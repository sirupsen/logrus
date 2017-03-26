package logrus_logentries

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/bsphere/le_go"
)

// LogentriesHook to send logs to a logging service compatible with the Logentries API.
type LogentriesHook struct {
	Token   string
	AppName string
}

// NewLogentriesHook creates a hook to be added to an instance of logger.
func NewLogentriesHook(token string, appName string) *LogentriesHook {
	return &LogentriesHook{token, appName}
}

// Fire is called when a log event is fired.
func (hook *LogentriesHook) Fire(entry *logrus.Entry) error {
	payload := fmt.Sprintf("%s: %s", hook.AppName, entry.Message)

	le, err := le_go.Connect(hook.Token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to logentries")
		return err
	}
	defer le.Close()

	le.Println(payload)

	return nil
}

// Levels returns the available logging levels.
func (hook *LogentriesHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
