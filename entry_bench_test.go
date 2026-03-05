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

func BenchmarkEntry_WithFields(b *testing.B) {
	fn := func() {}
	fnPtr := &fn

	tests := []struct {
		name string
		base logrus.Fields
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
			for i := 0; i < b.N; i++ {
				_ = e.WithFields(tc.fields)
			}
		})
	}
}
