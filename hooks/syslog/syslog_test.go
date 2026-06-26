// TinyGo currently (v0.41.1) doesn't fully support `log/syslog`;
// https://tinygo.org/docs/reference/lang-support/stdlib/#logsyslog
//go:build !windows && !nacl && !plan9 && !tinygo

package syslog_test

import (
	"io"
	"log/syslog"
	"testing"

	"github.com/sirupsen/logrus"
	lsyslog "github.com/sirupsen/logrus/hooks/syslog"
)

func TestLocalhostAddAndPrint(t *testing.T) {
	log := logrus.New()
	log.SetOutput(io.Discard)
	hook, err := lsyslog.NewSyslogHook("udp", "localhost:514", syslog.LOG_INFO, "")
	if err != nil {
		t.Errorf("Unable to connect to local syslog.")
	}

	log.Hooks.Add(hook)

	for _, level := range hook.Levels() {
		if len(log.Hooks[level]) != 1 {
			t.Errorf("SyslogHook was not added. The length of log.Hooks[%v]: %v", level, len(log.Hooks[level]))
		}
	}

	log.Info("Congratulations!")
}
