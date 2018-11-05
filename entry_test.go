package logrus_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/internal/testutils"
)

func TestEntryWithError(t *testing.T) {

	assert := assert.New(t)

	defer func() {
		logrus.ErrorKey = "error"
	}()

	err := fmt.Errorf("kaboom at layer %d", 4711)

	assert.Equal(err, logrus.WithError(err).Data["error"])

	logger := logrus.New()
	logger.Out = &bytes.Buffer{}
	entry := logrus.NewEntry(logger)

	assert.Equal(err, entry.WithError(err).Data["error"])

	logrus.ErrorKey = "err"

	assert.Equal(err, entry.WithError(err).Data["err"])

}

func TestEntryPanicln(t *testing.T) {
	errBoom := fmt.Errorf("boom time")
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{DisableTimestamp: true, DisableColors: true}
	outBuffer := bytes.Buffer{}

	defer func() {
		p := recover()
		assert.NotNil(t, p)

		switch pVal := p.(type) {
		case *logrus.Entry:
			assert.Equal(t, "kaboom", pVal.Message)
			assert.Equal(t, errBoom, pVal.Data["err"])
			assert.Regexp(t, `^level=panic msg=kaboom err="boom time"`, outBuffer.String())
		default:
			t.Fatalf("want type *Entry, got %T: %#v", pVal, pVal)
		}
	}()

	logger.Out = &outBuffer
	entry := logrus.NewEntry(logger)
	entry.WithField("err", errBoom).Panicln("kaboom")
}

// TestEntryFormattingError tests proper handling of any field formatting errors
func TestEntryFormattingError(t *testing.T) {
	logger := logrus.New()
	// Turning off timestamp makes comparing the result easier:
	logger.Formatter = &logrus.JSONFormatter{DisableTimestamp: true}
	// Functions can't be serialized, so they aren't allowed as a field value.
	badFieldValue := func() {}
	// This allows us to have multiple distinct tests, each with a new entry.
	// While it is really a good idea to break tests out into separate
	// independent units, we covered all the obvious cases with one multi-part
	// test. (Multi-part tests share an entry and test the entry state after
	// each change.)
	for description, independantTestCase := range map[string][]struct {
		Group          logrus.Fields
		ExpectedFields map[string][]string
	}{
		"Fields with no errors": {
			{
				Group: logrus.Fields{
					"foo": "bar",
				},
				ExpectedFields: map[string][]string{
					"foo": {"bar"},
				},
			},
		},
		"Fields with one error": {
			{
				Group: logrus.Fields{
					"foo": badFieldValue,
				},
				ExpectedFields: map[string][]string{
					"logrus_error": {`can not add field "foo"`},
				},
			},
		},
		"Fields with multiple errors": {
			{
				Group: logrus.Fields{
					"foo": badFieldValue,
					"bar": badFieldValue,
					"baz": badFieldValue,
				},
				ExpectedFields: map[string][]string{
					"logrus_error": {
						`can not add field "foo"`,
						`can not add field "bar"`,
						`can not add field "baz"`,
					},
				},
			},
		},
		"Mixed fields with and without errors": {
			{
				Group: logrus.Fields{
					"apple":  badFieldValue,
					"banana": "yellow",
					"carrot": badFieldValue,
					"daisy":  "gerber",
				},
				ExpectedFields: map[string][]string{
					"logrus_error": {
						`can not add field "apple"`,
						`can not add field "carrot"`,
					},
					"banana": {"yellow"},
					"daisy":  {"gerber"},
				},
			},
		},
		"No errors followed by multiple errors": {
			{
				Group: logrus.Fields{
					"foo": "bar",
				},
				ExpectedFields: map[string][]string{
					"foo": {"bar"},
				},
			},
			{
				Group: logrus.Fields{
					"Fred":   badFieldValue,
					"George": badFieldValue,
					"Ron":    badFieldValue,
					"Ginnie": badFieldValue,
				},
				ExpectedFields: map[string][]string{
					"foo": {"bar"},
					"logrus_error": {
						`can not add field "Fred"`,
						`can not add field "George"`,
						`can not add field "Ron"`,
						`can not add field "Ginnie"`,
					},
				},
			},
		},
		"Compound example": {
			{
				Group: logrus.Fields{
					"foo": "bar",
				},
				ExpectedFields: map[string][]string{
					"foo": {"bar"},
				},
			},
			{
				Group: logrus.Fields{
					"Fred":   badFieldValue,
					"George": badFieldValue,
					"Ron":    badFieldValue,
					"Ginnie": badFieldValue,
				},
				ExpectedFields: map[string][]string{
					"foo": {"bar"},
					"logrus_error": {
						`can not add field "Fred"`,
						`can not add field "George"`,
						`can not add field "Ron"`,
						`can not add field "Ginnie"`,
					},
				},
			},
			{
				Group: logrus.Fields{
					"six": badFieldValue,
				},
				ExpectedFields: map[string][]string{
					"foo": {"bar"},
					"logrus_error": {
						`can not add field "Fred"`,
						`can not add field "George"`,
						`can not add field "Ron"`,
						`can not add field "Ginnie"`,
						`can not add field "six"`,
					},
				},
			},
			{
				Group: logrus.Fields{
					"red": "green",
				},
				ExpectedFields: map[string][]string{
					"foo": {"bar"},
					"red": {"green"},
					"logrus_error": {
						`can not add field "Fred"`,
						`can not add field "George"`,
						`can not add field "Ron"`,
						`can not add field "Ginnie"`,
						`can not add field "six"`,
					},
				},
			},
			{
				Group: logrus.Fields{
					"seven": badFieldValue,
				},
				ExpectedFields: map[string][]string{
					"foo": {"bar"},
					"red": {"green"},
					"logrus_error": {
						`can not add field "Fred"`,
						`can not add field "George"`,
						`can not add field "Ron"`,
						`can not add field "Ginnie"`,
						`can not add field "six"`,
						`can not add field "seven"`,
					},
				},
			},
		},
	} {
		// Independant tests have a new Entry, and a new expected final state.
		entry := logrus.NewEntry(logger)
		var lastExpected map[string][]string
		for _, fieldGroupTestPart := range independantTestCase {
			// Set outBuffer as a new "file" to log to.
			outBuffer := &bytes.Buffer{}
			logger.Out = outBuffer

			// ******** SYSTEM UNDER TEST ********
			// Multiple calls to the WithFields() method is the primary
			// thing being tested:
			entry = entry.WithFields(fieldGroupTestPart.Group)
			// Everything below here is analizing the Entry state after
			// WithFields() is called:

			// Capture and parse logged output:
			entry.Info("baz")
			outputMap := make(map[string]string, len(fieldGroupTestPart.ExpectedFields)+2)
			if err := json.Unmarshal(outBuffer.Bytes(), &outputMap); err != nil {
				assert.Fail(t, fmt.Sprintf("Failure unmarshalling logger output, %#v testing %s from output %#v", err.Error(), description, outBuffer.String()))
			} else {
				// Remove level and msg created by logrus:
				delete(outputMap, "level")
				delete(outputMap, "msg")
				lastExpected = fieldGroupTestPart.ExpectedFields

				testutils.AssertMapOfStringToUnorderdStringsEqualf(t, ", ", fieldGroupTestPart.ExpectedFields, outputMap, "testing %s from map %#v", description, outputMap)
			}

		}
		// Prep for testing WithTime(): New outBuffer log "file"
		outBuffer := &bytes.Buffer{}
		logger.Out = outBuffer

		// ******** SYSTEM UNDER TEST ********
		// This is the call really being tested
		entry = entry.WithTime(time.Now())

		// Ensure changing the timestamp (with timestamps hidden) didn't change
		// the output.
		entry.Info("baz")
		outputMap := make(map[string]string, len(lastExpected))
		if err := json.Unmarshal(outBuffer.Bytes(), &outputMap); err != nil {
			assert.Fail(t, fmt.Sprintf("Failure unmarshalling logger output, %#v testing %s from output %#v", err.Error(), description, outBuffer.String()))
		} else {
			// Remove level and msg created by logrus:
			delete(outputMap, "level")
			delete(outputMap, "msg")
			testutils.AssertMapOfStringToUnorderdStringsEqualf(t, ", ", lastExpected, outputMap, "testing %s from map %#v", description, outputMap)
		}
	}
}

func TestEntryPanicf(t *testing.T) {
	errBoom := fmt.Errorf("boom again")

	defer func() {
		p := recover()
		assert.NotNil(t, p)

		switch pVal := p.(type) {
		case *logrus.Entry:
			assert.Equal(t, "kaboom true", pVal.Message)
			assert.Equal(t, errBoom, pVal.Data["err"])
		default:
			t.Fatalf("want type *Entry, got %T: %#v", pVal, pVal)
		}
	}()

	logger := logrus.New()
	logger.Out = &bytes.Buffer{}
	entry := logrus.NewEntry(logger)
	entry.WithField("err", errBoom).Panicf("kaboom %v", true)
}

const (
	badMessage   = "this is going to panic"
	panicMessage = "this is broken"
)

type panickyHook struct{}

func (p *panickyHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.InfoLevel}
}

func (p *panickyHook) Fire(entry *logrus.Entry) error {
	if entry.Message == badMessage {
		panic(panicMessage)
	}

	return nil
}

func TestEntryHooksPanic(t *testing.T) {
	logger := logrus.New()
	logger.Out = &bytes.Buffer{}
	logger.Level = logrus.InfoLevel
	logger.Hooks.Add(&panickyHook{})

	defer func() {
		p := recover()
		assert.NotNil(t, p)
		assert.Equal(t, panicMessage, p)

		entry := logrus.NewEntry(logger)
		entry.Info("another message")
	}()

	entry := logrus.NewEntry(logger)
	entry.Info(badMessage)
}
