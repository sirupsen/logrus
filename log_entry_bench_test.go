package logrus

import (
	"testing"
	"bytes"
)

func BenchmarkLogEntryWithFieldsNoLog(b *testing.B) {
	logger := NewLogger(InfoLevel)
	logger.Out = &bytes.Buffer{}
	entry := logger.Entry()
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
