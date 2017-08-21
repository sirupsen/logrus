package logrus

import (
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
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

type CallerHook struct {
	StackDepth int
}

func (hook *CallerHook) Fire(entry *Entry) error {
	entry.Data["caller"] = hook.caller(entry)
	return nil
}

func (hook *CallerHook) Levels() []Level {
	return []Level{
		PanicLevel,
		FatalLevel,
		ErrorLevel,
		WarnLevel,
		InfoLevel,
		DebugLevel,
	}
}

func (hook *CallerHook) caller(entry *Entry) string {
	if _, file, line, ok := runtime.Caller(hook.StackDepth); ok {
		return strings.Join([]string{filepath.Base(file), strconv.Itoa(line)}, ":")
	} else {
		return "this is not exist"
	}
}

func TestStackDepth(t *testing.T) {
	assert.Equal(t, true, strings.Contains("/usr/local/go/src/reflect/value.go", "reflect/value.go"))
	hook := &CallerHook{StackDepth: 8}
	//Be aware expectCaller's line number should equal Invoke(log, logList[i].Log, logList[i].Input...)'s line number
	expectCaller := "hook_test.go:218"
	//How to test panic and fatal?
	logList := []struct {
		Log          string
		ExpectCaller string
		Input        []interface{}
		ExpectMsg    string
		ExpectLevel  string
	}{
		{Log: "Debug", ExpectCaller: expectCaller, Input: []interface{}{"Hello1", "Debug"}, ExpectMsg: "Hello1Debug", ExpectLevel: DebugLevel.String()},
		{Log: "Debugf", ExpectCaller: expectCaller, Input: []interface{}{"Hello2 %s", "Debug"}, ExpectMsg: "Hello2 Debug", ExpectLevel: DebugLevel.String()},
		{Log: "Debugln", ExpectCaller: expectCaller, Input: []interface{}{"Hello3", "Debug"}, ExpectMsg: "Hello3 Debug", ExpectLevel: DebugLevel.String()},
		{Log: "Error", ExpectCaller: expectCaller, Input: []interface{}{"Hello1", "Error"}, ExpectMsg: "Hello1Error", ExpectLevel: ErrorLevel.String()},
		{Log: "Errorf", ExpectCaller: expectCaller, Input: []interface{}{"Hello2 %s", "Error"}, ExpectMsg: "Hello2 Error", ExpectLevel: ErrorLevel.String()},
		{Log: "Errorln", ExpectCaller: expectCaller, Input: []interface{}{"Hello3", "Error"}, ExpectMsg: "Hello3 Error", ExpectLevel: ErrorLevel.String()},
		{Log: "Info", ExpectCaller: expectCaller, Input: []interface{}{"Hello1", "Info"}, ExpectMsg: "Hello1Info", ExpectLevel: InfoLevel.String()},
		{Log: "Infof", ExpectCaller: expectCaller, Input: []interface{}{"Hello2 %s", "Info"}, ExpectMsg: "Hello2 Info", ExpectLevel: InfoLevel.String()},
		{Log: "Infoln", ExpectCaller: expectCaller, Input: []interface{}{"Hello3", "Info"}, ExpectMsg: "Hello3 Info", ExpectLevel: InfoLevel.String()},
		{Log: "Warning", ExpectCaller: expectCaller, Input: []interface{}{"Hello1", "Warning"}, ExpectMsg: "Hello1Warning", ExpectLevel: WarnLevel.String()},
		{Log: "Warningf", ExpectCaller: expectCaller, Input: []interface{}{"Hello2 %s", "Warning"}, ExpectMsg: "Hello2 Warning", ExpectLevel: WarnLevel.String()},
		{Log: "Warningln", ExpectCaller: expectCaller, Input: []interface{}{"Hello3", "Warning"}, ExpectMsg: "Hello3 Warning", ExpectLevel: WarnLevel.String()},
		{Log: "Warn", ExpectCaller: expectCaller, Input: []interface{}{"Hello1", "Warn"}, ExpectMsg: "Hello1Warn", ExpectLevel: WarnLevel.String()},
		{Log: "Warnf", ExpectCaller: expectCaller, Input: []interface{}{"Hello2 %s", "Warn"}, ExpectMsg: "Hello2 Warn", ExpectLevel: WarnLevel.String()},
		{Log: "Warnln", ExpectCaller: expectCaller, Input: []interface{}{"Hello3", "Warn"}, ExpectMsg: "Hello3 Warn", ExpectLevel: WarnLevel.String()},
		{Log: "Print", ExpectCaller: expectCaller, Input: []interface{}{"Hello1", "Print"}, ExpectMsg: "Hello1Print", ExpectLevel: InfoLevel.String()},
		{Log: "Printf", ExpectCaller: expectCaller, Input: []interface{}{"Hello2 %s", "Print"}, ExpectMsg: "Hello2 Print", ExpectLevel: InfoLevel.String()},
		{Log: "Println", ExpectCaller: expectCaller, Input: []interface{}{"Hello3", "Print"}, ExpectMsg: "Hello3 Print", ExpectLevel: InfoLevel.String()},
	}
	for i := range logList {
		LogAndAssertJSON(t, func(log *Logger) {
			log.Level = DebugLevel
			log.Hooks.Add(hook)
			Invoke(log, logList[i].Log, logList[i].Input...)
		}, func(fields Fields) {
			assert.Equal(t, logList[i].ExpectCaller, fields["caller"])
			assert.Equal(t, logList[i].ExpectMsg, fields["msg"])
			assert.Equal(t, logList[i].ExpectLevel, fields["level"])
		})
	}
	for i := range logList {
		LogAndAssertJSON(t, func(log *Logger) {
			log.Level = DebugLevel
			log.Hooks.Add(hook)
			Invoke(log.WithField("type", "check"), logList[i].Log, logList[i].Input...)
		}, func(fields Fields) {
			assert.Equal(t, logList[i].ExpectCaller, fields["caller"])
			assert.Equal(t, logList[i].ExpectMsg, fields["msg"])
			assert.Equal(t, logList[i].ExpectLevel, fields["level"])
		})
	}
}

func Invoke(any interface{}, name string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(any).MethodByName(name).Call(inputs)
}
