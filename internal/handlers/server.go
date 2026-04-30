package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/kamil7430/raspberry-voip/internal/state"
)

type server struct {
	state                                *state.State
	lastShowVerificationCodeRequest      time.Time
	lastShowVerificationCodeRequestMutex sync.Mutex
}

func NewHttpServer(state *state.State, addr string) *http.Server {
	s := server{
		state:                           state,
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
