package logrus_test

import (
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
)

func BenchmarkEntry_WithError(b *testing.B) {
	base := &logrus.Entry{Data: logrus.Fields{"a": 1}}
	errBoom := errors.New("boom")
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = base.WithError(errBoom)
	}
}

