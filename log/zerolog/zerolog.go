// Package zerolog adapts github.com/rs/zerolog to the harmoni/log.Logger interface.
package zerolog

import (
	harmonilog "github.com/harmonikit/harmoni/log"
	"github.com/rs/zerolog"
)

// Logger wraps zerolog.Logger as a harmoni/log.Logger.
type Logger struct {
	logger zerolog.Logger
	level  zerolog.Level
}

// NewLogger creates a harmoni Logger backed by a zerolog.Logger.
func NewLogger(logger zerolog.Logger, level zerolog.Level) *Logger {
	return &Logger{logger: logger, level: level}
}

// Log emits a log message at the given level.
func (l *Logger) Log(level harmonilog.Level, msg string, keysAndValues ...any) {
	zl := toZerologLevel(level)
	if zl < l.level {
		return
	}
	ev := l.logger.WithLevel(zl)
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			ev = ev.Interface(keyString(keysAndValues[i]), keysAndValues[i+1])
		}
	}
	ev.Msg(msg)
}

// With returns a new Logger with additional context.
func (l *Logger) With(keysAndValues ...any) harmonilog.Logger {
	ctx := l.logger.With()
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			ctx = ctx.Interface(keyString(keysAndValues[i]), keysAndValues[i+1])
		}
	}
	newLogger := ctx.Logger()
	return &Logger{logger: newLogger, level: l.level}
}

func toZerologLevel(l harmonilog.Level) zerolog.Level {
	switch l {
	case harmonilog.LevelDebug:
		return zerolog.DebugLevel
	case harmonilog.LevelInfo:
		return zerolog.InfoLevel
	case harmonilog.LevelWarn:
		return zerolog.WarnLevel
	case harmonilog.LevelError:
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

func keyString(k any) string {
	if s, ok := k.(string); ok {
		return s
	}
	return ""
}
