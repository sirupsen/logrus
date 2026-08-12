package slog_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	lslog "github.com/sirupsen/logrus/hooks/slog"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_basic(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.SetLevel(logrus.TraceLevel)

	s := slog.New(lslog.NewHandler(logger, nil))
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
			name: "trace",
			log: func(s *slog.Logger) {
				s.Log(context.Background(), slog.LevelDebug-4, "t")
			},
			wantLevel: logrus.TraceLevel,
		},
		{
			name: "fatal",
			log: func(s *slog.Logger) {
				s.Log(context.Background(), slog.LevelError+2, "f")
			},
			wantLevel: logrus.FatalLevel,
		},
		{
			name: "panic",
			log: func(s *slog.Logger) {
				s.Log(context.Background(), slog.LevelError+4, "p")
			},
			wantLevel: logrus.PanicLevel,
		},
		{
			name: "below trace",
			log: func(s *slog.Logger) {
				s.Log(context.Background(), slog.LevelDebug-5, "t")
			},
			wantLevel: logrus.TraceLevel,
		},
		{
			name: "between debug and info",
			log: func(s *slog.Logger) {
				s.Log(context.Background(), slog.LevelInfo-1, "d")
			},
			wantLevel: logrus.DebugLevel,
		},
		{
			name: "between info and warn",
			log: func(s *slog.Logger) {
				s.Log(context.Background(), slog.LevelWarn-1, "i")
			},
			wantLevel: logrus.InfoLevel,
		},
		{
			name: "between warn and error",
			log: func(s *slog.Logger) {
				s.Log(context.Background(), slog.LevelError-1, "w")
			},
			wantLevel: logrus.WarnLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, hook := test.NewNullLogger()
			logger.SetLevel(logrus.TraceLevel)
			s := slog.New(lslog.NewHandler(logger, nil))
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

	s := slog.New(lslog.NewHandler(logger, &lslog.HandlerOptions{
		LevelMapper: func(slog.Level) logrus.Level { return logrus.WarnLevel },
	}))
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

	s := slog.New(lslog.NewHandler(logger, nil))
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

	s := slog.New(lslog.NewHandler(logger, nil))
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

func TestHandler_groupAttrs(t *testing.T) {
	tests := []struct {
		doc   string
		attrs []any
		want  logrus.Fields
	}{
		{
			doc: "group",
			attrs: []any{
				slog.Group("g", slog.Int("n", 3), slog.String("s", "x")),
			},
			want: logrus.Fields{
				"g.n": int64(3),
				"g.s": "x",
			},
		},
		{
			doc: "empty key",
			attrs: []any{
				slog.Any("", "value"),
			},
			want: logrus.Fields{
				"": "value",
			},
		},
		{
			doc: "empty key in group",
			attrs: []any{
				slog.Group("group", slog.Any("", "value")),
			},
			want: logrus.Fields{
				"group": "value",
			},
		},
		{
			doc: "inline group",
			attrs: []any{
				slog.Group("", slog.String("key", "value")),
			},
			want: logrus.Fields{
				"key": "value",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.doc, func(t *testing.T) {
			logger, hook := test.NewNullLogger()
			logger.SetLevel(logrus.DebugLevel)

			s := slog.New(lslog.NewHandler(logger, nil))
			s.Info("msg", tc.attrs...)

			entry := hook.LastEntry()
			require.NotNil(t, entry)
			assert.Equal(t, tc.want, entry.Data)
		})
	}
}

func TestHandler_withEmptyGroupName(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	h := lslog.NewHandler(logger, nil)
	h2 := h.WithGroup("")
	if h2 != h {
		t.Error(`WithGroup("") should return the same handler`)
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

	h := lslog.NewHandler(logger, nil)
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

	s := slog.New(lslog.NewHandler(logger, nil))
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
	_ = lslog.NewHandler(nil, nil)
}

func TestHandler_withAttrsEmpty(t *testing.T) {
	logger := logrus.New()
	logger.Out = io.Discard
	h := lslog.NewHandler(logger, nil)
	if h.WithAttrs(nil) != h {
		t.Error("WithAttrs(nil) should return same handler")
	}
	if h.WithAttrs([]slog.Attr{}) != h {
		t.Error("WithAttrs(empty) should return same handler")
	}
}

// TestHandler_PreservesRecordTime verifies that Handle preserves the timestamp
// carried by the slog record.
func TestHandler_PreservesRecordTime(t *testing.T) {
	logger, hook := test.NewNullLogger()
	h := lslog.NewHandler(logger, nil)

	ts := time.Date(2026, 8, 12, 12, 34, 56, 789, time.UTC)
	record := slog.NewRecord(ts, slog.LevelInfo, "hello", 0)

	require.NoError(t, h.Handle(context.Background(), record))

	entry := hook.LastEntry()
	require.NotNil(t, entry)
	assert.Equal(t, ts, entry.Time)
}

// TestHandler_PreservesRecordCaller verifies that Handle preserves caller
// information carried by the slog record.
func TestHandler_PreservesRecordCaller(t *testing.T) {
	if runtime.Compiler == "tinygo" {
		// TinyGo currently (v0.41.1) doesn't support `runtime.Caller`;
		// https://tinygo.org/docs/reference/lang-support/stdlib/#logslog
		t.Log("SKIP: TinyGo does not support runtime.Caller") // no t.Skip on tinygo
		return
	}

	logger, hook := test.NewNullLogger()
	logger.SetReportCaller(true)

	h := lslog.NewHandler(logger, &lslog.HandlerOptions{
		AddSource: true,
	})

	var pcs [1]uintptr
	runtime.Callers(1, pcs[:])

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "hello", pcs[0])
	require.NoError(t, h.Handle(context.Background(), record))

	entry := hook.LastEntry()
	require.NotNil(t, entry)
	require.NotNil(t, entry.Caller)

	want, _ := runtime.CallersFrames(pcs[:]).Next()
	assert.Equal(t, want.Function, entry.Caller.Function)
	assert.Equal(t, want.File, entry.Caller.File)
	assert.Equal(t, want.Line, entry.Caller.Line)
}
