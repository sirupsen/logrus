package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/sirupsen/logrus"
)

func LogAndAssertJSON(t *testing.T, log func(*Logger), assertions func(fields Fields)) {
	var buffer bytes.Buffer
	var fields Fields

	logger := New()
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)

	log(logger)

	err := json.Unmarshal(buffer.Bytes(), &fields)
	require.Nil(t, err)

	assertions(fields)
}

func LogAndAssertText(t *testing.T, log func(*Logger), assertions func(fields map[string]string)) {
	var buffer bytes.Buffer

	logger := New()
	logger.Out = &buffer
	logger.Formatter = &TextFormatter{
		DisableColors: true,
	}

	log(logger)

	fields := make(map[string]string)
	for _, kv := range strings.Split(buffer.String(), " ") {
		if !strings.Contains(kv, "=") {
			continue
		}
		kvArr := strings.Split(kv, "=")
		key := strings.TrimSpace(kvArr[0])
		val := kvArr[1]
		if kvArr[1][0] == '"' {
			var err error
			val, err = strconv.Unquote(val)
			require.NoError(t, err)
		}
		fields[key] = val
	}
	assertions(fields)
}

// AssertMapOfStringToUnorderdStringsEqualf ensures actual matches expected,
// without regard to order of the seperator delimited parts of each key's value
func AssertMapOfStringToUnorderdStringsEqualf(t *testing.T, seperator string, expected map[string][]string, actual map[string]string, format string, msgParts ...interface{}) {
	message := fmt.Sprintf(format, msgParts...)
	AssertMapOfStringToUnorderdStringsEqual(t, seperator, expected, actual, message)
}

// AssertMapOfStringToUnorderdStringsEqual ensures actual matches expected,
// without regard to order of the seperator delimited parts of each key's value
func AssertMapOfStringToUnorderdStringsEqual(t *testing.T, seperator string, expected map[string][]string, actual map[string]string, msgParts ...interface{}) {
	message := fmt.Sprint(msgParts...)
	if message != "" {
		message = ": " + message
	}
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
			"Expected %#v to be an expected key%s",
			actualKey,
			message)
		assert.Equalf(
			t,
			0,
			foundTimes,
			"Expected %#v to only appear once%s",
			actualKey,
			message)
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
			"Expected key %#v (with value %#v) to have same parts as %#v%s",
			actualKey,
			actualValue,
			expected[actualKey],
			message)
	}
	// Make sure no expected keys were missing (or duplicate, already
	// checked for.)
	for key, foundTimes := range expectedKeysFound {
		assert.Equalf(
			t,
			1,
			foundTimes,
			"Expected key %#v to be found 1 time (not %#v times)%s",
			key,
			foundTimes,
			message)
	}
}
