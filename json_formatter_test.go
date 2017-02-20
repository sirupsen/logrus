package logrus

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestErrorNotLost(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("error", errors.New("wild walrus")))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["error"] != "wild walrus" {
		t.Fatal("Error field not set")
	}
}

func TestErrorNotLostOnFieldNotNamedError(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("omg", errors.New("wild walrus")))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["omg"] != "wild walrus" {
		t.Fatal("Error field not set")
	}
}

func TestFieldClashWithTime(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("time", "right now!"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
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

func TestTimestampInUTC(t *testing.T) {
	fiveHours := int(time.Hour/time.Second) * 5
	est := time.FixedZone("EST", fiveHours)
	time, err := time.ParseInLocation("2006-Jan-02", "2012-Jul-09", est)
	if err != nil {
		t.Fatal("Unable to Parse time/location", err)
	}

	entry := &Entry{}
	entry.Time = time

	formatter := &JSONFormatter{TimestampInUTC: true}
	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entryWithJSON := make(map[string]interface{})
	err = json.Unmarshal(b, &entryWithJSON)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}
	if entryWithJSON["time"] != "2012-07-08T19:00:00Z" {
		t.Fatal("Time was not converted to UTC")
	}
}

func TestFieldClashWithMsg(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("msg", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.msg"] != "something" {
		t.Fatal("fields.msg not set to original msg field")
	}
}

func TestFieldClashWithLevel(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.level"] != "something" {
		t.Fatal("fields.level not set to original level field")
	}
}

func TestJSONEntryEndsWithNewline(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	if b[len(b)-1] != '\n' {
		t.Fatal("Expected JSON log entry to end with a newline")
	}
}

func TestJSONMessageKey(t *testing.T) {
	formatter := &JSONFormatter{
		FieldMap: FieldMap{
			FieldKeyMsg: "message",
		},
	}

	b, err := formatter.Format(&Entry{Message: "oh hai"})
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !(strings.Contains(s, "message") && strings.Contains(s, "oh hai")) {
		t.Fatal("Expected JSON to format message key")
	}
}

func TestJSONLevelKey(t *testing.T) {
	formatter := &JSONFormatter{
		FieldMap: FieldMap{
			FieldKeyLevel: "somelevel",
		},
	}

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, "somelevel") {
		t.Fatal("Expected JSON to format level key")
	}
}

func TestJSONTimeKey(t *testing.T) {
	formatter := &JSONFormatter{
		FieldMap: FieldMap{
			FieldKeyTime: "timeywimey",
		},
	}

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, "timeywimey") {
		t.Fatal("Expected JSON to format time key")
	}
}

func TestJSONDisableTimestamp(t *testing.T) {
	formatter := &JSONFormatter{
		DisableTimestamp: true,
	}

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if strings.Contains(s, FieldKeyTime) {
		t.Error("Did not prevent timestamp", s)
	}
}

func TestJSONEnableTimestamp(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, FieldKeyTime) {
		t.Error("Timestamp not present", s)
	}
}
