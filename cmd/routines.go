package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/kamil7430/raspberry-voip/internal/display"
	"github.com/kamil7430/raspberry-voip/internal/handlers"
	"github.com/kamil7430/raspberry-voip/internal/state"
)

const httpServerAddr = ":2137"

func runHttpServer(state *state.State, d *display.DisplayController) {
	log.Println("Creating a web server instance...")
	server := handlers.NewHttpServer(state, httpServerAddr, d)
	defer server.Close()

	log.Println("Running the web server...")
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func runDisplayEventLoop(d *display.DisplayController) {
	log.Println("Starting display event loop...")
	d.EventLoop()
}
