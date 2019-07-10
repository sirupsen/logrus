package logatee

import (
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
)

var (
	logs = make(map[*logrus.Logger][]logrus.Entry)
	lock sync.RWMutex
)

type hook struct {
	tFunc  func() *testing.T
	logger *logrus.Logger
}

// Levels implements the logrus.Hook interface
// It returns all levels
func (h *hook) Levels() []logrus.Level {
	var lvls []logrus.Level

	// loop backwards through levels until the uint flips from 0s to 1s
	for l := uint32(logrus.TraceLevel); l <= uint32(logrus.TraceLevel); l-- {
		lvls = append(lvls, logrus.Level(l))
	}

	return lvls
}

// Fire implements the logrus.Hook interface
// It saves the entry to a list associated with the logger
func (h *hook) Fire(e *logrus.Entry) error {
	lock.Lock()
	defer lock.Unlock()
	logs[h.logger] = append(logs[h.logger], *e)
	return nil
}

// Logs returns all entries that were written through logger
func Logs(logger *logrus.Logger) []logrus.Entry {
	lock.RLock()
	defer lock.RUnlock()
	return logs[logger]
}

// Reset clears all entries associated with logger. Use it to
// reset tracking between tests
func Reset(logger *logrus.Logger) {
	lock.Lock()
	defer lock.Unlock()
	logs[logger] = nil
}
