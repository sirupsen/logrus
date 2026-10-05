package logrus_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestTextFormatterPreservesMessage(t *testing.T) {
	for _, mode := range []struct {
		name                                    string
		colors, disableTimestamp, fullTimestamp bool
	}{
		{name: "plain", disableTimestamp: true},
		{name: "colored_no_timestamp", colors: true, disableTimestamp: true},
		{name: "colored_relative_timestamp", colors: true},
		{name: "colored_full_timestamp", colors: true, fullTimestamp: true},
	} {
		for _, message := range []string{"message", "message\n", "message\n\n"} {
			t.Run(mode.name+"/"+strings.ReplaceAll(message, "\n", "_newline"), func(t *testing.T) {
				entry := &logrus.Entry{Logger: logrus.New(), Time: time.Unix(1, 0), Level: logrus.InfoLevel, Message: message, Data: logrus.Fields{"key": "value"}}
				formatter := &logrus.TextFormatter{ForceColors: mode.colors, DisableColors: !mode.colors, DisableTimestamp: mode.disableTimestamp, FullTimestamp: mode.fullTimestamp}
				first, err := formatter.Format(entry)
				if err != nil {
					t.Fatal(err)
				}
				first = bytes.Clone(first)
				if entry.Message != message {
					t.Errorf("Format changed Message: got %q, want %q", entry.Message, message)
				}
				second, err := formatter.Format(entry)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(first, second) {
					t.Errorf("repeated formatting changed output: %q != %q", first, second)
				}
				encoded, err := (&logrus.JSONFormatter{DisableTimestamp: true}).Format(entry)
				if err != nil {
					t.Fatal(err)
				}
				var data map[string]any
				if err = json.Unmarshal(encoded, &data); err != nil {
					t.Fatal(err)
				}
				if data["msg"] != message {
					t.Errorf("subsequent JSON lost original message: %q", data["msg"])
				}
				if mode.colors && !strings.Contains(string(first), fmt.Sprintf("%-44s ", strings.TrimSuffix(message, "\n"))) {
					t.Errorf("colored output did not retain existing newline handling: %q", first)
				}
			})
		}
	}
}

type coloredMessageHook struct {
	formatter logrus.TextFormatter
}

func (h *coloredMessageHook) Levels() []logrus.Level { return logrus.AllLevels }
func (h *coloredMessageHook) Fire(entry *logrus.Entry) error {
	_, err := h.formatter.Format(entry)
	return err
}

func TestTextFormatterHookPreservesJSONMessage(t *testing.T) {
	var output bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&output)
	logger.SetFormatter(&logrus.JSONFormatter{DisableTimestamp: true})
	logger.AddHook(&coloredMessageHook{formatter: logrus.TextFormatter{ForceColors: true, DisableTimestamp: true}})
	message := "message\n\n"
	logger.WithField("key", "value").Info(message)
	var data map[string]any
	if err := json.Unmarshal(output.Bytes(), &data); err != nil {
		t.Fatal(err)
	}
	if data["msg"] != message {
		t.Errorf("hook formatting changed JSON message: got %q, want %q", data["msg"], message)
	}
	if data["key"] != "value" {
		t.Errorf("hook formatting changed field: %v", data["key"])
	}
}
