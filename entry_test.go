package logrus

import (
	"bytes"
	"fmt"
	"regexp"
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
	// Turning off timestamp and color marks parsing the result easier:
	logger.Formatter = &TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	}
	// Regular Expressions used later
	keyValPartsPattern := regexp.MustCompile(`([^\s"=]*|"(?:[^\\"]|\\.)*")=([^\s"=]*|"(?:[^\\"]|\\.)*")(?:\s|$)`)
	commaSepPattern := regexp.MustCompile(", ")
	quotedCommaSepPattern := regexp.MustCompile(`", "`)
	// This allows us to have multiple distinct tests, each with a new entry.
	// While it is really a good idea to break tests out into separate
	// independent units, we covered all the obvious cases with one multi-part
	// test. (Multi-part tests share an entry and test the entry state after
	// each change.)
	for _, independantTestCase := range [][]struct {
		Group          Fields
		ExpectedFields map[string][]string
	}{
		// This is the first (and currently only) independant test case:
		{
			// These are the multiple parts of the single test case:
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
						`"can not add field \"Fred\""`,
						`"can not add field \"George\""`,
						`"can not add field \"Ron\""`,
						`"can not add field \"Ginnie\""`,
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
						`"can not add field \"Fred\""`,
						`"can not add field \"George\""`,
						`"can not add field \"Ron\""`,
						`"can not add field \"Ginnie\""`,
						`"can not add field \"six\""`,
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
						`"can not add field \"Fred\""`,
						`"can not add field \"George\""`,
						`"can not add field \"Ron\""`,
						`"can not add field \"Ginnie\""`,
						`"can not add field \"six\""`,
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
						`"can not add field \"Fred\""`,
						`"can not add field \"George\""`,
						`"can not add field \"Ron\""`,
						`"can not add field \"Ginnie\""`,
						`"can not add field \"six\""`,
						`"can not add field \"seven\""`,
					},
				},
			},
		},
	} {
		// Independant tests have a new Entry, and a new expected final state.
		entry := NewEntry(logger)
		finalFieldsStr := ""
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

			// Capture logged output:
			entry.Info("baz")
			finalFieldsStr = outBuffer.String()

			// Parse the output into key-value pairs.
			keyValPairs := keyValPartsPattern.FindAllStringSubmatch(outBuffer.String(), -1)
			assert.True(t, len(keyValPairs) > 0, "Expect at least one field")

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

			// A map of how many times we see each expected key. At the end
			// these should all be 1.
			expectedKeysFound := make(map[string]int, len(expectedFields))
			for key := range expectedFields {
				expectedKeysFound[key] = 0
			}

			// each key - value pair
			for _, pairMatch := range keyValPairs {
				assert.Len(t, pairMatch, 3, "Expected Regexp to only have exactly two matching capture groups")

				// Enforce no duplicate or unexpected keys
				foundTimes, expected := expectedKeysFound[pairMatch[1]]
				assert.Truef(
					t,
					expected,
					"Expected %#v to be an expected key from output %#v",
					pairMatch[1],
					outBuffer.String())
				assert.Equalf(
					t,
					0,
					foundTimes,
					"Expected %#v to only appear once from output %#v",
					pairMatch[1],
					outBuffer.String())
				expectedKeysFound[pairMatch[1]]++

				// If a value contains `, `, it must be in double quotes.
				// To sidestep the issue of string termination when reordering
				// we replace all instance of `, ` with `", "`  so we can simply
				// sort them, glue them back together and revert all `", "`s
				// back to `, `
				valParts := strings.Split(commaSepPattern.ReplaceAllString(pairMatch[2], `", "`), ", ")
				sort.Strings(valParts)
				sort.Strings(expectedFields[pairMatch[1]])

				// Glue the string slices back together for comparison.
				assert.Equal(
					t,
					quotedCommaSepPattern.ReplaceAllString(strings.Join(expectedFields[pairMatch[1]], ", "), ", "),
					quotedCommaSepPattern.ReplaceAllString(strings.Join(valParts, ", "), ", "),
					"Expected key %#v (with value %#v) from output %#v to have same parts as %#v",
					pairMatch[1],
					pairMatch[2],
					outBuffer.String(),
					expectedFields[pairMatch[1]])
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
					outBuffer.String())
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
		assert.Equalf(
			t,
			finalFieldsStr,
			outBuffer.String(),
			"WithTime() not expected to modify fields or errors: %#v, previously %#v",
			outBuffer.String(),
			finalFieldsStr)
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
