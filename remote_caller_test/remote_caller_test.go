// test remote caller from an external packages point of view
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"testing"

	log "github.com/Sirupsen/logrus"
)

var expectedFilePattern = regexp.MustCompile(`remote_caller_test\.go:\d+$`)

func TestVanilla(t *testing.T) {
	var buff bytes.Buffer
	var fields log.Fields

	log.SetOutput(&buff)
	log.SetFormatter(&log.JSONFormatter{})

	log.Info("dummy")

	err := json.Unmarshal(buff.Bytes(), &fields)
	if err != nil {
		t.Errorf("should have decoded message, got: %s", err)
	}

	caller := fields["caller"].(string)
	t.Logf("got caller: %s", caller)
	if !expectedFilePattern.MatchString(caller) {
		t.Errorf("found unexpected caller: %s", caller)
	}
}

func TestChainedField(t *testing.T) {
	var buff bytes.Buffer
	var fields log.Fields

	log.SetOutput(&buff)
	log.SetFormatter(&log.JSONFormatter{})

	log.WithField("foo", "bar").Info("baz")

	err := json.Unmarshal(buff.Bytes(), &fields)
	if err != nil {
		t.Errorf("should have decoded message, got: %s", err)
	}

	caller := fields["caller"].(string)
	t.Logf("got caller: %s", caller)
	if !expectedFilePattern.MatchString(caller) {
		t.Errorf("found unexpected caller: %s", caller)
	}
}

func TestReusedField(t *testing.T) {
	var buff bytes.Buffer
	var fields log.Fields

	log.SetOutput(&buff)
	log.SetFormatter(&log.JSONFormatter{})

	foolog := log.WithField("foo", "bar")
	foolog.Info("baz")
	buff.Reset()
	foolog.Info("baz2")

	err := json.Unmarshal(buff.Bytes(), &fields)
	if err != nil {
		t.Errorf("should have decoded message, got: %s", err)
	}

	caller := fields["caller"].(string)
	t.Logf("got caller: %s", caller)
	if !expectedFilePattern.MatchString(caller) {
		t.Errorf("found unexpected caller: %s", caller)
	}
}

func TestReusedMultipleField(t *testing.T) {
	var buff bytes.Buffer
	var fields log.Fields

	log.SetOutput(&buff)
	log.SetFormatter(&log.JSONFormatter{})

	foolog := log.WithFields(log.Fields{"foo": "bar", "q": 42})
	foolog.Info("baz")
	buff.Reset()
	foolog.Info("baz2")

	err := json.Unmarshal(buff.Bytes(), &fields)
	if err != nil {
		t.Errorf("should have decoded message, got: %s", err)
	}

	caller := fields["caller"].(string)
	t.Logf("got caller: %s", caller)
	if !expectedFilePattern.MatchString(caller) {
		t.Errorf("found unexpected caller: %s", caller)
	}
}

func TestMultipleChanins(t *testing.T) {
	var buff bytes.Buffer
	var fields log.Fields

	log.SetOutput(&buff)
	log.SetFormatter(&log.JSONFormatter{})

	foolog := log.WithFields(log.Fields{"foo": "bar", "q": 42}).WithField("such", "field")
	foolog.Info("baz")
	buff.Reset()
	foolog.Info("baz2")

	err := json.Unmarshal(buff.Bytes(), &fields)
	if err != nil {
		t.Errorf("should have decoded message, got: %s", err)
	}

	caller := fields["caller"].(string)
	t.Logf("got caller: %s", caller)
	if !expectedFilePattern.MatchString(caller) {
		t.Errorf("found unexpected caller: %s", caller)
	}
}

func TestErrorField(t *testing.T) {
	var buff bytes.Buffer
	var fields log.Fields

	log.SetOutput(&buff)
	log.SetFormatter(&log.JSONFormatter{})

	foolog := log.WithError(errors.New("wow, much error"))
	foolog.Info("baz")

	err := json.Unmarshal(buff.Bytes(), &fields)
	if err != nil {
		t.Errorf("should have decoded message, got: %s", err)
	}

	caller := fields["caller"].(string)
	t.Logf("got caller: %s", caller)
	if !expectedFilePattern.MatchString(caller) {
		t.Errorf("found unexpected caller: %s", caller)
	}
}

func TestNoCaller(t *testing.T) {
	var buff bytes.Buffer
	var fields log.Fields

	l := log.New()
	l.ShowCaller = false
	l.Out = &buff
	l.Formatter = &log.JSONFormatter{}

	l.Info("baz")

	err := json.Unmarshal(buff.Bytes(), &fields)
	if err != nil {
		t.Errorf("should have decoded message, got: %s", err)
	}

	if caller, exists := fields["caller"]; exists {
		t.Errorf("caller not expected, but found one: %s", caller)
	}

}
