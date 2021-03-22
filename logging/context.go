package logging

import "context"

type loggerContextKey struct{}

// FromContext returns the Logger value stored in ctx, if any. If no Logger can
// be found then a new one is returned.
func FromContext(ctx context.Context) Logger {
	log, ok := ctx.Value(loggerContextKey{}).(Logger)
	if !ok {
		return nil
	}

	return log
}

// NewContext returns a new Context that carries value e.
func NewContext(parent context.Context, log Logger) context.Context {
	if log == nil {
		return parent
	}

	return context.WithValue(parent, loggerContextKey{}, log)
}
