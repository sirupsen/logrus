package logrus

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/burke/ttyutils"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
)

type TextFormatter struct {
}

func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	var serialized []byte

	if ttyutils.IsTerminal(os.Stdout.Fd()) {
		levelText := strings.ToUpper(entry.Data["level"].(string))[0:4]

		levelColor := blue

		if entry.Data["level"] == "warning" {
			levelColor = yellow
		} else if entry.Data["level"] == "fatal" ||
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
				serialized = append(serialized, []byte(fmt.Sprintf("%s=%s ",
					f.ToString(key, false), f.ToString(value, true)))...)
			}
		}
	}

	return append(serialized, '\n'), nil
}

func (f *TextFormatter) AppendKeyValue(serialized []byte, key, value string) []byte {
	return append(serialized, []byte(fmt.Sprintf("%s=%s ",
		f.ToString(key, false), f.ToString(value, true)))...)
}

func (f *TextFormatter) ToString(value interface{}, escapeStrings bool) string {
	switch value.(type) {
	default:
		if escapeStrings {
			return fmt.Sprintf("'%s'", value)
		} else {
			return fmt.Sprintf("%s", value)
		}
	case int:
		return fmt.Sprintf("%s", strconv.Itoa(value.(int)))
	case uint64:
		return fmt.Sprintf("%s", strconv.FormatUint(value.(uint64), 10))
	case bool:
		return fmt.Sprintf("%s", strconv.FormatBool(value.(bool)))
	}
}
