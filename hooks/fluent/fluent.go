package fluent

import (
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"

	"github.com/Sirupsen/logrus"
)

var defaultLevels = []logrus.Level{
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
}

type fluentHook struct {
	host   string
	port   int
	levels []logrus.Level
}

func NewHook(host string, port int) *fluentHook {
	return &fluentHook{
		host:   host,
		port:   port,
		levels: defaultLevels,
	}
}

func getTagAndDel(entry *logrus.Entry) string {
	const key = "tag"
	var v interface{}
	var ok bool
	if v, ok = entry.Data[key]; !ok {
		return entry.Message
	}

	var val string
	if val, ok = v.(string); !ok {
		return entry.Message
	}
	delete(entry.Data, key)
	return val
}

func (hook *fluentHook) Fire(entry *logrus.Entry) error {
	logger, err := fluent.New(fluent.Config{
		FluentHost: hook.host,
		FluentPort: hook.port,
	})
	if err != nil {
		return err
	}
	defer logger.Close()

	tag := getTagAndDel(entry)
	err = logger.PostWithTime(tag, time.Now(), entry.Data)
	return err
}

func (hook *fluentHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *fluentHook) SetLevels(levels []logrus.Level) {
	hook.levels = levels
}
