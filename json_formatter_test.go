package logrus_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestErrorNotLost(t *testing.T) {
	formatter := &logrus.JSONFormatter{}

	b, err := formatter.Format(logrus.WithField("error", errors.New("wild walrus")))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]any)
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["error"] != "wild walrus" {
		t.Fatal("Error field not set")
	}
}

func TestErrorNotLostOnFieldNotNamedError(t *testing.T) {
	formatter := &logrus.JSONFormatter{}

	b, err := formatter.Format(logrus.WithField("omg", errors.New("wild walrus")))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]any)
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["omg"] != "wild walrus" {
		t.Fatal("Error field not set")
	}
}

func TestFieldClashWithTime(t *testing.T) {
	formatter := &logrus.JSONFormatter{}

	b, err := formatter.Format(logrus.WithField("time", "right now!"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]any)
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.time"] != "right now!" {
		t.Fatal("fields.time not set to original time field")
	}

	if entry["time"] != "0001-01-01T00:00:00Z" {
		t.Fatal("time field not set to current time, was: ", entry["time"])
	}
}

func TestFieldClashWithMsg(t *testing.T) {
	formatter := &logrus.JSONFormatter{}

	b, err := formatter.Format(logrus.WithField("msg", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]any)
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.msg"] != "something" {
		t.Fatal("fields.msg not set to original msg field")
	}
}

func TestFieldClashWithLevel(t *testing.T) {
	formatter := &logrus.JSONFormatter{}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]any)
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.level"] != "something" {
		t.Fatal("fields.level not set to original level field")
	}
}

func TestFieldClashWithRemappedFields(t *testing.T) {
	formatter := &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyMsg:   "@message",
		},
	}

	b, err := formatter.Format(logrus.WithFields(logrus.Fields{
		"@timestamp": "@timestamp",
		"@level":     "@level",
		"@message":   "@message",
		"timestamp":  "timestamp",
		"level":      "level",
		"msg":        "msg",
	}))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]any)
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	for _, field := range []string{"timestamp", "level", "msg"} {
		if entry[field] != field {
			t.Errorf("Expected field %v to be untouched; got %v", field, entry[field])
		}

		remappedKey := fmt.Sprintf("fields.%s", field)
		if remapped, ok := entry[remappedKey]; ok {
			t.Errorf("Expected %s to be empty; got %v", remappedKey, remapped)
		}
	}

	for _, field := range []string{"@timestamp", "@level", "@message"} {
		if entry[field] == field {
			t.Errorf("Expected field %v to be mapped to an Entry value", field)
		}

		remappedKey := fmt.Sprintf("fields.%s", field)
		if remapped, ok := entry[remappedKey]; ok {
			if remapped != field {
				t.Errorf("Expected field %v to be copied to %s; got %v", field, remappedKey, remapped)
			}
		} else {
			t.Errorf("Expected field %v to be copied to %s; was absent", field, remappedKey)
		}
	}
}

func TestFieldsInNestedDictionary(t *testing.T) {
	formatter := &logrus.JSONFormatter{
		DataKey: "args",
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"level": "level",
		"test":  "test",
	})
	logEntry.Level = logrus.InfoLevel

	b, err := formatter.Format(logEntry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]any)
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	args := entry["args"].(map[string]any)

	for _, field := range []string{"test", "level"} {
		if value, present := args[field]; !present || value != field {
			t.Errorf("Expected field %v to be present under 'args'; untouched", field)
		}
	}

	for _, field := range []string{"test", "fields.level"} {
		if _, present := entry[field]; present {
			t.Errorf("Expected field %v not to be present at top level", field)
		}
	}

	// with nested object, "level" shouldn't clash
	if entry["level"] != "info" {
		t.Errorf("Expected 'level' field to contain 'info'")
	}
}

func TestJSONEntryEndsWithNewline(t *testing.T) {
	formatter := &logrus.JSONFormatter{}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	if b[len(b)-1] != '\n' {
		t.Fatal("Expected JSON log entry to end with a newline")
	}
}

func TestJSONMessageKey(t *testing.T) {
	formatter := &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg: "message",
		},
	}

	b, err := formatter.Format(&logrus.Entry{Message: "oh hai"})
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, `"message":"oh hai"`) {
		t.Fatal("Expected JSON to format message key")
	}
}

func TestJSONLevelKey(t *testing.T) {
	formatter := &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "somelevel",
		},
	}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, "somelevel") {
		t.Fatal("Expected JSON to format level key")
	}
}

func TestJSONTimeKey(t *testing.T) {
	formatter := &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "timeywimey",
		},
	}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, "timeywimey") {
		t.Fatal("Expected JSON to format time key")
	}
}

func TestFieldDoesNotClashWithCaller(t *testing.T) {
	logrus.SetReportCaller(false)
	formatter := &logrus.JSONFormatter{}

	b, err := formatter.Format(logrus.WithField("func", "howdy pardner"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]any)
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["func"] != "howdy pardner" {
		t.Fatal("func field replaced when ReportCaller=false")
	}
}

func TestFieldClashWithCaller(t *testing.T) {
	logrus.SetReportCaller(true)
	formatter := &logrus.JSONFormatter{}
	e := logrus.WithField("func", "howdy pardner")
	e.Caller = &runtime.Frame{Function: "somefunc"}
	b, err := formatter.Format(e)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]any)
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.func"] != "howdy pardner" {
		t.Fatalf("fields.func not set to original func field when ReportCaller=true (got '%s')",
			entry["fields.func"])
	}

	if entry["func"] != "somefunc" {
		t.Fatalf("func not set as expected when ReportCaller=true (got '%s')",
			entry["func"])
	}

	logrus.SetReportCaller(false) // return to default value
}

func TestJSONDisableTimestamp(t *testing.T) {
	formatter := &logrus.JSONFormatter{
		DisableTimestamp: true,
	}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if strings.Contains(s, logrus.FieldKeyTime) {
		t.Error("Did not prevent timestamp", s)
	}
}

func TestJSONEnableTimestamp(t *testing.T) {
	formatter := &logrus.JSONFormatter{}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, logrus.FieldKeyTime) {
		t.Error("Timestamp not present", s)
	}
}

func TestJSONDisableHTMLEscape(t *testing.T) {
	formatter := &logrus.JSONFormatter{DisableHTMLEscape: true}

	b, err := formatter.Format(&logrus.Entry{Message: "& < >"})
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, "& < >") {
		t.Error("Message should not be HTML escaped", s)
	}
}

func TestJSONEnableHTMLEscape(t *testing.T) {
	formatter := &logrus.JSONFormatter{}

	b, err := formatter.Format(&logrus.Entry{Message: "& < >"})
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, `\u0026 \u003c \u003e`) {
		t.Error("Message should be HTML escaped", s)
	}
}
