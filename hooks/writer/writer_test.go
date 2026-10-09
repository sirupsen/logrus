package writer_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"github.com/stretchr/testify/assert"
)

func TestDifferentLevelsGoToDifferentWriters(t *testing.T) {
	var a, b bytes.Buffer

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	})
	log.SetOutput(io.Discard) // Send all logs to nowhere by default

	log.AddHook(&writer.Hook{
		Writer: &a,
		LogLevels: []logrus.Level{
			logrus.WarnLevel,
		},
	})
	log.AddHook(&writer.Hook{ // Send info and debug logs to stdout
		Writer: &b,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
		},
	})
	log.Warn("send to a")
	log.Info("send to b")

	assert.Equal(t, "level=warning msg=\"send to a\"\n", a.String())
	assert.Equal(t, "level=info msg=\"send to b\"\n", b.String())
}

type shortWriter struct{}

func (shortWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	return len(p) - 1, nil
}

type zeroWriter struct{}

func (zeroWriter) Write(p []byte) (int, error) {
	return 0, nil
}

type errWriter struct {
	err error
	n   int
}

func (w *errWriter) Write(p []byte) (int, error) {
	return w.n, w.err
}

func TestHookFire_ShortWrite(t *testing.T) {
	tests := []struct {
		name   string
		writer io.Writer
	}{
		{
			name:   "returns len-1 with nil error",
			writer: shortWriter{},
		},
		{
			name:   "returns 0 with nil error",
			writer: zeroWriter{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := &writer.Hook{
				Writer: tt.writer,
			}
			entry := logrus.New().WithField("test", "short_write")
			err := hook.Fire(entry)
			assert.ErrorIs(t, err, io.ErrShortWrite)
		})
	}
}

func TestHookFire_WriterErrorPreserved(t *testing.T) {
	customErr := io.ErrUnexpectedEOF
	tests := []struct {
		name   string
		writer io.Writer
	}{
		{
			name: "error with 0 written",
			writer: &errWriter{
				err: customErr,
				n:   0,
			},
		},
		{
			name: "error with short written",
			writer: &errWriter{
				err: customErr,
				n:   2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := &writer.Hook{
				Writer: tt.writer,
			}
			entry := logrus.New().WithField("test", "error_preserved")
			err := hook.Fire(entry)
			assert.ErrorIs(t, err, customErr)
		})
	}
}

func TestHookFire_SuccessfulWrite(t *testing.T) {
	var buf bytes.Buffer
	hook := &writer.Hook{
		Writer: &buf,
	}
	entry := logrus.New().WithField("test", "success")
	err := hook.Fire(entry)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf.String())
}

