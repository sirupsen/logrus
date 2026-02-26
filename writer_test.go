package logrus_test

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sirupsen/logrus"
)

func ExampleLogger_Writer_httpServer() {
	logger := logrus.New()
	w := logger.Writer()
	defer w.Close()

	srv := http.Server{
		// create a stdlib log.Logger that writes to
		// logrus.Logger.
		ErrorLog: log.New(w, "", 0),
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}

func ExampleLogger_Writer_stdlib() {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	// Use logrus for standard log output
	// Note that `log` here references stdlib's log
	// Not logrus imported under the name `log`.
	log.SetOutput(logger.Writer())
}

type bufferWithMu struct {
	buf *bytes.Buffer
	mu  sync.RWMutex
}

func (b *bufferWithMu) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *bufferWithMu) Read(p []byte) (int, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.buf.Read(p)
}

func (b *bufferWithMu) String() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.buf.String()
}

func TestWriterSplitNewlines(t *testing.T) {
	buf := &bufferWithMu{
		buf: bytes.NewBuffer(nil),
	}
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	}
	logger.SetOutput(buf)
	writer := logger.Writer()

	const logNum = 10

	for range logNum {
		_, err := writer.Write([]byte("bar\nfoo\n"))
		require.NoError(t, err, "writer.Write failed")
	}
	writer.Close()
	// Test is flaky because it writes in another goroutine,
	// we need to make sure to wait a bit so all write are done.
	time.Sleep(500 * time.Millisecond)

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	assert.Len(t, lines, logNum*2, "logger printed incorrect number of lines")
}

func TestWriterSplitsMax64KB(t *testing.T) {
	buf := &bufferWithMu{
		buf: bytes.NewBuffer(nil),
	}
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	}
	logger.SetOutput(buf)
	writer := logger.Writer()

	// write more than 64KB
	const bigWriteLen = bufio.MaxScanTokenSize + 100
	output := make([]byte, bigWriteLen)
	// lets not write zero bytes
	for i := range bigWriteLen {
		output[i] = 'A'
	}

	for range 3 {
		len, err := writer.Write(output)
		require.NoError(t, err, "writer.Write failed")
		assert.Equal(t, bigWriteLen, len, "bytes written")
	}
	writer.Close()
	// Test is flaky because it writes in another goroutine,
	// we need to make sure to wait a bit so all write are done.
	time.Sleep(500 * time.Millisecond)

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	// we should have 4 lines because we wrote more than 64 KB each time
	assert.Len(t, lines, 4, "logger printed incorrect number of lines")
}
