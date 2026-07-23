package logrus

import (
	"strings"
	"testing"
)

func TestTextFormatterNilEntry(t *testing.T) {
	f := &TextFormatter{DisableColors: true}
	_, err := f.Format(nil)
	if err == nil || !strings.Contains(err.Error(), "nil Entry") {
		t.Fatalf("got %v", err)
	}
}

func TestJSONFormatterNilEntry(t *testing.T) {
	f := &JSONFormatter{}
	_, err := f.Format(nil)
	if err == nil || !strings.Contains(err.Error(), "nil Entry") {
		t.Fatalf("got %v", err)
	}
}
