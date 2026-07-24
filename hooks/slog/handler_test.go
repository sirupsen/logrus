package slog_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	lslog "github.com/sirupsen/logrus/hooks/slog"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestHandler_basic(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.SetLevel(logrus.TraceLevel)

	s := slog.New(lslog.NewHandler(logger))
	s.Info("hello", "animal", "walrus")

	entry := hook.LastEntry()
	if entry == nil {
		t.Fatal("expected a log entry")
	}
	if entry.Level != logrus.InfoLevel {
		t.Errorf("level: got %v, want %v", entry.Level, logrus.InfoLevel)
	}
	if entry.Message != "hello" {
		t.Errorf("message: got %q, want %q", entry.Message, "hello")
	}
	if got, ok := entry.Data["animal"]; !ok || got != "walrus" {
		t.Errorf("field animal: got %v (%v), want walrus", got, ok)
	}
}

func TestHandler_levelMapping(t *testing.T) {
	tests := []struct {
		name      string
		log       func(*slog.Logger)
		wantLevel logrus.Level
	}{
		{
			name:      "debug",
			log:       func(s *slog.Logger) { s.Debug("d") },
			wantLevel: logrus.DebugLevel,
		},
		{
			name:      "info",
			log:       func(s *slog.Logger) { s.Info("i") },
			wantLevel: logrus.InfoLevel,
		},
		{
			name:      "warn",
			log:       func(s *slog.Logger) { s.Warn("w") },
			wantLevel: logrus.WarnLevel,
		},
		{
			name:      "error",
			log:       func(s *slog.Logger) { s.Error("e") },
			wantLevel: logrus.ErrorLevel,
		},
		{
			name: "trace via low level",
			log: func(s *slog.Logger) {
				s.Log(context.Background(), slog.LevelDebug-4, "t")
			},
			wantLevel: logrus.TraceLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, hook := test.NewNullLogger()
			logger.SetLevel(logrus.TraceLevel)
			s := slog.New(lslog.NewHandler(logger))
			tt.log(s)
			entry := hook.LastEntry()
			if entry == nil {
				t.Fatal("expected a log entry")
			}
			if entry.Level != tt.wantLevel {
				t.Errorf("level: got %v, want %v", entry.Level, tt.wantLevel)
			}
		})
	}
}

func TestHandler_customLevelMapper(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.SetLevel(logrus.TraceLevel)

	h := lslog.NewHandler(logger)
	h.LevelMapper = func(slog.Level) logrus.Level {
		return logrus.WarnLevel
	}
	s := slog.New(h)
	s.Info("mapped")

	entry := hook.LastEntry()
	if entry == nil {
		t.Fatal("expected a log entry")
	}
	if entry.Level != logrus.WarnLevel {
		t.Errorf("level: got %v, want %v", entry.Level, logrus.WarnLevel)
	}
}

func TestHandler_enabledRespectsLogrusLevel(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.SetLevel(logrus.WarnLevel)

	s := slog.New(lslog.NewHandler(logger))
	s.Info("suppressed")
	s.Warn("shown")

	entries := hook.AllEntries()
	if len(entries) != 1 {
		t.Fatalf("got %d entries, want 1", len(entries))
	}
	if entries[0].Message != "shown" {
		t.Errorf("message: got %q, want %q", entries[0].Message, "shown")
	}
}

func TestHandler_withAttrsAndGroups(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	s := slog.New(lslog.NewHandler(logger))
	s = s.With("root", 1)
	s = s.WithGroup("req")
	s = s.With("id", "abc")
	s.Info("done", "status", 200)

	entry := hook.LastEntry()
	if entry == nil {
		t.Fatal("expected a log entry")
	}
	if entry.Message != "done" {
		t.Errorf("message: got %q, want %q", entry.Message, "done")
	}
	// slog stores integer Attr values as int64.
	want := map[string]any{
		"root":       int64(1),
		"req.id":     "abc",
		"req.status": int64(200),
	}
	for k, v := range want {
		got, ok := entry.Data[k]
		if !ok {
			t.Errorf("missing field %q", k)
			continue
		}
		if got != v {
			t.Errorf("field %q: got %v (%T), want %v (%T)", k, got, got, v, v)
		}
	}
	if _, ok := entry.Data["id"]; ok {
		t.Error("ungrouped key 'id' should not be present")
	}
}

func TestHandler_inlineGroupAttr(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	s := slog.New(lslog.NewHandler(logger))
	s.Info("msg", slog.Group("g", slog.Int("n", 3), slog.String("s", "x")))

	entry := hook.LastEntry()
	if entry == nil {
		t.Fatal("expected a log entry")
	}
	if got := entry.Data["g.n"]; got != int64(3) {
		t.Errorf("g.n: got %v (%T), want int64(3)", got, got)
	}
	if got := entry.Data["g.s"]; got != "x" {
		t.Errorf("g.s: got %v, want x", got)
	}
}

func TestHandler_withEmptyGroupName(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	h := lslog.NewHandler(logger)
	h2 := h.WithGroup("")
	if h2 != h {
		t.Error("WithGroup(\"\") should return the same handler")
	}

	s := slog.New(h)
	s = s.WithGroup("")
	s.Info("msg", "k", "v")

	entry := hook.LastEntry()
	if entry == nil {
		t.Fatal("expected a log entry")
	}
	if got := entry.Data["k"]; got != "v" {
		t.Errorf("k: got %v, want v", got)
	}
}

func TestHandler_contextAndTime(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	type ctxKey struct{}
	ctx := context.WithValue(context.Background(), ctxKey{}, "val")
	ts := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	h := lslog.NewHandler(logger)
	s := slog.New(h)
	s.LogAttrs(ctx, slog.LevelInfo, "timed", slog.Time("when", ts))

	// Also log with an explicit record time via Handle path used by Logger.
	// Logger sets record.Time itself; verify WithTime was applied by checking
	// the entry is recent and context was stored.
	entry := hook.LastEntry()
	if entry == nil {
		t.Fatal("expected a log entry")
	}
	if entry.Context == nil {
		t.Fatal("expected entry context to be set")
	}
	if got := entry.Context.Value(ctxKey{}); got != "val" {
		t.Errorf("context value: got %v, want val", got)
	}
	if got := entry.Data["when"]; !got.(time.Time).Equal(ts) {
		t.Errorf("when: got %v, want %v", got, ts)
	}
	if entry.Time.IsZero() {
		t.Error("expected non-zero entry time")
	}
}

func TestHandler_outputIntegration(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logger.SetLevel(logrus.DebugLevel)

	s := slog.New(lslog.NewHandler(logger))
	s.With("chicken", "cluck").Error("error")

	got := strings.TrimSpace(buf.String())
	// TextFormatter field order is not guaranteed; check substrings.
	for _, want := range []string{"level=error", "msg=error", "chicken=cluck"} {
		if !strings.Contains(got, want) {
			t.Errorf("output %q missing %q", got, want)
		}
	}
}

func TestHandler_nilLoggerPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected panic for nil logger")
		}
	}()
	_ = lslog.NewHandler(nil)
}

func TestDefaultSlogLevelMapper(t *testing.T) {
	cases := []struct {
		in   slog.Level
		want logrus.Level
	}{
		{slog.LevelDebug - 1, logrus.TraceLevel},
		{slog.LevelDebug, logrus.DebugLevel},
		{slog.LevelInfo - 1, logrus.DebugLevel},
		{slog.LevelInfo, logrus.InfoLevel},
		{slog.LevelWarn - 1, logrus.InfoLevel},
		{slog.LevelWarn, logrus.WarnLevel},
		{slog.LevelError - 1, logrus.WarnLevel},
		{slog.LevelError, logrus.ErrorLevel},
		{slog.LevelError + 4, logrus.ErrorLevel},
	}
	for _, tc := range cases {
		got := lslog.DefaultSlogLevelMapper(tc.in)
		if got != tc.want {
			t.Errorf("DefaultSlogLevelMapper(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestHandler_withAttrsEmpty(t *testing.T) {
	logger := logrus.New()
	logger.Out = io.Discard
	h := lslog.NewHandler(logger)
	if h.WithAttrs(nil) != h {
		t.Error("WithAttrs(nil) should return same handler")
	}
	if h.WithAttrs([]slog.Attr{}) != h {
		t.Error("WithAttrs(empty) should return same handler")
	}
}
