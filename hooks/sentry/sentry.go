package logrus_sentry

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
)

const (
	timeout = 100 * time.Millisecond
)

var (
	severityMap = map[logrus.Level]raven.Severity{
		logrus.DebugLevel: raven.DEBUG,
		logrus.InfoLevel:  raven.INFO,
		logrus.WarnLevel:  raven.WARNING,
		logrus.ErrorLevel: raven.ERROR,
		logrus.FatalLevel: raven.FATAL,
		logrus.PanicLevel: raven.FATAL,
	}
)

func getAndDel(d logrus.Fields, key string) (string, bool) {
	var (
		ok  bool
		v   interface{}
		val string
	)
	if v, ok = d[key]; !ok {
		return "", false
	}

	if val, ok = v.(string); !ok {
		return "", false
	}
	delete(d, key)
	return val, true
}

// SentryHook delivers logs to a sentry server.
type SentryHook struct {
	client *raven.Client
	levels []logrus.Level
}

// NewSentryHook creates a hook to be added to an instance of logger
// and initializes the raven client.
func NewSentryHook(DSN string, levels []logrus.Level) (*SentryHook, error) {
	client, err := raven.NewClient(DSN, nil)
	if err != nil {
		return nil, err
	}
	return &SentryHook{client, levels}, nil
}

// Called when an event should be sent to sentry
// Special fields that sentry uses to give more information to the server
// are extracted from entry.Data (if they are found)
// These fields are: logger and server_name
func (hook *SentryHook) Fire(entry *logrus.Entry) error {
	packet := &raven.Packet{
		Message:   entry.Message,
		Timestamp: raven.Timestamp(entry.Time),
		Level:     severityMap[entry.Level],
		Platform:  "go",
	}

	d := entry.Data

	if logger, ok := getAndDel(d, "logger"); ok {
		packet.Logger = logger
	}
	if serverName, ok := getAndDel(d, "server_name"); ok {
		packet.ServerName = serverName
	}
	packet.Extra = map[string]interface{}(d)

	_, errCh := hook.client.Capture(packet, nil)
	timeoutCh := time.After(timeout)
	select {
	case err := <-errCh:
		return err
	case <-timeoutCh:
		return fmt.Errorf("no response from sentry server in %s", timeout)
	}
	return nil
}

// Levels returns the available logging levels.
func (hook *SentryHook) Levels() []logrus.Level {
	return hook.levels
}
