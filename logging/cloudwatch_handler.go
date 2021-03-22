package logging

import (
	"io"

	"golang.org/x/exp/slog"
)

// CloudWatchHandler is a handler that formats log messages in a way that is
// compatible with AWS CloudWatch.
type CloudWatchHandler struct {
	*slog.JSONHandler
}

// NewCloudWatchHandler returns a new CloudWatchHandler.
func NewCloudWatchHandler(w io.Writer, level slog.Level) *CloudWatchHandler {
	return &CloudWatchHandler{
		JSONHandler: slog.HandlerOptions{
			Level: level,
		}.NewJSONHandler(w),
	}
}
