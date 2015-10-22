package logrus_bugsnag

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/bugsnag/bugsnag-go"
)

type stackFrame struct {
	Method     string `json:"method"`
	File       string `json:"file"`
	LineNumber int    `json:"lineNumber"`
}

type exception struct {
	Message    string       `json:"message"`
	Stacktrace []stackFrame `json:"stacktrace"`
}
type notice struct {
	Events []struct {
		Exceptions []exception `json:"exceptions"`
	} `json:"events"`
}

func TestNoticeReceived(t *testing.T) {
	msg := make(chan exception, 1)
	expectedMsg := "foo"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var notice notice
		data, _ := ioutil.ReadAll(r.Body)
		if err := json.Unmarshal(data, &notice); err != nil {
			t.Error(err)
		}
		_ = r.Body.Close()

		msg <- notice.Events[0].Exceptions[0]
	}))
	defer ts.Close()

	hook := &bugsnagHook{}

	bugsnag.Configure(bugsnag.Configuration{
		Endpoint:     ts.URL,
		ReleaseStage: "production",
		APIKey:       "12345678901234567890123456789012",
		Synchronous:  true,
	})

	log := logrus.New()
	log.Hooks.Add(hook)

	log.WithFields(logrus.Fields{
		"error": errors.New(expectedMsg),
	}).Error("Bugsnag will not see this string")

	select {
	case received := <-msg:
		message := received.Message
		if message != expectedMsg {
			t.Errorf("Unexpected message received: %s", received)
		}
		if len(received.Stacktrace) < 1 {
			t.Error("Bugsnag error does not have a stack trace")
		}
		topFrame := received.Stacktrace[0]
		if topFrame.Method != "TestNoticeReceived" {
			t.Errorf("Unexpected method on top of call stack: '%s' (should be 'TestNoticeReceived')", topFrame.Method)
		}
	case <-time.After(time.Second):
		t.Error("Timed out; no notice received by Bugsnag API")
	}
}
