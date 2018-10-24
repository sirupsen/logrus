package testutils

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sirupsen/logrus"
)

// LogAndAssertJSON calls `log` function to log to a string buffer in JSON format,
// then passes the parsed output to `assertions` function for validation.
func LogAndAssertJSON(t *testing.T, log func(*logrus.Logger), assertions func(fields logrus.Fields)) {
	var buffer bytes.Buffer
	var fields logrus.Fields

	logger := logrus.New()
	logger.Out = &buffer
	logger.Formatter = new(logrus.JSONFormatter)

	log(logger)

	err := json.Unmarshal(buffer.Bytes(), &fields)
	require.Nil(t, err)

	assertions(fields)
}

// LogAndAssertText calls `log` function to log to a string buffer in text format,
// then passes the parsed output to `assertions` function for validation.
func LogAndAssertText(t *testing.T, log func(*logrus.Logger), assertions func(fields map[string]string)) {
	var buffer bytes.Buffer

	logger := logrus.New()
	logger.Out = &buffer
	logger.Formatter = &logrus.TextFormatter{
		DisableColors: true,
	}

	log(logger)

	fields := make(map[string]string)
	for _, kv := range strings.Split(buffer.String(), " ") {
		if !strings.Contains(kv, "=") {
			continue
		}
		kvArr := strings.Split(kv, "=")
		key := strings.TrimSpace(kvArr[0])
		val := kvArr[1]
		if kvArr[1][0] == '"' {
			var err error
			val, err = strconv.Unquote(val)
			require.NoError(t, err)
		}
		fields[key] = val
	}
	assertions(fields)
}
