package handlers

import (
	"net/http"

	"github.com/kamil7430/raspberry-voip/internal/state"
)

type server struct {
	state *state.State
}

func NewHttpServer(state *state.State, addr string) *http.Server {
	s := server{
		state,
	}

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/", s.configHandler)
	serveMux.HandleFunc("/save-config", s.saveConfigHandler)

	return &http.Server{
		Addr:    addr,
		Handler: serveMux,
	}
}
