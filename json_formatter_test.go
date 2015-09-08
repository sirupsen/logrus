package logrus

import (
	"bytes"
	"encoding/json"
	"errors"

	"testing"
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

func TestGolfsWithJsonFormatter(t *testing.T) {
	p := &Person{
		Name:  "Bruce",
		Alias: "Batman",
		Hideout: &Hideout{
			Name:        "JLU Tower",
			DimensionId: 52,
		},
	}

	jf := &JSONFormatter{}
	b, err := jf.Format(&Entry{
		Message: "the dark knight", Data: Fields{"hero": p}})
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	if bytes.Index(b, ([]byte)(`"hero.name":"Bruce"`)) < 0 {
		t.Fatalf(`missing "hero.name":"Bruce"`)
	}

	if bytes.Index(b, ([]byte)(`"hero.alias":"Batman"`)) < 0 {
		t.Fatalf(`missing "hero.alias":"Batman"`)
	}

	if bytes.Index(b, ([]byte)(`"hero.hideout.name":"JLU Tower"`)) < 0 {
		t.Fatalf(`missing "hero.hideout.name":"JLU Tower"`)
	}

	if bytes.Index(b, ([]byte)(`"hero.hideout.dimensionId":52`)) < 0 {
		t.Fatalf(`missing "hero.hideout.dimensionId":52`)
	}
}

func TestGolfsWithJsonFormatterAndNonGolfer(t *testing.T) {
	h := &Hideout{
		Name:        "JLU Tower",
		DimensionId: 52,
	}

	jf := &JSONFormatter{}
	b, err := jf.Format(&Entry{
		Message: "secret base", Data: Fields{"hideout": h}})
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	t.Log(string(b))

	if bytes.Index(b, ([]byte)(`"hideout":{"Name":"JLU Tower","DimensionId":52}`)) < 0 {
		t.Fatalf(`missing "hideout":{"Name":"JLU Tower","DimensionId":52}`)
	}
}
