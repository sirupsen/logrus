package logrus

import (
	"bytes"
	"testing"
)

func BenchmarkLogEntryWithFieldsNoLog(b *testing.B) {
	logger := NewLogger(InfoLevel)
	logger.Out = &bytes.Buffer{}
	entry := NewLogEntry(logger)
	for i := 0; i <= b.N; i++ {
		entry.AsDebug().WithField("test", "test").Write("message")
	}
}

func BenchmarkLegacyLoggerWithFieldsNoLog(b *testing.B) {
	logger := New()
	logger.SetLevel(InfoLevel)
	logger.Out = &bytes.Buffer{}
	for i := 0; i <= b.N; i++ {
		logger.WithField("test", "test").Debug("message")
	}
}

func BenchmarkLogEntryWithFieldsLogJSON(b *testing.B) {
	logger := NewLogger(DebugLevel)
	logger.Out = &bytes.Buffer{}
	logger.Formatter = new(JSONFormatter)
	entry := NewLogEntry(logger)
	for i := 0; i <= b.N; i++ {
		entry.AsDebug().WithField("test", "test").Write("message")
	}
}

func BenchmarkLegacyLoggerWithFieldsLogJSON(b *testing.B) {
	logger := New()
	logger.SetLevel(DebugLevel)
	logger.Out = &bytes.Buffer{}
	logger.Formatter = new(JSONFormatter)
	entry := logger.WithField("a", "b")
	for i := 0; i <= b.N; i++ {
		entry.WithField("test", "test").Debug("message")
	}
}

func BenchmarkLogEntryWithFieldsLogText(b *testing.B) {
	logger := NewLogger(DebugLevel)
	logger.Out = &bytes.Buffer{}
	entry := NewLogEntry(logger)
	for i := 0; i <= b.N; i++ {
		entry.AsDebug().WithField("test", "test").Write("message")
	}
}

func BenchmarkLegacyLoggerWithFieldsLogText(b *testing.B) {
	logger := New()
	logger.SetLevel(DebugLevel)
	logger.Out = &bytes.Buffer{}
	entry := logger.WithField("a", "b")
	for i := 0; i <= b.N; i++ {
		entry.WithField("test", "test").Debug("message")
	}
}
