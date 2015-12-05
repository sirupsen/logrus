package logstash

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogstashFormatter(t *testing.T) {
	assert := assert.New(t)

	lf := LogstashFormatter{Type: "abc"}

	fields := logrus.Fields{
		"message":      "def",
		"level":        "ijk",
		"type":         "lmn",
		"one":          1,
		"pi":           3.14,
		"bool":         true,
		"custom_error": errors.New("test custom error"),
	}

	entry := logrus.WithFields(fields).WithError(errors.New("test error"))
	entry.Message = "msg"
	entry.Level = logrus.InfoLevel

	b, _ := lf.Format(entry)

	var data map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	_ = dec.Decode(&data)

	// base fields
	assert.Equal(json.Number("1"), data["@version"])
	assert.NotEmpty(data["@timestamp"])
	assert.Equal("abc", data["type"])
	assert.Equal("msg", data["message"])
	assert.Equal("info", data["level"])

	// error fields
	assert.Equal("test error", data[logrus.ErrorKey])
	assert.Equal("test custom error", data["custom_error"])

	// substituted fields
	assert.Equal("def", data["fields.message"])
	assert.Equal("ijk", data["fields.level"])
	assert.Equal("lmn", data["fields.type"])

	// formats
	assert.Equal(json.Number("1"), data["one"])
	assert.Equal(json.Number("3.14"), data["pi"])
	assert.Equal(true, data["bool"])
}
