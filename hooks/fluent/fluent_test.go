package fluent

import (
	"bytes"
	"net"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
)

var message chan string

const (
	testHOST = "localhost"
)

func TestLogEntryMessageReceived(t *testing.T) {
	message = make(chan string, 1)
	port := startMockServer(t)
	hook := NewHook(testHOST, port)
	logger := logrus.New()
	logger.Hooks.Add(hook)

	logger.WithFields(logrus.Fields{
		"message": "message!",
		"tag":     "debug.test",
		"value":   "data",
	}).Error("hoge")

	received := <-message
	switch {
	case !strings.Contains(received, "\x94\xaadebug.test\xd2"):
		t.Errorf("message did not contain tag")
	case !strings.Contains(received, "value\xa4data"):
		t.Errorf("message did not contain value")
	case !strings.Contains(received, "\xa7message\xa8message!"):
		t.Errorf("message did not contain message")
	}
}

func startMockServer(t *testing.T) int {
	l, err := net.Listen("tcp", testHOST+":0")
	if err != nil {
		t.Errorf("Error listening:", err.Error())
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				t.Errorf("Error accepting:", err.Error())
			}
			go handleRequest(conn, l)
		}
	}()
	return l.Addr().(*net.TCPAddr).Port
}

func handleRequest(conn net.Conn, l net.Listener) {
	bf := new(bytes.Buffer)
	bf.ReadFrom(conn)
	conn.Close()
	message <- bf.String()
}
