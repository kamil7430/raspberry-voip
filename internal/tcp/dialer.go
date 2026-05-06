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

func Dial(addr string, state *state.State, d *display.DisplayController, a *audio.AudioHandler) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancelFunc := context.WithCancel(context.Background())
	if !state.TrySetConnectionContext(ctx, cancelFunc) {
		log.Fatal("couldn't set connection context")
	}
	defer cancelFunc()

	// our handshake with listener
	helloJson := helloMessage{
		DisplayName: state.GetDisplayName(),
	}
	err = json.NewEncoder(conn).Encode(helloJson)
	if err != nil {
		return err
	}

	var helloReceived helloMessage
	err = json.NewDecoder(conn).Decode(&helloReceived)
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
		// TODO: case reject button clicked
	}

	// main call loop
	go receiveAndPlay(conn, a)
	go sendFromAudioBuffer(conn, a)

	<-ctx.Done()

	d.CallFinishedChan <- &display.CallFinishedDetails{
		Time:   time.Now(),
		Reason: "Disconnected",
	}
	return nil
}
