package slog_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	lslog "github.com/sirupsen/logrus/hooks/slog"
)

func TestSlogHook(t *testing.T) {
	tests := []struct {
		name   string
		mapper lslog.LevelMapper
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
			mapper: func(logrus.Level) slog.Leveler {
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
			slogLogger := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{
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
			hook := lslog.NewSlogHook(slogLogger)
			hook.LevelMapper = tt.mapper
			log.AddHook(hook)
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

func TestSlogHook_error_propagates(t *testing.T) {
	stderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stderr = w
	t.Cleanup(func() {
		_ = r.Close()
	})

	slogLogger := slog.New(&errorHandler{})
	log := logrus.New()
	log.SetOutput(io.Discard)
	log.AddHook(lslog.NewSlogHook(slogLogger))
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

func TestSlogHook_source(t *testing.T) {
	buf := &bytes.Buffer{}
	slogLogger := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{
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
	log.AddHook(lslog.NewSlogHook(slogLogger))
	log.Info("info with source")
	got := strings.TrimSpace(buf.String())
	wantRE := regexp.MustCompile(`source=.*hooks[\\/]+slog[\\/]+slog_test\.go:\d+`)
	if !wantRE.MatchString(got) {
		t.Errorf("expected log to contain source attribute matching %q, got: %s", wantRE.String(), got)
	}
}
