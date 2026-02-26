//go:build go1.21
// +build go1.21

package slog

import (
	"log/slog"

	"github.com/sirupsen/logrus"
)

// LevelMapper maps a [github.com/sirupsen/logrus.Level] value to a
// [slog.Leveler] value. To change the default level mapping, for instance
// to allow mapping to custom or dynamic slog levels in your application, set
// [SlogHook.LevelMapper] to your own implementation of this function.
type LevelMapper func(logrus.Level) slog.Leveler

// SlogHook sends logs to slog.
type SlogHook struct {
	logger      *slog.Logger
	LevelMapper LevelMapper
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

func (h *SlogHook) toSlogLevel(level logrus.Level) slog.Leveler {
	if h.LevelMapper != nil {
		return h.LevelMapper(level)
	}
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
	lvl := h.toSlogLevel(entry.Level).Level()
	if !h.logger.Handler().Enabled(entry.Context, lvl) {
		return nil
	}
	attrs := make([]any, 0, len(entry.Data))
	for k, v := range entry.Data {
		attrs = append(attrs, slog.Any(k, v))
	}
	var pc uintptr
	if entry.Caller != nil {
		pc = entry.Caller.PC
	}
	r := slog.NewRecord(entry.Time, lvl, entry.Message, pc)
	r.Add(attrs...)
	return h.logger.Handler().Handle(entry.Context, r)
}
