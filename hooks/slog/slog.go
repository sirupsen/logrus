//go:build go1.21
// +build go1.21

package slog

import (
	"log/slog"

	"github.com/sirupsen/logrus"
)

// SlogHook sends logs to slog.
type SlogHook struct {
	logger *slog.Logger
}

var _ logrus.Hook = (*SlogHook)(nil)

// NewSlogHook creates a hook that sends logs to an existing slog Logger.
// This hook is intended to be used during transition from Logrus to slog,
// or as a shim between different parts of your application or different
// libraries that depend on different loggers.
//
// Example usage:
//
//	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
//	hook := NewSlogHook(logger)
func NewSlogHook(logger *slog.Logger) *SlogHook {
	return &SlogHook{
		logger: logger,
	}
}

func (*SlogHook) toSlogLevel(level logrus.Level) slog.Level {
	switch level {
	case logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel:
		return slog.LevelError
	case logrus.WarnLevel:
		return slog.LevelWarn
	case logrus.InfoLevel:
		return slog.LevelInfo
	case logrus.DebugLevel, logrus.TraceLevel:
		return slog.LevelDebug
	default:
		// Treat all unknown levels as errors
		return slog.LevelError
	}
}

// Levels always returns all levels, since slog allows controlling level
// enabling based on context.
func (h *SlogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire sends entry to the underlying slog logger. The Time and Caller fields
// of entry are ignored.
func (h *SlogHook) Fire(entry *logrus.Entry) error {
	attrs := make([]interface{}, 0, len(entry.Data))
	for k, v := range entry.Data {
		attrs = append(attrs, slog.Any(k, v))
	}
	h.logger.Log(
		entry.Context,
		h.toSlogLevel(entry.Level),
		entry.Message,
		attrs...,
	)
	return nil
}
