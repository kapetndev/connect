package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/crumbandbase/service-core-go/recovery"
	"github.com/crumbandbase/service-core-go/transport"
)

type picardError string

func (e picardError) Error() string {
	return string(e)
}

func (e picardError) RespondError(w http.ResponseWriter, r *http.Request) bool {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, e.Error())
	return true
}

func main() {
	recoveryFn := func(ctx context.Context, p interface{}) error {
		return picardError(fmt.Sprintf("recovered from panic: %v", p))
	}

	mw := transport.Chain(
		recovery.Handler(
			recovery.WithRecoveryContext(recoveryFn),
		),
		transport.WithJSONContentType,
		transport.WithHeaderValue("Cache-Control", "no-cache"),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", mw(transport.WithError(handler)))
	mux.HandleFunc("/healthz", mw(healthz))

	log.Println("server started on [::]:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func handler(http.ResponseWriter, *http.Request) error {
	panic("shut up, Wesley!")
	return nil
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
