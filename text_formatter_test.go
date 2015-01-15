package logrus

import (
	"bytes"
	"errors"

	"testing"
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
	checkQuoting(true, "/foobar")
	checkQuoting(true, "x y")
	checkQuoting(true, "x,y")
	checkQuoting(false, errors.New("invalid"))
	checkQuoting(true, errors.New("invalid argument"))
}

func TestTextPrint(t *testing.T) {
	tf := &TextFormatter{DisableColors: true}
	byts, _ := tf.Format(&Entry{Message: "msg content"})

	// make sure no leading or trailing spaces
	if string(byts) !=
		"time=\"0001-01-01T00:00:00Z\" level=panic msg=\"msg content\"\n" {
		t.Errorf("not expected: %q", string(byts))
	}
}

func TestColorPrint(t *testing.T) {
	tf := &TextFormatter{ForceColors: true}
	entry := WithField("testkey", "value")
	entry.Message = "msg content"
	byts, _ := tf.Format(entry)

	// make sure no leading or trailing spaces
	if string(byts) !=
		"\x1b[31mPANI\x1b[0m[0000] " +
			// length 44 plus one space
			"msg content                                  " +
			"\x1b[31mtestkey\x1b[0m=value\n" {
		t.Errorf("not expected: %q", string(byts))
	}
}
