package slog

import (
	"context"
	"log/slog"
	"maps"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
)

// SlogLevelMapper maps a [slog.Level] to a [logrus.Level].
//
// To change the default level mapping, for instance to map custom slog levels
// or to send [slog.LevelError] to [logrus.FatalLevel], set
// [Handler.LevelMapper] to your own implementation of this function.
type SlogLevelMapper func(slog.Level) logrus.Level

// Handler is a [slog.Handler] that writes records to a [logrus.Logger].
//
// It is intended for bridging libraries or application code that log via slog
// into an existing Logrus logger, for example during a gradual migration.
//
// Example usage:
//
//	logger := logrus.New()
//	slog.SetDefault(slog.New(NewHandler(logger)))
//	slog.Info("hello", "key", "value")
type Handler struct {
	logger *logrus.Logger

	// LevelMapper maps slog levels to logrus levels.
	// If nil, a default mapping is used (see [DefaultSlogLevelMapper]).
	LevelMapper SlogLevelMapper

	// fields holds attributes from prior WithAttrs calls, already resolved
	// under the group prefix that applied when they were added.
	fields logrus.Fields

	// groups is the current group name stack for attributes added later.
	groups []string
}

var _ slog.Handler = (*Handler)(nil)

// NewHandler creates a [slog.Handler] that writes to the provided logrus Logger.
// The provided logger must not be nil. NewHandler panics if logger is nil.
func NewHandler(logger *logrus.Logger) *Handler {
	if logger == nil {
		panic("cannot create handler from nil logger")
	}
	return &Handler{
		logger: logger,
		fields: make(logrus.Fields),
	}
}

// DefaultSlogLevelMapper is the default mapping from slog levels to logrus levels.
//
//	slog level >= Error  → ErrorLevel
//	slog level >= Warn   → WarnLevel
//	slog level >= Info   → InfoLevel
//	slog level >= Debug  → DebugLevel
//	otherwise            → TraceLevel
func DefaultSlogLevelMapper(level slog.Level) logrus.Level {
	switch {
	case level >= slog.LevelError:
		return logrus.ErrorLevel
	case level >= slog.LevelWarn:
		return logrus.WarnLevel
	case level >= slog.LevelInfo:
		return logrus.InfoLevel
	case level >= slog.LevelDebug:
		return logrus.DebugLevel
	default:
		return logrus.TraceLevel
	}
}

func (h *Handler) toLogrusLevel(level slog.Level) logrus.Level {
	if h.LevelMapper != nil {
		return h.LevelMapper(level)
	}
	return DefaultSlogLevelMapper(level)
}

// Enabled reports whether the handler handles records at the given level.
// It maps the slog level to a logrus level and consults the underlying logger.
func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return h.logger.IsLevelEnabled(h.toLogrusLevel(level))
}

// Handle converts the slog record into a logrus entry and logs it.
// Record time, context, message, and attributes (including those from
// WithAttrs/WithGroup) are preserved. Attributes are attached as logrus fields;
// group names are joined with "." as a key prefix (similar to [slog.TextHandler]).
func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	level := h.toLogrusLevel(record.Level)

	fields := maps.Clone(h.fields)
	if fields == nil {
		fields = make(logrus.Fields, record.NumAttrs())
	}
	record.Attrs(func(a slog.Attr) bool {
		appendAttr(fields, h.groups, a)
		return true
	})

	entry := logrus.NewEntry(h.logger).WithContext(ctx).WithTime(record.Time)
	if len(fields) > 0 {
		entry = entry.WithFields(fields)
	}
	entry.Log(level, record.Message)
	return nil
}

// WithAttrs returns a new Handler whose attributes consist of h's attributes
// followed by attrs.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := h.clone()
	h2.fields = maps.Clone(h.fields)
	if h2.fields == nil {
		h2.fields = make(logrus.Fields, len(attrs))
	}
	for _, a := range attrs {
		appendAttr(h2.fields, h.groups, a)
	}
	return h2
}

// WithGroup returns a new Handler with a group appended to the receiver's
// existing groups. All attributes added to the returned handler (via WithAttrs
// or Handle) are nested under the combined group names.
// If name is empty, WithGroup returns h unchanged.
func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := h.clone()
	h2.groups = append(slices.Clip(h.groups), name)
	return h2
}

func (h *Handler) clone() *Handler {
	return &Handler{
		logger:      h.logger,
		LevelMapper: h.LevelMapper,
		fields:      h.fields,
		groups:      h.groups,
	}
}

// appendAttr adds attr to fields, resolving it and applying group prefixes.
// Group attributes are flattened with "."-separated keys.
func appendAttr(fields logrus.Fields, groups []string, attr slog.Attr) {
	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) {
		return
	}
	if attr.Value.Kind() == slog.KindGroup {
		gs := attr.Value.Group()
		if len(gs) == 0 {
			return
		}
		g := groups
		if attr.Key != "" {
			g = append(slices.Clip(groups), attr.Key)
		}
		for _, a := range gs {
			appendAttr(fields, g, a)
		}
		return
	}
	if attr.Key == "" {
		return
	}
	key := attr.Key
	if len(groups) > 0 {
		key = strings.Join(groups, ".") + "." + key
	}
	fields[key] = attr.Value.Any()
}
