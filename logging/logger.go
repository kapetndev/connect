package logging

import (
	"context"
	"time"
)

const (
	// DurationKey is the default key name for a duration field.
	DurationKey = "duration"
	// MethodKey is the default key name for a method field.
	MethodKey = "method"
	// PathKey is the default key name for a path field.
	PathKey = "path"
	// PayloadKey is the default name for a payload field.
	PayloadKey = "jsonPayload"
	// StatusKey is the default name for a status field.
	StatusKey = "status"
	// TimeKey is the default nmme for a time field.
	TimeKey = "time"
	// TraceKey is the default name for a trace field.
	TraceKey = "trace"
)

// Fields represents a collection of log entry fields.
type Fields map[string]interface{}

// Logger does underlying logging work.
type Logger interface {
	// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
	Info(args ...interface{})
	// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
	Infoln(args ...interface{})
	// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
	Infof(format string, args ...interface{})
	// Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
	Warning(args ...interface{})
	// Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
	Warningln(args ...interface{})
	// Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
	Warningf(format string, args ...interface{})
	// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	Error(args ...interface{})
	// Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
	Errorln(args ...interface{})
	// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	Errorf(format string, args ...interface{})
	// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	// Implementations should call os.Exit() with a non-zero exit code.
	Fatal(args ...interface{})
	// Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
	// Implementations should call os.Exit() with a non-zero exit code.
	Fatalln(args ...interface{})
	// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	// Implementations should call os.Exit() with a non-zero exit code.
	Fatalf(format string, args ...interface{})
	// WithFields
	WithFields(Fields) Logger
	// WithField
	WithField(string, interface{}) Logger
}

func newRequestLogger(ctx context.Context, log Logger, funs []FieldFunc, method string, path string, start time.Time, timestampFormat string) Logger {
	fields := Fields{
		MethodKey: method,
		PathKey:   path,
		TimeKey:   start,
	}

	// Apply additional fields to be included in the log entry.
	for _, f := range funs {
		field, value := f(ctx)
		fields[field] = value
	}

	// If a deadline was set on the context and it has been exceeded then add
	// this to the log entry.
	if d, ok := ctx.Deadline(); ok {
		fields["request.deadline"] = d.Format(timestampFormat)
	}

	return log.WithFields(fields)
}
