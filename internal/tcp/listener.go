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
			log.Println(err)
			continue
		}

		ctx, cancelFunc := context.WithCancel(context.Background())
		if l.state.TrySetConnectionContext(ctx, cancelFunc) {
			go l.handleConnection(conn, ctx)
		} else {
			conn.Close()
		}
	}
}

func (l *Listener) handleConnection(conn net.Conn, ctx context.Context) {
	defer conn.Close()
	defer l.state.TerminateConnection()

	// our handshake with dialer
	var helloReceived helloMessage
	err := json.NewDecoder(conn).Decode(&helloReceived)
	if err != nil {
		log.Println(err)
		return
	}
	displayName := helloReceived.DisplayName

	helloJson := helloMessage{
		DisplayName: l.state.GetDisplayName(),
	}
	err = json.NewEncoder(conn).Encode(helloJson)
	if err != nil {
		log.Println(err)
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
		// TODO: case reject button clicked
		// TODO: case accept button clicked
	}

	// main call loop
	go sendFromAudioBuffer(conn, l.audio, ctx)
	go receiveAndPlay(conn, l.audio, ctx)

	<-ctx.Done()

	l.display.CallFinishedChan <- &display.CallFinishedDetails{
		Time:   time.Now(),
		Reason: "Disconnected",
	}
}
