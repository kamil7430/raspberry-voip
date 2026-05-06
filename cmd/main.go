package main

import (
	"log"

	"github.com/kamil7430/raspberry-voip/internal/audio"
	"github.com/kamil7430/raspberry-voip/internal/display"
	"github.com/kamil7430/raspberry-voip/internal/state"
)

const ( // TODO: correct?
	capCard    = 0
	capDevice  = 0
	playCard   = 0
	playDevice = 0
)

func main() {
	log.Println("Starting RPi VoIP!")

	s := state.NewState()
	d := display.NewDisplayController()
	a := audio.NewAudioHandler(capCard, capDevice, playCard, playDevice)

	go runHttpServer(s, d)
	go runDisplayEventLoop(d)
	runButtonHandler(s, d, a)
	runListener(s, d, a)
}
