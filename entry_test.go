package logrus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
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
	// This allows us to have multiple distinct tests, each with a new entry.
	// While it is really a good idea to break tests out into separate
	// independent units, we covered all the obvious cases with one multi-part
	// test. (Multi-part tests share an entry and test the entry state after
	// each change.)
	for _, independantTestCase := range [][]struct {
		Group          Fields
		ExpectedFields map[string][]string
	}{
		// Fields with no errors:
		{
			{
				Group: Fields{
					"foo": "bar",
				},
				ExpectedFields: map[string][]string{
					"foo": {"bar"},
				},
			},
		},
		// Fields with one error:
		{
			{
				Group: Fields{
					"foo": func() {},
				},
				ExpectedFields: map[string][]string{
					"logrus_error": {`can not add field "foo"`},
				},
			},
		},
		// Fields with multiple errors:
		{
			{
				Group: Fields{
					"foo": func() {},
					"bar": func() {},
					"baz": func() {},
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
		// Mixed fields with and without errors:
		{
			{
				Group: Fields{
					"apple":  func() {},
					"banana": "yellow",
					"carrot": func() {},
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
		// no errors followed by multiple errors:
		{
			{
				Group: Fields{
					"foo": "bar",
				},
				ExpectedFields: map[string][]string{
					"foo": {"bar"},
				},
			},
			{
				Group: Fields{
					"Fred":   func() {},
					"George": func() {},
					"Ron":    func() {},
					"Ginnie": func() {},
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
		// Compound example
		{
			{
				Group: Fields{
					"foo": "bar",
				},
				ExpectedFields: map[string][]string{
					"foo": {"bar"},
				},
			},
			{
				Group: Fields{
					"Fred":   func() {},
					"George": func() {},
					"Ron":    func() {},
					"Ginnie": func() {},
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
				Group: Fields{
					"six": func() {},
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
				Group: Fields{
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
				Group: Fields{
					"seven": func() {},
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
		entry := NewEntry(logger)
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
				assert.Fail(t, fmt.Sprintf("Failure unmarshalling logger output, %#v from output %#v", err.Error(), outBuffer.String()))
			} else {
				assert.True(t, len(outputMap) > 0, "Expect at least one field")

				// Inject level=info and msg=baz into our expected data, as they
				// come from the Info() call above:
				expectedFields := make(map[string][]string, len(fieldGroupTestPart.ExpectedFields)+2)
				for _, fieldSets := range []map[string][]string{
					map[string][]string{"level": {"info"}, "msg": {"baz"}},
					fieldGroupTestPart.ExpectedFields,
				} {
					for key, strValParts := range fieldSets {
						expectedFields[key] = strValParts
					}
				}
				lastExpected = expectedFields

				AssertMapOfStringToUnorderdStringsEqual(t, ", ", expectedFields, outputMap, outBuffer.String())
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
			assert.Fail(t, fmt.Sprintf("Failure unmarshalling logger output, %#v from output %#v", err.Error(), outBuffer.String()))
		} else {
			assert.True(t, len(outputMap) > 0, "Expect at least one field")
			AssertMapOfStringToUnorderdStringsEqual(t, ", ", lastExpected, outputMap, outBuffer.String())
		}
	}
}

func AssertMapOfStringToUnorderdStringsEqual(t *testing.T, seperator string, expected map[string][]string, actual map[string]string, output string) {
	// A map of how many times we see each expected key. At the end
	// these should all be 1.
	expectedKeysFound := make(map[string]int, len(expected))
	for key := range expected {
		expectedKeysFound[key] = 0
	}

	// each key - value pair
	for actualKey, actualValue := range actual {

		// Enforce no duplicate or unexpected keys
		foundTimes, wasExpected := expectedKeysFound[actualKey]
		assert.Truef(
			t,
			wasExpected,
			"Expected %#v to be an expected key from output %#v",
			actualKey,
			output)
		assert.Equalf(
			t,
			0,
			foundTimes,
			"Expected %#v to only appear once from output %#v",
			actualKey,
			output)
		expectedKeysFound[actualKey]++

		// Split the value on `seperator` so it can be sorted and reassembeled
		valParts := strings.Split(actualValue, seperator)
		sort.Strings(valParts)
		sort.Strings(expected[actualKey])

		// Glue the string slices back together for comparison.
		assert.Equal(
			t,
			strings.Join(expected[actualKey], seperator),
			strings.Join(valParts, seperator),
			"Expected key %#v (with value %#v) from output %#v to have same parts as %#v",
			actualKey,
			actualValue,
			output,
			expected[actualKey])
	}
	// Make sure no expected keys were missing (or duplicate, already
	// checked for.)
	for key, foundTimes := range expectedKeysFound {
		assert.Equalf(
			t,
			1,
			foundTimes,
			"Expected key %#v to be found 1 time (not %#v times) in output %#v",
			key,
			foundTimes,
			output)
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
