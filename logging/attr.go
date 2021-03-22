package logging

import (
	"context"

	"golang.org/x/exp/slog"
)

// AttrHandler is a function that can be used to add additional attributes to
// a log entry.
type AttrHandler func(context.Context) slog.Value

// NilValue is a slog.Value that represents a nil value.
var NilValue = slog.AnyValue(nil)

func severityValue(v slog.Value) slog.Value {
	level := v.Any().(slog.Level)

	switch {
	case level >= LevelCritical:
		return slog.StringValue("CRITICAL")
	case level >= LevelAlert:
		return slog.StringValue("ALERT")
	case level >= LevelEmergency:
		return slog.StringValue("EMERGENCY")
	case level >= LevelError:
		return slog.StringValue("ERROR")
	case level >= LevelWarning:
		return slog.StringValue("WARNING")
	case level >= LevelNotice:
		return slog.StringValue("NOTICE")
	case level >= LevelInfo:
		return slog.StringValue("INFO")
	case level >= LevelDebug:
		return slog.StringValue("DEBUG")
	case level >= LevelTrace:
		return slog.StringValue("TRACE")
	default:
		return slog.StringValue("DEFAULT")
	}
}
