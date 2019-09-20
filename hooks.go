package logrus

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
	if e == nil || len(e) == 0 {
		return ""
	}
	return e[0].Error()
}

func (e multiErr) Unwrap() error {
	if e == nil || len(e) == 0 {
		return nil
	}
	return e[0]
}

// Fire all the hooks for the passed level. Used by `entry.log` to fire
// appropriate hooks for a log entry. In case of an error, the first error
// encountered will be returned, but all hooks will fire.
func (hooks LevelHooks) Fire(level Level, entry *Entry) error {
	merr := make(multiErr, 0)
	for _, hook := range hooks[level] {
		if err := hook.Fire(entry); err != nil {
			if !entry.Logger.FireAllHooks {
				return err
			}
			merr = append(merr, err)
		}
	}
	if len(merr) > 0 {
		return merr
	}
	return nil
}
