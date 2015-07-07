package logrus

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryPanicln(t *testing.T) {
	errBoom := fmt.Errorf("boom time")

	defer func() {
		p := recover()
		assert.NotNil(t, p)

		switch pVal := p.(type) {
		case *Entry:
			assert.Equal(t, "kaboom", pVal.Message)
			assert.Equal(t, errBoom, pVal.Data["err"])
		default:
			t.Fatalf("want type *Entry, got %T: %#v", pVal, pVal)
		}
	}()

	logger := New()
	logger.Out = &bytes.Buffer{}
	entry := NewEntry(logger)
	entry.WithField("err", errBoom).Panicln("kaboom")
}

func TestEntryPanicf(t *testing.T) {
	errBoom := fmt.Errorf("boom again")

	defer func() {
		p := recover()
		assert.NotNil(t, p)

		switch pVal := p.(type) {
		case *Entry:
			assert.Equal(t, "kaboom true", pVal.Message)
			assert.Equal(t, errBoom, pVal.Data["err"])
		default:
			t.Fatalf("want type *Entry, got %T: %#v", pVal, pVal)
		}
	}()

	logger := New()
	logger.Out = &bytes.Buffer{}
	entry := NewEntry(logger)
	entry.WithField("err", errBoom).Panicf("kaboom %v", true)
}

func TestEntryLogLevel(t *testing.T) {
  	out := &bytes.Buffer{}
  	logger := New()
  	logger.Out = out
  	logger.Level = DebugLevel
  	entry := NewEntry(logger)
  	assert.Equal(t, DebugLevel, entry.Level)
  	entry.Level = WarnLevel
  	entry.Info("it should not be displayed")
  	assert.Equal(t, "", out.String())
  	entry.Warn("it should be displayed")
  	assert.Contains(t, out.String(), "it should be displayed")
}
