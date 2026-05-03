package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/kamil7430/raspberry-voip/internal/display"
	"github.com/kamil7430/raspberry-voip/web"
)

const showVerificationCodeTimeout = 5 * time.Second

type configPageData struct {
	DisplayName string
}

func (s *server) configHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metoda niedozwolona", http.StatusMethodNotAllowed)
		return
	}

	pageData := configPageData{
		DisplayName: s.state.GetDisplayName(),
	}

	err := web.Templates.ExecuteTemplate(w, "config.html", pageData)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *server) saveConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metoda niedozwolona", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Błąd przetwarzania formularza", http.StatusBadRequest)
	}

	verificationCode := r.FormValue("verificationCode")
	err := s.state.CheckAndConsumeVerificationCode(verificationCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	displayName := r.FormValue("displayName")
	err = s.state.SetDisplayName(displayName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("User changed display name: %s\n", displayName)

	http.Redirect(w, r, "/config", http.StatusSeeOther)
}

func (s *server) showVerificationCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metoda niedozwolona", http.StatusMethodNotAllowed)
		return
	}

	s.lastShowVerificationCodeRequestMutex.Lock()
	defer s.lastShowVerificationCodeRequestMutex.Unlock()

	timeNow := time.Now()

	if s.lastShowVerificationCodeRequest.Add(showVerificationCodeTimeout).After(timeNow) {
		http.Error(w, "Nowe żądanie wysłane zbyt szybko", http.StatusTooManyRequests)
		return
	}

	select {
	case s.display.ShowVerificationCodeChan <- &display.ShowVerificationCodeDetails{
		Time: timeNow,
		Code: s.state.CreateVerificationCode(),
	}:
		log.Println("Sent verification code show request to display")
	default:
		log.Fatal("The channel is not empty!")
	}

	s.lastShowVerificationCodeRequest = timeNow

	w.WriteHeader(http.StatusOK)
}
