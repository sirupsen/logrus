package logrus_caller

import (
	"github.com/Sirupsen/logrus"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type CallerHook struct {
}

func (hook *CallerHook) Fire(entry *logrus.Entry) error {
	entry.Data["caller"] = hook.caller()
	return nil
}

func (hook *CallerHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func (hook *CallerHook) caller() string {
	if _, file, line, ok := runtime.Caller(6); ok {
		return strings.Join([]string{filepath.Base(file), strconv.Itoa(line)}, ":")
	}
	// not sure what the convention should be here
	return ""
}
