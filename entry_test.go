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
	entry1 := &Entry{
		Logger: StandardLogger(),
		Time:   time.Now(),
		Level:  InfoLevel,
	}

	entry2 := entry1.WithFields(Fields{
		"key1": "value1",
		"key2": "value2",
	})

	// Check everything but Data is still the same
	assert.Equal(t, entry1.Logger, entry2.Logger)
	assert.Equal(t, entry1.Time, entry2.Time)
	assert.Equal(t, entry1.Level, entry2.Level)
	assert.Equal(t, entry1.Message, entry2.Message)
	assert.Equal(t, entry1.Buffer, entry2.Buffer)

	entry3 := entry2.WithFields(Fields{
		"key3": "value4",
		"key4": "value4",
	})
	// Check everything but Data is still the same
	assert.Equal(t, entry1.Logger, entry3.Logger)
	assert.Equal(t, entry1.Time, entry3.Time)
	assert.Equal(t, entry1.Level, entry3.Level)
	assert.Equal(t, entry1.Message, entry3.Message)
	assert.Equal(t, entry1.Buffer, entry3.Buffer)

}
