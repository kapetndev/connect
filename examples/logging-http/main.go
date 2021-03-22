package main

import (
	"context"
	"encoding/hex"
	"net/http"

	"github.com/sirupsen/logrus"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"

	"github.com/crumbandbase/service-core-go/logging"
	logging_logrus "github.com/crumbandbase/service-core-go/logging/logrus"
	"github.com/crumbandbase/service-core-go/transport"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	mw := transport.Chain(
		logging.RequestLogger(
			logging_logrus.New(logger),
			logging.WithField(traceFunc),
		),
		transport.WithJSONContentType,
		transport.WithHeaderValue("Cache-Control", "no-cache"),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", mw(transport.WithError(handler)))
	mux.HandleFunc("/healthz", mw(healthz))

	server := &ochttp.Handler{
		Handler: mux,
	}

	logger.Info("server started on [::]:8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) error {
	l := logging.FromContext(r.Context())
	l.Info("logging from request handler")

	return nil
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func traceFunc(ctx context.Context) (string, interface{}) {
	span := trace.FromContext(ctx)
	id := span.SpanContext().TraceID

	return logging.TraceKey, hex.EncodeToString(id[:])
}
