package slog_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"runtime"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	lslog "github.com/sirupsen/logrus/hooks/slog"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestHookPreservesCallerPosition(t *testing.T) {
	if runtime.Compiler == "tinygo" {
		t.Log("SKIP: TinyGo does not support runtime.Caller")
		return
	}
	source, _ := test.NewNullLogger()
	source.SetReportCaller(true)
	target, capture := test.NewNullLogger()
	source.AddHook(lslog.NewHook(slog.New(lslog.NewHandler(target, &lslog.HandlerOptions{AddSource: true})), nil))
	_, file, line, _ := runtime.Caller(0)
	source.Info("original")
	entry := capture.LastEntry()
	if entry == nil || entry.Caller == nil {
		t.Fatal("missing caller")
	}
	if entry.Caller.File != file || entry.Caller.Line != line+1 {
		t.Errorf("caller got %s:%d want %s:%d", entry.Caller.File, entry.Caller.Line, file, line+1)
	}
}

func TestHookJSONSourcePosition(t *testing.T) {
	if runtime.Compiler == "tinygo" {
		t.Log("SKIP: TinyGo does not support runtime.Caller")
		return
	}
	var output bytes.Buffer
	source, _ := test.NewNullLogger()
	source.SetReportCaller(true)
	source.AddHook(lslog.NewHook(slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{AddSource: true})), nil))
	_, file, line, _ := runtime.Caller(0)
	source.Info("original")
	var record struct {
		Source struct {
			File string
			Line int
		}
	}
	if err := json.Unmarshal(output.Bytes(), &record); err != nil {
		t.Fatal(err)
	}
	if record.Source.File != file || record.Source.Line != line+1 {
		t.Errorf("source got %s:%d want %s:%d", record.Source.File, record.Source.Line, file, line+1)
	}
}

func TestHookWithoutCallerPC(t *testing.T) {
	for _, caller := range []*runtime.Frame{nil, {File: "manual.go", Line: 42}} {
		var output bytes.Buffer
		logger := logrus.New()
		hook := lslog.NewHook(slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{AddSource: true})), nil)
		if err := hook.Fire(&logrus.Entry{Logger: logger, Data: logrus.Fields{"key": "value"}, Time: time.Unix(1, 0), Level: logrus.InfoLevel, Message: "original", Caller: caller}); err != nil {
			t.Fatal(err)
		}
		var record struct {
			Msg    string
			Key    string
			Source struct {
				File string
				Line int
			}
		}
		if err := json.Unmarshal(output.Bytes(), &record); err != nil {
			t.Fatal(err)
		}
		if record.Msg != "original" || record.Key != "value" || record.Source.File != "" || record.Source.Line != 0 {
			t.Errorf("unexpected record: %+v", record)
		}
	}
}
