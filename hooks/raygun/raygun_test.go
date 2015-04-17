package raygun

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/sditools/goraygun"
)

type customErr struct {
	msg string
}

func (e *customErr) Error() string {
	return e.msg
}

const (
	testAPIKey    = "abcxyz"
	expectedClass = "github.com/Sirupsen/logrus/hooks/raygun"
	expectedMsg   = "oh no some error occured."
	unintendedMsg = "Airbrake will not see this string"
)

var (
	entryCh = make(chan goraygun.Entry, 1)
)

// TestLogEntryMessageReceived checks if invoking Logrus' log.Error
// method causes an XML payload containing the log entry message is received
// by a HTTP server emulating an Airbrake-compatible endpoint.
func TestLogEntryMessageReceived(t *testing.T) {
	log := logrus.New()
	ts := startRaygunServer(t)
	defer ts.Close()

	hook := NewHook(ts.URL, testAPIKey, true)
	log.Hooks.Add(hook)

	log.Error(expectedMsg)

	select {
	case received := <-entryCh:
		if received.Details.Error.Message != expectedMsg {
			t.Errorf("Unexpected message received: %s", received)
		}
	case <-time.After(time.Second):
		t.Error("Timed out; no notice received by Raygun API")
	}
}

// TestLogEntryMessageReceived confirms that, when passing an error type using
// logrus.Fields, a HTTP server emulating an Airbrake endpoint receives the
// error message returned by the Error() method on the error interface
// rather than the logrus.Entry.Message string.
func TestLogEntryWithErrorReceived(t *testing.T) {
	log := logrus.New()
	ts := startRaygunServer(t)
	defer ts.Close()

	hook := NewHook(ts.URL, testAPIKey, true)
	log.Hooks.Add(hook)

	log.WithFields(logrus.Fields{
		"error": &customErr{expectedMsg},
	}).Error(unintendedMsg)

	select {
	case received := <-entryCh:
		if received.Details.Error.Message != expectedMsg {
			t.Errorf("Unexpected message received: %s", received.Details.Error.Message)
		}
		if received.Details.Error.ClassName != expectedClass {
			t.Errorf("Unexpected error class: %s", received.Details.Error.ClassName)
		}
	case <-time.After(time.Second):
		t.Error("Timed out; no notice received by Airbrake API")
	}
}

func startRaygunServer(t *testing.T) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var entry goraygun.Entry
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			t.Error(err)
		}
		r.Body.Close()

		entryCh <- entry
	}))

	return ts
}
