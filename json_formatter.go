package logrus

import (
	"encoding/json"
	"fmt"
	"time"
)

type JSONFormatter struct {
}

type jsonEntry struct {
	Time  time.Time `json:"time"`
	Level string    `json:"level"`
	Msg   string    `json:"msg"`
	Data  Fields    `json:"data,omitempty"`
}

func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	jsonEntry := jsonEntry{
		Time: entry.Time,
		Msg:  entry.Msg,
		Data: entry.Data,
	}

	switch entry.Level {
	case Debug:
		jsonEntry.Level = "debug"
	case Info:
		jsonEntry.Level = "info"
	case Warn:
		jsonEntry.Level = "warn"
	case Error:
		jsonEntry.Level = "error"
	case Fatal:
		jsonEntry.Level = "fatal"
	case Panic:
		jsonEntry.Level = "panic"
	}

	serialized, err := json.Marshal(jsonEntry)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal data to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}

func (f *JSONFormatter) Unformat(buffer []byte) (*Entry, error) {
	var jsonEntry jsonEntry

	err := json.Unmarshal(buffer, &jsonEntry)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal entry to JSON, %v", err)
	}

	entry := Entry{
		Time: jsonEntry.Time,
		Msg:  jsonEntry.Msg,
		Data: jsonEntry.Data,
	}

	switch jsonEntry.Level {
	case "debug":
		entry.Level = Debug
	case "info":
		entry.Level = Info
	case "warn":
		entry.Level = Warn
	case "error":
		entry.Level = Error
	case "fatal":
		entry.Level = Fatal
	case "panic":
		entry.Level = Panic
	}

	return &entry, nil
}
