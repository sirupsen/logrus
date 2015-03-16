package logrus_bugsnag

import (
	"github.com/Sirupsen/logrus"
	"github.com/bugsnag/bugsnag-go"
)

// BugsnagHook sends exceptions to an exception-tracking service compatible
// with the Bugsnag API. Before using this hook, you must call
// bugsnag.Configure().
//
// Entries that trigger an Error, Fatal or Panic should now include an "error"
// field to send to Bugsnag
type BugsnagHook struct{}

// Fire forwards an error to Bugsnag. Given a logrus.Entry, it extracts the
// implicitly-required "error" field and sends it off.
func (hook *BugsnagHook) Fire(entry *logrus.Entry) error {
	if entry.Data["error"] == nil {
		entry.Logger.WithFields(logrus.Fields{
			"source": "bugsnag",
		}).Warn("Exceptions sent to Bugsnag must have an 'error' key with the error")
		return nil
	}

	err, ok := entry.Data["error"].(error)
	if !ok {
		entry.Logger.WithFields(logrus.Fields{
			"source": "bugsnag",
		}).Warn("Exceptions sent to Bugsnag must have an `error` key of type `error`")
		return nil
	}

	bugsnagErr := bugsnag.Notify(err)
	if bugsnagErr != nil {
		entry.Logger.WithFields(logrus.Fields{
			"source": "bugsnag",
			"error":  bugsnagErr,
		}).Warn("Failed to send error to Bugsnag")
	}

	return nil
}

// Levels enumerates the log levels on which the error should be forwarded to
// bugsnag: everything at or above the "Error" level.
func (hook *BugsnagHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
