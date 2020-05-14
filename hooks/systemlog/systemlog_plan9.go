// +build plan9

// FIX IT: this is a stub to allow build on plan9.
//         Need to be implemented according plan9 system facilities.

package lachesis_log

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
