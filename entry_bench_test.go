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

func BenchmarkEntry_WithField_Chain(b *testing.B) {
	base := &logrus.Entry{Data: logrus.Fields{"a": 1}}
	errBoom := errors.New("boom")
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		e := base

		e = e.WithField("k0", 0)
		e = e.WithField("k1", 1)
		e = e.WithField("k2", 2)
		e = e.WithField("k3", 3)
		e = e.WithError(errBoom)
		_ = e
	}
}