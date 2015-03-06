package airbrake

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
)

type notice struct {
	Error struct {
		Message string `xml:"message"`
	} `xml:"error"`
}

func TestNoticeReceived(t *testing.T) {
	msg := make(chan string, 1)
	expectedMsg := "foo"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var notice notice
		if err := xml.NewDecoder(r.Body).Decode(&notice); err != nil {
			t.Error(err)
		}
		r.Body.Close()

		msg <- notice.Error.Message
	}))
	defer ts.Close()

	hook := NewHook(ts.URL, "foo", "production")

	log := logrus.New()
	log.Hooks.Add(hook)
	log.Error(expectedMsg)

	select {
	case received := <-msg:
		if received != expectedMsg {
			t.Errorf("Unexpected message received: %s", received)
		}
	case <-time.After(time.Second):
		t.Error("Timed out; no notice received by Airbrake API")
	}
}
