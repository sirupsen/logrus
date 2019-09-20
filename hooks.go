package logrus

import (
	"fmt"
	"strings"
)

// A hook to be fired when logging on the logging levels returned from
// `Levels()` on your implementation of the interface. Note that this is not
// fired in a goroutine or a channel with workers, you should handle such
// functionality yourself if your call is non-blocking and you don't wish for
// the logging calls for levels returned from `Levels()` to block.
type Hook interface {
	Levels() []Level
	Fire(*Entry) error
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

// multiErr may contain multiple errors.
type multiErr []error

func (e multiErr) Error() string {
	if len(e) == 1 {
		return e[0].Error()
	}
	str := make([]string, len(e))
	for i, err := range e.Errors() {
		str[i] = fmt.Sprintf("Error #%d: %s", i+1, err)
	}
	return strings.Join(str, "\n")
}

func (e multiErr) Unwrap() error {
	if e == nil || len(e) == 0 {
		return nil
	}
	return e[0]
}

func (e multiErr) Errors() []error {
	return e
}

// Fire all the hooks for the passed level. Used by `entry.log` to fire
// appropriate hooks for a log entry. By default, when an error occurs, further
// firing of hooks is aborted. To fire all hooks, even in case of an error,
// set the Logger.FireAllHooks flag to true, in which case a returned error
// will be a composite of multiple errors which may be inspected with a type
// assertion. For example:
//
//  type errorGroup interface {
//      Errors() []error
//  }
//  for _, e := range err.(errorGroup).Errors() {
//      // inspect individual errors here
//   }
func (hooks LevelHooks) Fire(level Level, entry *Entry) error {
	var merr multiErr
	for _, hook := range hooks[level] {
		if err := hook.Fire(entry); err != nil {
			if !entry.Logger.FireAllHooks {
				return err
			}
			merr = append(merr, err)
		}
	}
	return merr
}
