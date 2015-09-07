package logrus

import (
	"fmt"
	"os"
	"sync"
)

// A hook to be fired when logging on the logging levels returned from
// `Levels()` on your implementation of the interface.
type Hook interface {
	Levels() []Level
	Fire(*Entry) error
}

// A hook that should be fired asynchronously. The Async() method is
// a no-op that simply distinguishes asynchronous hooks from regular,
// synchronous ones.
type AsyncHook interface {
	Hook
	Async()
}

// Internal type for storing the hooks on a logger instance.
type LevelHooks map[Level][]Hook

// Add a hook to an instance of logger. This is called with
// `log.Hooks.Add(new(MyHook))` where `MyHook` implements the `Hook` interface.
func (hooks LevelHooks) Add(hook Hook) {
	for _, level := range hook.Levels() {
		hooks[level] = append(hooks[level], hook)
	}
}

func hookFailed(entry *Entry, err error) {
	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	fmt.Fprintf(os.Stderr, "Failed to fire hook: %v\n", err)
}

// Fire all the hooks for the passed level. Used by `entry.log` to fire
// appropriate hooks for a log entry.
func (hooks LevelHooks) Fire(level Level, entry *Entry, done chan<- struct{}) {
	var wg sync.WaitGroup
	wg.Add(len(hooks[level]))
	for _, hook := range hooks[level] {
		if _, ok := hook.(AsyncHook); ok {
			go func(h Hook) {
				err := h.Fire(entry)
				if err != nil {
					hookFailed(entry, err)
				}
				wg.Done()
			}(hook)
		} else {
			err := hook.Fire(entry)
			if err != nil {
				hookFailed(entry, err)
			}
			wg.Done()
		}
	}
	go func() {
		wg.Wait()
		done <- struct{}{}
		close(done)
	}()
}
