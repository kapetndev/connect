package recovery

import (
	"context"
	"fmt"
	"net/http"

	"github.com/crumbandbase/service-core-go/transport"
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

				err := recoverHandler(r.Context(), p, o.recovery)

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

func recoverHandler(ctx context.Context, p interface{}, fn RecoveryContextFunc) error {
	err := fn(ctx, p)

	if err == nil {
		return fmt.Errorf("internal server error: %v", p)
	}

	return err
}
