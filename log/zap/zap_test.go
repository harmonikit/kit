package zap_test

import (
	"bytes"
	"testing"

	harmonilog "github.com/harmonikit/harmoni/log"
	zaplog "github.com/harmonikit/kit/log/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger_Log(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{MessageKey: "msg", LevelKey: "level"}),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)
	logger := zaplog.NewLogger(zap.New(core), zapcore.DebugLevel)

	logger.Log(harmonilog.LevelInfo, "hello", "key", "value")

	out := buf.String()
	if out == "" {
		t.Fatal("expected log output")
	}
}

func TestLogger_With(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{MessageKey: "msg", LevelKey: "level"}),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)
	base := zaplog.NewLogger(zap.New(core), zapcore.DebugLevel)

	derived := base.With("component", "server")
	derived.Log(harmonilog.LevelInfo, "started")

	out := buf.String()
	if out == "" {
		t.Fatal("expected log output")
	}
}

func TestLogger_LevelFilter(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{MessageKey: "msg", LevelKey: "level"}),
		zapcore.AddSync(&buf),
		zapcore.ErrorLevel, // only error level passes
	)
	logger := zaplog.NewLogger(zap.New(core), zapcore.WarnLevel)

	logger.Log(harmonilog.LevelInfo, "should be filtered") // below warn level
	if buf.String() != "" {
		t.Fatal("info message should have been filtered")
	}
}

func TestLogger_Interface(t *testing.T) {
	var _ harmonilog.Logger = zaplog.NewLogger(zap.NewNop(), zapcore.DebugLevel)
}