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

func TestHasFieldsWithJsonFormatter(t *testing.T) {
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

	if bytes.Index(b, ([]byte)(`"name":"Bruce"`)) < 0 {
		t.Fatalf(`missing "name":"Bruce"`)
	}

	if bytes.Index(b, ([]byte)(`"alias":"Batman"`)) < 0 {
		t.Fatalf(`missing "alias":"Batman"`)
	}

	if bytes.Index(b, ([]byte)(`"hideout.name":"JLU Tower"`)) < 0 {
		t.Fatalf(`missing "hideout.name":"JLU Tower"`)
	}

	if bytes.Index(b, ([]byte)(`"hideout.dimensionId":52`)) < 0 {
		t.Fatalf(`missing "hideout.dimensionId":52`)
	}
}

func TestHasTypeFieldsExceptWithJsonFormatter(t *testing.T) {
	p := &Person{
		Name:  "Bruce",
		Alias: "Batman",
		Hideout: &Hideout{
			Name:        "JLU Tower",
			DimensionId: 52,
			except:      []string{"dimensionId"},
		},
		useTypeFields: true,
		except:        []string{"name"},
	}

	jf := &JSONFormatter{}
	b, err := jf.Format(&Entry{
		Message: "the dark knight", Data: Fields{"hero": p}})
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	if bytes.Index(b, ([]byte)(`"name":"Bruce"`)) >= 0 {
		t.Fatalf(`has "name":"Bruce"`)
	}

	if bytes.Index(b, ([]byte)(`"alias":"Batman"`)) < 0 {
		t.Fatalf(`missing "alias":"Batman"`)
	}

	if bytes.Index(b, ([]byte)(`"hideout.name":"JLU Tower"`)) < 0 {
		t.Fatalf(`missing "hideout.name":"JLU Tower"`)
	}

	if bytes.Index(b, ([]byte)(`"hideout.dimensionId":52`)) >= 0 {
		t.Fatalf(`has "hideout.dimensionId":52`)
	}
}
