package logrus_papertrail

import (
	"fmt"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stvp/go-udp-testing"
)

func TestWritingToUDP(t *testing.T) {
	port := 16661
	udp.SetAddr(fmt.Sprintf(":%d", port))

	hook, err := NewPapertrailHook("localhost", port, "test")
	if err != nil {
		t.Errorf("Unable to connect to local UDP server.")
	}

	log := logrus.New()
	log.Hooks.Add(hook)

	udp.ShouldReceive(t, "foo", func() {
		log.Info("foo")
	})
}

func TestUseHostname(t *testing.T) {
	port := 16661

	hostname, _ := os.Hostname()
	udp.SetAddr(fmt.Sprintf(":%d", port))

	hook, err := NewPapertrailHook("localhost", port, "test")
	if err != nil {
		t.Errorf("Unable to connect to local UDP server.")
	}
	hook.UseHostname()
	if hook.Hostname != hostname {
		t.Errorf("Expected %s got %s", hostname, hook.Hostname)
	}
}
