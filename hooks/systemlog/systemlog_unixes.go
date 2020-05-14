// +build !windows,!nacl,!plan9

package systemlog

import (
	"fmt"
	"log/syslog"
	"os"

	"github.com/sirupsen/logrus"
)

// SystemlogHook to send logs via syslog.
type SystemlogHook struct {
	Writer        *syslog.Writer
}

// Creates a hook to be added to an instance of logger. This is called with
// `hook, err := NewSystemlogHook("udp", "localhost:514", "")`
// `if err == nil { log.Hooks.Add(hook) }`
func NewSystemlogHook(network, raddr string, tag string) (*SystemlogHook, error) {
	w, err := syslog.Dial(network, raddr, syslog.LOG_INFO, tag)
	return &SystemlogHook{w,}, err
}

func (hook *SystemlogHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case logrus.PanicLevel:
		return hook.Writer.Crit(line)
	case logrus.FatalLevel:
		return hook.Writer.Crit(line)
	case logrus.ErrorLevel:
		return hook.Writer.Err(line)
	case logrus.WarnLevel:
		return hook.Writer.Warning(line)
	case logrus.InfoLevel:
		return hook.Writer.Info(line)
	case logrus.DebugLevel, logrus.TraceLevel:
		return hook.Writer.Debug(line)
	default:
		return nil
	}
}

func (hook *SystemlogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
