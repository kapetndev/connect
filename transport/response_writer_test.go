package transport_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kapetndev/connect/transport"
)

func TestResponseWriter(t *testing.T) {
	t.Parallel()

	t.Run("captures the status code returned with the response", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := transport.NewResponseWriter(w)

		// Capture the status code.
		rw.WriteHeader(http.StatusOK)

		if rw.StatusCode() != http.StatusOK {
			t.Errorf("status codes are not equal: %d != %d", rw.StatusCode(), http.StatusOK)
		}
	})

	t.Run("captures the payload returned with the response", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := transport.NewResponseWriter(w)

		// Capture the payload.
		rw.Write([]byte("hello, world"))

		expectedPayload := []byte("hello, world")
		if !bytes.Equal(rw.Payload(), expectedPayload) {
			t.Errorf("payloads are not equal: %s != %s", rw.Payload(), expectedPayload)
		}
	})
}
