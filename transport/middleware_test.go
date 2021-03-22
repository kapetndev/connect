package transport_test

import (
	stderrors "errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crumbandbase/expect"
	"github.com/crumbandbase/service-core-go/testdata/errors"
	"github.com/crumbandbase/service-core-go/transport"
)

type calledResponder func(http.ResponseWriter, *http.Request) bool

func setupHandler(mw transport.Middleware) calledResponder {
	return func(w http.ResponseWriter, r *http.Request) (called bool) {
		handler := func(http.ResponseWriter, *http.Request) {
			called = true
		}

		handlerFunc := mw(handler)
		handlerFunc(w, r)
		return
	}
}

func TestChain(t *testing.T) {
	t.Parallel()

	t.Run("invokes the handler directly when no middlewares are given", func(t *testing.T) {
		handlerFunc := setupHandler(transport.Chain())

		w := httptest.NewRecorder()
		r := &http.Request{}

		called := handlerFunc(w, r)
		expect.Equal(t, called, true, "handler not invoked")
	})

	t.Run("passes through middlewares to the handler when middlewares are given", func(t *testing.T) {
		handlerFunc := setupHandler(transport.Chain(func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			}
		}))

		w := httptest.NewRecorder()
		r := &http.Request{}

		called := handlerFunc(w, r)
		expect.Equal(t, called, true, "handler not invoked")
	})
}

func TestWithContentType(t *testing.T) {
	t.Run("adds specificed content type to the response", func(t *testing.T) {
		handlerFunc := setupHandler(transport.WithJSONContentType)

		w := httptest.NewRecorder()
		r := &http.Request{}

		called := handlerFunc(w, r)
		expect.Equal(t, called, true, "handler not invoked")
		expect.Equal(t, w.Header()["Content-Type"][0], "application/json", "content type missing or incorrect")
	})
}

func TestWithHeaderValue(t *testing.T) {
	t.Run("add an arbitrary header to the response", func(t *testing.T) {
		handlerFunc := setupHandler(transport.WithHeaderValue("hello", "world"))

		w := httptest.NewRecorder()
		r := &http.Request{}

		called := handlerFunc(w, r)
		expect.Equal(t, called, true, "handler not invoked")
		expect.Equal(t, w.Header()["Hello"][0], "world", "header missing or incorrect")
	})
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
			b, err := ioutil.ReadAll(w.Body)
			if err != nil {
				return err
			}

			return stderrors.New(string(b))
		}

		return nil
	}
}

func TestWithError(t *testing.T) {
	t.Run("does not return an error when the request is successful", func(t *testing.T) {
		handlerFunc := setupErrorHander(nil)

		w := httptest.NewRecorder()
		r := &http.Request{}

		err := handlerFunc(w, r)
		expect.Equal(t, err, nil, "error is not <nil>")
	})

	t.Run("returns a default error when the request fails and the error does not have a responder", func(t *testing.T) {
		handlerFunc := setupErrorHander(stderrors.New("bad thing happened"))

		w := httptest.NewRecorder()
		r := &http.Request{}

		err := handlerFunc(w, r)
		expect.Equal(t, err.Error(), "bad thing happened\n", "error has incorrect value")
	})

	t.Run("returns an error when the request fails", func(t *testing.T) {
		handlerFunc := setupErrorHander(errors.New("custom bad thing happened"))

		w := httptest.NewRecorder()
		r := &http.Request{}

		err := handlerFunc(w, r)
		expect.Equal(t, err.Error(), "custom bad thing happened\n", "error has incorrect value")
	})
}
