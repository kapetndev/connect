package logging

import (
	"net/http"
	"time"

	"github.com/crumbandbase/service-core-go/transport"
)

// RequestLogger returns a middleware that instantiates a request scoped logger
// and injects it into the request context to be used within subsequent
// middlewares and handlers.
func RequestLogger(log Logger, opts ...Option) transport.Middleware {
	o := applyOptions(opts)
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			ctx := r.Context()
			path := getPath(r)

			// Configure the logger passed into the middleware. This function will
			// always return a new logger and is therefore thread safe.
			entry := newRequestLogger(ctx, log, o.fields, r.Method, path, startTime, o.timestampFormat)

			// Wrap the response writer so we may capture the status code and payload
			// from the handler.
			rw := transport.NewResponseWriter(w)

			// Invoke the hander and log the response.
			next.ServeHTTP(rw, r.WithContext(NewContext(ctx, entry)))

			// Suppress request logs matching some pattern.
			if o.shouldDiscard(ctx, path, nil) {
				return
			}

			logRequest(entry, rw, time.Since(startTime))
		}
	}
}

func getPath(r *http.Request) string {
	if r.URL == nil {
		return ""
	}
	return r.URL.Path
}

func logRequest(log Logger, w *transport.ResponseWriter, duration time.Duration) {
	statusCode := w.StatusCode()

	log = log.WithFields(Fields{
		StatusKey:   statusCode,
		DurationKey: duration,
	})

	// Everything with a status code of 400 or greater is considered an error.
	if statusCode < 400 {
		entryWithResponseFields(log, w.Payload(), PayloadKey).Info()
	} else {
		entryWithResponseFields(log, w.Payload(), PayloadKey).Error()
	}
}

func entryWithResponseFields(log Logger, msg []byte, key string) Logger {
	if msg != nil {
		return log.WithField(key, byteSliceMarshallable(msg))
	}

	return log
}

// byteSliceMarshallable is a wrapper type allowing us to embed a JSON object
// within a log entry. Without this the logger will return the raw bytes.
type byteSliceMarshallable []byte

// MarshalJSON simply returns the underlying byte slice.
func (b byteSliceMarshallable) MarshalJSON() ([]byte, error) {
	return b, nil
}
