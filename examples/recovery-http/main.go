package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/kapetndev/connect/recovery"
	"github.com/kapetndev/connect/transport"
)

func main() {
	recoveryFn := func(ctx context.Context, p interface{}) error {
		return fmt.Errorf("recovered from panic: %v", p)
	}

	mw := transport.Chain(
		recovery.Handler(
			recovery.WithRecoveryContext(recoveryFn),
		),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", mw(handler))
	mux.HandleFunc("/healthz", mw(healthz))

	log.Println("server started on [::]:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func handler(http.ResponseWriter, *http.Request) {
	panic("something bad happened")
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
