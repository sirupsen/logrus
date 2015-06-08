package defaultfields

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestDefaultFieldLogged(t *testing.T) {

	var buffer bytes.Buffer
	var fields logrus.Fields

	log := logrus.New()
	log.Out = &buffer
	log.Formatter = new(logrus.JSONFormatter)

	hook := NewDefaultFieldsHook([]logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})

	defaultKey := "default"
	expected := "Test Value"
	hook.AddDefaultField(defaultKey, expected)
	log.Hooks.Add(hook)

	log.WithField("another", "value").Info("A message")

	err := json.Unmarshal(buffer.Bytes(), &fields)
	if err != nil {
		t.Error("Cannot unmarshal JSON logged output")
	}

	actual, ok := fields[defaultKey]
	if !ok {
		t.Error("Default key not in logged output")
	}

	if expected != actual {
		t.Error("Default value does not match expected:", expected, actual)
	}
}
