package logrus_sentry

import (
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
)

const (
	message     = "error message"
	server_name = "testserver.internal"
	logger_name = "test.logger"
)

func getTestLogger() *logrus.Logger {
	l := logrus.New()
	l.Out = ioutil.Discard
	return l
}

type resultPacket struct {
	raven.Packet
	Stacktrace raven.Stacktrace `json:"sentry.interfaces.Stacktrace"`
}

func WithTestDSN(t *testing.T, tf func(string, <-chan *resultPacket)) {
	pch := make(chan *resultPacket, 1)
	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		contentType := req.Header.Get("Content-Type")
		var bodyReader io.Reader = req.Body
		if contentType == "application/octet-stream" {
			bodyReader = base64.NewDecoder(base64.StdEncoding, bodyReader)
			bodyReader, _ = zlib.NewReader(bodyReader)
		}
		d := json.NewDecoder(bodyReader)
		p := &resultPacket{}
		err := d.Decode(p)
		if err != nil {
			t.Fatal(err.Error())
		}

		pch <- p
	}))
	defer s.Close()

	fragments := strings.SplitN(s.URL, "://", 2)
	dsn := fmt.Sprintf(
		"%s://public:secret@%s/sentry/project-id",
		fragments[0],
		fragments[1],
	)
	tf(dsn, pch)
}

func TestSpecialFields(t *testing.T) {
	WithTestDSN(t, func(dsn string, pch <-chan *resultPacket) {
		logger := getTestLogger()

		hook, err := NewSentryHook(dsn, []logrus.Level{
			logrus.ErrorLevel,
		})

		if err != nil {
			t.Fatal(err.Error())
		}
		logger.Hooks.Add(hook)

		req, _ := http.NewRequest("GET", "url", nil)
		logger.WithFields(logrus.Fields{
			"server_name":  server_name,
			"logger":       logger_name,
			"http_request": req,
		}).Error(message)

		packet := <-pch
		if packet.Logger != logger_name {
			t.Errorf("logger should have been %s, was %s", logger_name, packet.Logger)
		}

		if packet.ServerName != server_name {
			t.Errorf("server_name should have been %s, was %s", server_name, packet.ServerName)
		}
	})
}

func TestSentryHandler(t *testing.T) {
	WithTestDSN(t, func(dsn string, pch <-chan *resultPacket) {
		logger := getTestLogger()
		hook, err := NewSentryHook(dsn, []logrus.Level{
			logrus.ErrorLevel,
		})
		if err != nil {
			t.Fatal(err.Error())
		}
		logger.Hooks.Add(hook)

		logger.Error(message)
		packet := <-pch
		if packet.Message != message {
			t.Errorf("message should have been %s, was %s", message, packet.Message)
		}
	})
}

func TestSentryStacktrace(t *testing.T) {
	WithTestDSN(t, func(dsn string, pch <-chan *resultPacket) {
		logger := getTestLogger()
		hook, err := NewSentryHook(dsn, []logrus.Level{
			logrus.ErrorLevel,
			logrus.InfoLevel,
		})
		hook.SetStacktraceLevel(logrus.ErrorLevel)
		if err != nil {
			t.Fatal(err.Error())
		}
		logger.Hooks.Add(hook)

		logger.Error(message) // this is the call that the last frame of stacktrace should capture
		packet := <-pch
		if packet.Message != message {
			t.Errorf("message should have been %s, was %s", message, packet.Message)
		}
		stacktraceSize := len(packet.Stacktrace.Frames)
		if stacktraceSize == 0 {
			t.Error("Stacktrace should not be empty")
		}
		lastFrame := packet.Stacktrace.Frames[stacktraceSize-1]
		expectedFilename := "github.com/Sirupsen/logrus/hooks/sentry/sentry_test.go"
		if lastFrame.Filename != expectedFilename {
			t.Errorf("File name should have been %s, was %s", expectedFilename, lastFrame.Filename)
		}
		expectedLineno := 129
		if lastFrame.Lineno != expectedLineno {
			t.Errorf("Line number should have been %s, was %s", expectedLineno, lastFrame.Lineno)
		}

		logger.Info(message)
		packet = <-pch
		if packet.Message != message {
			t.Errorf("message should have been %s, was %s", message, packet.Message)
		}
		stacktraceSize = len(packet.Stacktrace.Frames)
		if stacktraceSize > 0 {
			t.Error("Stacktrace should be empty")
		}
	})
}
