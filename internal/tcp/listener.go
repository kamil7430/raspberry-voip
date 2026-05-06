package tcp

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/kamil7430/raspberry-voip/internal/audio"
	"github.com/kamil7430/raspberry-voip/internal/display"
	"github.com/kamil7430/raspberry-voip/internal/state"
)

const listenerAddr = ":8080"

type Listener struct {
	state   *state.State
	display *display.DisplayController
	audio   *audio.AudioHandler
}

func NewListener(state *state.State, d *display.DisplayController, a *audio.AudioHandler) Listener {
	return Listener{
		state:   state,
		display: d,
		audio:   a,
	}
}

func (l *Listener) Listen() {
	listener, err := net.Listen("tcp", listenerAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancelFunc := context.WithCancel(nil)
		if l.state.TrySetConnectionContext(ctx, cancelFunc) {
			go l.handleConnection(conn, ctx)
		} else {
			conn.Close()
		}
	}
}

func (l *Listener) handleConnection(conn net.Conn, ctx context.Context) {
	defer conn.Close()

	// our handshake with dialer
	buffer := make([]byte, bufferSize)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Printf("Error reading from connection: %s\n", err)
		return
	}
	var helloReceived helloMessage
	err = json.Unmarshal(buffer, &helloReceived)
	if err != nil {
		log.Printf("Error unmarshalling the hello message: %s\n", err)
		return
	}
	displayName := helloReceived.DisplayName

	helloJson, err := json.Marshal(helloMessage{
		DisplayName: l.state.GetDisplayName(),
	})
	if err != nil {
		log.Printf("Error marshalling hello message: %s\n", err)
		return
	}
	_, err = conn.Write(helloJson)
	if err != nil {
		log.Printf("Error sending hello message: %s\n", err)
		return
	}

	l.display.IncomingCallChan <- &display.IncomingCallDetails{
		DisplayName: displayName,
	}

	// wait for pickup
	timeoutTicker := time.NewTicker(dialingTime)
	defer timeoutTicker.Stop()
	select {
	case <-timeoutTicker.C:
		l.display.CallFinishedChan <- &display.CallFinishedDetails{
			Time:   time.Now(),
			Reason: "Dial timeout",
		}
		return
		// case reject button clicked
		// case accept button clicked
	}

	// main call loop
	for {
		select {
		case <-ctx.Done():
			l.display.CallFinishedChan <- &display.CallFinishedDetails{
				Time:   time.Now(),
				Reason: "Disconnected",
			}
			return
		default:
			sendFromAudioBuffer(conn, l.audio)
			receiveAndPlay(conn, l.audio)
		}
	}
}
