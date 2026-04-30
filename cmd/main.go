package main

import (
	"log"

	"github.com/kamil7430/raspberry-voip/internal/display"
	"github.com/kamil7430/raspberry-voip/internal/state"
)

func main() {
	log.Println("Starting RPi VoIP!")

	s := state.NewState()
	d := display.NewDisplayController()

	go runHttpServer(&s, &d)
	runDisplayEventLoop(&d)
}
