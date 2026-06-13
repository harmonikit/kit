package zerolog_test

import (
	"bytes"
	"testing"

	harmonilog "github.com/harmonikit/harmoni/log"
	zerologlog "github.com/harmonikit/kit/log/zerolog"
	"github.com/rs/zerolog"
)

func TestLogger_Log(t *testing.T) {
	var buf bytes.Buffer
	logger := zerologlog.NewLogger(zerolog.New(&buf), zerolog.DebugLevel)

	logger.Log(harmonilog.LevelInfo, "hello", "key", "value")

	out := buf.String()
	if out == "" {
		t.Fatal("expected log output")
	}
}

func TestLogger_With(t *testing.T) {
	var buf bytes.Buffer
	base := zerologlog.NewLogger(zerolog.New(&buf), zerolog.DebugLevel)

	derived := base.With("component", "server")
	derived.Log(harmonilog.LevelInfo, "started")

	out := buf.String()
	if out == "" {
		t.Fatal("expected log output")
	}
}

func TestLogger_LevelFilter(t *testing.T) {
	var buf bytes.Buffer
	logger := zerologlog.NewLogger(zerolog.New(&buf), zerolog.ErrorLevel)

	logger.Log(harmonilog.LevelInfo, "should be filtered")
	if buf.String() != "" {
		t.Fatal("info message should have been filtered")
	}
}

func TestLogger_Interface(t *testing.T) {
	var _ harmonilog.Logger = zerologlog.NewLogger(zerolog.New(nil), zerolog.DebugLevel)
}