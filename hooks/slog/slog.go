package slog

import (
	"context"
	"log/slog"
	"runtime"

	"github.com/sirupsen/logrus"
)

// LevelMapper maps a [logrus.Level] to a [slog.Leveler].
//
// To change the default level mapping, for instance to allow mapping to custom
// or dynamic slog levels in your application, set [SlogHook.LevelMapper]
// to your own implementation of this function.
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
// The provided logger must not be nil. NewSlogHook panics if logger is nil.
//
// Example usage:
//
//	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
//	hook := NewSlogHook(logger)
func NewSlogHook(logger *slog.Logger) *SlogHook {
	if logger == nil {
		panic("cannot create hook from nil logger")
	}
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

// Fire forwards the provided logrus Entry to the underlying slog.Logger's
// Handler, mapping it to a slog.Record. Time and caller information are
// preserved when available, and Entry.Data is converted to attributes.
// If Entry.Context is nil, context.Background() is used.
func (h *SlogHook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		ctx = context.Background()
	}
	lvl := h.toSlogLevel(entry.Level).Level()
	attrs := make([]slog.Attr, 0, len(entry.Data))
	for k, v := range entry.Data {
		attrs = append(attrs, slog.Any(k, v))
	}
	var pcs [1]uintptr
	// skip 8 callers to get to the original logrus caller
	runtime.Callers(8, pcs[:])
	r := slog.NewRecord(entry.Time, lvl, entry.Message, pcs[0])
	r.AddAttrs(attrs...)
	return h.logger.Handler().Handle(ctx, r)
}
