package logrus_sourcefile

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"runtime"
	"strings"
)

type SourceFileHook struct {
	LogLevel logrus.Level
}

func (hook *SourceFileHook) Fire(entry *logrus.Entry) (_ error) {
	for skip := 4; skip < 9; skip++ {
		_, file, line, _ := runtime.Caller(skip)
		split := strings.Split(file, "/")
		if l := len(split); l > 1 {
			pkg := split[l-2]
			if pkg != "logrus" {
				file = fmt.Sprintf("%s/%s:%d", split[l-2], split[l-1], line)
				// set source_file field
				entry.Data["source_file"] = file
				return
			}
		}
	}

	return
}

func (hook *SourceFileHook) Levels() []logrus.Level {
	levels := make([]logrus.Level, hook.LogLevel+1)
	for i, _ := range levels {
		levels[i] = logrus.Level(i)
	}
	return levels
}
