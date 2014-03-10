package logrus

import (
	"fmt"
	"os"
	"sort"
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

	levelText := strings.ToUpper(entry.Data["level"].(string))[0:4]
	levelColor := blue

	if entry.Data["level"] == "warning" {
		levelColor = yellow
	} else if entry.Data["level"] == "fatal" ||
		entry.Data["level"] == "panic" {
		levelColor = red
	}

	if ttyutils.IsTerminal(os.Stdout.Fd()) {
		serialized = append(serialized, []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m[%04d] %-45s ", levelColor, levelText, miniTS(), entry.Data["msg"]))...)
	}

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

	return append(serialized, '\n'), nil
}
