package logrus_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

type benchStringer string

func (s benchStringer) String() string { return string(s) }

var numericFields = logrus.Fields{
	"i":   int(42),
	"i64": int64(-1234567890),
	"u":   uint(99),
	"u64": uint64(18446744073709551615),
	"f32": float32(3.1415927),
	"f64": float64(-1.2345e6),
}

var boolFields = logrus.Fields{
	"t": true,
	"f": false,
	"x": true,
	"y": false,
}

var stringerFields = logrus.Fields{
	"s1": benchStringer("alpha"),
	"s2": benchStringer("beta"),
	"s3": benchStringer("gamma-delta"), // includes '-' (still unquoted)
	"s4": benchStringer("needs quote"), // includes space -> quoted path
}

// smallFields is a small size data set for benchmarking
var smallFields = logrus.Fields{
	"foo":   "bar",
	"baz":   "qux",
	"one":   "two",
	"three": "four",
}

// largeFields is a large size data set for benchmarking
var largeFields = logrus.Fields{
	"foo":       "bar",
	"baz":       "qux",
	"one":       "two",
	"three":     "four",
	"five":      "six",
	"seven":     "eight",
	"nine":      "ten",
	"eleven":    "twelve",
	"thirteen":  "fourteen",
	"fifteen":   "sixteen",
	"seventeen": "eighteen",
	"nineteen":  "twenty",
	"a":         "b",
	"c":         "d",
	"e":         "f",
	"g":         "h",
	"i":         "j",
	"k":         "l",
	"m":         "n",
	"o":         "p",
	"q":         "r",
	"s":         "t",
	"u":         "v",
	"w":         "x",
	"y":         "z",
	"this":      "will",
	"make":      "thirty",
	"entries":   "yeah",
}

var errorFields = logrus.Fields{
	"foo": fmt.Errorf("bar"),
	"baz": fmt.Errorf("qux"),
}

func BenchmarkNumericTextFormatter(b *testing.B) {
	doBenchmark(b, &logrus.TextFormatter{DisableColors: true}, numericFields)
}

func BenchmarkBoolTextFormatter(b *testing.B) {
	doBenchmark(b, &logrus.TextFormatter{DisableColors: true}, boolFields)
}

func BenchmarkStringerTextFormatter(b *testing.B) {
	doBenchmark(b, &logrus.TextFormatter{DisableColors: true}, stringerFields)
}

func BenchmarkErrorTextFormatter(b *testing.B) {
	doBenchmark(b, &logrus.TextFormatter{DisableColors: true}, errorFields)
}

func BenchmarkSmallTextFormatter(b *testing.B) {
	doBenchmark(b, &logrus.TextFormatter{DisableColors: true}, smallFields)
}

func BenchmarkLargeTextFormatter(b *testing.B) {
	doBenchmark(b, &logrus.TextFormatter{DisableColors: true}, largeFields)
}

func BenchmarkSmallColoredTextFormatter(b *testing.B) {
	doBenchmark(b, &logrus.TextFormatter{ForceColors: true}, smallFields)
}

func BenchmarkLargeColoredTextFormatter(b *testing.B) {
	doBenchmark(b, &logrus.TextFormatter{ForceColors: true}, largeFields)
}

func BenchmarkSmallJSONFormatter(b *testing.B) {
	doBenchmark(b, &logrus.JSONFormatter{}, smallFields)
}

func BenchmarkLargeJSONFormatter(b *testing.B) {
	doBenchmark(b, &logrus.JSONFormatter{}, largeFields)
}

var sink []byte

func doBenchmark(b *testing.B, formatter logrus.Formatter, fields logrus.Fields) {
	logger := logrus.New()

	entry := &logrus.Entry{
		Time:    time.Time{},
		Level:   logrus.InfoLevel,
		Message: "message",
		Data:    fields,
		Logger:  logger,
	}

	// Warm once to determine output size and validate.
	d, err := formatter.Format(entry)
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(d)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d, err = formatter.Format(entry)
		if err != nil {
			b.Fatal(err)
		}
	}
	sink = d
}
