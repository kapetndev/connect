package transport_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kapetndev/connect/transport"
)

type calledResponder func(http.ResponseWriter, *http.Request) bool

func setupHandler(mw transport.Middleware) calledResponder {
	return func(w http.ResponseWriter, r *http.Request) (called bool) {
		handler := func(http.ResponseWriter, *http.Request) {
			called = true
		}

		mw(handler)(w, r)
		return
	}
}

func TestChain(t *testing.T) {
	t.Parallel()

	t.Run("invokes the handler directly when no middlewares are given", func(t *testing.T) {
		handlerFunc := setupHandler(transport.Chain())

		w := httptest.NewRecorder()
		r := &http.Request{}

		if called := handlerFunc(w, r); !called {
			t.Error("handler not invoked")
		}
	})

	t.Run("passes through middlewares to the handler when middlewares are given", func(t *testing.T) {
		handlerFunc := setupHandler(transport.Chain(func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			}
		}))

		w := httptest.NewRecorder()
		r := &http.Request{}

		if called := handlerFunc(w, r); !called {
			t.Error("handler not invoked")
		}
	})
}

func TestWithContentType(t *testing.T) {
	t.Parallel()

	t.Run("adds specificed content type to the response", func(t *testing.T) {
		handlerFunc := setupHandler(transport.WithJSONContentType)

		w := httptest.NewRecorder()
		r := &http.Request{}

		if called := handlerFunc(w, r); !called {
			t.Error("handler not invoked")
		}

		expectedContentType := "application/json"
		if w.Header()["Content-Type"][0] != expectedContentType {
			t.Errorf("content types are not equal: %s != %s", w.Header()["Content-Type"][0], expectedContentType)
		}
	})
}

func TestWithHeaderValue(t *testing.T) {
	t.Parallel()

	t.Run("add an arbitrary header to the response", func(t *testing.T) {
		handlerFunc := setupHandler(transport.WithHeaderValue("hello", "world"))

		w := httptest.NewRecorder()
		r := &http.Request{}

		if called := handlerFunc(w, r); !called {
			t.Error("handler not invoked")
		}

		expectedHeaderValue := "world"
		if w.Header()["Hello"][0] != expectedHeaderValue {
			t.Errorf("header values are not equal: %s != %s", w.Header()["Hello"][0], expectedHeaderValue)
		}
	})
}

type customError string

func (e customError) Error() string {
	return string(e)
}

func (e customError) RespondError(w http.ResponseWriter, r *http.Request) bool {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, string(e))
	return true
}

type errorResponder func(*httptest.ResponseRecorder, *http.Request) error

func setupErrorHander(err error) errorResponder {
	return func(w *httptest.ResponseRecorder, r *http.Request) error {
		handler := func(w http.ResponseWriter, r *http.Request) error {
			return err
		}

		handlerFunc := transport.WithError(handler)
		handlerFunc(w, r)

		if w.Code != http.StatusOK {
			b, err := io.ReadAll(w.Body)
			if err != nil {
				return err
			}

			return errors.New(string(b))
		}

		return nil
	}
}

func TestWithError(t *testing.T) {
	t.Parallel()

	t.Run("does not return an error when the request is successful", func(t *testing.T) {
		handlerFunc := setupErrorHander(error(nil))

		w := httptest.NewRecorder()
		r := &http.Request{}

		if err := handlerFunc(w, r); err != nil {
			t.Error("error is not <nil>")
		}
	})

	t.Run("returns a default error when the request fails and the error does not have a responder", func(t *testing.T) {
		handlerFunc := setupErrorHander(errors.New("something bad happened"))

		w := httptest.NewRecorder()
		r := &http.Request{}

		expectedErrorMesssage := errorMessage("something bad happened")
		if err := handlerFunc(w, r); err.Error() != expectedErrorMesssage {
			t.Errorf("messages are not equal: %s != %s", err.Error(), expectedErrorMesssage)
		}
	})

	t.Run("returns an error when the request fails and the error does have a responder", func(t *testing.T) {
		handlerFunc := setupErrorHander(customError("something custom and bad happened"))

		w := httptest.NewRecorder()
		r := &http.Request{}

		expectedErrorMesssage := errorMessage("something custom and bad happened")
		if err := handlerFunc(w, r); err.Error() != expectedErrorMesssage {
			t.Errorf("messages are not equal: %s != %s", err.Error(), expectedErrorMesssage)
		}
	})
}

func errorMessage(message string) string {
	return message + "\n"
}
