package logrus

type Hook interface {
	Levels() []Level
	Fire(*Entry) error
}

type levelHooks map[Level][]Hook

func (hooks levelHooks) Add(hook Hook) {
	for _, level := range hook.Levels() {
		hooks[level] = append(hooks[level], hook)
	}
}

func (hooks levelHooks) Fire(level Level, entry *Entry) error {
	for _, hook := range hooks[level] {
		if err := hook.Fire(entry); err != nil {
			return err
		}
	}

	return nil
}
