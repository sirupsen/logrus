// Package slog provides adapters between Logrus and log/slog.
//
// [Handler] forwards slog records to a Logrus logger, while [Hook] forwards
// Logrus entries to a slog logger. They can be used independently to support
// incremental migration between the two logging APIs.
package slog

import (
	"context"
	"log/slog"

	"github.com/sirupsen/logrus"
)

// slog levels corresponding to Logrus levels.
const (
	slogLevelTrace = slog.LevelDebug - 4
	slogLevelDebug = slog.LevelDebug
	slogLevelInfo  = slog.LevelInfo
	slogLevelWarn  = slog.LevelWarn
	slogLevelError = slog.LevelError
	slogLevelFatal = slog.LevelError + 2
	slogLevelPanic = slog.LevelError + 4
)

// Hook sends Logrus entries to slog.
//
// By default, Logrus levels map to their corresponding slog levels. Trace,
// Fatal, and Panic use custom slog levels below Debug and above Error,
// respectively, preserving their relative severity:
//
//   - [logrus.TraceLevel] -> [slog.LevelDebug] - 4
//   - [logrus.DebugLevel] -> [slog.LevelDebug]
//   - [logrus.InfoLevel] -> [slog.LevelInfo]
//   - [logrus.WarnLevel] -> [slog.LevelWarn]
//   - [logrus.ErrorLevel] -> [slog.LevelError]
//   - [logrus.FatalLevel] -> [slog.LevelError] + 2
//   - [logrus.PanicLevel] -> [slog.LevelError] + 4
//
// Set [Hook.LevelMapper] to customize this mapping.
type Hook struct {
	logger *slog.Logger

	// LevelMapper maps Logrus levels to slog levels. If nil, the default
	// mapping is used. Set it to customize level mapping, for example to
	// support custom or dynamic slog levels.
	LevelMapper func(logrus.Level) slog.Leveler
}

var _ logrus.Hook = (*Hook)(nil)

// NewHook creates a hook that sends logs to an existing slog Logger.
// This hook is intended to be used during transition from Logrus to slog,
// or as a shim between different parts of your application or different
// libraries that depend on different loggers.
//
// The provided logger must not be nil. NewHook panics if logger is nil.
func NewHook(logger *slog.Logger) *Hook {
	if logger == nil {
		panic("cannot create hook from nil logger")
	}
	return &Hook{
		logger: logger,
	}
}

// toSlogLevel maps a Logrus level using LevelMapper or the default mapping.
func (h *Hook) toSlogLevel(level logrus.Level) slog.Leveler {
	if h.LevelMapper != nil {
		return h.LevelMapper(level)
	}
	switch level {
	case logrus.PanicLevel:
		return slogLevelPanic
	case logrus.FatalLevel:
		return slogLevelFatal
	case logrus.ErrorLevel:
		return slogLevelError
	case logrus.WarnLevel:
		return slogLevelWarn
	case logrus.InfoLevel:
		return slogLevelInfo
	case logrus.DebugLevel:
		return slogLevelDebug
	case logrus.TraceLevel:
		return slogLevelTrace
	default:
		// Treat all unknown levels as errors
		return slogLevelError
	}
}

// Levels always returns all levels, since slog allows controlling level
// enabling based on context.
func (h *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire forwards the provided logrus Entry to the underlying slog.Logger's
// Handler, mapping it to a slog.Record. Time and caller information are
// preserved when available, and Entry.Data is converted to attributes.
// If Entry.Context is nil, context.Background() is used.
func (h *Hook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		ctx = context.Background()
	}
	lvl := h.toSlogLevel(entry.Level).Level()
	handler := h.logger.Handler()
	if !handler.Enabled(ctx, lvl) {
		return nil
	}
	attrs := make([]slog.Attr, 0, len(entry.Data))
	for k, v := range entry.Data {
		attrs = append(attrs, slog.Any(k, v))
	}
	var pc uintptr
	if entry.Caller != nil {
		pc = entry.Caller.PC
	}
	r := slog.NewRecord(entry.Time, lvl, entry.Message, pc)
	r.AddAttrs(attrs...)
	return handler.Handle(ctx, r)
}
