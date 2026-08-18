package logrus_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type noopHook struct{}

func (noopHook) Levels() []logrus.Level   { return logrus.AllLevels }
func (noopHook) Fire(*logrus.Entry) error { return nil }

type nopBufferPool struct{}

func (nopBufferPool) Get() *bytes.Buffer { return new(bytes.Buffer) }
func (nopBufferPool) Put(*bytes.Buffer)  {}

type nopFormatter struct{}

func (nopFormatter) Format(*logrus.Entry) ([]byte, error) { return nil, nil }

type funcError func()

func (funcError) Error() string { return "boom" }

type contextKeyType string

func TestEntryWithError(t *testing.T) {
	expErr := fmt.Errorf("kaboom at layer %d", 4711)
	logger, hook := test.NewNullLogger()

	logger.WithError(expErr).Error("default key")

	tmpKey := logrus.ErrorKey
	logrus.ErrorKey = "err" //nolint:reassign // ignore "reassigning variable ErrorKey in other package logrus (reassign)"
	t.Cleanup(func() {
		logrus.ErrorKey = tmpKey //nolint:reassign // ignore "reassigning variable ErrorKey in other package logrus (reassign)"
	})

	logger.WithError(expErr).Error("custom key")

	require.Len(t, hook.Entries, 2)
	assert.Equal(t, expErr, hook.Entries[0].Data["error"])
	assert.Equal(t, expErr, hook.Entries[1].Data["err"])
}

// TestEntryErrorFieldValues verifies that error fields are accepted, rejected,
// and formatted consistently across WithError, WithField, and WithFields.
//
// relates to https://github.com/sirupsen/logrus/pull/1577#discussion_r3799339160
func TestEntryErrorFieldValues(t *testing.T) {
	errFunc := funcError(func() {})
	fn := func() {}

	values := []struct {
		name  string
		value any
		valid bool
	}{
		{name: "error", value: errors.New("boom"), valid: true},
		{name: "function error", value: errFunc, valid: true},
		{name: "function error pointer", value: &errFunc, valid: true},
		{name: "string", value: "boom", valid: true},
		{name: "function", value: fn, valid: false},
		{name: "function pointer", value: &fn, valid: false},
	}

	tests := []struct {
		doc   string
		apply func(*logrus.Logger, any) *logrus.Entry
	}{
		{
			doc: "WithError",
			apply: func(logger *logrus.Logger, value any) *logrus.Entry {
				return logger.WithError(value.(error))
			},
		},
		{
			doc: "WithField",
			apply: func(logger *logrus.Logger, value any) *logrus.Entry {
				return logger.WithField(logrus.ErrorKey, value)
			},
		},
		{
			doc: "WithFields",
			apply: func(logger *logrus.Logger, value any) *logrus.Entry {
				return logger.WithFields(logrus.Fields{logrus.ErrorKey: value})
			},
		},
	}

	for _, tc := range tests {
		for _, val := range values {
			t.Run(tc.doc+"("+val.name+")", func(t *testing.T) {
				if _, ok := val.value.(error); tc.doc == "WithError" && !ok {
					skip(t, "WithError only accepts error values")
					return
				}

				var buf bytes.Buffer

				logger := logrus.New()
				logger.SetOutput(&buf)
				logger.SetFormatter(&logrus.JSONFormatter{})

				tc.apply(logger, val.value).Error("test")

				var data map[string]any
				require.NoError(t, json.Unmarshal(buf.Bytes(), &data))

				if val.valid {
					assert.Equal(t, "boom", data[logrus.ErrorKey])
					assert.NotContains(t, data, logrus.FieldKeyLogrusError)
				} else {
					assert.NotContains(t, data, logrus.ErrorKey)
					assert.Contains(t, data, logrus.FieldKeyLogrusError)
				}
			})
		}
	}
}

func TestEntryWithContext(t *testing.T) {
	assert := assert.New(t)
	var contextKey contextKeyType = "foo"
	ctx := context.WithValue(context.Background(), contextKey, "bar")

	assert.Equal(ctx, logrus.WithContext(ctx).Context)

	logger := logrus.New()
	logger.SetOutput(io.Discard)
	entry := logrus.NewEntry(logger)

	assert.Equal(ctx, entry.WithContext(ctx).Context)
}

func TestEntryWithContextPreservesData(t *testing.T) {
	logger, hook := test.NewNullLogger()

	// Initialize a parent Entry object with a key/value field.
	parentEntry := logger.WithField("parentKey", "parentValue")

	// Create two child Entry objects from the parent in different contexts.
	var contextKey1 contextKeyType = "foo"
	ctx1 := context.WithValue(context.Background(), contextKey1, "bar")
	childEntry1 := parentEntry.WithContext(ctx1)
	assert.Equal(t, ctx1, childEntry1.Context)

	var contextKey2 contextKeyType = "bar"
	ctx2 := context.WithValue(context.Background(), contextKey2, "baz")
	childEntry2 := parentEntry.WithContext(ctx2)
	assert.Equal(t, ctx2, childEntry2.Context)
	assert.NotEqual(t, ctx1, ctx2)

	childEntry1.Info("child 1")
	childEntry2.Info("child 2")
	parentEntry.Info("parent")

	require.Len(t, hook.Entries, 3)

	// Ensure that data set in the parent Entry are preserved in both children
	// and when the parent is reused.
	assert.Equal(t, "parentValue", hook.Entries[0].Data["parentKey"])
	assert.Equal(t, "parentValue", hook.Entries[1].Data["parentKey"])
	assert.Equal(t, "parentValue", hook.Entries[2].Data["parentKey"])
}

func TestEntryWithTimePreservesData(t *testing.T) {
	logger, hook := test.NewNullLogger()

	// Initialize a parent Entry object with a key/value field.
	parentEntry := logger.WithField("parentKey", "parentValue")

	// Create two child Entry objects from the parent with two different times.
	childEntry1 := parentEntry.WithTime(time.Now().AddDate(0, 0, 1))
	childEntry2 := parentEntry.WithTime(time.Now().AddDate(0, 0, 2))

	childEntry1.Info("child 1")
	childEntry2.Info("child 2")
	parentEntry.Info("parent")

	require.Len(t, hook.Entries, 3)

	// Ensure that data set in the parent Entry are preserved in both children
	// and when the parent is reused.
	assert.Equal(t, "parentValue", hook.Entries[0].Data["parentKey"])
	assert.Equal(t, "parentValue", hook.Entries[1].Data["parentKey"])
	assert.Equal(t, "parentValue", hook.Entries[2].Data["parentKey"])
}

// TestEntryLogPanicLevelDoesNotPanic verifies that generic Log methods treat
// PanicLevel as severity only and do not trigger panic behavior.
//
// Regression and related history:
//   - https://github.com/sirupsen/logrus/pull/65
//   - https://github.com/sirupsen/logrus/pull/1283
func TestEntryLogPanicLevelDoesNotPanic(t *testing.T) {
	tests := []struct {
		doc  string
		log  func(*logrus.Entry)
		want string
	}{
		{
			doc:  "Log",
			log:  func(e *logrus.Entry) { e.Log(logrus.PanicLevel, "kaboom") },
			want: "kaboom",
		},
		{
			doc:  "Logf",
			log:  func(e *logrus.Entry) { e.Logf(logrus.PanicLevel, "kaboom %d", 42) },
			want: "kaboom 42",
		},
		{
			doc:  "Logln",
			log:  func(e *logrus.Entry) { e.Logln(logrus.PanicLevel, "kaboom", 42) },
			want: "kaboom 42",
		},
	}

	for _, tc := range tests {
		t.Run(tc.doc, func(t *testing.T) {
			logger, hook := test.NewNullLogger()

			require.NotPanics(t, func() {
				tc.log(logrus.NewEntry(logger))
			})

			got := hook.LastEntry()
			require.NotNil(t, got)
			assert.Equal(t, logrus.PanicLevel, got.Level)
			assert.Equal(t, tc.want, got.Message)
		})
	}
}

func TestEntryPanicln(t *testing.T) {
	errBoom := fmt.Errorf("boom time")

	defer func() {
		p := recover()
		assert.NotNil(t, p)

		switch pVal := p.(type) {
		case *logrus.Entry:
			assert.Equal(t, "kaboom", pVal.Message)
			assert.Equal(t, errBoom, pVal.Data["err"])
		default:
			t.Fatalf("want type *Entry, got %T: %#v", pVal, pVal)
		}
	}()

	logger := logrus.New()
	logger.SetOutput(io.Discard)
	entry := logrus.NewEntry(logger)
	entry.WithField("err", errBoom).Panicln("kaboom")
}

func TestEntryPanicf(t *testing.T) {
	errBoom := fmt.Errorf("boom again")

	defer func() {
		p := recover()
		assert.NotNil(t, p)

		switch pVal := p.(type) {
		case *logrus.Entry:
			assert.Equal(t, "kaboom true", pVal.Message)
			assert.Equal(t, errBoom, pVal.Data["err"])
		default:
			t.Fatalf("want type *Entry, got %T: %#v", pVal, pVal)
		}
	}()

	logger := logrus.New()
	logger.SetOutput(io.Discard)
	entry := logrus.NewEntry(logger)
	entry.WithField("err", errBoom).Panicf("kaboom %v", true)
}

func TestEntryPanic(t *testing.T) {
	errBoom := fmt.Errorf("boom again")

	defer func() {
		p := recover()
		assert.NotNil(t, p)

		switch pVal := p.(type) {
		case *logrus.Entry:
			assert.Equal(t, "kaboom", pVal.Message)
			assert.Equal(t, errBoom, pVal.Data["err"])
		default:
			t.Fatalf("want type *Entry, got %T: %#v", pVal, pVal)
		}
	}()

	logger := logrus.New()
	logger.SetOutput(io.Discard)
	entry := logrus.NewEntry(logger)
	entry.WithField("err", errBoom).Panic("kaboom")
}

const (
	badMessage   = "this is going to panic"
	panicMessage = "this is broken"
)

type panickyHook struct{}

func (p *panickyHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.InfoLevel}
}

func (p *panickyHook) Fire(entry *logrus.Entry) error {
	if entry.Message == badMessage {
		panic(panicMessage)
	}

	return nil
}

func TestEntryHooksPanic(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	logger.SetLevel(logrus.InfoLevel)
	logger.AddHook(&panickyHook{})

	defer func() {
		p := recover()
		assert.NotNil(t, p)
		assert.Equal(t, panicMessage, p)

		entry := logrus.NewEntry(logger)
		entry.Info("another message")
	}()

	entry := logrus.NewEntry(logger)
	entry.Info(badMessage)
}

// TestEntryDerivationPreservesCaller verifies that derived entries retain
// explicitly set caller information.
func TestEntryDerivationPreservesCaller(t *testing.T) {
	entry := logrus.NewEntry(logrus.New())
	entry.Caller = &runtime.Frame{
		Function: "example.function",
		File:     "example.go",
		Line:     42,
	}

	tests := []struct {
		doc    string
		derive func(*logrus.Entry) *logrus.Entry
	}{
		{
			doc:    "Dup",
			derive: func(entry *logrus.Entry) *logrus.Entry { return entry.Dup() },
		},
		{
			doc:    "WithContext",
			derive: func(entry *logrus.Entry) *logrus.Entry { return entry.WithContext(context.Background()) },
		},
		{
			doc:    "WithError",
			derive: func(entry *logrus.Entry) *logrus.Entry { return entry.WithError(errors.New("boom")) },
		},
		{
			doc:    "WithField",
			derive: func(entry *logrus.Entry) *logrus.Entry { return entry.WithField("foo", "bar") },
		},
		{
			doc:    "WithFields",
			derive: func(entry *logrus.Entry) *logrus.Entry { return entry.WithFields(logrus.Fields{"foo": "bar"}) },
		},
		{
			doc:    "WithTime",
			derive: func(entry *logrus.Entry) *logrus.Entry { return entry.WithTime(time.Now()) },
		},
	}

	for _, tc := range tests {
		t.Run(tc.doc, func(t *testing.T) {
			got := tc.derive(entry)

			require.NotNil(t, got.Caller)
			assert.Equal(t, "example.function", got.Caller.Function)
			assert.Equal(t, "example.go", got.Caller.File)
			assert.Equal(t, 42, got.Caller.Line)
		})
	}
}

func TestEntryWithIncorrectField(t *testing.T) {
	var buf bytes.Buffer

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(&buf)

	fn := func() {}
	eWithFunc := logger.WithFields(logrus.Fields{"func": fn})
	eWithFuncPtr := logger.WithFields(logrus.Fields{"funcPtr": &fn})

	assert.Equal(t, `skipping unsupported field "func"`, getErr(t, &buf, eWithFunc))
	assert.Equal(t, `skipping unsupported field "funcPtr"`, getErr(t, &buf, eWithFuncPtr))

	eWithFunc = eWithFunc.WithField("not_a_func", "it is a string")
	eWithFuncPtr = eWithFuncPtr.WithField("not_a_func", "it is a string")

	assert.Equal(t, `skipping unsupported field "func"`, getErr(t, &buf, eWithFunc))
	assert.Equal(t, `skipping unsupported field "funcPtr"`, getErr(t, &buf, eWithFuncPtr))

	eWithFunc = eWithFunc.WithTime(time.Now())
	eWithFuncPtr = eWithFuncPtr.WithTime(time.Now())

	assert.Equal(t, `skipping unsupported field "func"`, getErr(t, &buf, eWithFunc))
	assert.Equal(t, `skipping unsupported field "funcPtr"`, getErr(t, &buf, eWithFuncPtr))
}

func getErr(t *testing.T, buf *bytes.Buffer, entry *logrus.Entry) string {
	t.Helper()

	buf.Reset()
	entry.Info("test")

	var data map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &data))

	got, _ := data[logrus.FieldKeyLogrusError].(string)
	return got
}

func TestEntryLogfLevel(t *testing.T) {
	var buffer bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buffer)
	logger.SetLevel(logrus.InfoLevel)
	entry := logrus.NewEntry(logger)

	entry.Logf(logrus.DebugLevel, "%s", "debug")
	assert.NotContains(t, buffer.String(), "debug")

	entry.Logf(logrus.WarnLevel, "%s", "warn")
	assert.Contains(t, buffer.String(), "warn")
}

func TestEntryLoggerMutationRace(t *testing.T) {
	tests := []struct {
		doc    string
		mutate func(*logrus.Logger)
	}{
		{doc: "AddHook", mutate: func(l *logrus.Logger) { l.AddHook(noopHook{}) }},
		{doc: "SetBufferPool", mutate: func(l *logrus.Logger) { l.SetBufferPool(nopBufferPool{}) }},
		{doc: "SetFormatter", mutate: func(l *logrus.Logger) { l.SetFormatter(&logrus.TextFormatter{}) }},
		{doc: "SetLevel", mutate: func(l *logrus.Logger) { l.SetLevel(logrus.InfoLevel) }},
		{doc: "SetOutput", mutate: func(l *logrus.Logger) { l.SetOutput(io.Discard) }},
		{doc: "SetReportCaller", mutate: func(l *logrus.Logger) { l.SetReportCaller(true) }},
		{doc: "ReplaceHooks_withHookPresent", mutate: func(l *logrus.Logger) {
			// Replace with a fresh map each time to maximize mutation.
			h := make(logrus.LevelHooks)
			for _, lvl := range logrus.AllLevels {
				h[lvl] = []logrus.Hook{noopHook{}}
			}
			l.ReplaceHooks(h)
		}},
	}

	for _, tc := range tests {
		t.Run(tc.doc, func(t *testing.T) {
			runEntryLoggerRace(t, tc.mutate)
		})
	}
}

func runEntryLoggerRace(t *testing.T, mutate func(logger *logrus.Logger)) {
	t.Helper()

	logger := logrus.New()
	logger.SetOutput(io.Discard)
	entry := logrus.NewEntry(logger)

	const n = 100

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		for range n {
			_, _ = entry.Bytes()
		}
	}()

	go func() {
		defer wg.Done()
		for range n {
			entry.Info("should not race")
		}
	}()

	go func() {
		defer wg.Done()
		for range n {
			mutate(logger)
		}
	}()

	go func() {
		defer wg.Done()
		for range n {
			entry.Info("should not race")
		}
	}()

	wg.Wait()
}

// reentrantValue is a type whose MarshalJSON method triggers another log call,
// which would deadlock if the logger mutex is held during formatting.
type reentrantValue struct {
	logger *logrus.Logger
}

func (r reentrantValue) MarshalJSON() ([]byte, error) {
	r.logger.Info("reentrant log from MarshalJSON")
	return []byte(`"reentrant"`), nil
}

// TestEntryReentrantLoggingDeadlock verifies that logging from within a field's
// MarshalJSON (or similar serialization callback) does not deadlock.
// This is a regression test for https://github.com/sirupsen/logrus/issues/1448.
func TestEntryReentrantLoggingDeadlock(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.JSONFormatter{})

	done := make(chan struct{})
	go func() {
		defer close(done)
		logger.WithFields(logrus.Fields{
			"key": reentrantValue{logger: logger},
		}).Info("outer log message")
	}()

	select {
	case <-done:
		// Success: the log call completed without deadlocking.
		output := buf.String()
		assert.Contains(t, output, "outer log message")
		assert.Contains(t, output, "reentrant log from MarshalJSON")
		assert.Contains(t, output, `"key":"reentrant"`)
	case <-time.After(5 * time.Second):
		t.Fatal("deadlock detected: reentrant logging from MarshalJSON blocked for 5 seconds")
	}
}

// TestEntryWithFieldsThenBranch verifies that an Entry with fields can be
// logged, used as the parent of a derived Entry, and reused without inheriting
// fields added to the derived Entry.
func TestEntryWithFieldsThenBranch(t *testing.T) {
	logger, hook := test.NewNullLogger()

	entry := logger.
		WithField("a", 1).
		WithField("b", 2)

	entry.Info("parent")

	require.Len(t, hook.Entries, 1)
	assert.Equal(t, logrus.Fields{
		"a": 1,
		"b": 2,
	}, hook.Entries[0].Data)

	child := entry.WithField("c", 3)
	child.Info("child")

	require.Len(t, hook.Entries, 2)
	assert.Equal(t, logrus.Fields{
		"a": 1,
		"b": 2,
		"c": 3,
	}, hook.Entries[1].Data)

	entry.Info("parent again")

	require.Len(t, hook.Entries, 3)
	assert.Equal(t, logrus.Fields{
		"a": 1,
		"b": 2,
	}, hook.Entries[2].Data)
}

// TestEntryDataIsMutable verifies that Entry.Data exposes materialized fields,
// can be modified directly, and remains isolated between derived entries.
func TestEntryDataIsMutable(t *testing.T) {
	logger, hook := test.NewNullLogger()

	// Fields added to an entry are exposed immediately through Data.
	parent := logger.WithFields(logrus.Fields{
		"foo": "bar",
		"one": 1,
	})

	assert.Equal(t, "bar", parent.Data["foo"])
	assert.Equal(t, 1, parent.Data["one"])

	// Derived entries inherit the parent's fields without sharing its Data map.
	child := parent.WithFields(logrus.Fields{
		"two": 2,
	})

	assert.Equal(t, "bar", child.Data["foo"])
	assert.Equal(t, 1, child.Data["one"])
	assert.Equal(t, 2, child.Data["two"])

	// Data is public and mutable. Changes to the child must not affect its
	// parent, and newly added fields must be preserved.
	child.Data["foo"] = "changed"
	child.Data["three"] = 3

	assert.Equal(t, "bar", parent.Data["foo"])
	assert.Equal(t, 1, parent.Data["one"])
	assert.NotContains(t, parent.Data, "three")

	// Logging must use the currently exposed Data and must not overwrite direct
	// mutations by re-materializing fields from internal state.
	parent.Info("parent")
	child.Info("child")

	require.Len(t, hook.Entries, 2)

	assert.Equal(t, logrus.Fields{
		"foo": "bar",
		"one": 1,
	}, hook.Entries[0].Data)

	assert.Equal(t, logrus.Fields{
		"foo":   "changed",
		"one":   1,
		"two":   2,
		"three": 3,
	}, hook.Entries[1].Data)
}
