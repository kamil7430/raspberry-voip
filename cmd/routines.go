package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/kamil7430/raspberry-voip/internal/audio"
	"github.com/kamil7430/raspberry-voip/internal/buttons"
	"github.com/kamil7430/raspberry-voip/internal/display"
	"github.com/kamil7430/raspberry-voip/internal/handlers"
	"github.com/kamil7430/raspberry-voip/internal/state"
	"github.com/kamil7430/raspberry-voip/internal/tcp"
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

func runListener(s *state.State, d *display.DisplayController, a *audio.AudioHandler) {
	log.Println("Starting tcp listener...")
	listener := tcp.NewListener(s, d, a)
	listener.Listen()
}

func runButtonHandler(s *state.State) {
	log.Println("Starting button handler...")
	buttonHandler := buttons.NewConcreteButtonHandler(s)
	buttonHandler.Start()
}
