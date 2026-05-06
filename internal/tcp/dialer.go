package tcp

import (
	"context"
	"encoding/json"
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

	helloJson, err := json.Marshal(helloMessage{
		DisplayName: state.GetDisplayName(),
	})
	if err != nil {
		return err
	}
	_, err = conn.Write(helloJson)
	if err != nil {
		return err
	}

	buffer := make([]byte, bufferSize)
	_, err = conn.Read(buffer)
	if err != nil {
		return err
	}
	var helloReceived helloMessage
	err = json.Unmarshal(buffer, &helloReceived)
	if err != nil {
		return err
	}
	displayName := helloReceived.DisplayName

	d.DialingChan <- &display.DialingDetails{
		DisplayName: displayName,
	}

	// TODO: wait for pickup

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
