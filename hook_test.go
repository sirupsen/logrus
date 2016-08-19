package logrus

import (
	"io/ioutil"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestHook struct {
	Fired bool
}

func (hook *TestHook) Fire(entry *Entry) error {
	hook.Fired = true
	return nil
}

func (hook *TestHook) Levels() []Level {
	return []Level{
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
		assert.Equal(t, hook.Fired, false)

		log.Print("test")
	}, func(fields Fields) {
		assert.Equal(t, hook.Fired, true)
	})
}

type ModifyHook struct {
}

func (hook *ModifyHook) Fire(entry *Entry) error {
	entry.Data["wow"] = "whale"
	return nil
}

func (hook *ModifyHook) Levels() []Level {
	return []Level{
		DebugLevel,
		InfoLevel,
		WarnLevel,
		ErrorLevel,
		FatalLevel,
		PanicLevel,
	}
}

type SynchronizedModifyHook struct{}

func (hook *SynchronizedModifyHook) Levels() []Level {
	return AllLevels
}

func (hook *SynchronizedModifyHook) Fire(entry *Entry) error {
	entry.SetField("wow", "whale")
	return nil
}

func TestHookCanModifyEntry(t *testing.T) {
	hook := new(ModifyHook)

	LogAndAssertJSON(t, func(log *Logger) {
		log.Hooks.Add(hook)
		log.WithField("wow", "elephant").Print("test")
	}, func(fields Fields) {
		assert.Equal(t, fields["wow"], "whale")
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
		assert.Equal(t, fields["wow"], "whale")
		assert.Equal(t, hook2.Fired, true)
	})
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
		assert.Equal(t, hook.Fired, false)
	})
}

func TestErrorHookShouldFireOnError(t *testing.T) {
	hook := new(ErrorHook)

	LogAndAssertJSON(t, func(log *Logger) {
		log.Hooks.Add(hook)
		log.Error("test")
	}, func(fields Fields) {
		assert.Equal(t, hook.Fired, true)
	})
}

func TestEntryDataRaceWithHook(t *testing.T) {
	hook := new(SynchronizedModifyHook)

	logger := New()
	logger.Hooks.Add(hook)
	logger.Out = ioutil.Discard

	defer func() {
		assert.Nil(t, recover())
	}()

	entry := logger.WithFields(Fields{"a": "b"})

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			switch i % 2 {
			case 0:
				entry.Info("Hello")
			case 1:
				entry.WithFields(Fields{
					"something": "else",
				}).Info("World")
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
