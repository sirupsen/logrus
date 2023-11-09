//go:build go1.21
// +build go1.21

package slog

import (
	"bytes"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSlogHook(t *testing.T) {
	tests := []struct {
		name   string
		mapper LevelMapper
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
			hook := NewSlogHook(slogLogger)
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
