package slog_test

import (
	"log/slog"
	"testing"

	"github.com/sirupsen/logrus"
	lslog "github.com/sirupsen/logrus/hooks/slog"
)

// TestLevel verifies the default Logrus-to-slog level mapping.
func TestLevel(t *testing.T) {
	tests := []struct {
		name  string
		level logrus.Level
		want  slog.Level
	}{
		{
			name:  "panic",
			level: logrus.PanicLevel,
			want:  slog.LevelError + 4,
		},
		{
			name:  "fatal",
			level: logrus.FatalLevel,
			want:  slog.LevelError + 2,
		},
		{
			name:  "error",
			level: logrus.ErrorLevel,
			want:  slog.LevelError,
		},
		{
			name:  "warn",
			level: logrus.WarnLevel,
			want:  slog.LevelWarn,
		},
		{
			name:  "info",
			level: logrus.InfoLevel,
			want:  slog.LevelInfo,
		},
		{
			name:  "debug",
			level: logrus.DebugLevel,
			want:  slog.LevelDebug,
		},
		{
			name:  "trace",
			level: logrus.TraceLevel,
			want:  slog.LevelDebug - 4,
		},
		{
			name:  "unknown",
			level: logrus.Level(255),
			want:  slog.LevelError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := lslog.Level(tc.level).Level(); got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

// TestSlogLevel verifies the default slog-to-Logrus level mapping.
func TestSlogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
		want  logrus.Level
	}{
		{
			name:  "above panic",
			level: slog.LevelError + 5,
			want:  logrus.PanicLevel,
		},
		{
			name:  "panic",
			level: slog.LevelError + 4,
			want:  logrus.PanicLevel,
		},
		{
			name:  "between fatal and panic",
			level: slog.LevelError + 3,
			want:  logrus.FatalLevel,
		},
		{
			name:  "fatal",
			level: slog.LevelError + 2,
			want:  logrus.FatalLevel,
		},
		{
			name:  "between error and fatal",
			level: slog.LevelError + 1,
			want:  logrus.ErrorLevel,
		},
		{
			name:  "error",
			level: slog.LevelError,
			want:  logrus.ErrorLevel,
		},
		{
			name:  "between warn and error",
			level: slog.LevelError - 1,
			want:  logrus.WarnLevel,
		},
		{
			name:  "warn",
			level: slog.LevelWarn,
			want:  logrus.WarnLevel,
		},
		{
			name:  "between info and warn",
			level: slog.LevelWarn - 1,
			want:  logrus.InfoLevel,
		},
		{
			name:  "info",
			level: slog.LevelInfo,
			want:  logrus.InfoLevel,
		},
		{
			name:  "between debug and info",
			level: slog.LevelInfo - 1,
			want:  logrus.DebugLevel,
		},
		{
			name:  "debug",
			level: slog.LevelDebug,
			want:  logrus.DebugLevel,
		},
		{
			name:  "between trace and debug",
			level: slog.LevelDebug - 1,
			want:  logrus.TraceLevel,
		},
		{
			name:  "trace",
			level: slog.LevelDebug - 4,
			want:  logrus.TraceLevel,
		},
		{
			name:  "below trace",
			level: slog.LevelDebug - 5,
			want:  logrus.TraceLevel,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := lslog.SlogLevel(tc.level).Level(); got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
