package airbrake

import (
	"errors"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/tobi/airbrake-go"
)

// Set graylog.BufSize = <value> _before_ calling NewHook
var BufSize uint = 1024

// AirbrakeHook to send exceptions to an exception-tracking service compatible
// with the Airbrake API.
type airbrakeHook struct {
	APIKey      string
	Endpoint    string
	Environment string
	entryChan   chan *logrus.Entry
}

func NewHook(endpoint, apiKey, env string) *airbrakeHook {
	hook := &airbrakeHook{
		APIKey:      apiKey,
		Endpoint:    endpoint,
		Environment: env,
		entryChan:   make(chan *logrus.Entry, BufSize),
	}
	go hook.fire()
	return hook
}

func (hook *airbrakeHook) Fire(entry *logrus.Entry) error {
	hook.entryChan <- entry
	return nil
}

// fire sends errors to airbrake when an entry is available on entryChan
func (hook *airbrakeHook) fire() {
	for {
		entry := <-hook.entryChan
		airbrake.ApiKey = hook.APIKey
		airbrake.Endpoint = hook.Endpoint
		airbrake.Environment = hook.Environment

		var notifyErr error
		err, ok := entry.Data["error"].(error)
		if ok {
			notifyErr = err
		} else {
			notifyErr = errors.New(entry.Message)
		}

		if err = airbrake.Notify(notifyErr); err != nil {
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
