package logrus

import (
	"bytes"
	"errors"

	"testing"
)

func TestQuoting(t *testing.T) {
	tf := new(TextFormatter)

	checkQuoting := func(q bool, value interface{}) {
		b, _ := tf.Format(WithField("test", value))
		idx := bytes.LastIndex(b, []byte{'='})
		cont := bytes.Contains(b[idx:], []byte{'"'})
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
