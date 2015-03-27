package airbrake

import (
	"errors"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/airbrake/gobrake"
)

// Set graylog.BufSize = <value> _before_ calling NewHook
var BufSize uint = 1024

// AirbrakeHook to send exceptions to an exception-tracking service compatible
// with the Airbrake API.
type airbrakeHook struct {
	Airbrake   *gobrake.Notifier
	noticeChan chan *gobrake.Notice
}

func NewHook(projectID int64, apiKey, env string) *airbrakeHook {
	airbrake := gobrake.NewNotifier(projectID, apiKey)
	airbrake.SetContext("environment", env)
	hook := &airbrakeHook{
		Airbrake:   airbrake,
		noticeChan: make(chan *gobrake.Notice, BufSize),
	}
	go hook.fire()
	return hook
}

func (hook *airbrakeHook) Fire(entry *logrus.Entry) error {
	var notifyErr error
	err, ok := entry.Data["error"].(error)
	if ok {
		notifyErr = err
	} else {
		notifyErr = errors.New(entry.Message)
	}
	notice := hook.Airbrake.Notice(notifyErr, nil, 3)
	hook.noticeChan <- notice
	return nil
}

// fire sends errors to airbrake when an entry is available on entryChan
func (hook *airbrakeHook) fire() {
	for {
		notice := <-hook.noticeChan

		if err := hook.Airbrake.SendNotice(notice); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to send error to Airbrake: %v\n", err)
		}

	}
}

func (hook *airbrakeHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
