package tcp

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"time"

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

	// our handshake with listener
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

	// wait for pickup
	timeoutTicker := time.NewTicker(dialingTime)
	defer timeoutTicker.Stop()
	select {
	case <-timeoutTicker.C:
		d.CallFinishedChan <- &display.CallFinishedDetails{
			Time:   time.Now(),
			Reason: "Dial timeout",
		}
		return nil
		// case reject button clicked
	}

	// main call loop
	for {
		select {
		case <-ctx.Done():
			d.CallFinishedChan <- &display.CallFinishedDetails{
				Time:   time.Now(),
				Reason: "Disconnected",
			}
			return nil
		default:
			receive(conn)
			send(conn)
		}
	}
}
