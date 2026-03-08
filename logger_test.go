package logrus_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

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
