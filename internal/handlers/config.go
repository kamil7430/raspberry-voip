package handlers

import (
	"fmt"
	"net/http"
)

func configHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world!")
}
