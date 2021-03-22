package recovery_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kapetndev/connect/recovery"
)

var (
	goodHTTPRequest     = []byte("good")
	panicHTTPRequest    = []byte("panic")
	nilPanicHTTPRequest = []byte("nilPanic")
)

type errorResponder func(*httptest.ResponseRecorder, *http.Request) error

func handler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	returnPanics(string(body))
}

func newRequest(t *testing.T, body []byte) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to create request: %s", err)
	}
	return req
}

func setupRecoveryHandler(t *testing.T, opts ...recovery.Option) errorResponder {
	return func(w *httptest.ResponseRecorder, r *http.Request) error {
		mw := recovery.Handler(opts...)
		handlerFunc := mw(handler)

		handlerFunc(w, r)

		if w.Code != http.StatusOK {
			b, err := io.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}

			return errors.New(string(b))
		}

		return nil
	}
}

func TestHandler_Default(t *testing.T) {
	t.Parallel()

	t.Run("does not invoke the recovery function when the request is successful", func(t *testing.T) {
		handlerFunc := setupRecoveryHandler(t)

		w := httptest.NewRecorder()
		r := newRequest(t, goodHTTPRequest)

		if err := handlerFunc(w, r); err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}
	})

	t.Run("recovers and returns an error when the request panics", func(t *testing.T) {
		handlerFunc := setupRecoveryHandler(t)

		w := httptest.NewRecorder()
		r := newRequest(t, panicHTTPRequest)

		expectedErrorMessage := errorMessage("internal server error", panicMessage)
		if err := handlerFunc(w, r); err.Error() != expectedErrorMessage {
			t.Errorf("messages are not equal: %s != %s", err, expectedErrorMessage)
		}
	})

	t.Run("recovers and returns an error when the request causes a nil panic", func(t *testing.T) {
		handlerFunc := setupRecoveryHandler(t)

		w := httptest.NewRecorder()
		r := newRequest(t, nilPanicHTTPRequest)

		expectedErrorMessage := errorMessage("internal server error", "<nil>")
		if err := handlerFunc(w, r); err.Error() != expectedErrorMessage {
			t.Errorf("messages are not equal: %s != %s", err, expectedErrorMessage)
		}
	})
}

func setupOverrideRecoveryHandler(t *testing.T) errorResponder {
	recoveryFunc := func(ctx context.Context, p interface{}) error {
		return fmt.Errorf("panic triggered: %v", p)
	}

	return setupRecoveryHandler(t, recovery.WithRecoveryContext(recoveryFunc))
}

func TestHandler_Override(t *testing.T) {
	t.Parallel()

	t.Run("does not invoke the recovery function when the request is successful", func(t *testing.T) {
		handlerFunc := setupOverrideRecoveryHandler(t)

		w := httptest.NewRecorder()
		r := newRequest(t, goodHTTPRequest)

		if err := handlerFunc(w, r); err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}
	})

	t.Run("recovers and returns an error when the request panics", func(t *testing.T) {
		handlerFunc := setupOverrideRecoveryHandler(t)

		w := httptest.NewRecorder()
		r := newRequest(t, panicHTTPRequest)

		expectedErrorMessage := errorMessage("panic triggered", panicMessage)
		if err := handlerFunc(w, r); err.Error() != expectedErrorMessage {
			t.Errorf("messages are not equal: %s != %s", err, expectedErrorMessage)
		}
	})

	t.Run("recovers and returns an error when the request causes a nil panic", func(t *testing.T) {
		handlerFunc := setupOverrideRecoveryHandler(t)

		w := httptest.NewRecorder()
		r := newRequest(t, nilPanicHTTPRequest)

		expectedErrorMessage := errorMessage("panic triggered", "<nil>")
		if err := handlerFunc(w, r); err.Error() != expectedErrorMessage {
			t.Errorf("messages are not equal: %s != %s", err, expectedErrorMessage)
		}
	})
}

func errorMessage(prefix, message string) string {
	return prefix + ": " + message + "\n"
}
