package logrus

import (
	"encoding/json"
	"fmt"
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

	timeKey := f.TimeKey
	if timeKey == "" {
		timeKey = "time"
	}

	messageKey := f.MessageKey
	if messageKey == "" {
		messageKey = "msg"
	}

	levelKey := f.LevelKey
	if levelKey == "" {
		levelKey = "level"
	}

	data[timeKey] = entry.Time.Format(timestampFormat)
	data[messageKey] = entry.Message
	data[levelKey] = entry.Level.String()

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
