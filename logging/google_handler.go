package logging

import (
	"context"
	"io"

	"golang.org/x/exp/slog"
)

// Google Cloud Logging specific attributes.
// https://cloud.google.com/logging/docs/agent/logging/configuration#process-payload
const (
	googleCloudLabelsKey         = "logging.googleapis.com/labels"
	googleCloudMessageKey        = "message"
	googleCloudMethodKey         = "requestMethod"
	googleCloudSeverityKey       = "severity"
	googleCloudSourceLocationKey = "logging.googleapis.com/sourceLocation"
	googleCloudSpanKey           = "logging.googleapis.com/spanId"
	googleCloudTraceKey          = "logging.googleapis.com/trace"
	googleCloudResponseSizeKey   = "responseSize"
)

// GoogleCloudHandler is a handler that formats log messages in a way that is
// compatible with Google Cloud Logging.
type GoogleCloudHandler struct {
	handler      slog.Handler
	SpanHandler  AttrHandler
	TraceHandler AttrHandler
}

// NewGoogleCloudHandler returns a new
func NewGoogleCloudHandler(w io.Writer, level slog.Level) *GoogleCloudHandler {
	return &GoogleCloudHandler{
		handler: slog.HandlerOptions{
			Level: level,

			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if len(a.Value.String()) == 0 {
					a.Key = "" // Drop empty attributes.
					return a
				}

				switch a.Key {
				case slog.LevelKey:
					a.Key = googleCloudSeverityKey
					a.Value = severityValue(a.Value)
				case slog.MessageKey:
					a.Key = googleCloudMessageKey
				case slog.SourceKey:
					a.Key = googleCloudSourceLocationKey
				case MethodKey:
					a.Key = googleCloudMethodKey
				}

				return a
			},
		}.NewJSONHandler(w),
	}
}

// Enabled reports whether the handler handles records at the given level. The
// handler ignores records whose level is lower.
func (h *GoogleCloudHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle formats its argument Record as a JSON object on a single line.
func (h *GoogleCloudHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := make([]slog.Attr, 0, r.NumAttrs())
	httpRequest := make([]slog.Attr, 0)

	// Separate out the HTTP request attributes.
	r.Attrs(func(a slog.Attr) {
		switch a.Key {
		case MethodKey, StatusKey:
			httpRequest = append(httpRequest, a)
		default:
			attrs = append(attrs, a)
		}
	})

	attrs = append(attrs, slog.Group("httpRequest", httpRequest...))

	if h.SpanHandler != nil {
		if span := h.SpanHandler(ctx); span != NilValue {
			attrs = append(attrs, slog.Any(googleCloudSpanKey, span))
		}
	}
	if h.TraceHandler != nil {
		if trace := h.TraceHandler(ctx); trace != NilValue {
			attrs = append(attrs, slog.Any(googleCloudTraceKey, trace))
		}
	}

	// Create a new record with the attributes we want to keep.
	record := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	record.AddAttrs(attrs...)

	h.handler.Handle(ctx, record)
	return nil
}

// WithAttrs returns a new GoogleCloudHandler whose attributes consists of h's
// attributes followed by attrs.
func (h *GoogleCloudHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &GoogleCloudHandler{
		handler: h.handler.WithAttrs(attrs),
	}
}

// WithGroup returns a new GoogleCloudHandler whose attributes consists of h's
// attributes followed by a group with the given name.
func (h *GoogleCloudHandler) WithGroup(name string) slog.Handler {
	return &GoogleCloudHandler{
		handler: h.handler.WithGroup(name),
	}
}

// WithLabels returns a new GoogleCloudHandler whose attributes consists of h's
// attributes followed by the given labels.
func (h *GoogleCloudHandler) WithLabels(labels map[string]string) slog.Handler {
	return &GoogleCloudHandler{
		handler: h.handler.WithAttrs([]slog.Attr{
			slog.Any(googleCloudLabelsKey, labels),
		}),
	}
}
