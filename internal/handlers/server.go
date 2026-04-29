package handlers

import "net/http"

func NewServer(addr string) *http.Server {
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/", configHandler)

	return &http.Server{
		Addr:    addr,
		Handler: serveMux,
	}
}
