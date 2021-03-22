package logging

import (
	"context"
	"os"
	"sync/atomic"

	"golang.org/x/exp/slog"
)

// Level denotes the severity of a log entry.
const (
	LevelTrace     = slog.Level(-8)
	LevelDebug     = slog.LevelDebug
	LevelInfo      = slog.LevelInfo
	LevelNotice    = slog.Level(2)
	LevelWarning   = slog.LevelWarn
	LevelError     = slog.LevelError
	LevelEmergency = slog.Level(12)
	LevelAlert     = slog.Level(16)
	LevelCritical  = slog.Level(20)
)

// Extended logger attribute keys.
const (
	DeadlineKey = "deadline"
	DurationKey = "duration"
	MethodKey   = "method"
	PathKey     = "path"
	ResponseKey = "jsonPayload"
	StatusKey   = "status"
)

// LeveledLogger is a logger that logs messages at a specific level.
type LeveledLogger struct {
	logger *slog.Logger
}

var defaultLogger atomic.Value

func init() {
	defaultLogger.Store(New(slog.NewJSONHandler(os.Stdout)))
}

// Default returns the default leveled logger.
func Default() *LeveledLogger {
	return defaultLogger.Load().(*LeveledLogger)
}

// SetDefault sets the default leveled logger.
func SetDefault(l *LeveledLogger) {
	defaultLogger.Store(l)
}

// New returns a new leveled logger that logs messages to the given handler.
func New(h slog.Handler) *LeveledLogger {
	return &LeveledLogger{
		logger: slog.New(h),
	}
}

// Trace logs a message at the trace level.
func (l *LeveledLogger) Trace(ctx context.Context, msg string, attrs ...any) {
	l.log(ctx, LevelTrace, msg, attrs...)
}

// Debug logs a message at the debug level.
func (l *LeveledLogger) Debug(ctx context.Context, msg string, attrs ...any) {
	l.log(ctx, LevelDebug, msg, attrs...)
}

// Info logs a message at the info level.
func (l *LeveledLogger) Info(ctx context.Context, msg string, attrs ...any) {
	l.log(ctx, LevelInfo, msg, attrs...)
}

// Notice logs a message at the notice level.
func (l *LeveledLogger) Notice(ctx context.Context, msg string, attrs ...any) {
	l.log(ctx, LevelNotice, msg, attrs...)
}

// Warning logs a message at the warning level.
func (l *LeveledLogger) Warning(ctx context.Context, msg string, attrs ...any) {
	l.log(ctx, LevelWarning, msg, attrs...)
}

// Error logs a message at the error level.
func (l *LeveledLogger) Error(ctx context.Context, msg string, attrs ...any) {
	l.log(ctx, LevelError, msg, attrs...)
}

// Emergency logs a message at the emergency level.
func (l *LeveledLogger) Emergency(ctx context.Context, msg string, attrs ...any) {
	l.log(ctx, LevelEmergency, msg, attrs...)
}

// Alert logs a message at the alert level.
func (l *LeveledLogger) Alert(ctx context.Context, msg string, attrs ...any) {
	l.log(ctx, LevelAlert, msg, attrs...)
}

// Critical logs a message at the critical level.
func (l *LeveledLogger) Critical(ctx context.Context, msg string, attrs ...any) {
	l.log(ctx, LevelCritical, msg, attrs...)
}

// Log logs a message at the specified level.
func (l *LeveledLogger) log(ctx context.Context, level slog.Level, msg string, attrs ...any) {
	// If a deadline was set on the context and it has been exceeded then add
	// this to the log entry.
	if d, ok := ctx.Deadline(); ok {
		attrs = append(attrs, slog.Time(DeadlineKey, d))
	}

	l.logger.Log(ctx, level, msg, attrs...)
}
