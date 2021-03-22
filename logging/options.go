package logging

import (
	"context"
	"os"

	"golang.org/x/exp/slog"
)

var defaultOptions = options{
	handler:       slog.NewTextHandler(os.Stdout),
	shouldDiscard: permitAllRequestLogs,
}

// options describe the full set of options that may be configured to influence
// the output of the logger.
type options struct {
	handler       slog.Handler
	shouldDiscard FilterFunc
}

// Option is a function that can configure one or more logging options.
type Option func(*options)

// FilterFunc customises the behaviour used to determine log suppression.
type FilterFunc func(context.Context, string, error) bool

// WithHandler returns a logging option to customise the handler used to
// output log entries.
func WithHandler(h slog.Handler) Option {
	return func(o *options) {
		o.handler = h
	}
}

// WithFilter returns a logging option to suppress log entries based on the
// provided filter function.
func WithFilter(f FilterFunc) Option {
	return func(o *options) {
		o.shouldDiscard = f
	}
}

func permitAllRequestLogs(context.Context, string, error) bool {
	return false
}

func applyOptions(opts []Option) options {
	cfg := defaultOptions
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}
