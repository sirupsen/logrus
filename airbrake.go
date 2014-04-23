package logrus

import (
	"github.com/tobi/airbrake-go"
)

// AirbrakeHook to send exceptions to an exception-tracking service compatible
// with the Airbrake API. You must set:
// * airbrake.Endpoint
// * airbrake.ApiKey
// * airbrake.Environment (only sends exceptions when set to "production")
//
// Before using this hook, to send exceptions. Entries that trigger an Error,
// Fatal or Panic should now include an "Error" field to send to Airbrake.
type AirbrakeHook struct{}

func (hook *AirbrakeHook) Fire(entry *Entry) error {
	if entry.Data["error"] == nil {
		entry.Logger.WithFields(Fields{
			"source":   "airbrake",
			"endpoint": airbrake.Endpoint,
		}).Warn("Exceptions sent to Airbrake must have an 'error' key with the error")
		return nil
	}

	err := airbrake.Notify(entry.Data["error"].(error))
	if err != nil {
		entry.Logger.WithFields(Fields{
			"source":   "airbrake",
			"endpoint": airbrake.Endpoint,
		}).Warn("Failed to send error to Airbrake")
	}

	return nil
}

func (hook *AirbrakeHook) Levels() []Level {
	return []Level{
		Error,
		Fatal,
		Panic,
	}
}
