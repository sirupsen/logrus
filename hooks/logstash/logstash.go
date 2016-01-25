// +build !windows,!nacl,!plan9

package logrus_logstash

import (
	"net"

	"github.com/Sirupsen/logrus"
	logrus_logstash_fmt "github.com/Sirupsen/logrus/formatters/logstash"
)

// LogstashHook represents a connection to Logstash.
type LogstashHook struct {
	conn net.Conn
}

// NewLogstashHook creates a hook to a Logstash instance that is available at
// `Address` through `Protocol` (one of these: "tcp", "udp").
func NewLogstashHook(Protocol, Address string) (*LogstashHook, error) {
	conn, err := net.Dial(Protocol, Address)
	if err != nil {
		return nil, err
	}
	return &LogstashHook{conn: conn}, nil
}

func (h *LogstashHook) Fire(entry *logrus.Entry) error {
	formatter := logrus_logstash_fmt.LogstashFormatter{}
	dataBytes, err := formatter.Format(entry)
	if err != nil {
		return err
	}
	if _, err = h.conn.Write(dataBytes); err != nil {
		return err
	}
	return nil
}

func (h *LogstashHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
