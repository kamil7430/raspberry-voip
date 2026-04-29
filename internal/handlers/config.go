package handlers

import (
	"log"
	"net/http"

	"github.com/kamil7430/raspberry-voip/web"
)

func (s *server) configHandler(w http.ResponseWriter, r *http.Request) {
	err := web.Templates.Execute(w, "config.html")
	if err != nil {
		log.Fatal(err)
	}
}
