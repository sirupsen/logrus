package logrus_sourcefile

import (
	"fmt"
	"github.com/flowhealth/logrus"
	"runtime"
	"strings"
)

type SourceFileHook struct {
	LogLevel logrus.Level
}

func (hook *SourceFileHook) Fire(entry *logrus.Entry) error {
	skip := 4
	if len(entry.Data) == 0 {
		skip = 6
	}

	// set source_file field
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		split := strings.Split(file, "/")
		if l := len(split); l > 2 {
			file = fmt.Sprintf("%s/%s:%d", split[l-2], split[l-1], line)
		}
		entry.Data["source_file"] = file
	}

	return nil
}

func (hook *SourceFileHook) Levels() []logrus.Level {
	levels := make([]logrus.Level, hook.LogLevel+1)
	for i, _ := range levels {
		levels[i] = logrus.Level(i)
	}
	return levels
}
