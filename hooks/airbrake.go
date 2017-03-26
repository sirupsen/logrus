package hook

import (
	"github.com/Sirupsen/logrus"
	"github.com/tobi/airbrake-go"
)

type airbrakeHook struct{}

// AirbrakeHook to send errors to an exception-tracking service
// compatible with the Airbrake API. Errors will only be sent when env
// is set to "production".
//
//		// Sends errors
//		airhook := hook.AirbrakeHook(..., ..., "production")
//		// Doesn't send errors
//		airhook := hook.AirbrakeHook(..., ..., "dev")
//
// Entries that are invoked for an Error, Fatal or Panic should now include
// an "error" field containing a value of type `error`, which will be sent
// to Airbrake:
//
//		airhook := hook.AirbrakeHook(..., ..., "production")
//		log.Hooks.Add(airhook)
//		// only `err` will be sent to airbrake
//		log.WithField("error": err).Panic("what the hell have you built?!")
//
// The arguments will set global vars in the airbrake client, thus only
// one instance of this hook should be created.
func AirbrakeHook(endpoint, apiKey, env string) logrus.Hook {
	airbrake.Endpoint = endpoint
	airbrake.ApiKey = apiKey
	airbrake.Environment = env
	return &airbrakeHook{}
}

func (hook *airbrakeHook) Fire(entry *logrus.Entry) error {
	if entry.Data["error"] == nil {
		entry.Logger.WithFields(logrus.Fields{
			"source":   "airbrake",
			"endpoint": airbrake.Endpoint,
		}).Warn("Entries sent to Airbrake must have a field 'error' set containing the error")
		return nil
	}

	err := airbrake.Notify(entry.Data["error"].(error))
	if err != nil {
		entry.Logger.WithFields(logrus.Fields{
			"source":   "airbrake",
			"endpoint": airbrake.Endpoint,
		}).Warn("Failed to send error to Airbrake")
	}

	return nil
}

func (hook *airbrakeHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.Error,
		logrus.Fatal,
		logrus.Panic,
	}
}
