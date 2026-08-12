package logrus_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// TestLoggerLogPanicLevelDoesNotPanic verifies that generic Log methods treat
// PanicLevel as severity only and do not trigger panic behavior.
//
// Regression and related history:
//   - https://github.com/sirupsen/logrus/pull/65
//   - https://github.com/sirupsen/logrus/pull/1283
func TestLoggerLogPanicLevelDoesNotPanic(t *testing.T) {
	tests := []struct {
		doc  string
		log  func(*logrus.Logger)
		want string
	}{
		{
			doc:  "Log",
			log:  func(logger *logrus.Logger) { logger.Log(logrus.PanicLevel, "boom") },
			want: "boom",
		},
		{
			doc:  "Logf",
			log:  func(logger *logrus.Logger) { logger.Logf(logrus.PanicLevel, "boom %d", 42) },
			want: "boom 42",
		},
		{
			doc:  "Logln",
			log:  func(logger *logrus.Logger) { logger.Logln(logrus.PanicLevel, "boom", 42) },
			want: "boom 42",
		},
		{
			doc: "LogFn",
			log: func(logger *logrus.Logger) {
				logger.LogFn(logrus.PanicLevel, func() []any { return []any{"boom", 42} })
			},
			want: "boom42",
		},
	}

	for _, tc := range tests {
		t.Run(tc.doc, func(t *testing.T) {
			logger, hook := test.NewNullLogger()

			require.NotPanics(t, func() {
				tc.log(logger)
			})

			got := hook.LastEntry()
			require.NotNil(t, got)
			assert.Equal(t, logrus.PanicLevel, got.Level)
			assert.Equal(t, tc.want, got.Message)
		})
	}
}
