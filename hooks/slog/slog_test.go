package slog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"testing/slogtest"

	"github.com/sirupsen/logrus"
	lslog "github.com/sirupsen/logrus/hooks/slog"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHook(t *testing.T) {
	tests := []struct {
		name   string
		mapper func(logrus.Level) slog.Level
		fn     func(*logrus.Logger)
		want   []string
	}{
		{
			name: "defaults",
			fn: func(log *logrus.Logger) {
				log.Info("info")
			},
			want: []string{
				"level=INFO msg=info",
			},
		},
		{
			name: "with fields",
			fn: func(log *logrus.Logger) {
				log.WithFields(logrus.Fields{
					"chicken": "cluck",
				}).Error("error")
			},
			want: []string{
				"level=ERROR msg=error chicken=cluck",
			},
		},
		{
			name: "level mapper",
			mapper: func(logrus.Level) slog.Level {
				return slog.LevelInfo
			},
			fn: func(log *logrus.Logger) {
				log.WithFields(logrus.Fields{
					"chicken": "cluck",
				}).Error("error")
			},
			want: []string{
				"level=INFO msg=error chicken=cluck",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			slogger := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{
				// Remove timestamps from logs, for easier comparison
				ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						return slog.Attr{}
					}
					return a
				},
			}))
			log := logrus.New()
			log.Out = io.Discard
			log.AddHook(lslog.NewHook(slogger, &lslog.HookOptions{
				LevelMapper: tt.mapper,
			}))
			tt.fn(log)
			got := strings.Split(strings.TrimSpace(buf.String()), "\n")
			if len(got) != len(tt.want) {
				t.Errorf("Got %d log lines, expected %d", len(got), len(tt.want))
				return
			}
			for i, line := range got {
				if line != tt.want[i] {
					t.Errorf("line %d differs from expectation.\n Got: %s\nWant: %s", i, line, tt.want[i])
				}
			}
		})
	}
}

type errorHandler struct{}

var _ slog.Handler = (*errorHandler)(nil)

func (h *errorHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *errorHandler) Handle(context.Context, slog.Record) error {
	return errors.New("boom")
}

func (h *errorHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h *errorHandler) WithGroup(string) slog.Handler {
	return h
}

func TestHook_error_propagates(t *testing.T) {
	stderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stderr = w
	t.Cleanup(func() {
		_ = r.Close()
	})

	slogger := slog.New(&errorHandler{})
	log := logrus.New()
	log.SetOutput(io.Discard)
	log.AddHook(lslog.NewHook(slogger, nil))
	log.WithField("key", "value").Error("test error")

	// Restore stderr before closing the pipe writer to avoid leaving os.Stderr
	// pointing at a closed file descriptor.
	os.Stderr = stderr
	_ = w.Close()
	gotStderr, _ := io.ReadAll(r)
	if !bytes.Contains(gotStderr, []byte("boom")) {
		t.Errorf("expected stderr to contain 'boom', got: %s", string(gotStderr))
	}
}

func TestHook_source(t *testing.T) {
	if runtime.Compiler == "tinygo" {
		// TinyGo currently (v0.41.1) doesn't support `runtime.Caller`;
		// https://tinygo.org/docs/reference/lang-support/stdlib/#logslog
		t.Log("SKIP: TinyGo does not support runtime.Caller") // no t.Skip on tinygo
		return
	}
	buf := &bytes.Buffer{}
	slogger := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
		AddSource: true,
	}))
	log := logrus.New()
	log.Out = io.Discard
	log.ReportCaller = true
	log.AddHook(lslog.NewHook(slogger, nil))
	log.Info("info with source")
	got := strings.TrimSpace(buf.String())
	wantRE := regexp.MustCompile(`source=.*hooks[\\/]+slog[\\/]+slog_test\.go:\d+`)
	if !wantRE.MatchString(got) {
		t.Errorf("expected log to contain source attribute matching %q, got: %s", wantRE.String(), got)
	}
}

func TestHookLevelMapping(t *testing.T) {
	tests := []struct {
		level logrus.Level
		want  string
	}{
		{level: logrus.PanicLevel, want: "ERROR+4"},
		{level: logrus.FatalLevel, want: "ERROR+2"},
		{level: logrus.ErrorLevel, want: "ERROR"},
		{level: logrus.WarnLevel, want: "WARN"},
		{level: logrus.InfoLevel, want: "INFO"},
		{level: logrus.DebugLevel, want: "DEBUG"},
		{level: logrus.TraceLevel, want: "DEBUG-4"},
	}

	for _, tc := range tests {
		t.Run(tc.level.String(), func(t *testing.T) {
			var buf bytes.Buffer
			hook := lslog.NewHook(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
				Level: slog.LevelDebug - 4, // slogLevelTrace
				ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
					if attr.Key == slog.TimeKey {
						return slog.Attr{}
					}
					return attr
				},
			})), nil)

			logger := logrus.New()
			logger.SetOutput(io.Discard)
			entry := logrus.NewEntry(logger)
			entry.Level = tc.level
			entry.Message = "message"

			require.NoError(t, hook.Fire(entry))
			assert.Contains(t, buf.String(), "level="+tc.want+" ")
		})
	}
}

// TestHookHandler_RecursionGuard verifies that the slog Hook and Handler can be
// composed without recursively forwarding a record back through Logrus
// loggers it has already traversed.
func TestHookHandler_RecursionGuard(t *testing.T) {
	t.Run("same logger", func(t *testing.T) {
		logger, hook := test.NewNullLogger()

		slogger := slog.New(lslog.NewHandler(logger, nil))
		logger.AddHook(lslog.NewHook(slogger, nil))

		logger.Info("hello")

		entries := hook.AllEntries()
		require.Len(t, entries, 1)
		assert.Equal(t, "hello", entries[0].Message)
	})

	t.Run("different logger", func(t *testing.T) {
		loggerA, hookA := test.NewNullLogger()
		loggerB, hookB := test.NewNullLogger()

		slogger := slog.New(lslog.NewHandler(loggerB, nil))
		loggerA.AddHook(lslog.NewHook(slogger, nil))

		loggerA.Info("hello")

		entriesA := hookA.AllEntries()
		require.Len(t, entriesA, 1)
		assert.Equal(t, "hello", entriesA[0].Message)

		entriesB := hookB.AllEntries()
		require.Len(t, entriesB, 1)
		assert.Equal(t, "hello", entriesB[0].Message)
	})

	t.Run("multiple loggers", func(t *testing.T) {
		loggerA, hookA := test.NewNullLogger()
		loggerB, hookB := test.NewNullLogger()

		sloggerA := slog.New(lslog.NewHandler(loggerA, nil))
		sloggerB := slog.New(lslog.NewHandler(loggerB, nil))

		loggerA.AddHook(lslog.NewHook(sloggerB, nil))
		loggerB.AddHook(lslog.NewHook(sloggerA, nil))

		loggerA.Info("hello")

		entriesA := hookA.AllEntries()
		require.Len(t, entriesA, 1)
		assert.Equal(t, "hello", entriesA[0].Message)

		entriesB := hookB.AllEntries()
		require.Len(t, entriesB, 1)
		assert.Equal(t, "hello", entriesB[0].Message)
	})
}

// TestHookHandler_RecursionGuardPreservesContext verifies that adding recursion
// tracking does not discard context values carried by a Logrus entry.
func TestHookHandler_RecursionGuardPreservesContext(t *testing.T) {
	type ctxKey struct{}

	logger, hook := test.NewNullLogger()
	slogger := slog.New(lslog.NewHandler(logger, nil))
	logger.AddHook(lslog.NewHook(slogger, nil))

	ctx := context.WithValue(context.Background(), ctxKey{}, "value")
	logger.WithContext(ctx).Info("hello")

	entries := hook.AllEntries()
	require.Len(t, entries, 1)
	require.NotNil(t, entries[0].Context)
	assert.Equal(t, "value", entries[0].Context.Value(ctxKey{}))
}

func TestHandler_slogtest(t *testing.T) {
	// slogtest cases the Handler does not pass, with the failure they report;
	// they are skipped.
	exceptions := map[string]string{
		// logrus replaces a zero Record.Time with time.Now().
		"zero-time": `unexpected key "time": a Handler should ignore a zero Record.Time`,

		// The handler flattens groups into "."-joined keys ("G.c"), while slogtest
		// expects each group as a nested map.
		"groups":                    `missing group "G": a Handler should handle Group attributes`,
		"WithGroup":                 `missing group "G": a Handler should handle the WithGroup method`,
		"multi-With":                `missing group "G": a Handler should handle multiple WithGroup and WithAttr calls`,
		"empty-group-record":        `missing group "G": a Handler should not output groups if there are no attributes`,
		"nested-empty-group-record": `missing group "G": a Handler should not output nested groups if there are no attributes`, // Go 1.25+
		"resolve-groups":            `missing group "G": a Handler should call Resolve on attribute values in groups`,
		"resolve-WithAttrs-groups":  `missing group "G": a Handler should call Resolve on attribute values in groups from WithAttrs`,
	}

	var buf bytes.Buffer

	newHandler := func(t *testing.T) slog.Handler {
		if reason, ok := exceptions[path.Base(t.Name())]; ok {
			t.Skip("known failure: " + reason)
		}
		buf.Reset()
		l := logrus.New()
		l.SetOutput(&buf)
		l.SetFormatter(&logrus.JSONFormatter{})
		return lslog.NewHandler(l, nil)
	}

	result := func(t *testing.T) map[string]any {
		var m map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &m), "output: %q", buf.Bytes())
		return m
	}

	slogtest.Run(t, newHandler, result)
}
