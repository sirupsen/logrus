// TinyGo doesn't implement b.SetParallelism.
//go:build !tinygo

package logrus_test

import (
	"io"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func BenchmarkDummyLogger(b *testing.B) {
	nullf, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0o666)
	if err != nil {
		b.Fatalf("%v", err)
	}
	defer nullf.Close()
	doLoggerBenchmark(b, nullf, &logrus.TextFormatter{DisableColors: true}, smallFields)
}

func BenchmarkDummyLoggerNoLock(b *testing.B) {
	nullf, err := os.OpenFile(os.DevNull, os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		b.Fatalf("%v", err)
	}
	defer nullf.Close()
	doLoggerBenchmarkNoLock(b, nullf, &logrus.TextFormatter{DisableColors: true}, smallFields)
}

func doLoggerBenchmark(b *testing.B, out *os.File, formatter logrus.Formatter, fields logrus.Fields) {
	logger := logrus.Logger{
		Out:       out,
		Level:     logrus.InfoLevel,
		Formatter: formatter,
	}

	b.RunParallel(func(pb *testing.PB) {
		entry := logger.WithFields(fields) // new entry per goroutine
		for pb.Next() {
			entry.Info("aaa")
		}
	})
}

func doLoggerBenchmarkNoLock(b *testing.B, out *os.File, formatter logrus.Formatter, fields logrus.Fields) {
	logger := logrus.Logger{
		Out:       out,
		Level:     logrus.InfoLevel,
		Formatter: formatter,
	}
	logger.SetNoLock()

	b.RunParallel(func(pb *testing.PB) {
		entry := logger.WithFields(fields) // new entry per goroutine
		for pb.Next() {
			entry.Info("aaa")
		}
	})
}

type nopFormatter struct{}

func (nopFormatter) Format(*logrus.Entry) ([]byte, error) {
	return nil, nil
}

func BenchmarkLoggerLog(b *testing.B) {
	logger := logrus.New()
	logger.SetFormatter(nopFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(io.Discard)

	b.ReportAllocs()

	b.Run("disabled_level", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			logger.Log(logrus.DebugLevel, "test")
		}
	})
	b.Run("enabled_log", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			logger.Log(logrus.WarnLevel, "test")
		}
	})
	b.Run("enabled_logln", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			logger.Logln(logrus.WarnLevel, "test")
		}
	})
}

func BenchmarkLoggerJSONFormatter(b *testing.B) {
	doLoggerBenchmarkWithFormatter(b, &logrus.JSONFormatter{})
}

func BenchmarkLoggerTextFormatter(b *testing.B) {
	doLoggerBenchmarkWithFormatter(b, &logrus.TextFormatter{})
}

func doLoggerBenchmarkWithFormatter(b *testing.B, f logrus.Formatter) {
	b.SetParallelism(100)
	log := logrus.New()
	log.Formatter = f
	log.Out = io.Discard
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.
				WithField("foo1", "bar1").
				WithField("foo2", "bar2").
				Info("this is a dummy log")
		}
	})
}
