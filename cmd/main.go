package main

import (
	"log"

	"github.com/kamil7430/raspberry-voip/internal/audio"
	"github.com/kamil7430/raspberry-voip/internal/display"
	"github.com/kamil7430/raspberry-voip/internal/state"
)

func main() {
	log.Println("Starting RPi VoIP!")

	s := state.NewState()
	d := display.NewDisplayController()
	a := audio.NewAudioHandler()

	go runHttpServer(s, d)
	go runDisplayEventLoop(d)
	runButtonHandler(s, d, a)
	runListener(s, d, a)
}
