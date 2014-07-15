package logrus

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func LogAndAssertText(t *testing.T, log func(*Logger), assertions func(*Entry)) {
	var buffer bytes.Buffer

	logger := New()
	logger.Out = &buffer
	formatter := new(TextFormatter)
	logger.Formatter = formatter

	log(logger)

	entry, err := formatter.Unformat(buffer.Bytes())
	assert.Nil(t, err)

	if assert.NotNil(t, entry) {
		assertions(entry)
	}
}

func TestTextPrint(t *testing.T) {
	LogAndAssertText(t, func(log *Logger) {
		log.Print("test")
	}, func(e *Entry) {
		assert.Equal(t, e.Msg, "test", "Entry: %v", e)
		assert.Equal(t, e.Level, Info)
	})
}

func TestTextMultiData(t *testing.T) {
	LogAndAssertText(t, func(log *Logger) {
		log.WithField("wow", "pink elephant").WithField("such", "big whale").Print("test with spaces")
	}, func(e *Entry) {
		assert.Equal(t, e.Msg, "test with spaces", "Entry: %v", e)
		assert.Equal(t, e.Level, Info)
		assert.Equal(t, e.Data["wow"], "pink elephant")
		assert.Equal(t, e.Data["such"], "big whale")
	})
}
