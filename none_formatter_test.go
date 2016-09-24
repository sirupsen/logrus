package logrus

import "testing"

// TestNoneFormat should return nil,nil for ignoring log output
func TestNoneFormat(t *testing.T) {
	nof := &NoneFormatter{}

	b, _ := nof.Format(WithField("test", "anything"))
	if b != nil {
		t.Errorf("none formatter returned an actual formatter?")
	}
}
