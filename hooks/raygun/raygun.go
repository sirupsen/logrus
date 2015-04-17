package raygun

import (
	"errors"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/sditools/goraygun"
)

type raygunHook struct {
	client *goraygun.Client
}

func NewHook(Endpoint string, ApiKey string, Enabled bool) *raygunHook {
	client := goraygun.Init(goraygun.Settings{
		ApiKey:   ApiKey,
		Endpoint: Endpoint,
		Enabled:  Enabled,
	}, goraygun.Entry{})
	return &raygunHook{client}
}

func (hook *raygunHook) Fire(logEntry *logrus.Entry) error {
	// Start with a copy of the default entry
	raygunEntry := hook.client.Entry

	if request, ok := logEntry.Data["request"]; ok {
		raygunEntry.Details.Request.Populate(*(request.(*http.Request)))
	}

	var reportErr error
	if err, ok := logEntry.Data["error"]; ok {
		reportErr = err.(error)
	} else {
		reportErr = errors.New(logEntry.Message)
	}

	hook.client.Report(reportErr, raygunEntry)

	return nil
}

func (hook *raygunHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
