package recovery

import (
	"fmt"
	"net/http"

	"github.com/kapetndev/connect/transport"
)

// Handler returns a middleware for panic recovery.
func Handler(opts ...Option) transport.Middleware {
	o := applyOptions(opts)
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			panicked := true

			defer func() {
				var p interface{}

				if p = recover(); p == nil && !panicked {
					return
				}

				err := o.recovery(r.Context(), p)
				if err == nil {
					err = fmt.Errorf("internal server error: %v", p)
				}

				res, ok := err.(transport.ErrorResponder)
				if ok && res.RespondError(w, r) {
					return
				}

				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, err.Error())
			}()

			// If the call to the handler does not result in a panic then the code
			// following it will be executed. Therefore if the handler panics then
			// `panicked` will not be set to false. This is essential to be able to
			// detect nil panics from the handler.
			next.ServeHTTP(w, r)
			panicked = false
		}
	}
}
