package tcp

import (
	"context"
	"log"
	"net"

	"github.com/kamil7430/raspberry-voip/internal/display"
	"github.com/kamil7430/raspberry-voip/internal/state"
)

const listenerAddr = ":8080"

type Listener struct {
	state   *state.State
	display *display.DisplayController
}

func NewListener(state *state.State, d *display.DisplayController) Listener {
	return Listener{
		state:   state,
		display: d,
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
	for {
		select {
		case <-ctx.Done():
			l.display.CallFinishedChan <- &display.CallFinishedDetails{}
			conn.Close()
			return
		default:
			send()
			receive()
		}
	}
}
