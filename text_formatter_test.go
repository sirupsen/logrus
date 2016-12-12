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

func TestTimestampLocationOverride(t *testing.T) {
	moment := time.Unix(1472181897, 0).UTC()

	zone, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal("Unable to load Eastern Time: ", err)
	}
	local := moment.In(zone)
	body := WithField("dont", "care")
	body.Time = local

	customFormatter := &TextFormatter{Location: time.UTC}
	b, err := customFormatter.Format(body)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	if !strings.Contains(string(b), "2016-08-26T03:24:57Z") {
		t.Errorf("Expected a zulu timestamp: %s", b)
	}
}

// TODO add tests for sorting etc., this requires a parser for the text
// formatter output.
