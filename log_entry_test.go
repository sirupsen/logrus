package logrus

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func logEntryAndAssertJSON(t *testing.T, loggerLevel Level, log func(*LogEntry), assertions func(Fields, *LogEntry)) {
	t.Helper()
	var buffer bytes.Buffer
	var fields Fields

	logger := NewLogger(loggerLevel)
	logger.SetLevel(loggerLevel)
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)

	entry := logger.Entry()

	log(entry)

	json.Unmarshal(buffer.Bytes(), &fields)

	assertions(fields, entry)
}

func TestEntryLogging(t *testing.T) {
	testCases := []struct {
		title                 string
		loggerLevel           Level
		entryLevel            Level
		expectedLevelAfterLog Level
		message               string
		shouldLog             bool
	}{
		{
			title:                 "entry_with_the_same_level_as_log_level_should_log",
			loggerLevel:           DebugLevel,
			entryLevel:            DebugLevel,
			expectedLevelAfterLog: DebugLevel,
			message:               "log me",
			shouldLog:             true,
		},
		{
			title:                 "entry_with_the_level_lower_than_the_log_level_should_log",
			loggerLevel:           DebugLevel,
			entryLevel:            InfoLevel,
			expectedLevelAfterLog: DebugLevel,
			message:               "log me",
			shouldLog:             true,
		},
		{
			title:                 "entry_with_the_level_higher_than_the_log_level_should_not_log",
			loggerLevel:           InfoLevel,
			entryLevel:            DebugLevel,
			expectedLevelAfterLog: InfoLevel,
			message:               "log me",
			shouldLog:             false,
		},
	}
	assert := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			logEntryAndAssertJSON(t, tc.loggerLevel,
				func(entry *LogEntry) {
					entry.AsLevel(tc.entryLevel).Write(tc.message)
				},
				func(fields Fields, entry *LogEntry) {
					assert.Equal(tc.expectedLevelAfterLog, entry.Level)
					assert.Equal(tc.expectedLevelAfterLog, entry.Logger.Level)
					msg, ok := fields["msg"]
					if tc.shouldLog {
						if !ok {
							t.Error("Failed to retrieve the message. Nothing was logged")
						}
						if logged, ok := checkLoggedMessage(tc.message, msg); !ok {
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
			entry:    logger.Entry(),
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
	if logged, ok := checkLoggedMessage(expected, msg); !ok {
		t.Errorf("expected %s, received '%v'", expected, logged)
	}
}

func checkLoggedMessage(expected string, actual interface{}) (string, bool) {
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
