package logrus_test

import (
	"bytes"
	"testing"

	harmonilog "github.com/harmonikit/harmoni/log"
	logruslog "github.com/harmonikit/kit/log/logrus"
	"github.com/sirupsen/logrus"
)

func TestLogger_Log(t *testing.T) {
	var buf bytes.Buffer
	l := logrus.New()
	l.SetOutput(&buf)
	l.SetLevel(logrus.DebugLevel)
	logger := logruslog.NewLogger(l, logrus.DebugLevel)

	logger.Log(harmonilog.LevelInfo, "hello", "key", "value")

	out := buf.String()
	if out == "" {
		t.Fatal("expected log output")
	}
}

func TestLogger_With(t *testing.T) {
	var buf bytes.Buffer
	l := logrus.New()
	l.SetOutput(&buf)
	l.SetLevel(logrus.DebugLevel)
	base := logruslog.NewLogger(l, logrus.DebugLevel)

	derived := base.With("component", "server")
	derived.Log(harmonilog.LevelInfo, "started")

	out := buf.String()
	if out == "" {
		t.Fatal("expected log output")
	}
}

func TestLogger_LevelFilter(t *testing.T) {
	var buf bytes.Buffer
	l := logrus.New()
	l.SetOutput(&buf)
	l.SetLevel(logrus.DebugLevel)
	logger := logruslog.NewLogger(l, logrus.ErrorLevel)

	logger.Log(harmonilog.LevelInfo, "should be filtered")
	if buf.String() != "" {
		t.Fatal("info message should have been filtered")
	}
}

func TestLogger_Interface(t *testing.T) {
	var _ harmonilog.Logger = logruslog.NewLogger(logrus.New(), logrus.DebugLevel)
}