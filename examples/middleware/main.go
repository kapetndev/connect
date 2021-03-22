package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/crumbandbase/errors"
	"github.com/crumbandbase/service-core-go/transport"
	"github.com/crumbandbase/service-core-go/transport/sentinels"
)

type badRequestError struct {
	err error
}

func (e badRequestError) Error() string {
	return "bad request: " + e.err.Error()
}

func (e badRequestError) RespondError(w http.ResponseWriter, r *http.Request) bool {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, e.Error())
	return true
}

func main() {
	mw := transport.Chain(
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

func handler(w http.ResponseWriter, r *http.Request) error {
	starship := struct {
		Name string
	}{
		Name: "Enterprise",
	}

	if err := json.NewEncoder(w).Encode(starship); err != nil {
		return badRequestError{errors.Wrap(sentinels.MarshalError, err)}
	}

	return nil
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
