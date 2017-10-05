package logrus

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"encoding/json"
)
func writeLogAndAssertJSON(loggerLevel Level, log func(*LogWriter), assertions func(Fields, *LogWriter)) {
	var buffer bytes.Buffer
	var fields Fields

	logger := NewLogger(loggerLevel)
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)

	log(logger)

	json.Unmarshal(buffer.Bytes(), &fields)

	assertions(fields, logger)
}

func writeLogAndAssertText(t *testing.T, loggerLevel Level, log func(*LogWriter), assertions func(Fields, *LogWriter)) {
	t.Helper()

	var buffer bytes.Buffer
	logger := NewLogger(loggerLevel)
	logger.Out = &buffer
	logger.Formatter = &TextFormatter{
		DisableColors: true,
	}

	log(logger)

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

	assertions(fields, logger)
}

func TestLogWriterDebugText(t *testing.T) {
	testCases := []struct {
		title       string
		loggerLevel Level
		message     string
		shouldLog   bool
	}{
		{
			title:       "logging_with_the_same_level_as_log_level_should_log",
			loggerLevel: DebugLevel,
			message:     "message",
			shouldLog:   true,
		},
		{
			title:       "logging_with_the_level_higher_than_the_log_level_should_not_log",
			loggerLevel: InfoLevel,
			message:     "message",
			shouldLog:   false,
		},
	}
	assrt := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			writeLogAndAssertText(t, tc.loggerLevel,
				func(lw *LogWriter) {
					lw.Debug(tc.message)
				},
				func(fields Fields, lw *LogWriter) {
					assrt.Equal(tc.loggerLevel, lw.Level)
					msg, ok := fields["msg"]
					if tc.shouldLog {
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

func TestLogWriterDebugJSON(t *testing.T) {
	testCases := []struct {
		title       string
		loggerLevel Level
		message     string
		shouldLog   bool
	}{
		{
			title:       "logging_with_the_same_level_as_log_level_should_log",
			loggerLevel: DebugLevel,
			message:     "message",
			shouldLog:   true,
		},
		{
			title:       "logging_with_the_level_higher_than_the_log_level_should_not_log",
			loggerLevel: InfoLevel,
			message:     "message",
			shouldLog:   false,
		},
	}
	assrt := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			writeLogAndAssertJSON(tc.loggerLevel,
				func(lw *LogWriter) {
					lw.Debug(tc.message)
				},
				func(fields Fields, lw *LogWriter) {
					assrt.Equal(tc.loggerLevel, lw.Level)
					msg, ok := fields["msg"]
					if tc.shouldLog {
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

func TestLogWriterInfoText(t *testing.T) {
	testCases := []struct {
		title       string
		loggerLevel Level
		message     string
		shouldLog   bool
	}{
		{
			title:       "logging_with_the_same_level_as_log_level_should_log",
			loggerLevel: DebugLevel,
			message:     "message",
			shouldLog:   true,
		},
		{
			title:       "logging_with_the_level_higher_than_the_log_level_should_not_log",
			loggerLevel: WarnLevel,
			message:     "message",
			shouldLog:   false,
		},
	}
	assrt := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			writeLogAndAssertText(t, tc.loggerLevel,
				func(lw *LogWriter) {
					lw.Info(tc.message)
				},
				func(fields Fields, lw *LogWriter) {
					assrt.Equal(tc.loggerLevel, lw.Level)
					msg, ok := fields["msg"]
					if tc.shouldLog {
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

func TestLogWriterInfoJSON(t *testing.T) {
	testCases := []struct {
		title       string
		loggerLevel Level
		message     string
		shouldLog   bool
	}{
		{
			title:       "logging_with_the_same_level_as_log_level_should_log",
			loggerLevel: DebugLevel,
			message:     "message",
			shouldLog:   true,
		},
		{
			title:       "logging_with_the_level_higher_than_the_log_level_should_not_log",
			loggerLevel: WarnLevel,
			message:     "message",
			shouldLog:   false,
		},
	}
	assrt := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			writeLogAndAssertJSON(tc.loggerLevel,
				func(lw *LogWriter) {
					lw.Info(tc.message)
				},
				func(fields Fields, lw *LogWriter) {
					assrt.Equal(tc.loggerLevel, lw.Level)
					msg, ok := fields["msg"]
					if tc.shouldLog {
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
