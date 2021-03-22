package logging

import (
	"context"
	"time"
)

var defaultOptions = options{
	durationFieldValue: defaultDurationFunc,
	shouldDiscard:      permitAllRequestLogs,
	timestampFormat:    time.RFC3339,
}

// options describe the full set of options that may be configured to influence
// the output of the logger.
type options struct {
	durationFieldValue DurationFieldFunc
	fields             []FieldFunc
	shouldDiscard      FilterFunc
	timestampFormat    string
}

// Option is a function that can configure one or more logging options.
type Option func(*options)

// DurationFieldFunc customises the log field used for request durations.
type DurationFieldFunc func(time.Duration) (string, interface{})

// FieldFunc customises the set of standard log entry fields.
type FieldFunc func(context.Context) (string, interface{})

// FilterFunc customises the behaviour used to determine log suppression.
type FilterFunc func(context.Context, string, error) bool

// WithDurationField returns a logging option to customise the log field used
// for request durations.
func WithDurationField(f DurationFieldFunc) Option {
	return func(o *options) {
		o.durationFieldValue = f
	}
}

// WithField returns a logging option to customise a loggers standard
// fields.
func WithField(f FieldFunc) Option {
	return func(o *options) {
		o.fields = append(o.fields, f)
	}
}

// WithFilter returns a logging option to customise a loggers behaviour,
// allowing certain logs to be suppressed given some condition.
func WithFilter(f FilterFunc) Option {
	return func(o *options) {
		o.shouldDiscard = f
	}
}

// WithTimestampFormat returns a logging option to customise a loggers default
// timestmap format.
func WithTimestampFormat(f string) Option {
	return func(o *options) {
		o.timestampFormat = f
	}
}

func defaultDurationFunc(duration time.Duration) (string, interface{}) {
	return DurationKey, duration
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
