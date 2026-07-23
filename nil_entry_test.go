package logrus

import "testing"

func TestNilEntryWithField(t *testing.T) {
	var e *Entry
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panicked: %v", r)
		}
	}()
	if e.WithField("a", 1) != nil {
		t.Fatal("want nil")
	}
	if e.WithFields(Fields{"a": 1}) != nil {
		t.Fatal("want nil")
	}
}
