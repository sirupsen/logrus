package logrus

import (
	"bytes"
	"fmt"
	"github.com/kortschak/zalgo"
)

type ZalgoFormatter struct {
	victim *TextFormatter
	pain   *bytes.Buffer
	z      *zalgo.Corrupter
}

func NewZalgoFormatterrrrrr() *ZalgoFormatter {
	pain := bytes.NewBuffer(nil)
	z := zalgo.NewCorrupter(pain)

	z.Zalgo = func(n int, z *zalgo.Corrupter) {
		z.Up += 0.1
		z.Middle += complex(0.01, 0.01)
		z.Down += complex(real(z.Down)*0.1, 0)
	}

	return &ZalgoFormatter{
		victim: &TextFormatter{},
		pain:   pain,
		z:      z,
	}
}

func (zal *ZalgoFormatter) Format(entry *Entry) ([]byte, error) {
	zal.pain.Reset()

	zal.z.Up = complex(0, 0.2)
	zal.z.Middle = complex(0, 0.2)
	zal.z.Down = complex(0.001, 0.3)

	victimsWish := entry.Data["msg"].(string)

	_, _ = fmt.Fprint(zal.z, victimsWish)

	victimsReality := zal.pain.String()

	entry.Data["msg"] = victimsReality
	return zal.victim.Format(entry)
}
