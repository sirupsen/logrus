package logatee

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func New(t *testing.T) *logrus.Logger {
	return NewFunc(func() *testing.T {
		return t
	})
}

func NewFunc(tFunc func() *testing.T) *logrus.Logger {
	log := logrus.New()
	log.Out = &testWriter{tFunc}
	log.Level = logrus.TraceLevel
	log.Hooks.Add(&hook{
		tFunc:  tFunc,
		logger: log,
	})

	return log
}

type testWriter struct {
	tFunc func() *testing.T
}

func (w *testWriter) Write(b []byte) (int, error) {
	w.tFunc().Helper()
	w.tFunc().Log(string(b))

	return len(b), nil
}
