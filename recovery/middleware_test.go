package recovery_test

import (
	"bytes"
	"context"
	stderrors "errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crumbandbase/expect"
	"github.com/crumbandbase/service-core-go/recovery"
	"github.com/crumbandbase/service-core-go/testdata/errors"
)

var (
	goodHTTPRequest     = []byte("good")
	panicHTTPRequest    = []byte("panic")
	nilPanicHTTPRequest = []byte("nilPanic")
)

type errorResponder func(*httptest.ResponseRecorder, *http.Request) error

func handler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	switch string(body) {
	case "panic":
		panic("very bad thing happened")
	case "nilPanic":
		panic(nil)
	}
}

func newRequest(t *testing.T, body []byte) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	return req
}

func setupRecoveryHandler(t *testing.T, opts ...recovery.Option) errorResponder {
	return func(w *httptest.ResponseRecorder, r *http.Request) error {
		mw := recovery.Handler(opts...)
		handlerFunc := mw(handler)

		handlerFunc(w, r)

		if w.Code != http.StatusOK {
			b, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}

			return stderrors.New(string(b))
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

		err := handlerFunc(w, r)
		expect.Equal(t, err, nil, "error was not <nil>")
	})

	t.Run("recovers and returns an error when the request panics", func(t *testing.T) {
		handlerFunc := setupRecoveryHandler(t)

		w := httptest.NewRecorder()
		r := newRequest(t, panicHTTPRequest)

		err := handlerFunc(w, r)
		expect.Equal(t, err.Error(), "internal server error: very bad thing happened\n", "messages are not equal")
	})

	t.Run("recovers and returns an error when the request causes a nil panic", func(t *testing.T) {
		handlerFunc := setupRecoveryHandler(t)

		w := httptest.NewRecorder()
		r := newRequest(t, nilPanicHTTPRequest)

		err := handlerFunc(w, r)
		expect.Equal(t, err.Error(), "internal server error: <nil>\n", "messages are not equal")
	})
}

func setupOverrideRecoveryHandler(t *testing.T) errorResponder {
	recoveryFunc := func(ctx context.Context, p interface{}) error {
		return errors.New("panic triggered: %v", p)
	}

	return setupRecoveryHandler(t, recovery.WithRecoveryContext(recoveryFunc))
}

func TestHandler_Override(t *testing.T) {
	t.Parallel()

	t.Run("does not invoke the recovery function when the request is successful", func(t *testing.T) {
		handlerFunc := setupOverrideRecoveryHandler(t)

		w := httptest.NewRecorder()
		r := newRequest(t, goodHTTPRequest)

		err := handlerFunc(w, r)
		expect.Equal(t, err, nil, "error was not <nil>")
	})

	t.Run("recovers and returns an error when the request panics", func(t *testing.T) {
		handlerFunc := setupOverrideRecoveryHandler(t)

		w := httptest.NewRecorder()
		r := newRequest(t, panicHTTPRequest)

		err := handlerFunc(w, r)
		expect.Equal(t, err.Error(), "panic triggered: very bad thing happened\n", "messages are not equal")
	})

	t.Run("recovers and returns an error when the request causes a nil panic", func(t *testing.T) {
		handlerFunc := setupOverrideRecoveryHandler(t)

		w := httptest.NewRecorder()
		r := newRequest(t, nilPanicHTTPRequest)

		err := handlerFunc(w, r)
		expect.Equal(t, err.Error(), "panic triggered: <nil>\n", "messages are not equal")
	})
}
