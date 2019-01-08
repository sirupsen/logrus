// +build windows nacl plan9

package syslog

import (
	"github.com/sirupsen/logrus"
)

// SyslogHook to do nothing on syslogless systems.
type SyslogHook struct {
}

// Creates a stub hook to be added to an instance of logger on syslogless systems.
func NewSyslogHook(network, raddr string, priority Priority, tag string) (*SyslogHook, error) {
	return &SyslogHook{}, nil
}

func (hook *SyslogHook) Fire(entry *logrus.Entry) error {
	return nil
}

func (hook *SyslogHook) Levels() []logrus.Level {
	return make([]logrus.Level, 0, 0)
}
