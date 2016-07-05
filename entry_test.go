package logrus

import (
	"bytes"
	"fmt"
	"testing"

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

func TestEntryFatal(t *testing.T) {
	assert := assert.New(t)

	exiter := new(exiter)
	logger := New()
	logger.Out = &bytes.Buffer{}
	logger.Exit = exiter.Exit
	entry := NewEntry(logger)

	entry.Fatal("kaboom")

	assert.Equal(1, exiter.exitCode)
}

func TestEntryFatalf(t *testing.T) {
	assert := assert.New(t)

	exiter := new(exiter)
	logger := New()
	logger.Out = &bytes.Buffer{}
	logger.Exit = exiter.Exit
	entry := NewEntry(logger)

	entry.Fatalf("kaboom")

	assert.Equal(1, exiter.exitCode)
}

func TestEntryFatalln(t *testing.T) {
	assert := assert.New(t)

	exiter := new(exiter)
	logger := New()
	logger.Out = &bytes.Buffer{}
	logger.Exit = exiter.Exit
	entry := NewEntry(logger)

	entry.Fatalln("kaboom")

	assert.Equal(1, exiter.exitCode)
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
