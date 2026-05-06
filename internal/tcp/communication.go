package tcp

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/kamil7430/raspberry-voip/internal/audio"
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
			if err != nil {
				log.Printf("Error sending voice: %s\n", err)
				return
			}
		}
	}
}
