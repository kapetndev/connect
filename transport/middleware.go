package transport

import (
	"fmt"
	"net/http"
)

// Middleware is a function type alias representing a single item in a
// middleware stack.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain returns a function taking a HandlerFunc. This function iterates
// through a variadic list of HandlerFunc, representing the middleware stack,
// passing to each a reference to the previous function.
func Chain(fns ...Middleware) Middleware {
	return func(route http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			handler := route
			for _, fn := range fns {
				handler = fn(handler)
			}

			handler(w, r)
		}
	}
}

// HeaderValue is a HTTP middleware setting an aritrary header value.
type HeaderValue struct {
	handler     http.HandlerFunc
	headerName  string
	headerValue string
}

func (h HeaderValue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(http.CanonicalHeaderKey(h.headerName), h.headerValue)
	h.handler.ServeHTTP(w, r)
}

// WithJSONContentType returns a middlewate that sets the Content-Type HTTP
// header on the response with a value indicating the payload is JSON.
func WithJSONContentType(next http.HandlerFunc) http.HandlerFunc {
	return WithHeaderValue("content-type", "application/json")(next)
}

// WithHeaderValue returns a middleware that sets an arbitrary HTTP header on
// the response with a given value.
func WithHeaderValue(key string, value string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(http.CanonicalHeaderKey(key), value)
			next.ServeHTTP(w, r)
		}
	}
}

// ErrorResponder describes how to write an error message to `w`.
type ErrorResponder interface {
	RespondError(w http.ResponseWriter, r *http.Request) bool
}

// ErrorHandlerFunc describes a HTTP handler that returns an error.
type ErrorHandlerFunc func(http.ResponseWriter, *http.Request) error

// WithError is a wrapper around a handler function that delegates error
// reporting to the returned error value.
func WithError(h ErrorHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			return
		}

		res, ok := err.(ErrorResponder)
		if ok && res.RespondError(w, r) {
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
	}
}
