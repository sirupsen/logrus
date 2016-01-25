package logrus_logstash

import (
	"bufio"
	"encoding/json"
	"net"
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestLogMessage(t *testing.T) {
	serverIsUp := make(chan bool)
	acceptConn := make(chan bool)
	finishTesting := make(chan bool)
	logFields := logrus.Fields{
		"Ninja":  "Razen",
		"Weapon": "NinjaStars",
	}

	go func(t *testing.T) {
		l, err := net.Listen("tcp", "localhost:32769")
		if err != nil {
			t.Error(err)
		}
		serverIsUp <- true
		defer l.Close()
		conn, err := l.Accept()
		acceptConn <- true
		if err != nil {
			t.Error(err)
		}
		message, err := bufio.NewReader(conn).ReadString('\n')
		var fields logrus.Fields
		json.Unmarshal([]byte(message), &fields)
		for k, v := range logFields {
			if tv, ok := fields[k]; !ok {
				t.Errorf("Expected to have the '%s' field but got none", k)
			} else if tv != v {
				t.Errorf("Expected '%s' to be set to '%s' but got '%s'", k, v, tv)
			}
		}
		finishTesting <- true
	}(t)

	<-serverIsUp
	hook, err := NewLogstashHook("tcp", "localhost:32769")
	if err != nil {
		t.Error(err)
	}
	logger := logrus.New()
	logger.Hooks.Add(hook)
	ctx := logger.WithFields(logFields)
	ctx.Info("my message")
	<-acceptConn
	<-finishTesting
}
