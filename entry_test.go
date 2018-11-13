package logrus_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/sirupsen/logrus"
	. "github.com/sirupsen/logrus/internal/testutils"
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
	logger := New()
	logger.Formatter = &TextFormatter{DisableTimestamp: true, DisableColors: true}
	outBuffer := bytes.Buffer{}

	defer func() {
		p := recover()
		assert.NotNil(t, p)

		switch pVal := p.(type) {
		case *Entry:
			assert.Equal(t, "kaboom", pVal.Message)
			assert.Equal(t, errBoom, pVal.Data["err"])
			assert.Regexp(t, `^level=panic msg=kaboom err="boom time"`, outBuffer.String())
		default:
			t.Fatalf("want type *Entry, got %T: %#v", pVal, pVal)
		}
	}()

	logger.Out = &outBuffer
	entry := NewEntry(logger)
	entry.WithField("err", errBoom).Panicln("kaboom")
}

// TestEntryFormattingError tests proper handling of any field formatting errors
func TestEntryFormattingError(t *testing.T) {
	logger := New()
	// Turning off timestamp makes comparing the result easier:
	logger.Formatter = &JSONFormatter{DisableTimestamp: true}
	// Functions can't be serialized, so they aren't allowed as a field value.
	badFieldValue := func() {}
	// makeAssertion converts a string with comma seperated parts to an
	// UnorderedStringParts assertion
	makeAssertion := func(expected string) StringFieldAssertion {
		expectedParsed := NewUnorderedStringParts(expected, ", ")
		return StringFieldAssertion(func(actual string, key string, message string) bool {
			if message != "" {
				message = " " + message
			}
			return expectedParsed.AssertEqualString(t, actual, fmt.Sprintf(" for key %#v%s", key, message))
		})
	}
	type testCaseStep struct {
		Group                   Fields
		ExpectedFieldValueParts map[string]StringFieldAssertion
	}
	type testCase struct {
		Description string
		Steps       []testCaseStep
	}
	for _, independantTestCase := range []testCase{
		testCase{
			Description: "Fields with no errors",
			Steps: []testCaseStep{
				testCaseStep{
					Group: Fields{
						"foo": "bar",
					},
					ExpectedFieldValueParts: map[string]StringFieldAssertion{
						"foo": makeAssertion("bar"),
					},
				},
			},
		},
		testCase{
			Description: "Fields with one error",
			Steps: []testCaseStep{
				testCaseStep{
					Group: Fields{
						"foo": badFieldValue,
					},
					ExpectedFieldValueParts: map[string]StringFieldAssertion{
						"logrus_error": makeAssertion(`can not add field "foo"`),
					},
				},
			},
		},
		testCase{
			Description: "Fields with multiple errors",
			Steps: []testCaseStep{
				testCaseStep{
					Group: Fields{
						"foo": badFieldValue,
						"bar": badFieldValue,
						"baz": badFieldValue,
					},
					ExpectedFieldValueParts: map[string]StringFieldAssertion{
						"logrus_error": makeAssertion(`can not add field "foo", can not add field "bar", can not add field "baz"`),
					},
				},
			},
		},
		testCase{
			Description: "Mixed fields with and without errors",
			Steps: []testCaseStep{
				testCaseStep{
					Group: Fields{
						"apple":  badFieldValue,
						"banana": "yellow",
						"carrot": badFieldValue,
						"daisy":  "gerber",
					},
					ExpectedFieldValueParts: map[string]StringFieldAssertion{
						"logrus_error": makeAssertion(`can not add field "apple", can not add field "carrot"`),
						"banana":       makeAssertion("yellow"),
						"daisy":        makeAssertion("gerber"),
					},
				},
			},
		},
		testCase{
			Description: "No errors followed by multiple errors",
			Steps: []testCaseStep{
				testCaseStep{
					Group: Fields{
						"foo": "bar",
					},
					ExpectedFieldValueParts: map[string]StringFieldAssertion{
						"foo": makeAssertion("bar"),
					},
				},
				testCaseStep{
					Group: Fields{
						"Fred":   badFieldValue,
						"George": badFieldValue,
						"Ron":    badFieldValue,
						"Ginnie": badFieldValue,
					},
					ExpectedFieldValueParts: map[string]StringFieldAssertion{
						"foo":          makeAssertion("bar"),
						"logrus_error": makeAssertion(`can not add field "Fred", can not add field "George", can not add field "Ron", can not add field "Ginnie"`),
					},
				},
			},
		},
		testCase{
			Description: "Compound example",
			Steps: []testCaseStep{
				testCaseStep{
					Group: Fields{
						"foo": "bar",
					},
					ExpectedFieldValueParts: map[string]StringFieldAssertion{
						"foo": makeAssertion("bar"),
					},
				},
				testCaseStep{
					Group: Fields{
						"Fred":   badFieldValue,
						"George": badFieldValue,
						"Ron":    badFieldValue,
						"Ginnie": badFieldValue,
					},
					ExpectedFieldValueParts: map[string]StringFieldAssertion{
						"foo":          makeAssertion("bar"),
						"logrus_error": makeAssertion(`can not add field "Fred", can not add field "George", can not add field "Ron", can not add field "Ginnie"`),
					},
				},
				testCaseStep{
					Group: Fields{
						"six": badFieldValue,
					},
					ExpectedFieldValueParts: map[string]StringFieldAssertion{
						"foo":          makeAssertion("bar"),
						"logrus_error": makeAssertion(`can not add field "Fred", can not add field "George", can not add field "Ron", can not add field "Ginnie", can not add field "six"`),
					},
				},
				testCaseStep{
					Group: Fields{
						"red": "green",
					},
					ExpectedFieldValueParts: map[string]StringFieldAssertion{
						"foo":          makeAssertion("bar"),
						"red":          makeAssertion("green"),
						"logrus_error": makeAssertion(`can not add field "Fred", can not add field "George", can not add field "Ron", can not add field "Ginnie", can not add field "six"`),
					},
				},
				testCaseStep{
					Group: Fields{
						"seven": badFieldValue,
					},
					ExpectedFieldValueParts: map[string]StringFieldAssertion{
						"foo":          makeAssertion("bar"),
						"red":          makeAssertion("green"),
						"logrus_error": makeAssertion(`can not add field "Fred", can not add field "George", can not add field "Ron", can not add field "Ginnie", can not add field "six", can not add field "seven"`),
					},
				},
			},
		},
	} {
		// Independant tests have a new Entry, and a new expected final state.
		entry := NewEntry(logger)
		var lastExpected map[string]StringFieldAssertion
		for _, fieldGroupTestPart := range independantTestCase.Steps {
			// Set outBuffer as a new "file" to log to.
			outBuffer := &bytes.Buffer{}
			logger.Out = outBuffer

			// Multiple calls to the WithFields() method is the primary
			// thing being tested:
			entry = entry.WithFields(fieldGroupTestPart.Group)

			// Everything below here is analizing the Entry state after
			// WithFields() is called:

			// Capture and parse logged output:
			entry.Info("baz")
			outputMap := make(map[string]string, len(fieldGroupTestPart.ExpectedFieldValueParts)+2)
			if err := json.Unmarshal(outBuffer.Bytes(), &outputMap); err != nil {
				assert.Fail(t, fmt.Sprintf("Failure unmarshalling logger output, %#v testing %s from output %#v", err.Error(), independantTestCase.Description, outBuffer.String()))
			} else {
				// Remove level and msg created by logrus:
				delete(outputMap, "level")
				delete(outputMap, "msg")
				lastExpected = fieldGroupTestPart.ExpectedFieldValueParts

				ApplyAssertsToMapOfStringf(t, fieldGroupTestPart.ExpectedFieldValueParts, outputMap, "testing %s from map %#v", independantTestCase.Description, outputMap)
			}

		}
		// Prep for testing WithTime(): New outBuffer log "file"
		outBuffer := &bytes.Buffer{}
		logger.Out = outBuffer

		// This is the call really being tested
		entry = entry.WithTime(time.Now())

		// Ensure changing the timestamp (with timestamps hidden) didn't change
		// the output.
		entry.Info("baz")
		outputMap := make(map[string]string, len(lastExpected))
		if err := json.Unmarshal(outBuffer.Bytes(), &outputMap); err != nil {
			assert.Fail(t, fmt.Sprintf("Failure unmarshalling logger output, %#v testing %s from output %#v", err.Error(), independantTestCase.Description, outBuffer.String()))
		} else {
			// Remove level and msg created by logrus:
			delete(outputMap, "level")
			delete(outputMap, "msg")
			ApplyAssertsToMapOfStringf(t, lastExpected, outputMap, "testing %s from map %#v", independantTestCase.Description, outputMap)
		}
	}
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
