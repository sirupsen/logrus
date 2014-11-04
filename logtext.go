package logtext

// A Logrus hook that adds filename/line number/stack trace to our log outputs

import (
	"github.com/Sirupsen/logrus"
	"runtime"
)

// Log depth is how many levels to ascend to find where the actual log call occurred
// while debugOnly sets whether or not stack traces should be printed outside of
// debug prints
type Logtext struct{
	Formatter logrus.Formatter
	LogDepth int
	DebugOnly bool
}

func NewLogtext(formatter logrus.Formatter, debugOnly bool) *Logtext {
	return &Logtext{
		LogDepth: 4,
		Formatter: formatter,
		DebugOnly: debugOnly,
	}
}

// Creates a hook to be added to an instance of logger. This is called with
func (hook *Logtext) Format(entry *logrus.Entry) ([]byte, error) {

	if _, file, line, ok := runtime.Caller(hook.LogDepth); ok {
		entry.Data["line"] = line
		entry.Data["file"] = file
	}

	if !hook.DebugOnly || entry.Level == logrus.DebugLevel {
		stack := getTrace()
		entry.Data["stack"] = stack
	}

	return hook.Formatter.Format(entry)
}

// handles getting the stack trace and returns it as a string
func getTrace() string {
	stack := make([]byte, 2048)
	size := runtime.Stack(stack, false)
	return string(stack[:size])
}
