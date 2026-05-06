package tcp

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/kamil7430/raspberry-voip/internal/audio"
	"github.com/kamil7430/raspberry-voip/internal/state"
)

const (
	bufferSize  = 1024
	dialingTime = 15 * time.Second
)

func receiveAndPlay(conn net.Conn, audio *audio.AudioHandler, ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}

		receiveBuffer := make([]byte, bufferSize)

		n, err := conn.Read(receiveBuffer)
		if err != nil {
			log.Printf("Error reading from connection: %s\n", err)
			return
		}

		select {
		case <-ctx.Done():
			return
		case audio.In <- receiveBuffer[:n]:
			log.Printf("TCP send voice chunk")
		}
	}
}

func sendFromAudioBuffer(conn net.Conn, audio *audio.AudioHandler, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case out, ok := <-audio.Out:
			if !ok {
				return
			}

			_, err := conn.Write(out)
			log.Printf("TCP send voice chunk")
			if err != nil {
				log.Printf("Error sending voice: %s\n", err)
				return
			}
		}
	}
}

func handleRejectButtonClick(conn net.Conn, state *state.State, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case clickTime := <-state.RejectButtonClickChan:
			if clickTime.Add(500 * time.Millisecond).After(time.Now()) {
				_ = json.NewEncoder(conn).Encode(&finishCallMessage{Rejected: true})
				state.TerminateConnection()
			}
		}
	}
}

func listenForCallFinish(conn net.Conn, state *state.State) {
	err := json.NewDecoder(conn).Decode(&finishCallMessage{})
	if err != nil {
		log.Println(err)
	}
	state.TerminateConnection()
}
