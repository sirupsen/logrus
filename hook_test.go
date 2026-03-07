package logrus_test

import (
	"fmt"
	"io"
	"maps"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "github.com/sirupsen/logrus/internal/testutils"
)

// RecordingFormatter is a test helper that implements Formatter and
// records information about the last formatted Entry.
//
// On every call to Format, it increments Calls and stores a shallow copy
// of entry.Data in EntryData, overwriting any previous value. The formatted
// output is the raw entry.Message as bytes.
type RecordingFormatter struct {
	Calls     atomic.Int64
	EntryData Fields
}

func (f *RecordingFormatter) Format(entry *Entry) ([]byte, error) {
	f.Calls.Add(1)
	f.EntryData = maps.Clone(entry.Data)
	return []byte(entry.Message), nil
}

type TestHook struct {
	Fired bool
}

func (hook *TestHook) Fire(entry *Entry) error {
	hook.Fired = true
	return nil
}

func (hook *TestHook) Levels() []Level {
	return []Level{
		TraceLevel,
		DebugLevel,
		InfoLevel,
		WarnLevel,
		ErrorLevel,
		FatalLevel,
		PanicLevel,
	}
}

func TestHookFires(t *testing.T) {
	hook := new(TestHook)

	LogAndAssertJSON(t, func(log *Logger) {
		log.Hooks.Add(hook)
		assert.False(t, hook.Fired)

		log.Print("test")
	}, func(fields Fields) {
		assert.True(t, hook.Fired)
	})
}

type ModifyHook struct {
	Calls atomic.Int64
}

func (hook *ModifyHook) Fire(entry *Entry) error {
	hook.Calls.Add(1)
	entry.Data["wow"] = "whale"
	return nil
}

func (hook *ModifyHook) Levels() []Level {
	return []Level{
		TraceLevel,
		DebugLevel,
		InfoLevel,
		WarnLevel,
		ErrorLevel,
		FatalLevel,
		PanicLevel,
	}
}

func TestHookCanModifyEntry(t *testing.T) {
	hook := new(ModifyHook)

	LogAndAssertJSON(t, func(log *Logger) {
		log.Hooks.Add(hook)
		log.WithField("wow", "elephant").Print("test")
	}, func(fields Fields) {
		assert.Equal(t, "whale", fields["wow"])
	})
}

func TestCanFireMultipleHooks(t *testing.T) {
	hook1 := new(ModifyHook)
	hook2 := new(TestHook)

	LogAndAssertJSON(t, func(log *Logger) {
		log.Hooks.Add(hook1)
		log.Hooks.Add(hook2)

		log.WithField("wow", "elephant").Print("test")
	}, func(fields Fields) {
		assert.Equal(t, "whale", fields["wow"])
		assert.True(t, hook2.Fired)
	})
}

type SingleLevelModifyHook struct {
	ModifyHook
}

func (h *SingleLevelModifyHook) Levels() []Level {
	return []Level{InfoLevel}
}

// TestHookEntryIsPristine tests that each log gets a pristine copy of Entry,
// and changes from modifying hooks are not persisted.
//
// Regression test for https://github.com/sirupsen/logrus/issues/795
func TestHookEntryIsPristine(t *testing.T) {
	formatter := &RecordingFormatter{}
	hook := &SingleLevelModifyHook{}
	l := New()
	l.SetOutput(io.Discard)
	l.SetFormatter(formatter)
	l.AddHook(hook)

	// Initial message should have a pristine copy of Entry.
	l.Error("first")
	assert.Equal(t, int64(0), hook.Calls.Load())
	assert.Equal(t, int64(1), formatter.Calls.Load())
	require.Empty(t, formatter.EntryData)

	// Info message modifies data through SingleLevelModifyHook
	l.Info("second")
	assert.Equal(t, int64(1), hook.Calls.Load())
	assert.Equal(t, int64(2), formatter.Calls.Load())
	require.Equal(t, Fields{"wow": "whale"}, formatter.EntryData)

	// Should have a pristine copy of Entry.
	l.Error("third")
	assert.Equal(t, int64(1), hook.Calls.Load())
	assert.Equal(t, int64(3), formatter.Calls.Load())
	require.Empty(t, formatter.EntryData)
}

type ErrorHook struct {
	Fired bool
}

func (hook *ErrorHook) Fire(entry *Entry) error {
	hook.Fired = true
	return nil
}

func (hook *ErrorHook) Levels() []Level {
	return []Level{
		ErrorLevel,
	}
}

func TestErrorHookShouldntFireOnInfo(t *testing.T) {
	hook := new(ErrorHook)

	LogAndAssertJSON(t, func(log *Logger) {
		log.Hooks.Add(hook)
		log.Info("test")
	}, func(fields Fields) {
		assert.False(t, hook.Fired)
	})
}

func TestErrorHookShouldFireOnError(t *testing.T) {
	hook := new(ErrorHook)

	LogAndAssertJSON(t, func(log *Logger) {
		log.Hooks.Add(hook)
		log.Error("test")
	}, func(fields Fields) {
		assert.True(t, hook.Fired)
	})
}

func TestAddHookRace(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	hook := new(ErrorHook)
	LogAndAssertJSON(t, func(log *Logger) {
		go func() {
			defer wg.Done()
			log.AddHook(hook)
		}()
		go func() {
			defer wg.Done()
			log.Error("test")
		}()
		wg.Wait()
	}, func(fields Fields) {
		// the line may have been logged
		// before the hook was added, so we can't
		// actually assert on the hook
	})
}

func TestAddHookRace2(t *testing.T) {
	// Test modifies the standard-logger; restore it afterward.
	stdLogger := StandardLogger()
	oldOut := stdLogger.Out
	oldHooks := stdLogger.ReplaceHooks(make(LevelHooks))
	t.Cleanup(func() {
		stdLogger.SetOutput(oldOut)
		stdLogger.ReplaceHooks(oldHooks)
	})
	stdLogger.SetOutput(io.Discard)

	for i := range 3 {
		testname := fmt.Sprintf("Test %d", i)
		t.Run(testname, func(t *testing.T) {
			t.Parallel()

			_ = test.NewGlobal()
			Info(testname)
		})
	}
}

type HookCallFunc struct {
	F func()
}

func (h *HookCallFunc) Levels() []Level {
	return AllLevels
}

func (h *HookCallFunc) Fire(e *Entry) error {
	h.F()
	return nil
}

func TestHookFireOrder(t *testing.T) {
	checkers := []string{}
	h := LevelHooks{}
	h.Add(&HookCallFunc{F: func() { checkers = append(checkers, "first hook") }})
	h.Add(&HookCallFunc{F: func() { checkers = append(checkers, "second hook") }})
	h.Add(&HookCallFunc{F: func() { checkers = append(checkers, "third hook") }})

	if err := h.Fire(InfoLevel, &Entry{}); err != nil {
		t.Error("unexpected error:", err)
	}
	require.Equal(t, []string{"first hook", "second hook", "third hook"}, checkers)
}
