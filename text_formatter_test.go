package logrus

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestQuoting(t *testing.T) {
	tf := &TextFormatter{DisableColors: true}

	checkQuoting := func(q bool, value interface{}) {
		b, _ := tf.Format(WithField("test", value))
		idx := bytes.Index(b, ([]byte)("test="))
		cont := bytes.Contains(b[idx+5:], []byte{'"'})
		if cont != q {
			if q {
				t.Errorf("quoting expected for: %#v", value)
			} else {
				t.Errorf("quoting not expected for: %#v", value)
			}
		}
	}

	checkQuoting(false, "abcd")
	checkQuoting(false, "v1.0")
	checkQuoting(false, "1234567890")
	checkQuoting(true, "/foobar")
	checkQuoting(true, "x y")
	checkQuoting(true, "x,y")
	checkQuoting(false, errors.New("invalid"))
	checkQuoting(true, errors.New("invalid argument"))
}

func TestTimestampFormat(t *testing.T) {
	checkTimeStr := func(format string) {
		customFormatter := &TextFormatter{DisableColors: true, TimestampFormat: format}
		customStr, _ := customFormatter.Format(WithField("test", "test"))
		timeStart := bytes.Index(customStr, ([]byte)("time="))
		timeEnd := bytes.Index(customStr, ([]byte)("level="))
		timeStr := customStr[timeStart+5 : timeEnd-1]
		if timeStr[0] == '"' && timeStr[len(timeStr)-1] == '"' {
			timeStr = timeStr[1 : len(timeStr)-1]
		}
		if format == "" {
			format = time.RFC3339
		}
		_, e := time.Parse(format, (string)(timeStr))
		if e != nil {
			t.Errorf("time string \"%s\" did not match provided time format \"%s\": %s", timeStr, format, e)
		}
	}

	checkTimeStr("2006-01-02T15:04:05.000000000Z07:00")
	checkTimeStr("Mon Jan _2 15:04:05 2006")
	checkTimeStr("")
}

func TestDisableLevelTruncation(t *testing.T) {
	entry := &Entry{
		Time:    time.Now(),
		Message: "testing",
	}
	keys := []string{}
	timestampFormat := "Mon Jan 2 15:04:05 -0700 MST 2006"
	checkDisableTruncation := func(disabled bool, level Level) {
		tf := &TextFormatter{DisableLevelTruncation: disabled}
		var b bytes.Buffer
		entry.Level = level
		tf.printColored(&b, entry, keys, timestampFormat)
		logLine := (&b).String()
		if disabled {
			expected := strings.ToUpper(level.String())
			if !strings.Contains(logLine, expected) {
				t.Errorf("level string expected to be %s when truncation disabled", expected)
			}
		} else {
			expected := strings.ToUpper(level.String())
			if len(level.String()) > 4 {
				if strings.Contains(logLine, expected) {
					t.Errorf("level string %s expected to be truncated to %s when truncation is enabled", expected, expected[0:4])
				}
			} else {
				if !strings.Contains(logLine, expected) {
					t.Errorf("level string expected to be %s when truncation is enabled and level string is below truncation threshold", expected)
				}
			}
		}
	}

	checkDisableTruncation(true, DebugLevel)
	checkDisableTruncation(true, InfoLevel)
	checkDisableTruncation(false, ErrorLevel)
	checkDisableTruncation(false, InfoLevel)
}

func TestDisableTimestampWithColoredOutput(t *testing.T) {
	tf := &TextFormatter{DisableTimestamp: true, ForceColors: true}

	b, _ := tf.Format(WithField("test", "test"))
	if strings.Contains(string(b), "[0000]") {
		t.Error("timestamp not expected when DisableTimestamp is true")
	}
}

// TODO add tests for sorting etc., this requires a parser for the text
// formatter output.
