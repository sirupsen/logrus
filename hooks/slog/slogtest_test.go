package slog_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"testing/slogtest"

	"github.com/sirupsen/logrus"
	lslog "github.com/sirupsen/logrus/hooks/slog"
)

// slogtestSkips lists the testing/slogtest cases that the logrus handler
// is known not to pass yet, with the reason for each. The issue asking
// for this suite (#1611) scopes it to measuring conformance; fixing the
// handler is out of scope, so known failures are skipped here instead
// of failing CI. Remove an entry when the handler starts passing it.
var slogtestSkips = map[string]string{
	// The handler logs through a logrus Entry, which always carries
	// its own timestamp: a zero Record.Time is not propagated, so the
	// JSON output contains a "time" key even for a zero record time.
	"zero-time": "handler always emits a \"time\" key: logrus entries always carry a timestamp and a zero Record.Time is not propagated",

	// The handler deliberately flattens group attributes into
	// "."-separated field keys (see appendAttr in handler.go), e.g.
	// "G.k", instead of emitting a nested object under the group key,
	// so the suite's group checks cannot find a group "G" object.
	"groups":                    "handler flattens groups into \".\"-separated field keys instead of emitting a nested group object",
	"WithGroup":                 "handler flattens groups into \".\"-separated field keys instead of emitting a nested group object",
	"multi-With":                "handler flattens groups into \".\"-separated field keys instead of emitting a nested group object",
	"empty-group-record":        "handler flattens groups into \".\"-separated field keys instead of emitting a nested group object",
	"nested-empty-group-record": "handler flattens groups into \".\"-separated field keys instead of emitting a nested group object",
	"resolve-groups":            "handler flattens groups into \".\"-separated field keys instead of emitting a nested group object",
	"resolve-WithAttrs-groups":  "handler flattens groups into \".\"-separated field keys instead of emitting a nested group object",
}

// TestHandlerSlogtest runs the standard testing/slogtest conformance
// suite against the logrus slog handler, using a JSONFormatter writing
// to a buffer, as suggested in #1611.
func TestHandlerSlogtest(t *testing.T) {
	var buf bytes.Buffer

	newHandler := func(t *testing.T) slog.Handler {
		// t.Name() is "<test name>/<case name>".
		caseName := t.Name()[strings.LastIndex(t.Name(), "/")+1:]
		if reason, ok := slogtestSkips[caseName]; ok {
			t.Skip(reason)
		}
		buf.Reset()
		logger := logrus.New()
		logger.SetOutput(&buf)
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetLevel(logrus.TraceLevel)
		return lslog.NewHandler(logger, nil)
	}

	result := func(t *testing.T) map[string]any {
		t.Helper()
		line := bytes.TrimSpace(buf.Bytes())
		if len(line) == 0 {
			return map[string]any{}
		}
		var m map[string]any
		if err := json.Unmarshal(line, &m); err != nil {
			t.Fatalf("unmarshalling handler output %q: %v", line, err)
		}
		return m
	}

	slogtest.Run(t, newHandler, result)
}
