package logrus

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEntryWithError(t *testing.T) {

	assert := assert.New(t)

	defer func() {
		ErrorKey = "error"
	}()

	err := fmt.Errorf("kaboom at layer %d", 4711)

	assert.Equal(err, WithError(err).Data["error"])

	logger := New()
	logger.Out = &bytes.Buffer{}
	entry := NewEntry(logger)

	assert.Equal(err, entry.WithError(err).Data["error"])

	ErrorKey = "err"

	assert.Equal(err, entry.WithError(err).Data["err"])

}

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

func TestEntryWithFields(t *testing.T) {
	logger := New()
	entry := &Entry{
		Time:    time.Date(2001, 2, 3, 4, 5, 6, 7, time.UTC),
		Message: "The message",
		Level:   DebugLevel,
		Data:    Fields{"foo": "bar"},
		Logger:  logger,
		Buffer:  bytes.NewBuffer([]byte{}),
	}
	newEntry := entry.WithFields(Fields{"baz": 42, "one": "more"})
	assert.Equal(t, entry.Time, newEntry.Time)
	assert.Equal(t, entry.Message, newEntry.Message)
	assert.Equal(t, entry.Level, newEntry.Level)
	assert.Equal(t, entry.Logger, newEntry.Logger)
	assert.NotEqual(t, entry.Buffer, newEntry.Buffer)

	value, ok := newEntry.Data["foo"]
	assert.True(t, ok)
	assert.Equal(t, "bar", value)

	value, ok = newEntry.Data["baz"]
	assert.True(t, ok)
	assert.Equal(t, 42, value)

	value, ok = newEntry.Data["one"]
	assert.True(t, ok)
	assert.Equal(t, "more", value)
}
