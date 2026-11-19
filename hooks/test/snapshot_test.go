package test_test

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
)

func TestHookAccessorsReturnIndependentEntries(t *testing.T) {
	for _, accessor := range []string{"last", "all"} {
		t.Run(accessor, func(t *testing.T) {
			logger, hook := logtest.NewNullLogger()
			logger.WithField("key", "original").Warn("original message")
			get := func() *logrus.Entry {
				if accessor == "last" {
					return hook.LastEntry()
				}
				return hook.AllEntries()[0]
			}
			first := get()
			second := get()
			first.Message = "changed"
			first.Level = logrus.DebugLevel
			first.Data["key"] = "changed"
			first.Data["extra"] = true
			for _, entry := range []*logrus.Entry{second, get()} {
				if entry.Message != "original message" || entry.Level != logrus.WarnLevel {
					t.Errorf("accessor changed recorded entry metadata: message=%q level=%v", entry.Message, entry.Level)
				}
				if entry.Data["key"] != "original" || len(entry.Data) != 1 {
					t.Errorf("accessor changed recorded fields: %#v", entry.Data)
				}
			}
			if hook.Entries[0].Message != "original message" || hook.Entries[0].Data["key"] != "original" {
				t.Error("accessor modified exported recorded entries")
			}
			hook.Reset()
			logger.Info("new message")
			if second.Message != "original message" {
				t.Error("reset and logging changed the previous snapshot")
			}
		})
	}
}

func TestHookAccessorEmptyAndNilData(t *testing.T) {
	hook := new(logtest.Hook)
	if hook.LastEntry() != nil || len(hook.AllEntries()) != 0 {
		t.Fatal("empty hook returned an entry")
	}
	entry := logrus.Entry{Message: "message", Level: logrus.WarnLevel}
	if err := hook.Fire(&entry); err != nil {
		t.Fatal(err)
	}
	if got := hook.LastEntry(); got.Data != nil || got.Message != entry.Message || got.Level != entry.Level {
		t.Errorf("nil-data entry was not preserved: %#v", got)
	}
	if got := hook.AllEntries()[0]; got.Data != nil || got.Message != entry.Message || got.Level != entry.Level {
		t.Errorf("nil-data entry was not preserved: %#v", got)
	}
}

func TestHookAccessorsPreserveEventMetadata(t *testing.T) {
	logger, hook := logtest.NewNullLogger()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	caller := &runtime.Frame{Function: "example.run", File: "example.go", Line: 12}
	timestamp := time.Date(2026, time.January, 1, 2, 3, 4, 0, time.UTC)
	entry := logrus.Entry{Logger: logger, Context: ctx, Caller: caller, Time: timestamp, Level: logrus.ErrorLevel, Message: "event", Data: logrus.Fields{"key": "value"}}
	if err := hook.Fire(&entry); err != nil {
		t.Fatal(err)
	}
	for _, got := range []*logrus.Entry{hook.LastEntry(), hook.AllEntries()[0]} {
		if got.Logger != logger || got.Context != ctx || got.Caller != caller || got.Time != timestamp || got.Level != entry.Level || got.Message != entry.Message {
			t.Errorf("accessor lost event metadata: %#v", got)
		}
	}
}

func TestHookAccessorsConcurrentSnapshots(t *testing.T) {
	for _, accessor := range []string{"last", "all"} {
		t.Run(accessor, func(t *testing.T) {
			logger, hook := logtest.NewNullLogger()
			logger.Info("original")
			var wg sync.WaitGroup
			for range 4 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for range 100 {
						var got *logrus.Entry
						if accessor == "last" {
							got = hook.LastEntry()
						} else {
							got = hook.AllEntries()[0]
						}
						got.Message = "local edit"
					}
				}()
			}
			wg.Wait()
			if got := hook.LastEntry().Message; got != "original" {
				t.Errorf("recorded message = %q", got)
			}
		})
	}
}
