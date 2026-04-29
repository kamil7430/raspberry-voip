package handlers

import (
	"log"
	"net/http"

	"github.com/kamil7430/raspberry-voip/web"
)

type configPageData struct {
	DisplayName string
}

func (s *server) configHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metoda niedozwolona", http.StatusMethodNotAllowed)
		return
	}

	pageData := configPageData{
		s.state.GetDisplayName(),
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

	displayName := r.FormValue("displayName")
	err := s.state.SetDisplayName(displayName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("User changed display name: %s\n", displayName)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
