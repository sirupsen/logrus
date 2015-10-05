package logrus

import (
	"bytes"
	"errors"
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

func TestGolfsWithTextFormatter(t *testing.T) {
	p := &Person{
		Name:  "Bruce",
		Alias: "Batman",
		Hideout: &Hideout{
			Name:        "JLU Tower",
			DimensionId: 52,
		},
	}

	tf := &TextFormatter{DisableColors: true}
	b, _ := tf.Format(&Entry{
		Message: "the dark knight", Data: Fields{"hero": p}})

	if bytes.Index(b, ([]byte)("hero.name=Bruce")) < 0 {
		t.Fatalf("missing hero.name=Bruce")
	}

	if bytes.Index(b, ([]byte)("hero.alias=Batman")) < 0 {
		t.Fatalf("missing hero.alias=Batman")
	}

	if bytes.Index(b, ([]byte)(`hero.hideout.name="JLU Tower"`)) < 0 {
		t.Fatalf(`missing hero.hideout.name="JLU Tower"`)
	}

	if bytes.Index(b, ([]byte)("hero.hideout.dimensionId=52")) < 0 {
		t.Fatalf("missing hero.hideout.dimensionId=52")
	}
}

func TestGolfsWithTextFormatterAndNonGolfer(t *testing.T) {
	h := &Hideout{
		Name:        "JLU Tower",
		DimensionId: 52,
	}

	tf := &TextFormatter{DisableColors: true}
	b, _ := tf.Format(&Entry{
		Message: "secret base", Data: Fields{"hideout": h}})
	t.Log(string(b))

	if bytes.Index(b, ([]byte)("hideout=&{JLU Tower 52}")) < 0 {
		t.Fatalf("missing hideout={JLU Tower 52}")
	}
}

// TODO add tests for sorting etc., this requires a parser for the text
// formatter output.
