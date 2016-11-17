package logrus

import (
	"encoding/json"
	"fmt"
)

type fieldKey string

const (
	DefaultKeyMsg   = "msg"
	DefaultKeyLevel = "level"
	DefaultKeyTime  = "time"
)

type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
	MessageKey      string
	LevelKey        string
	TimeKey         string
}

func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	data := make(Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	prefixFieldClashes(data)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = DefaultTimestampFormat
	}

	data[f.resolveKey(f.TimeKey, DefaultKeyTime)] = entry.Time.Format(timestampFormat)
	data[f.resolveKey(f.MessageKey, DefaultKeyMsg)] = entry.Message
	data[f.resolveKey(f.LevelKey, DefaultKeyLevel)] = entry.Level.String()

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}

func (f *JSONFormatter) resolveKey(key, defaultKey string) string {
	if len(key) > 0 {
		return key
	}
	return defaultKey
}
