package logrus

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestIssue1597DisableTimestampKeepsUserTimeField(t *testing.T) {
	entry := &Entry{
		Logger:  New(),
		Time:    time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
		Level:   InfoLevel,
		Message: "event",
		Data:    Fields{"time": "from application"},
	}

	t.Run("json default time key", func(t *testing.T) {
		f := &JSONFormatter{DisableTimestamp: true}
		b, err := f.Format(entry)
		if err != nil {
			t.Fatal(err)
		}
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			t.Fatal(err)
		}
		if got := m["time"]; got != "from application" {
			t.Errorf("time = %v, want %q (output: %s)", got, "from application", b)
		}
		if _, ok := m["fields.time"]; ok {
			t.Errorf("unexpected fields.time key (output: %s)", b)
		}
	})

	t.Run("json custom time key", func(t *testing.T) {
		f := &JSONFormatter{
			DisableTimestamp: true,
			FieldMap:         FieldMap{FieldKeyTime: "@timestamp"},
		}
		e := *entry
		e.Data = Fields{"@timestamp": "from application"}
		b, err := f.Format(&e)
		if err != nil {
			t.Fatal(err)
		}
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			t.Fatal(err)
		}
		if got := m["@timestamp"]; got != "from application" {
			t.Errorf("@timestamp = %v, want %q (output: %s)", got, "from application", b)
		}
		if _, ok := m["fields.@timestamp"]; ok {
			t.Errorf("unexpected fields.@timestamp key (output: %s)", b)
		}
	})

	t.Run("text default time key", func(t *testing.T) {
		f := &TextFormatter{DisableTimestamp: true, DisableColors: true}
		b, err := f.Format(entry)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Contains(b, []byte("time=\"from application\"")) {
			t.Errorf("want user time field preserved, got %q", b)
		}
		if bytes.Contains(b, []byte("fields.time")) {
			t.Errorf("unexpected fields.time in %q", b)
		}
	})

	t.Run("timestamps enabled still prefixes collision", func(t *testing.T) {
		f := &JSONFormatter{}
		b, err := f.Format(entry)
		if err != nil {
			t.Fatal(err)
		}
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			t.Fatal(err)
		}
		if got := m["fields.time"]; got != "from application" {
			t.Errorf("fields.time = %v, want %q (output: %s)", got, "from application", b)
		}
		if _, ok := m["time"]; !ok {
			t.Errorf("want automatic time key (output: %s)", b)
		}
	})
}
