package transport

import "net/http"

// ResponseWriter is used by a HTTP handler to construct a HTTP response. Both
// the status code and payload are captured by this type.
type ResponseWriter struct {
	http.ResponseWriter

	// Captured values.
	statusCode int
	payload    []byte
}

// NewResponseWriter returns a new ResponseWriter.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

// WriteHeader sends a HTTP response header with the provided status code.
func (w *ResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Write writes the data to the connection as part of a HTTP reply.
func (w *ResponseWriter) Write(payload []byte) (int, error) {
	w.payload = payload
	return w.ResponseWriter.Write(payload)
}

// StatusCode returns the status code last written to the writer.
func (w *ResponseWriter) StatusCode() int {
	return w.statusCode
}

// Payload returns the payload written to the writer.
func (w *ResponseWriter) Payload() []byte {
	return w.payload
}
