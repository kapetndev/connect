package transport_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crumbandbase/expect"
	"github.com/crumbandbase/service-core-go/transport"
)

func TestResponseWriter(t *testing.T) {
	t.Run("captures the status code returned with the response", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := transport.NewResponseWriter(w)

		// Capture the status code.
		rw.WriteHeader(http.StatusOK)

		expect.Equal(t, rw.StatusCode(), http.StatusOK, "status has incorrect value")
	})

	t.Run("captures the payload returned with the response", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := transport.NewResponseWriter(w)

		// Capture the payload.
		rw.Write([]byte("hello, world"))

		expect.Equal(t, rw.Payload(), []byte("hello, world"), "payload has incorrect value")
	})
}
