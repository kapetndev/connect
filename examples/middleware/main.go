package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/kapetndev/connect/transport"
)

func main() {
	mw := transport.Chain(
		transport.WithJSONContentType,
		transport.WithHeaderValue("Cache-Control", "no-cache"),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", mw(transport.WithError(handler)))
	mux.HandleFunc("/error", mw(transport.WithError(handlerWithError)))
	mux.HandleFunc("/healthz", mw(healthz))

	log.Println("server started on [::]:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func handler(w http.ResponseWriter, _ *http.Request) error {
	return json.NewEncoder(w).Encode(map[string]string{"hello": "world"})
}

func handlerWithError(w http.ResponseWriter, r *http.Request) error {
	return errors.New("something bad happened")
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
