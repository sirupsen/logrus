package systemlog

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLocalhostAddAndPrint(t *testing.T) {
	log := logrus.New()
	hook, err := NewSystemlogHook("udp", "localhost:514", "systemlog-test")

	if err != nil {
		t.Errorf("Unable to connect to local syslog.")
	}

	log.Hooks.Add(hook)

	for _, level := range hook.Levels() {
		if len(log.Hooks[level]) != 1 {
			t.Errorf("SystemlogHook was not added. The length of log.Hooks[%v]: %v", level, len(log.Hooks[level]))
		}
	}

	log.Info("Congratulations!")
}
