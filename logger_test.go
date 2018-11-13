package logrus_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/sirupsen/logrus"
)

func TestFieldValueError(t *testing.T) {
	buf := &bytes.Buffer{}
	l := &Logger{
		Out:       buf,
		Formatter: new(JSONFormatter),
		Hooks:     make(LevelHooks),
		Level:     DebugLevel,
	}
	l.WithField("func", func() {}).Info("test")
	var data map[string]interface{}
	json.Unmarshal(buf.Bytes(), &data)
	_, ok := data[FieldKeyLogrusError]
	require.True(t, ok)
}

func TestNoFieldValueError(t *testing.T) {
	buf := &bytes.Buffer{}
	l := &Logger{
		Out:       buf,
		Formatter: new(JSONFormatter),
		Hooks:     make(LevelHooks),
		Level:     DebugLevel,
	}
	l.WithField("str", "str").Info("test")
	var data map[string]interface{}
	json.Unmarshal(buf.Bytes(), &data)
	_, ok := data[FieldKeyLogrusError]
	require.False(t, ok)
}
