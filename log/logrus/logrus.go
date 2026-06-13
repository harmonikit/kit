// Package logrus adapts github.com/sirupsen/logrus to the harmoni/log.Logger interface.
package logrus

import (
	harmonilog "github.com/harmonikit/harmoni/log"
	"github.com/sirupsen/logrus"
)

// Logger wraps *logrus.Logger as a harmoni/log.Logger.
type Logger struct {
	logger *logrus.Logger
	level  logrus.Level
}

// NewLogger creates a harmoni Logger backed by a logrus.Logger.
func NewLogger(logger *logrus.Logger, level logrus.Level) *Logger {
	return &Logger{logger: logger, level: level}
}

// Log emits a log message at the given level.
func (l *Logger) Log(level harmonilog.Level, msg string, keysAndValues ...any) {
	lvl := toLogrusLevel(level)
	if lvl > l.level {
		return
	}
	entry := l.logger.WithFields(toLogrusFields(keysAndValues))
	entry.Log(lvl, msg)
}

// With returns a new Logger with additional context.
func (l *Logger) With(keysAndValues ...any) harmonilog.Logger {
	newLogger := &logrus.Logger{
		Out:       l.logger.Out,
		Formatter: l.logger.Formatter,
		Hooks:     l.logger.Hooks,
		Level:     l.logger.Level,
		ExitFunc:  l.logger.ExitFunc,
	}
	// Add context as persistent fields via a new entry.
	entry := l.logger.WithFields(toLogrusFields(keysAndValues))
	newLogger.Formatter = entry.Logger.Formatter
	return &entryLogger{entry: entry, level: l.level}
}

// entryLogger wraps a *logrus.Entry for the With() return value.
type entryLogger struct {
	entry *logrus.Entry
	level logrus.Level
}

func (el *entryLogger) Log(level harmonilog.Level, msg string, keysAndValues ...any) {
	lvl := toLogrusLevel(level)
	if lvl > el.level {
		return
	}
	el.entry.Log(lvl, msg)
}

func (el *entryLogger) With(keysAndValues ...any) harmonilog.Logger {
	return &entryLogger{
		entry: el.entry.WithFields(toLogrusFields(keysAndValues)),
		level: el.level,
	}
}

func toLogrusLevel(l harmonilog.Level) logrus.Level {
	switch l {
	case harmonilog.LevelDebug:
		return logrus.DebugLevel
	case harmonilog.LevelInfo:
		return logrus.InfoLevel
	case harmonilog.LevelWarn:
		return logrus.WarnLevel
	case harmonilog.LevelError:
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}

func toLogrusFields(kv []any) logrus.Fields {
	fields := logrus.Fields{}
	for i := 0; i < len(kv); i += 2 {
		if i+1 < len(kv) {
			fields[keyString(kv[i])] = kv[i+1]
		}
	}
	return fields
}

func keyString(k any) string {
	if s, ok := k.(string); ok {
		return s
	}
	return ""
}
