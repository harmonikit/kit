// Package zap adapts go.uber.org/zap to the harmoni/log.Logger interface.
package zap

import (
	harmonilog "github.com/harmonikit/harmoni/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps *zap.Logger as a harmoni/log.Logger.
type Logger struct {
	logger *zap.Logger
	level  zapcore.LevelEnabler
}

// NewLogger creates a harmoni Logger backed by a zap.Logger.
// The level parameter controls which messages are logged.
func NewLogger(logger *zap.Logger, level zapcore.LevelEnabler) *Logger {
	return &Logger{logger: logger, level: level}
}

// Log emits a log message at the given level.
func (l *Logger) Log(level harmonilog.Level, msg string, keysAndValues ...any) {
	zapLevel := toZapLevel(level)
	if !l.level.Enabled(zapLevel) {
		return
	}
	fields := toZapFields(keysAndValues)
	l.logger.Log(zapLevel, msg, fields...)
}

// With returns a new Logger with additional context.
func (l *Logger) With(keysAndValues ...any) harmonilog.Logger {
	return &Logger{
		logger: l.logger.With(toZapFields(keysAndValues)...),
		level:  l.level,
	}
}

// toZapLevel maps a harmoni Level to a zapcore.Level.
func toZapLevel(l harmonilog.Level) zapcore.Level {
	switch l {
	case harmonilog.LevelDebug:
		return zapcore.DebugLevel
	case harmonilog.LevelInfo:
		return zapcore.InfoLevel
	case harmonilog.LevelWarn:
		return zapcore.WarnLevel
	case harmonilog.LevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// toZapFields converts alternating key-value pairs to zap.Field slices.
func toZapFields(kv []any) []zap.Field {
	if len(kv) == 0 {
		return nil
	}
	fields := make([]zap.Field, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		if i+1 < len(kv) {
			fields = append(fields, zap.Any(keyString(kv[i]), kv[i+1]))
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
