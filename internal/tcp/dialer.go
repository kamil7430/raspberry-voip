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
		return err
	}
	defer state.TerminateConnection()

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

	answeredChan := make(chan bool)
	go func() {
		var answer callAnswerMessage
		err := json.NewDecoder(conn).Decode(&answer)
		if err != nil {
			log.Println("Error decoding answer")
		}
		answeredChan <- answer.Answered
	}()

	shouldProceed := false
	for !shouldProceed {
		select {
		case <-timeoutTicker.C:
			d.CallFinishedChan <- &display.CallFinishedDetails{
				Time:   time.Now(),
				Reason: "Dial timeout",
			}
			return nil
		case clickTime := <-state.RejectButtonClickChan:
			if clickTime.Add(500 * time.Millisecond).After(time.Now()) {
				d.CallFinishedChan <- &display.CallFinishedDetails{
					Time:   time.Now(),
					Reason: "Dial cancelled",
				}
				return nil
			}
		case answer := <-answeredChan:
			if answer {
				shouldProceed = true
			} else {
				d.CallFinishedChan <- &display.CallFinishedDetails{
					Time:   time.Now(),
					Reason: "Dial rejected",
				}
				return nil
			}
		}
	}

	d.InCallChan <- &display.InCallDetails{
		DisplayName: displayName,
		CallStart:   time.Now(),
	}
	a.Start(ctx)

	// main call loop
	go receiveAndPlay(conn, a, ctx)
	go sendFromAudioBuffer(conn, a, ctx)
	go handleRejectButtonClick(conn, state, ctx)
	go listenForCallFinish(conn, state)

	<-ctx.Done()

	d.CallFinishedChan <- &display.CallFinishedDetails{
		Time:   time.Now(),
		Reason: "Disconnected",
	}
	return nil
}
