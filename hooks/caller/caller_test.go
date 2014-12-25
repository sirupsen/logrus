package logrus_caller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Sirupsen/logrus"
)

func LogAndAssertJSON(t *testing.T, log func(*logrus.Logger), assertions func(fields logrus.Fields)) {
	var buffer bytes.Buffer
	var fields logrus.Fields

	logger := logrus.New()
	logger.Hooks.Add(&CallerHook{})
	logger.Out = &buffer
	logger.Formatter = new(logrus.JSONFormatter)

	log(logger)

	err := json.Unmarshal(buffer.Bytes(), &fields)
	if err != nil {
		t.Error("Error unmarshaling log entry")
	}

	assertions(fields)
}

func TestCaller(t *testing.T) {
	LogAndAssertJSON(t, func(logger *logrus.Logger) {
		logger.Info("Hello World")
	}, func(fields logrus.Fields) {
		expected := "caller_test.go:33"

		if fields["caller"] != expected {
			t.Error(fmt.Sprintf("Caller was %s, expected %s", fields["caller"], expected))
		}
	})
}
