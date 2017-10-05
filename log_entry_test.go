package logrus

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"errors"
	"github.com/stretchr/testify/assert"
)

func logEntryAndAssertJSON(loggerLevel Level, log func(*LogEntry), assertions func(Fields, *LogEntry)) {
	var buffer bytes.Buffer
	var fields Fields

	logger := NewLogger(loggerLevel)
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)

	entry := NewLogEntry(logger)

	log(entry)

	json.Unmarshal(buffer.Bytes(), &fields)

	assertions(fields, entry)
}

func logEntryAndAssertText(t *testing.T, loggerLevel Level, log func(*LogEntry), assertions func(Fields, *LogEntry)) {
	t.Helper()

	var buffer bytes.Buffer
	logger := NewLogger(loggerLevel)
	logger.Out = &buffer
	logger.Formatter = &TextFormatter{
		DisableColors: true,
	}

	entry := NewLogEntry(logger)

	log(entry)

	fields := make(Fields)
	for _, kv := range strings.Split(strings.TrimSpace(buffer.String()), " ") {
		if !strings.Contains(kv, "=") {
			continue
		}
		kvArr := strings.Split(kv, "=")

		key := strings.TrimSpace(kvArr[0])
		val := kvArr[1]
		if kvArr[1][0] == '"' {
			var err error
			val, err = strconv.Unquote(val)
			assert.NoError(t, err)
		}
		fields[key] = val
	}

	assertions(fields, entry)
}

func TestEntryLoggingWithTextFormatter(t *testing.T) {
	testCases := []struct {
		title       string
		loggerLevel Level
		entryLevel  Level
		message     string
	}{
		{
			title:       "entry_with_the_same_level_as_log_level_should_log",
			loggerLevel: DebugLevel,
			entryLevel:  DebugLevel,
			message:     "message",
		},
		{
			title:       "entry_with_the_level_lower_than_the_log_level_should_log",
			loggerLevel: DebugLevel,
			entryLevel:  InfoLevel,
			message:     "message",
		},
		{
			title:       "entry_with_the_level_higher_than_the_log_level_should_not_log",
			loggerLevel: InfoLevel,
			entryLevel:  DebugLevel,
			message:     "message",
		},
	}
	assrt := assert.New(t)
	for _, tc := range testCases {
		shouldLog := tc.entryLevel <= tc.loggerLevel
		t.Run(tc.title, func(t *testing.T) {
			logEntryAndAssertText(t, tc.loggerLevel,
				func(entry *LogEntry) {
					entry.AsLevel(tc.entryLevel).Write(tc.message)
				},
				func(fields Fields, entry *LogEntry) {
					assrt.Equal(tc.loggerLevel, entry.Level)
					assrt.Equal(tc.loggerLevel, entry.Logger.Level)
					msg, ok := fields["msg"]
					if shouldLog {
						if !ok {
							t.Error("Failed to retrieve the message. Nothing was logged")
						}
						if logged, ok := checkLoggedField(tc.message, msg); !ok {
							t.Errorf("expected %s, received '%v'", tc.message, logged)
						}
						return
					}
					if ok {
						t.Errorf("we shouldn't have logged anything but the output was %v", fields)
					}
				})
		})
	}
}

func TestEntryLoggingWithJSONFormatter(t *testing.T) {
	testCases := []struct {
		title       string
		loggerLevel Level
		entryLevel  Level
		message     string
	}{
		{
			title:       "entry_with_the_same_level_as_log_level_should_log",
			loggerLevel: DebugLevel,
			entryLevel:  DebugLevel,
			message:     "log me",
		},
		{
			title:       "entry_with_the_level_lower_than_the_log_level_should_log",
			loggerLevel: DebugLevel,
			entryLevel:  InfoLevel,
			message:     "log me",
		},
		{
			title:       "entry_with_the_level_higher_than_the_log_level_should_not_log",
			loggerLevel: InfoLevel,
			entryLevel:  DebugLevel,
			message:     "log me",
		},
	}
	assrt := assert.New(t)

	for _, tc := range testCases {
		shouldLog := tc.entryLevel <= tc.loggerLevel
		t.Run(tc.title, func(t *testing.T) {
			logEntryAndAssertJSON(tc.loggerLevel,
				func(entry *LogEntry) {
					entry.AsLevel(tc.entryLevel).Write(tc.message)
				},
				func(fields Fields, entry *LogEntry) {
					assrt.Equal(tc.loggerLevel, entry.Level)
					assrt.Equal(tc.loggerLevel, entry.Logger.Level)
					msg, ok := fields["msg"]
					if shouldLog {
						if !ok {
							t.Error("Failed to retrieve the message. Nothing was logged")
						}
						if logged, ok := checkLoggedField(tc.message, msg); !ok {
							t.Errorf("expected %s, received '%v'", tc.message, logged)
						}
						return
					}
					if ok {
						t.Errorf("we shouldn't have logged anything but the output was %v", fields)
					}
				})
		})
	}
}

func TestWithField(t *testing.T) {
	testCases := []struct {
		title       string
		loggerLevel Level
		entryLevel  Level
		message     string
	}{
		{
			title:       "entry_with_the_same_level_as_log_level_should_log",
			loggerLevel: DebugLevel,
			entryLevel:  DebugLevel,
			message:     "log me",
		},
		{
			title:       "entry_with_the_level_lower_than_the_log_level_should_log",
			loggerLevel: DebugLevel,
			entryLevel:  InfoLevel,
			message:     "log me",
		},
		{
			title:       "entry_with_the_level_higher_than_the_log_level_should_not_log",
			loggerLevel: InfoLevel,
			entryLevel:  DebugLevel,
			message:     "log me",
		},
	}
	assrt := assert.New(t)
	const fieldKey = "field"
	for _, tc := range testCases {
		shouldLog := tc.entryLevel <= tc.loggerLevel
		t.Run(tc.title, func(t *testing.T) {
			logEntryAndAssertJSON(tc.loggerLevel,
				func(entry *LogEntry) {
					entry.AsLevel(tc.entryLevel).WithField(fieldKey, tc.title).Write(tc.message)
				},
				func(fields Fields, entry *LogEntry) {
					assrt.Equal(tc.loggerLevel, entry.Level)
					assrt.Equal(tc.loggerLevel, entry.Logger.Level)
					msg, msgOk := fields["msg"]
					field, fieldOk := fields[fieldKey]
					if shouldLog {
						if !msgOk {
							t.Error("Failed to retrieve the message. Nothing was logged")
						}
						if logged, ok := checkLoggedField(tc.message, msg); !ok {
							t.Errorf("expected %s, received '%v'", tc.message, logged)
						}

						if logged, ok := checkLoggedField(tc.title, field); !ok {
							t.Errorf("expected '%v' to be %s, but it was '%v'", fieldKey, tc.title, logged)
						}
						return
					}
					if msgOk || fieldOk {
						t.Errorf("we shouldn't have logged anything but the output was %v", fields)
					}
				})
		})
	}
}

func TestWithFields(t *testing.T) {
	testCases := []struct {
		title          string
		message        string
		originalFields Fields
		logFields      Fields
		entryLevel     Level
		loggerLevel    Level
		shouldLog      bool
	}{
		{
			title:          "original_fields_should_not_change_after_logging",
			message:        "message",
			originalFields: Fields{"original_key": "original_value"},
			logFields:      Fields{"log_field_key": "log_field_value"},
			entryLevel:     DebugLevel,
			loggerLevel:    DebugLevel,
			shouldLog:      true,
		},
		{
			title:          "original_fields_should_get_logged_even_if_log_fields_are_empty",
			message:        "message",
			originalFields: Fields{"original_key": "original_value"},
			logFields:      Fields{},
			entryLevel:     DebugLevel,
			loggerLevel:    DebugLevel,
			shouldLog:      true,
		},
		{
			title:          "no_fields_should_get_logged_if_entry_level_is_higher_than_logger_level",
			message:        "message",
			originalFields: Fields{"original_key": "original_value"},
			logFields:      Fields{"log_field_key": "log_field_value"},
			entryLevel:     DebugLevel,
			loggerLevel:    InfoLevel,
			shouldLog:      false,
		},
	}
	assrt := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var buffer bytes.Buffer
			var fields Fields

			logger := NewLogger(tc.loggerLevel)
			logger.Out = &buffer
			logger.Formatter = new(JSONFormatter)
			entry := logger.WithFields(tc.originalFields)
			entry.AsLevel(tc.entryLevel).WithFields(tc.logFields).Write(tc.message)

			json.Unmarshal(buffer.Bytes(), &fields)

			assrt.Equal(entry.Data, tc.originalFields)
			assrt.Equal(tc.loggerLevel, entry.Level)

			if tc.shouldLog {
				msg, ok := fields["msg"]
				if !ok {
					t.Error("Failed to retrieve the message. Nothing was logged")
				}

				if logged, ok := checkLoggedField(tc.message, msg); !ok {
					t.Errorf("expected %s, received '%v'", tc.message, logged)
				}

				assertFields(t, tc.logFields, fields)
				assertFields(t, tc.originalFields, fields)

				return
			}
			if len(fields) > 0 {
				t.Errorf("we shouldn't have logged anything but the output was %v", fields)
			}
		})
	}
}

func TestWithError(t *testing.T) {
	testCases := []struct {
		title          string
		message        string
		originalFields Fields
		err            error
		entryLevel     Level
		loggerLevel    Level
		shouldLog      bool
	}{
		{
			title:          "original_fields_and_error_should_get_logged",
			message:        "message",
			originalFields: Fields{"original_key": "original_value"},
			err:            errors.New("test error"),
			entryLevel:     DebugLevel,
			loggerLevel:    DebugLevel,
			shouldLog:      true,
		},
		{
			title:          "no_fields_should_get_logged_if_entry_level_is_higher_than_logger_level",
			message:        "message",
			originalFields: Fields{"original_key": "original_value"},
			err:            errors.New("test error"),
			entryLevel:     DebugLevel,
			loggerLevel:    InfoLevel,
			shouldLog:      false,
		},
	}
	assrt := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var buffer bytes.Buffer
			var fields Fields

			logger := NewLogger(tc.loggerLevel)
			logger.Out = &buffer
			logger.Formatter = new(JSONFormatter)
			entry := logger.WithFields(tc.originalFields)
			cloned := entry.AsLevel(tc.entryLevel).WithError(tc.err)
			cloned.Write(tc.message)

			json.Unmarshal(buffer.Bytes(), &fields)

			assrt.Equal(entry.Data, tc.originalFields)
			assrt.Equal(tc.loggerLevel, entry.Level)

			if tc.shouldLog {
				msg, ok := fields["msg"]
				if !ok {
					t.Error("Failed to retrieve the message. Nothing was logged")
				}

				if logged, ok := checkLoggedField(tc.message, msg); !ok {
					t.Errorf("expected %s, received '%v'", tc.message, logged)
				}

				assertFields(t, Fields{ErrorKey: tc.err.Error()}, fields)
				assertFields(t, tc.originalFields, fields)

				return
			}
			assertFields(t, tc.originalFields, cloned.Data)
			if len(fields) > 0 {
				t.Errorf("we shouldn't have logged anything but the output was %v", fields)
			}
		})
	}
}

func TestEntryInstantiation(t *testing.T) {
	logger := NewLogger(InfoLevel)
	buf := &bytes.Buffer{}
	logger.Out = buf

	testCases := []struct {
		title    string
		entry    *LogEntry
		writable bool
		message  string
	}{
		{
			title:    "invalid_log_entry_should_not_log",
			entry:    &LogEntry{},
			writable: false,
		},
		{
			title: "an_entry_with_manually_set_level_should_not_log",
			entry: &LogEntry{
				Level: DebugLevel,
			},
			writable: false,
		},
		{
			title:    "valid_entry_should_log",
			entry:    NewLogEntry(logger),
			writable: true,
			message:  "some cool stuff",
		},
		{
			title:    "valid_entry_created_by_logger_should_log",
			entry:    NewLogEntry(logger),
			writable: true,
			message:  "some cool stuff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			buf.Reset()
			tc.entry.Write(tc.message)
			if tc.writable {
				assertLoggedMessage(t, tc.message, buf)
				return
			}

			assertEmptyOutput(t, buf)
		})
	}
}

func assertFields(t *testing.T, expected Fields, actual Fields) {
	t.Helper()
	for fieldKey, expectedFieldValue := range expected {
		field, ok := actual[fieldKey]
		if !ok {
			t.Errorf("Filed '%v' has not been logged", fieldKey)
			continue
		}
		if field != expectedFieldValue {
			t.Errorf("expected [%v] to be %s, but it was '%v'", fieldKey, expectedFieldValue, field)
		}
	}
}

func assertEmptyOutput(t *testing.T, output *bytes.Buffer) {
	t.Helper()
	if output.Len() > 0 {
		fields := inspectJsonOutput(t, output)
		t.Errorf("we shouldn't have logged anything but the output was %v", fields)
	}
}

func assertLoggedMessage(t *testing.T, expected string, output *bytes.Buffer) {
	t.Helper()
	fields := inspectJsonOutput(t, output)
	msg, ok := fields["msg"]
	if !ok {
		t.Error("Failed to retrieve the message. Nothing was logged")
	}
	if logged, ok := checkLoggedField(expected, msg); !ok {
		t.Errorf("expected %s, received '%v'", expected, logged)
	}
}

func checkLoggedField(expected string, actual interface{}) (string, bool) {
	logged, ok := actual.(string)
	return logged, ok && logged == expected
}

func inspectJsonOutput(t *testing.T, buffer *bytes.Buffer) Fields {
	t.Helper()
	var fields Fields
	err := json.Unmarshal(buffer.Bytes(), &fields)
	if err != nil {
		t.Errorf("Failed to unmarshal the log output %s", err)
		return nil
	}
	return fields
}
