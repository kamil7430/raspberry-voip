package tcp

import (
	"context"
	"log"
	"net"

	"github.com/kamil7430/raspberry-voip/internal/display"
	"github.com/kamil7430/raspberry-voip/internal/state"
)

func Dial(addr string, state *state.State, d *display.DisplayController) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancelFunc := context.WithCancel(nil)
	if !state.TrySetConnectionContext(ctx, cancelFunc) {
		log.Fatal("couldn't set connection context")
	}

	// TODO: dialing -- waiting for pick up

	for {
		select {
		case <-ctx.Done():
			d.CallFinishedChan <- &display.CallFinishedDetails{}
			return nil
		default:
			receive()
			send()
		}
	}
}
