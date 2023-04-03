package main

import (
	"net/http"
	"os"

	"golang.org/x/exp/slog"

	"github.com/kapetndev/connect/logging"
	"github.com/kapetndev/connect/transport"
)

func main() {
	h := logging.NewGoogleCloudHandler(os.Stdout, slog.LevelDebug)

	mw := transport.Chain(
		logging.RequestLogger(
			logging.WithHandler(h),
		),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", mw(handler))
	mux.HandleFunc("/healthz", mw(healthz))

	logger := slog.New(h)

	logger.Info("server started on [::]:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Error("failed to serve: " + err.Error())
	}
}

func handler(_ http.ResponseWriter, r *http.Request) {
	l := logging.FromContext(r.Context())
	l.Info(r.Context(), "logging from request handler")
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
