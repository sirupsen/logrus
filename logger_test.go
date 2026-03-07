package logrus_test

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldValueError(t *testing.T) {
	var buf bytes.Buffer
	l := &logrus.Logger{
		Out:       &buf,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
	l.WithField("func", func() {}).Info("test")
	var data map[string]any
	err := json.Unmarshal(buf.Bytes(), &data)
	require.NoError(t, err, "output:\n%s", buf.String())
	_, ok := data[logrus.FieldKeyLogrusError]
	require.True(t, ok, `cannot find expected "logrus_error" field: %v`, data)
}

func TestNoFieldValueError(t *testing.T) {
	var buf bytes.Buffer
	l := &logrus.Logger{
		Out:       &buf,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
	l.WithField("str", "str").Info("test")
	var data map[string]any
	err := json.Unmarshal(buf.Bytes(), &data)
	require.NoError(t, err, "output:\n%s", buf.String())
	_, ok := data[logrus.FieldKeyLogrusError]
	require.False(t, ok)
}

func TestWarningAndWarninglnFormatting(t *testing.T) {
	var buf bytes.Buffer
	l := &logrus.Logger{
		Out: &buf,
		Formatter: &logrus.TextFormatter{
			DisableColors:          true,
			DisableTimestamp:       true,
			DisableLevelTruncation: true,
		},
		Hooks: make(logrus.LevelHooks),
		Level: logrus.DebugLevel,
	}
	l.Warning("hello", "world")
	expected := "level=warning msg=helloworld\n"
	assert.Equal(t, expected, buf.String())

	buf.Reset()
	l.Warningln("hello", "world")

	expected = "level=warning msg=\"hello world\"\n"
	assert.Equal(t, expected, buf.String())
}

type testBufferPool struct {
	buffers []*bytes.Buffer
	get     int
}

func (p *testBufferPool) Get() *bytes.Buffer {
	p.get++
	return new(bytes.Buffer)
}

func (p *testBufferPool) Put(buf *bytes.Buffer) {
	p.buffers = append(p.buffers, buf)
}

func TestLogger_SetBufferPool(t *testing.T) {
	l := logrus.New()
	l.SetOutput(io.Discard)

	pool := new(testBufferPool)
	l.SetBufferPool(pool)

	l.Info("test")

	assert.Equal(t, 1, pool.get, "Logger.SetBufferPool(): The BufferPool.Get() must be called")
	assert.Len(t, pool.buffers, 1, "Logger.SetBufferPool(): The BufferPool.Put() must be called")
}
