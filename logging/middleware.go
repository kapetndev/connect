package logging

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/exp/slog"

	"github.com/kapetndev/connect/transport"
)

// RequestLogger returns a middleware that instantiates a request scoped logger
// and injects it into the request context to be used within subsequent
// middlewares and handlers.
func RequestLogger(opts ...Option) transport.Middleware {
	o := applyOptions(opts)
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			ctx := r.Context()

			// Configure the logger passed into the middleware.
			logger := New(o.handler)

			// Wrap the response writer so we may capture the status code and payload
			// from the handler.
			rw := transport.NewResponseWriter(w)

			// Invoke the hander and log the response.
			next.ServeHTTP(rw, r.WithContext(NewContext(ctx, logger)))

			// Suppress request logs matching some pattern.
			if o.shouldDiscard(ctx, r.URL.Path, nil) {
				return
			}

			// Log the request/response.
			o.handler.Handle(ctx, newRequestRecord(ctx, startTime, rw, r))
		}
	}
}

func newRequestRecord(ctx context.Context, t time.Time, rw *transport.ResponseWriter, r *http.Request) slog.Record {
	statusCode := rw.StatusCode()

	level := slog.LevelInfo
	if statusCode >= http.StatusBadRequest {
		level = slog.LevelError
	}

	record := newCommonRecord(ctx, level, t, r.Method, r.URL.Path)

	record.AddAttrs(slog.Int(StatusKey, statusCode))

	// If the response includes a payload then add it to the log entry. This
	// assumes that the payload is a JSON object.
	if rw.Payload() != nil {
		record.AddAttrs(slog.Any(ResponseKey, byteSliceMarshallable(rw.Payload())))
	}

	return record
}

// byteSliceMarshallable is a wrapper type allowing us to embed a JSON object
// within a log entry. Without this the logger will return the raw bytes.
type byteSliceMarshallable []byte

// MarshalJSON simply returns the underlying byte slice.
func (b byteSliceMarshallable) MarshalJSON() ([]byte, error) {
	return b, nil
}
