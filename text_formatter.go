package logrus

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
)

func init() {
	baseTimestamp = time.Now()
}

func miniTS() int {
	return int(time.Since(baseTimestamp) / time.Second)
}

type TextFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool
}

func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	var serialized []byte

	if f.ForceColors || IsTerminal() {
		levelText := strings.ToUpper(entry.Data["level"].(string))[0:4]

		levelColor := blue

		if entry.Data["level"] == "warning" {
			levelColor = yellow
		} else if entry.Data["level"] == "error" ||
			entry.Data["level"] == "fatal" ||
			entry.Data["level"] == "panic" {
			levelColor = red
		}

		serialized = append(serialized, []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m[%04d] %-45s ", levelColor, levelText, miniTS(), entry.Data["msg"]))...)

		keys := make([]string, 0)
		for k, _ := range entry.Data {
			if k != "level" && k != "time" && k != "msg" {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		first := true
		for _, k := range keys {
			v := entry.Data[k]
			if first {
				first = false
			} else {
				serialized = append(serialized, ' ')
			}
			serialized = append(serialized, []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m=%v", levelColor, k, v))...)
		}
	} else {
		serialized = f.AppendKeyValue(serialized, "time", entry.Data["time"].(string))
		serialized = f.AppendKeyValue(serialized, "level", entry.Data["level"].(string))
		serialized = f.AppendKeyValue(serialized, "msg", entry.Data["msg"].(string))

		for key, value := range entry.Data {
			if key != "time" && key != "level" && key != "msg" {
				serialized = f.AppendKeyValue(serialized, key, value)
			}
		}
	}

	return append(serialized, '\n'), nil
}

func (f *TextFormatter) AppendKeyValue(serialized []byte, key, value interface{}) []byte {
	if _, ok := value.(string); ok {
		return append(serialized, []byte(fmt.Sprintf("%v=%q ", key, value))...)
	} else {
		return append(serialized, []byte(fmt.Sprintf("%v=%v ", key, value))...)
	}
}
