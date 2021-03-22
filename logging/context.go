package logging

import "context"

type loggerContextKey struct{}

// FromContext returns the LeveledLogger value stored in ctx, if any. If no
// LeveledLogger can be found then a default logger is returned.
func FromContext(ctx context.Context) *LeveledLogger {
	logger, ok := ctx.Value(loggerContextKey{}).(*LeveledLogger)
	if !ok {
		return Default()
	}

	return logger
}

// NewContext returns a new Context that carries a LeveledLogger.
func NewContext(parent context.Context, logger *LeveledLogger) context.Context {
	return context.WithValue(parent, loggerContextKey{}, logger)
}
