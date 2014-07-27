package logrus

import (
	"encoding/json"
	"fmt"
)

type JSONFormatter struct {
	*internalFormatter
}

func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	f.prefixFieldClashes(entry)

	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
