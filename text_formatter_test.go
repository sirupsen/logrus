package logrus

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func LogAndAssertText(t *testing.T, log func(*Logger), assertions func(Fields)) {
	var buffer bytes.Buffer

	logger := New()
	logger.Out = &buffer
	formatter := new(TextFormatter)
	logger.Formatter = formatter

	log(logger)

	entry, err := formatter.Unformat(buffer.Bytes())
	assert.Nil(t, err)

	if assert.NotNil(t, entry) {
		assertions(entry.Data)
	}
}

func TestTextPrint(t *testing.T) {
	LogAndAssertText(t, func(log *Logger) {
		log.Print("test")
	}, func(fields Fields) {
		assert.Equal(t, fields["msg"], "test")
		assert.Equal(t, fields["level"], "info")
	})
}

func TestTextMultiData(t *testing.T) {
	LogAndAssertText(t, func(log *Logger) {
		log.WithField("wow", "pink elephant").WithField("such", "big whale").Print("test with spaces")
	}, func(fields Fields) {
		assert.Equal(t, fields["msg"], "test with spaces")
		assert.Equal(t, fields["level"], "info")
		assert.Equal(t, fields["wow"], "pink elephant")
		assert.Equal(t, fields["such"], "big whale")
	})
}
