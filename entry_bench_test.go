package logrus_test

import (
	"errors"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
)

func BenchmarkEntry_WithError(b *testing.B) {
	base := &logrus.Entry{Data: logrus.Fields{"a": 1}}
	errBoom := errors.New("boom")
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = base.WithError(errBoom)
	}
}

func BenchmarkEntry_WithField_Chain(b *testing.B) {
	base := &logrus.Entry{Data: logrus.Fields{"a": 1}}
	errBoom := errors.New("boom")
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		e := base

		e = e.WithField("k0", 0)
		e = e.WithField("k1", 1)
		e = e.WithField("k2", 2)
		e = e.WithField("k3", 3)
		e = e.WithError(errBoom)
		_ = e
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
			for range b.N {
				_ = e.WithFields(tc.fields)
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
