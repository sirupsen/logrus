// +build nacl

// FIX IT: this is a stub to allow build on nacl.
//         Need to be implemented according nacl system facilities.

package systemlog

import (
	"github.com/sirupsen/logrus"
)

// SystemlogHook to do nothing on syslogless systems.
type SystemlogHook struct {
}

// Creates a stub hook to be added to an instance of logger on syslogless systems.
func NewSystemlogHook(network, raddr string, tag string) (*SystemlogHook, error) {
	return &SystemlogHook{}, nil
}

func (hook *SystemlogHook) Fire(entry *logrus.Entry) error {
	return nil
}

func (hook *SystemlogHook) Levels() []logrus.Level {
	return make([]logrus.Level, 0, 0)
}
