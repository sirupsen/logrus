package logrus

import (
	"regexp"
	"testing"
)

var str = "[INFO] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."

var strs = []string{
	"[INFO] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[INFO] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[INFO] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[INFO] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[INFO] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[WARN] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[ERROR] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[FATAL] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[INFO] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[DEBUG] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[DEBUG] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[DEBUG] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[DEBUG] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[DEBUG] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[DEBUG] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[INFO] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[INFO] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"[INFO] Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
}

func TestRawRegexp(t *testing.T) {

	r := regexp.MustCompile(`^\[\w+\]`)
	b := r.Find([]byte(str))
	l := string(b)[1 : len(b)-1]

	if l != "INFO" {
		t.Errorf("INFO did not return InfoLevel: %#v", l)
	}

}

func TestRegexpParser(t *testing.T) {

	pr := RegexpParser{}
	pr.prefixRegex()
	l, err := pr.Parse(&str)

	if err != nil {
		t.Errorf("RegexpParse error: %v", err)
	}
	if l != InfoLevel {
		t.Errorf("INFO did not return InfoLevel: %#v", l)
	}

}

func BenchmarkFilterMatchPrefix(b *testing.B) {
	pr := RegexpParser{}
	pr.prefixRegex()
	for n := 0; n < b.N; n++ {
		_, _ = pr.Parse(&str)
	}
}

func BenchmarkFilterMatchPrefixRange(b *testing.B) {

	pr := RegexpParser{}
	pr.prefixRegex()
	for n := 0; n < b.N; n++ {
		_, _ = pr.Parse(&strs[n%len(strs)])
	}
}

func TestPrefixStrCmp(t *testing.T) {
	pp := PrefixStrCmp{}
	l, err := pp.Parse(&str)
	if err != nil {
		t.Errorf("PrefixStrCmp Parse error: %v", err)
	}
	if l != InfoLevel {
		t.Errorf("INFO did not return InfoLevel: %#v", l)
	}
}

func BenchmarkPrefixStrCmp(b *testing.B) {
	pp := PrefixStrCmp{}
	for n := 0; n < b.N; n++ {
		_, _ = pp.Parse(&str)
	}
}

func BenchmarkPrefixStrCmpRange(b *testing.B) {
	pp := PrefixStrCmp{}
	for n := 0; n < b.N; n++ {
		_, _ = pp.Parse(&strs[n%len(strs)])
	}
}
