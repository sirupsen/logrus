package logrus_test

import (
	"errors"
	"io"
	"runtime"
	"testing"

	"github.com/sirupsen/logrus"
)

// BenchmarkEntry_NewEntry measures the cost of creating entries that are not
// subsequently logged.
func BenchmarkEntry_NewEntry(b *testing.B) {
	logger := logrus.New()

	b.Run("empty", func(b *testing.B) {
		b.ReportAllocs()
		var entry *logrus.Entry
		for range b.N {
			entry = logrus.NewEntry(logger)
		}
		runtime.KeepAlive(entry)
	})

	b.Run("with_field", func(b *testing.B) {
		b.ReportAllocs()
		var entry *logrus.Entry
		for range b.N {
			entry = logrus.NewEntry(logger).WithField("key", "value")
		}
		runtime.KeepAlive(entry)
	})
}

// BenchmarkEntry_WithError compares the ways an error can be added to an Entry.
func BenchmarkEntry_WithError(b *testing.B) {
	logger := logrus.New()
	base := logrus.NewEntry(logger).WithField("a", 1)
	errBoom := errors.New("boom")

	b.Run("WithError", func(b *testing.B) {
		b.ReportAllocs()
		var entry *logrus.Entry
		for range b.N {
			entry = base.WithError(errBoom)
		}
		runtime.KeepAlive(entry)
	})

	b.Run("WithField", func(b *testing.B) {
		b.ReportAllocs()
		var entry *logrus.Entry
		for range b.N {
			entry = base.WithField(logrus.ErrorKey, errBoom)
		}
		runtime.KeepAlive(entry)
	})

	b.Run("WithFields", func(b *testing.B) {
		b.ReportAllocs()
		var entry *logrus.Entry
		for range b.N {
			entry = base.WithFields(logrus.Fields{logrus.ErrorKey: errBoom})
		}
		runtime.KeepAlive(entry)
	})
}

func BenchmarkEntry_WithField_Chain(b *testing.B) {
	base := &logrus.Entry{Data: logrus.Fields{"a": 1}}
	errBoom := errors.New("boom")

	b.ReportAllocs()
	b.ResetTimer()

	var result *logrus.Entry
	for range b.N {
		e := base

		e = e.WithField("k0", 0)
		e = e.WithField("k1", 1)
		e = e.WithField("k2", 2)
		e = e.WithField("k3", 3)
		e = e.WithError(errBoom)

		result = e
	}
	runtime.KeepAlive(result)
}

// BenchmarkEntry_WithField_Chain_Disabled measures the cost of constructing a
// chain of fields when the resulting log call is disabled by the log level.
func BenchmarkEntry_WithField_Chain_Disabled(b *testing.B) {
	logger := logrus.New()
	logger.SetFormatter(nopFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	base := logrus.NewEntry(logger).WithField("a", 1)
	errBoom := errors.New("boom")

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		base.
			WithField("k0", 0).
			WithField("k1", 1).
			WithField("k2", 2).
			WithField("k3", 3).
			WithError(errBoom).
			Debug("message")
	}
}

// BenchmarkEntry_WithField_Chain_Enabled measures the cost of constructing and
// logging a chain of fields when the log level is enabled.
func BenchmarkEntry_WithField_Chain_Enabled(b *testing.B) {
	logger := logrus.New()
	logger.SetFormatter(nopFormatter{})
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(io.Discard)

	base := logrus.NewEntry(logger).WithField("a", 1)
	errBoom := errors.New("boom")

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		base.
			WithField("k0", 0).
			WithField("k1", 1).
			WithField("k2", 2).
			WithField("k3", 3).
			WithError(errBoom).
			Debug("message")
	}
}

func BenchmarkEntry_WithFields(b *testing.B) {
	fn := func() {}
	fnPtr := &fn

	tests := []struct {
		name   string
		base   logrus.Fields
		fields logrus.Fields
	}{
		{
			name:   "valid_fields_only",
			base:   logrus.Fields{"a": 1, "b": "two"},
			fields: logrus.Fields{"c": 3, "d": "four"},
		},
		{
			name:   "contains_func",
			base:   logrus.Fields{"a": 1},
			fields: logrus.Fields{"bad": fn},
		},
		{
			name:   "contains_func_ptr",
			base:   logrus.Fields{"a": 1},
			fields: logrus.Fields{"bad": fnPtr},
		},
		{
			name:   "mixed_valid_invalid",
			base:   logrus.Fields{"a": 1, "b": 2},
			fields: logrus.Fields{"c": 3, "bad": fn, "d": 4},
		},
		{
			name:   "larger_map",
			base:   logrus.Fields{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7, "h": 8, "i": 9, "j": 10},
			fields: logrus.Fields{"k": 11, "l": 12, "m": 13, "n": 14, "o": 15},
		},
	}

	for _, tc := range tests {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			e := &logrus.Entry{Data: tc.base}
			b.ResetTimer()
			var result *logrus.Entry
			for range b.N {
				result = e.WithFields(tc.fields)
			}
			runtime.KeepAlive(result)
		})
	}
}

// BenchmarkEntry_WithFields_Reused measures logging through an Entry that was
// constructed once with fields and reused across multiple log calls.
func BenchmarkEntry_WithFields_Reused(b *testing.B) {
	tests := []struct {
		name  string
		entry func(*logrus.Logger) *logrus.Entry
	}{
		{
			name: "small",
			entry: func(logger *logrus.Logger) *logrus.Entry {
				return logger.WithFields(smallFields)
			},
		},
		{
			name: "small_duplicates",
			entry: func(logger *logrus.Logger) *logrus.Entry {
				return logger.
					WithFields(smallFields).
					WithFields(smallFields).
					WithFields(smallFields)
			},
		},
		{
			name: "large",
			entry: func(logger *logrus.Logger) *logrus.Entry {
				return logger.WithFields(largeFields)
			},
		},
		{
			name: "large_duplicates",
			entry: func(logger *logrus.Logger) *logrus.Entry {
				return logger.
					WithFields(largeFields).
					WithFields(largeFields).
					WithFields(largeFields)
			},
		},
	}

	for _, tc := range tests {
		b.Run(tc.name, func(b *testing.B) {
			logger := logrus.New()
			logger.SetFormatter(nopFormatter{})
			logger.SetOutput(io.Discard)

			entry := tc.entry(logger)

			b.ReportAllocs()
			b.ResetTimer()

			for range b.N {
				entry.Info("message")
			}
		})
	}
}

func benchmarkEntryInfo(b *testing.B, reportCaller bool) {
	// JSONFormatter is used intentionally to measure realistic end-to-end
	// ReportCaller overhead (Entry.log + caller field formatting),
	// not getCaller() in isolation.
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(reportCaller)
	logger.SetLevel(logrus.InfoLevel) // ensure Info is enabled
	logger.SetOutput(io.Discard)

	entry := logrus.NewEntry(logger)

	// getCaller has a package-level sync.Once; exclude initialization from the benchmark.
	entry.Info("warmup")
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		entry.Info("test message")
	}
}

func BenchmarkEntry_ReportCaller_NoCaller(b *testing.B)   { benchmarkEntryInfo(b, false) }
func BenchmarkEntry_ReportCaller_WithCaller(b *testing.B) { benchmarkEntryInfo(b, true) }

//go:noinline
func caller4(entry *logrus.Entry) { caller3(entry) }

//go:noinline
func caller3(entry *logrus.Entry) { caller2(entry) }

//go:noinline
func caller2(entry *logrus.Entry) { caller1(entry) }

//go:noinline
func caller1(entry *logrus.Entry) { entry.Info("test message") }

// benchmarkEntryReportCallerDepth4 simulates a wrapper call site.
// It does not increase getCaller() scan depth (which stops at the first
// non-logrus frame), but ensures ReportCaller overhead is stable with
// wrapper layers.
func benchmarkEntryReportCallerDepth4(b *testing.B, reportCaller bool) {
	// JSONFormatter is used intentionally to measure realistic end-to-end
	// ReportCaller overhead (Entry.log + caller field formatting),
	// not getCaller() in isolation.
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(reportCaller)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(io.Discard)

	entry := logrus.NewEntry(logger)

	// getCaller has a package-level sync.Once; exclude initialization from the benchmark.
	entry.Info("warmup")
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		caller4(entry)
	}
}

func BenchmarkEntry_ReportCaller_NoCaller_Depth4(b *testing.B) {
	benchmarkEntryReportCallerDepth4(b, false)
}

func BenchmarkEntry_ReportCaller_WithCaller_Depth4(b *testing.B) {
	benchmarkEntryReportCallerDepth4(b, true)
}
