package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/sirupsen/logrus"
)

// LogAndAssertJSON calls `log` function to log to a string buffer in JSON format,
// then passes the parsed output to `assertions` function for validation.
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

// StringFieldAssertion asserts some rules about a value of a struct field or
// map value
type StringFieldAssertion func(string, string, string) bool

// ApplyAssertsToMapOfStringf applies a map of assertion functions to all values
// in actual. Asserts that assertions and actual have the same set of keys.
func ApplyAssertsToMapOfStringf(t *testing.T, assertions map[string]StringFieldAssertion, actual map[string]string, msg string, args ...interface{}) bool {
	message := fmt.Sprintf(msg, args...)
	return ApplyAssertsToMapOfString(t, assertions, actual, message)
}

// ApplyAssertsToMapOfString applies a map of assertion functions to all values
// in actual. Asserts that assertions and actual have the same set of keys.
func ApplyAssertsToMapOfString(t *testing.T, assertions map[string]StringFieldAssertion, actual map[string]string, msgAndArgs ...interface{}) bool {
	message := fmt.Sprint(msgAndArgs...)
	messageSuffix := message
	if messageSuffix != "" {
		messageSuffix = ": " + messageSuffix
	}
	// A map of which expected key we saw. At the end
	// these should all be true.
	expectedKeysFound := make(map[string]bool, len(assertions))
	for key := range assertions {
		expectedKeysFound[key] = false
	}

	result := true
	// each key - value pair
	for actualKey, actualValue := range actual {

		// Enforce no unexpected keys
		_, wasExpected := expectedKeysFound[actualKey]
		assert.Truef(
			t,
			wasExpected,
			"Expected %#v to be an expected key%s",
			actualKey,
			messageSuffix)
		result = result && wasExpected
		expectedKeysFound[actualKey] = true

		// Apply the assertion
		result = result && assertions[actualKey](actualValue, actualKey, message)
	}
	// Make sure no expected keys were missing.
	for key, found := range expectedKeysFound {
		assert.Truef(
			t,
			found,
			"Expected key %#v to be present%s",
			key,
			messageSuffix)
		result = result && found
	}
	return result
}
