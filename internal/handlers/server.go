package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/kamil7430/raspberry-voip/internal/display"
	"github.com/kamil7430/raspberry-voip/internal/state"
)

type server struct {
	state                                *state.State
	display                              *display.DisplayController
	lastShowVerificationCodeRequest      time.Time
	lastShowVerificationCodeRequestMutex sync.Mutex
}

func NewHttpServer(state *state.State, addr string, d *display.DisplayController) *http.Server {
	s := server{
		state:                           state,
		display:                         d,
		lastShowVerificationCodeRequest: time.Now(),
	}

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/config", s.configHandler)
	serveMux.HandleFunc("/save-config", s.saveConfigHandler)
	serveMux.HandleFunc("/show-verification-code", s.showVerificationCode)

	return &http.Server{
		Addr:    addr,
		Handler: serveMux,
	}
}
