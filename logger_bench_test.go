package logrus_test

import (
	"io"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func BenchmarkDummyLogger(b *testing.B) {
	nullf, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	if err != nil {
		b.Fatalf("%v", err)
	}
	defer nullf.Close()
	doLoggerBenchmark(b, nullf, &logrus.TextFormatter{DisableColors: true}, smallFields)
}

func BenchmarkDummyLoggerNoLock(b *testing.B) {
	nullf, err := os.OpenFile(os.DevNull, os.O_WRONLY|os.O_APPEND, 0666)
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
	entry := logger.WithFields(fields)
	b.RunParallel(func(pb *testing.PB) {
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
	entry := logger.WithFields(fields)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			entry.Info("aaa")
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
