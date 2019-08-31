package logrus

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestFormattingWithFieldFilters(t *testing.T) {
	textFormatter := &TextFormatter{FieldFilters: []FieldFilter{IgnoreEmptyValueFieldFilter}}
	jsonFormatter := &JSONFormatter{FieldFilters: []FieldFilter{IgnoreEmptyValueFieldFilter}}

	entry := &Entry {
		Message: "Testing field filtering",
		Time: time.Now(),
		Level: InfoLevel,
		Data: Fields{
			"wantedFieldKey":                  "value",
			"unwantedFieldKeyWithEmptyStringValue": "",
			"unwantedFieldKeyWithNilValue":    nil,
		},
	}

	testFilteringOfFields := func(formatters ...Formatter) {
		for _, formatter := range formatters {
			b, err := formatter.Format(entry)
			require.NoError(t, err)
			require.True(t, strings.Contains(string(b), "wantedFieldKey=value"))
			require.False(t, strings.Contains(string(b), "unwantedFieldKeyWithEmptyStringValue"))
			require.False(t, strings.Contains(string(b), "unwantedFieldKeyWithNilValue"))
		}
	}

	testFilteringOfFields(textFormatter, jsonFormatter)
}
