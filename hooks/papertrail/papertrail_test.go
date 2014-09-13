package logrus_papertrail

import (
	"net"
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestWritingToUDP(t *testing.T) {
	log := logrus.New()
	port := 16661

	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("127.0.0.1"),
	}

	c, err := net.ListenUDP("udp", &addr)
	if err != nil {
		t.Fatalf("ListenUDP failed: %v", err)
	}
	defer c.Close()

	hook, err := NewPapertrailHook("localhost", port, "test")
	if err != nil {
		t.Errorf("Unable to connect to local UDP server.")
	}

	log.Hooks.Add(hook)
	log.Info("Today was a good day.")

	var buf = make([]byte, 1500)
	n, _, err := c.ReadFromUDP(buf)

	if err != nil {
		t.Fatalf("Error reading data from local UDP server")
	}

	if n <= 0 {
		t.Errorf("Nothing written to local UDP server.")
	}
}
