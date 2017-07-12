package source_file

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLocalhostAddAndPrint(t *testing.T) {
	log := logrus.New()

	hook := CodeLineHook{LogLevel: logrus.DebugLevel}
	log.Hooks.Add(&hook)

	for _, level := range hook.Levels() {
		if len(log.Hooks[level]) != 1 {
			t.Errorf("CodeLineHook was not added. The length of log.Hooks[%v]: %v", level, len(log.Hooks[level]))
		}
	}

	log.Info("Congratulations!")
}
