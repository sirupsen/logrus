package logrus_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/sirupsen/logrus"
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

const (
	badMessage   = "this is going to panic"
	panicMessage = "this is broken"
)

type panickyHook struct{}

func (p *panickyHook) Levels() []Level {
	return []Level{InfoLevel}
}

func (p *panickyHook) Fire(entry *Entry) error {
	if entry.Message == badMessage {
		panic(panicMessage)
	}

	return nil
}

func TestEntryHooksPanic(t *testing.T) {
	logger := New()
	logger.Out = &bytes.Buffer{}
	logger.Level = InfoLevel
	logger.Hooks.Add(&panickyHook{})

	defer func() {
		p := recover()
		assert.NotNil(t, p)
		assert.Equal(t, panicMessage, p)

		entry := NewEntry(logger)
		entry.Info("another message")
	}()

	entry := NewEntry(logger)
	entry.Info(badMessage)
}

func TestEntryWithIncorrectField(t *testing.T) {
	fn := func() {}

	e := Entry{}
	eWithFunc := e.WithFields(Fields{"func": fn})
	eWithFuncPtr := e.WithFields(Fields{"funcPtr": &fn})

	assertEntryErrorField(t, eWithFunc, `can not add field "func"`)
	assertEntryErrorField(t, eWithFuncPtr, `can not add field "funcPtr"`)

	eWithFunc = eWithFunc.WithField("not_a_func", "it is a string")
	eWithFuncPtr = eWithFuncPtr.WithField("not_a_func", "it is a string")

	assertEntryErrorField(t, eWithFunc, `can not add field "func"`)
	assertEntryErrorField(t, eWithFuncPtr, `can not add field "funcPtr"`)

	eWithFunc = eWithFunc.WithTime(time.Now())
	eWithFuncPtr = eWithFuncPtr.WithTime(time.Now())

	assertEntryErrorField(t, eWithFunc, `can not add field "func"`)
	assertEntryErrorField(t, eWithFuncPtr, `can not add field "funcPtr"`)
}

func getEntryField(e *Entry, fiedlName string) (interface{}, error) {
	myFormatter := JSONFormatter{}

	jsonBuff, err := myFormatter.Format(e)
	if err != nil {
		return nil, err
	}
	var anyHash map[string]interface{}
	err = json.Unmarshal(jsonBuff, &anyHash)
	if err != nil {
		return nil, err
	}
	if val, ok := anyHash[fiedlName]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("Entry has no field named %#v", fiedlName)
}

func assertEntryErrorField(t *testing.T, e *Entry, expectedErrorText string) bool {
	assert := assert.New(t)
	actualErrorText := ""

	fieldValue, err := getEntryField(e, "logrus_error")
	if err != nil {
		return assert.Fail(fmt.Sprintf("Entry is missing the error field resulting in error: %s", err.Error()))
	}
	if fieldText, ok := fieldValue.(string); ok {
		actualErrorText = fieldText
	} else {
		return assert.Fail("Entry's error field is not a string")
	}
	return assert.Equalf(expectedErrorText, actualErrorText, "entry error field did not have the expected value.")
}
