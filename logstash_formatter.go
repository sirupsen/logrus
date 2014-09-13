package logrus

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// Formatter generates json in logstash format.
// Logstash site: http://logstash.net/
type LogstashFormatter struct {
	Type             string // if not empty use for logstash type field.
	FileLineLogLevel Level  // level on wich filename and linenumber will write to log. Be careful it's expensive.
}

func (f *LogstashFormatter) Format(entry *Entry) ([]byte, error) {
	skip := 6 // caller skip number, default for logger
	if len(entry.Data) == 5 {
		skip = 4 // for entry
	}

	entry.Data["@version"] = 1
	entry.Data["@timestamp"] = entry.Time.Format(time.RFC3339)

	// set message field
	_, ok := entry.Data["message"]
	if ok {
		entry.Data["fields.message"] = entry.Data["message"]
	}
	entry.Data["message"] = entry.Message

	// set level field
	_, ok = entry.Data["level"]
	if ok {
		entry.Data["fields.level"] = entry.Data["level"]
	}
	entry.Data["level"] = entry.Level.String()

	// set type field
	if f.Type != "" {
		_, ok = entry.Data["type"]
		if ok {
			entry.Data["fields.type"] = entry.Data["type"]
		}
		entry.Data["type"] = f.Type
	}

	// set file and line fields
	if f.FileLineLogLevel >= entry.Level {
		_, ok = entry.Data["source_file"]
		if ok {
			entry.Data["fields.source_file"] = entry.Data["source_file"]
		}
		_, file, line, ok := runtime.Caller(skip)
		if ok {
			split := strings.Split(file, "/")
			if l := len(split); l > 2 {
				file = fmt.Sprintf("%s/%s:%d", split[l-2], split[l-1], line)
			}
			entry.Data["source_file"] = file
		}
	}

	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
