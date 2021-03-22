package logging

import (
	"context"
	"time"

	"golang.org/x/exp/slog"
)

func newCommonRecord(ctx context.Context, level slog.Level, t time.Time, method, path string) slog.Record {
	duration := time.Since(t)

	record := slog.NewRecord(t, level, "", 0)
	record.AddAttrs(
		slog.Duration(DurationKey, duration),
		slog.String(MethodKey, method),
		slog.String(PathKey, path),
	)

	// If a deadline was set on the context add this to the log entry.
	if d, ok := ctx.Deadline(); ok {
		record.AddAttrs(slog.Time(DeadlineKey, d))
	}

	return record
}
