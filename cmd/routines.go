package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/kamil7430/raspberry-voip/internal/handlers"
)

func runHttpServer() {
	log.Println("Creating a web server instance...")
	server := handlers.NewServer(":2137")
	defer server.Close()

	log.Println("Running the web server...")
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
