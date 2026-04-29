package main

import (
	"log"

	"github.com/kamil7430/raspberry-voip/internal/state"
)

func main() {
	log.Println("Starting RPi VoIP!")

	s := state.NewState()

	runHttpServer(&s) // will be a go routine
}
