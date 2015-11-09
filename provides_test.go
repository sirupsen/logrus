package logrus

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

type fubarMsgAndData struct {
	msg  string
	data map[string]interface{}
}

func (f *fubarMsgAndData) String() string {
	return "FUBAR"
}

type notFubarMsgAndData fubarMsgAndData

func (nf *notFubarMsgAndData) GetLogMessage() string {
	return nf.msg
}

func (nf *notFubarMsgAndData) GetLogData() map[string]interface{} {
	return nf.data
}

func TestProvidesLogMessageAndDataFubar(t *testing.T) {
	f := &fubarMsgAndData{
		msg: "NAY",
		data: map[string]interface{}{
			"hello": "world",
			"good":  "bye",
		},
	}
	assertProvidesLogMessageAndDataFubar(f, ErrorLevel, t)
	assertProvidesLogMessageAndDataFubar(f, WarnLevel, t)
	assertProvidesLogMessageAndDataFubar(f, InfoLevel, t)
	assertProvidesLogMessageAndDataFubar(f, DebugLevel, t)
}

func TestProvidesLogMessageAndData(t *testing.T) {
	f := &notFubarMsgAndData{
		msg: "YAY",
		data: map[string]interface{}{
			"hello": "world",
			"good":  "bye",
		},
	}
	assertProvidesLogMessageAndData(f, ErrorLevel, t)
	assertProvidesLogMessageAndData(f, WarnLevel, t)
	assertProvidesLogMessageAndData(f, InfoLevel, t)
	assertProvidesLogMessageAndData(f, DebugLevel, t)
}

func assertProvidesLogMessageAndDataFubar(
	obj interface{}, lvl Level, t *testing.T) {
	b := &bytes.Buffer{}
	l := New()
	l.Level = DebugLevel
	l.Formatter = &TextFormatter{DisableColors: true}
	l.Out = b
	switch lvl {
	case ErrorLevel:
		l.Error(obj)
	case WarnLevel:
		l.Warn(obj)
	case InfoLevel:
		l.Info(obj)
	case DebugLevel:
		l.Debug(obj)
	}
	s := b.String()
	exp := fmt.Sprintf("level=%s msg=FUBAR", lvl)
	if !strings.Contains(s, exp) {
		t.Fatalf("s is `%s`, expected `%s`", s, exp)
	}
}

func assertProvidesLogMessageAndData(
	obj interface{}, lvl Level, t *testing.T) {
	b := &bytes.Buffer{}
	l := New()
	l.Level = DebugLevel
	l.Formatter = &TextFormatter{DisableColors: true}
	l.Out = b
	switch lvl {
	case ErrorLevel:
		l.Error(obj)
	case WarnLevel:
		l.Warn(obj)
	case InfoLevel:
		l.Info(obj)
	case DebugLevel:
		l.Debug(obj)
	}
	s := b.String()
	exp := fmt.Sprintf("level=%s msg=YAY good=bye hello=world", lvl)
	if !strings.Contains(s, exp) {
		t.Fatalf("s is `%s`, expected `%s`", s, exp)
	}
}
