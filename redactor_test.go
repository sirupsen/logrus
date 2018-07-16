package logrus

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testSentence = "this is a test sentence"

func TestNoRedactor(t *testing.T) {
	logger := New()
	var buffer bytes.Buffer
	logger.Out = &buffer
	logger.Warning(testSentence)
	assert.Equal(t, bytes.Contains(buffer.Bytes(), []byte("test")), true, "test has been redacted but it shouldn't")
}

func TestRedactorOnText(t *testing.T) {
	logger := New()
	var buffer bytes.Buffer
	logger.Out = &buffer
	logger.SetRedactor(func(in []byte) []byte {
		return bytes.Replace(in, []byte("test"), []byte("****"), -1)
	})
	logger.Warning(testSentence)
	assert.Equal(t, bytes.Contains(buffer.Bytes(), []byte("test")), false, "test has not been redacted but it shouldn't")
}

func TestRedactorOnField(t *testing.T) {
	logger := New()
	var buffer bytes.Buffer
	logger.Out = &buffer
	logger.SetRedactor(func(in []byte) []byte {
		return bytes.Replace(in, []byte("test"), []byte("****"), -1)
	})
	logger.WithField("field_one", testSentence).Warning(testSentence)
	assert.Equal(t, bytes.Contains(buffer.Bytes(), []byte("test")), false, "test has not been redacted but it shouldn't")
}
