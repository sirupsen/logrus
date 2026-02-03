package test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAllHooks(t *testing.T) {
	assert := assert.New(t)

	logger, hook := NewNullLogger()
	assert.Nil(hook.LastEntry())
	assert.Empty(hook.Entries)

	logger.Error("Hello error")
	assert.Equal(logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal("Hello error", hook.LastEntry().Message)
	assert.Len(hook.Entries, 1)

	logger.Warn("Hello warning")
	assert.Equal(logrus.WarnLevel, hook.LastEntry().Level)
	assert.Equal("Hello warning", hook.LastEntry().Message)
	assert.Len(hook.Entries, 2)

	hook.Reset()
	assert.Nil(hook.LastEntry())
	assert.Empty(hook.Entries)

	hook = NewGlobal()

	logrus.Error("Hello error")
	assert.Equal(logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal("Hello error", hook.LastEntry().Message)
	assert.Len(hook.Entries, 1)
}

func TestLoggingWithHooksRace(t *testing.T) {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	unlocker := r.Intn(100)

	assert := assert.New(t)
	logger, hook := NewNullLogger()

	var wgOne, wgAll sync.WaitGroup
	wgOne.Add(1)
	wgAll.Add(100)

	for i := 0; i < 100; i++ {
		go func(i int) {
			logger.Info("info")
			wgAll.Done()
			if i == unlocker {
				wgOne.Done()
			}
		}(i)
	}

	wgOne.Wait()

	assert.Equal(logrus.InfoLevel, hook.LastEntry().Level)
	assert.Equal("info", hook.LastEntry().Message)

	wgAll.Wait()

	entries := hook.AllEntries()
	assert.Len(entries, 100)
}

// nolint:staticcheck // linter assumes logger.Fatal exits, resulting in false SA4006 warnings.
func TestFatalWithAlternateExit(t *testing.T) {
	assert := assert.New(t)

	logger, hook := NewNullLogger()
	logger.ExitFunc = func(code int) {}

	logger.Fatal("something went very wrong")
	assert.Equal(logrus.FatalLevel, hook.LastEntry().Level)
	assert.Equal("something went very wrong", hook.LastEntry().Message)
	assert.Len(hook.Entries, 1)
}

func TestNewLocal(t *testing.T) {
	assert := assert.New(t)
	logger := logrus.New()

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			logger.Info("info")
			wg.Done()
		}()
	}

	hook := NewLocal(logger)
	assert.NotNil(hook)
}
