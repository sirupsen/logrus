package logrus_honeybadger

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/agonzalezro/goneybadger"
)

// HoneybadgerHook to send exceptions to an exception-tracking service compatible
// with the Honeybadger API.
type HoneybadgerHook struct {
	apiKey, env string
	timeout     time.Duration
}

func NewHoneyBadgerHook(apiKey, env string, timeout time.Duration) *HoneybadgerHook {
	return &HoneybadgerHook{
		apiKey:  apiKey,
		env:     env,
		timeout: timeout,
	}
}

func (h *HoneybadgerHook) Fire(entry *logrus.Entry) error {
	var notifyErr string

	if err, ok := entry.Data["error"].(error); ok {
		notifyErr = err.Error()
	} else {
		notifyErr = entry.Message
	}

	gb := goneybadger.NewWithTimeout(h.apiKey, h.env, h.timeout)
	err := gb.Notify(notifyErr)
	if err != nil {
		entry.Logger.WithFields(logrus.Fields{
			"source": "honeybadger",
			"error":  err,
		}).Warn("Failed to send error to HoneyBadger")
	}

	return nil
}

func (h *HoneybadgerHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
