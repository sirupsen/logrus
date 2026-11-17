//go:build !tinygo

// TinyGo currently (v0.41.1) doesn't support `runtime.Caller`;
// https://tinygo.org/docs/reference/lang-support/stdlib/#logslog
package slog_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"path"
	"testing"
	"testing/slogtest"

	"github.com/sirupsen/logrus"
	lslog "github.com/sirupsen/logrus/hooks/slog"
	"github.com/stretchr/testify/require"
)

func TestHandler_slogtest(t *testing.T) {
	// slogtest cases the Handler does not pass, with the failure they report;
	// they are skipped.
	exceptions := map[string]string{
		// logrus replaces a zero Record.Time with time.Now().
		"zero-time": `unexpected key "time": a Handler should ignore a zero Record.Time`,

		// The handler flattens groups into "."-joined keys ("G.c"), while slogtest
		// expects each group as a nested map.
		"groups":                    `missing group "G": a Handler should handle Group attributes`,
		"WithGroup":                 `missing group "G": a Handler should handle the WithGroup method`,
		"multi-With":                `missing group "G": a Handler should handle multiple WithGroup and WithAttr calls`,
		"empty-group-record":        `missing group "G": a Handler should not output groups if there are no attributes`,
		"nested-empty-group-record": `missing group "G": a Handler should not output nested groups if there are no attributes`, // Go 1.25+
		"resolve-groups":            `missing group "G": a Handler should call Resolve on attribute values in groups`,
		"resolve-WithAttrs-groups":  `missing group "G": a Handler should call Resolve on attribute values in groups from WithAttrs`,
	}

	var buf bytes.Buffer

	newHandler := func(t *testing.T) slog.Handler {
		if reason, ok := exceptions[path.Base(t.Name())]; ok {
			t.Skip("known failure: " + reason)
		}
		buf.Reset()
		l := logrus.New()
		l.SetOutput(&buf)
		l.SetFormatter(&logrus.JSONFormatter{})
		return lslog.NewHandler(l, nil)
	}

	result := func(t *testing.T) map[string]any {
		var m map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &m), "output: %q", buf.Bytes())
		return m
	}

	slogtest.Run(t, newHandler, result)
}
