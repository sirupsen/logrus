package logrus

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func LogAndAssertJSON(t *testing.T, log func(*Logger), assertions func(fields Fields)) {
	var buffer bytes.Buffer
	var fields Fields

	logger := New()
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)

	log(logger)

	err := json.Unmarshal(buffer.Bytes(), &fields)
	assert.Nil(t, err)

	assertions(fields)
}

func TestPrint(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Print("test")
	}, func(fields Fields) {
		assert.Equal(t, fields["msg"], "test")
		assert.Equal(t, fields["level"], "info")
	})
}

func TestPrintf(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Printf("%s", "test")
	}, func(fields Fields) {
		assert.Equal(t, fields["msg"], "test")
		assert.Equal(t, fields["level"], "info")
	})
}

func TestPrintln(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Println("test")
	}, func(fields Fields) {
		assert.Equal(t, fields["msg"], "test")
		assert.Equal(t, fields["level"], "info")
	})
}

func TestInfo(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Info("test")
	}, func(fields Fields) {
		assert.Equal(t, fields["msg"], "test")
		assert.Equal(t, fields["level"], "info")
	})
}

func TestWarn(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Warn("test")
	}, func(fields Fields) {
		assert.Equal(t, fields["msg"], "test")
		assert.Equal(t, fields["level"], "warning")
	})
}
