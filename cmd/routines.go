package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/kamil7430/raspberry-voip/internal/handlers"
)

const httpServerAddr = ":2137"

func runHttpServer() {
	log.Println("Creating a web server instance...")
	server := handlers.NewServer(httpServerAddr)
	defer server.Close()

	log.Println("Running the web server...")
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
