package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/kamil7430/raspberry-voip/internal/handlers"
	"github.com/kamil7430/raspberry-voip/internal/state"
)

const httpServerAddr = ":2137"

func runHttpServer(state *state.State) {
	log.Println("Creating a web server instance...")
	server := handlers.NewHttpServer(state, httpServerAddr)
	defer server.Close()

	log.Println("Running the web server...")
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
