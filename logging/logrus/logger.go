package logrus

import (
	"github.com/sirupsen/logrus"

	"github.com/crumbandbase/service-core-go/logging"
)

// Logger is an adapter around the Logrus logging package allowing it to be
// generalised.
type Logger struct {
	*logrus.Entry
}

// New returns a new Logger
func New(logger *logrus.Logger) *Logger {
	return &Logger{logrus.NewEntry(logger)}
}

// V checks if the log level of the logger is greater than the given level. It
// satisfies the service-core-go logger interface.
func (l *Logger) V(level int) bool {
	return l.Entry.Logger.IsLevelEnabled(logrus.Level(uint32(level)))
}

// WithFields adds a map of fields to the log entry.
func (l *Logger) WithFields(fields logging.Fields) logging.Logger {
	return &Logger{l.Entry.WithFields(logrus.Fields(fields))}
}

// WithField adds a single field to the log entry.
func (l *Logger) WithField(key string, value interface{}) logging.Logger {
	return &Logger{l.Entry.WithField(key, value)}
}
