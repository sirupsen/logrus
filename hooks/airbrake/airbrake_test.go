package airbrake

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/airbrake/gobrake"
)

type customErr struct {
	msg string
}

func (e *customErr) Error() string {
	return e.msg
}

const (
	testAPIKey    = "abcxyz"
	testEnv       = "development"
	expectedClass = "*airbrake.customErr"
	expectedMsg   = "foo"
	unintendedMsg = "Airbrake will not see this string"
)

var (
	noticeError = make(chan *gobrake.Error, 1)
)

// TestLogEntryMessageReceived checks if invoking Logrus' log.Error
// method causes an XML payload containing the log entry message is received
// by a HTTP server emulating an Airbrake-compatible endpoint.
func TestLogEntryMessageReceived(t *testing.T) {
	log := logrus.New()
	hook := newTestHook()
	log.Hooks.Add(hook)

	log.Error(expectedMsg)

	select {
	case received := <-noticeError:
		if received.Message != expectedMsg {
			t.Errorf("Unexpected message received: %s", received.Message)
		}
	case <-time.After(time.Second):
		t.Error("Timed out; no notice received by Airbrake API")
	}
}

// TestLogEntryMessageReceived confirms that, when passing an error type using
// logrus.Fields, a HTTP server emulating an Airbrake endpoint receives the
// error message returned by the Error() method on the error interface
// rather than the logrus.Entry.Message string.
func TestLogEntryWithErrorReceived(t *testing.T) {
	log := logrus.New()
	hook := newTestHook()
	log.Hooks.Add(hook)

	log.WithFields(logrus.Fields{
		"error": &customErr{expectedMsg},
	}).Error(unintendedMsg)

	select {
	case received := <-noticeError:
		if received.Message != expectedMsg {
			t.Errorf("Unexpected message received: %s", received.Message)
		}
		if received.Type != expectedClass {
			t.Errorf("Unexpected error class: %s", received.Type)
		}
	case <-time.After(time.Second):
		t.Error("Timed out; no notice received by Airbrake API")
	}
}

// TestLogEntryWithNonErrorTypeNotReceived confirms that, when passing a
// non-error type using logrus.Fields, a HTTP server emulating an Airbrake
// endpoint receives the logrus.Entry.Message string.
//
// Only error types are supported when setting the 'error' field using
// logrus.WithFields().
func TestLogEntryWithNonErrorTypeNotReceived(t *testing.T) {
	log := logrus.New()
	hook := newTestHook()
	log.Hooks.Add(hook)

	log.WithFields(logrus.Fields{
		"error": expectedMsg,
	}).Error(unintendedMsg)

	select {
	case received := <-noticeError:
		if received.Message != unintendedMsg {
			t.Errorf("Unexpected message received: %s", received.Message)
		}
	case <-time.After(time.Second):
		t.Error("Timed out; no notice received by Airbrake API")
	}
}

// Returns a new airbrakeHook with the test server proxied
func newTestHook() *airbrakeHook {
	// Make a http.Client with the transport
	httpClient := &http.Client{Transport: &FakeRoundTripper{}}

	hook := NewHook(123, testAPIKey, "production")
	hook.Airbrake.Client = httpClient
	return hook
}

// gobrake API doesn't allow to override endpoint, we need a http.Roundtripper
type FakeRoundTripper struct {
}

func (rt *FakeRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	notice := &gobrake.Notice{}
	err = json.Unmarshal(b, notice)
	noticeError <- notice.Errors[0]

	res := &http.Response{
		StatusCode: http.StatusCreated,
		Body:       ioutil.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
	}
	return res, nil
}
