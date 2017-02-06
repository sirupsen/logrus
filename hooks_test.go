package logrus

import (
	"testing"
)

func TestLevelHooks_Add(t *testing.T) {
	const hooksAddLength = 1000
	hookList := make([]Hook, 0, hooksAddLength)
	for i := 0; i < hooksAddLength; i++ {
		hookList = append(hookList, new(TestHook))
	}

	hooks := NewLevelHooks()
	for _, h := range hookList {
		go func(hooks *LevelHooks, h Hook) { hooks.Add(h) }(hooks, h) // raise data race.
	}

	hooks2 := NewLevelHooks()
	for _, h := range hookList {
		hooks2.Add(h) // no data race.
	}
}

func TestLevelHooks_Fire(t *testing.T) {
	const hooksFireLength = 1000
	entryList := make([]Entry, 0, hooksFireLength)
	for i := 0; i < hooksFireLength; i++ {
		entryList = append(entryList, Entry{})
	}

	hooks := NewLevelHooks()
	for _, e := range entryList {
		go func(hooks *LevelHooks, e Entry) { hooks.Fire(ErrorLevel, &e) }(hooks, e)
	}
}
