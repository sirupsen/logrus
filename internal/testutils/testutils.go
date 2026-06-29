package testutils

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/require"
)

func LogAndAssertJSON(t *testing.T, log func(*logrus.Logger), assertions func(logrus.Fields)) {
	var buffer bytes.Buffer
	logger := logrus.New()
	logger.Out = &buffer
	logger.Formatter = new(logrus.JSONFormatter)

	log(logger)

	var fields logrus.Fields
	err := json.Unmarshal(buffer.Bytes(), &fields)
	require.NoError(t, err)

	assertions(fields)
}

func LogAndAssertText(t *testing.T, log func(*logrus.Logger), assertions func(fields map[string]string)) {
	var buffer bytes.Buffer
	logger := logrus.New()
	logger.Out = &buffer
	logger.Formatter = &logrus.TextFormatter{
		DisableColors: true,
	}

	log(logger)

	fields := make(map[string]string)
	for _, kv := range strings.Fields(strings.TrimRight(buffer.String(), "\n")) {
		key, val, ok := strings.Cut(kv, "=")
		if !ok {
			continue
		}
		if strings.HasPrefix(val, `"`) {
			var err error
			val, err = strconv.Unquote(val)
			require.NoError(t, err)
		}
		key = strings.TrimSpace(key)
		fields[key] = val
	}
	assertions(fields)
}
