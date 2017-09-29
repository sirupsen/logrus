package logrus

import (
	"testing"

	"bytes"
	"encoding/json"
)

func TestEntryInstantiation(t *testing.T) {
	logger := NewLogger()
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
			entry:    logger.AsDebug(),
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
	if logged, ok := msg.(string); !ok || logged != expected {
		t.Errorf("expected %s, received '%v'", expected, logged)
	}
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
