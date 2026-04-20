package logrus_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type contextKeyType string

func TestEntryWithError(t *testing.T) {
	expErr := fmt.Errorf("kaboom at layer %d", 4711)
	assert.Equal(t, expErr, logrus.WithError(expErr).Data["error"])

	logger := logrus.New()
	logger.SetOutput(io.Discard)
	entry := logrus.NewEntry(logger)

	assert.Equal(t, expErr, entry.WithError(expErr).Data["error"])

	tmpKey := logrus.ErrorKey
	logrus.ErrorKey = "err" //nolint:reassign // ignore "reassigning variable ErrorKey in other package logrus (reassign)"
	t.Cleanup(func() {
		logrus.ErrorKey = tmpKey //nolint:reassign // ignore "reassigning variable ErrorKey in other package logrus (reassign)"
	})

	assert.Equal(t, expErr, entry.WithError(expErr).Data["err"])
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

func TestEntryWithContextCopiesData(t *testing.T) {
	assert := assert.New(t)

	// Initialize a parent Entry object with a key/value set in its Data map
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	parentEntry := logrus.NewEntry(logger).WithField("parentKey", "parentValue")

	// Create two children Entry objects from the parent in different contexts
	var contextKey1 contextKeyType = "foo"
	ctx1 := context.WithValue(context.Background(), contextKey1, "bar")
	childEntry1 := parentEntry.WithContext(ctx1)
	assert.Equal(ctx1, childEntry1.Context)

	var contextKey2 contextKeyType = "bar"
	ctx2 := context.WithValue(context.Background(), contextKey2, "baz")
	childEntry2 := parentEntry.WithContext(ctx2)
	assert.Equal(ctx2, childEntry2.Context)
	assert.NotEqual(ctx1, ctx2)

	// Ensure that data set in the parent Entry are preserved to both children
	assert.Equal("parentValue", childEntry1.Data["parentKey"])
	assert.Equal("parentValue", childEntry2.Data["parentKey"])

	// Modify data stored in the child entry
	childEntry1.Data["childKey"] = "childValue"

	// Verify that data is successfully stored in the child it was set on
	val, exists := childEntry1.Data["childKey"]
	assert.True(exists)
	assert.Equal("childValue", val)

	// Verify that the data change to child 1 has not affected its sibling
	val, exists = childEntry2.Data["childKey"]
	assert.False(exists)
	assert.Empty(val)

	// Verify that the data change to child 1 has not affected its parent
	val, exists = parentEntry.Data["childKey"]
	assert.False(exists)
	assert.Empty(val)
}

func TestEntryWithTimeCopiesData(t *testing.T) {
	assert := assert.New(t)

	// Initialize a parent Entry object with a key/value set in its Data map
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	parentEntry := logrus.NewEntry(logger).WithField("parentKey", "parentValue")

	// Create two children Entry objects from the parent with two different times
	childEntry1 := parentEntry.WithTime(time.Now().AddDate(0, 0, 1))
	childEntry2 := parentEntry.WithTime(time.Now().AddDate(0, 0, 2))

	// Ensure that data set in the parent Entry are preserved to both children
	assert.Equal("parentValue", childEntry1.Data["parentKey"])
	assert.Equal("parentValue", childEntry2.Data["parentKey"])

	// Modify data stored in the child entry
	childEntry1.Data["childKey"] = "childValue"

	// Verify that data is successfully stored in the child it was set on
	val, exists := childEntry1.Data["childKey"]
	assert.True(exists)
	assert.Equal("childValue", val)

	// Verify that the data change to child 1 has not affected its sibling
	val, exists = childEntry2.Data["childKey"]
	assert.False(exists)
	assert.Empty(val)

	// Verify that the data change to child 1 has not affected its parent
	val, exists = parentEntry.Data["childKey"]
	assert.False(exists)
	assert.Empty(val)
}

func TestEntryDupCopiesMetadata(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)

	buffer := bytes.NewBufferString("buffered")
	caller := &runtime.Frame{Function: "example.fn", File: "example.go", Line: 42}
	ctx := context.WithValue(context.Background(), contextKeyType("dup"), "value")

	original := logrus.NewEntry(logger)
	original.Data["animal"] = "walrus"
	original.Time = time.Unix(1700000000, 0)
	original.Level = logrus.WarnLevel
	original.Caller = caller
	original.Message = "duplicated"
	original.Buffer = buffer
	original.Context = ctx

	dup := original.Dup()

	require.NotSame(t, original, dup)
	assert.Equal(t, original.Logger, dup.Logger)
	assert.Equal(t, original.Time, dup.Time)
	assert.Equal(t, original.Level, dup.Level)
	assert.Equal(t, original.Caller, dup.Caller)
	assert.Equal(t, original.Message, dup.Message)
	assert.Equal(t, original.Buffer, dup.Buffer)
	assert.Equal(t, original.Context, dup.Context)
	assert.Equal(t, original.Data, dup.Data)

	dup.Data["animal"] = "otter"
	assert.Equal(t, "walrus", original.Data["animal"])
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

func TestEntryWithIncorrectField(t *testing.T) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(io.Discard)
	entry := logrus.NewEntry(logger)

	fn := func() {}
	eWithFunc := entry.WithFields(logrus.Fields{"func": fn})
	eWithFuncPtr := entry.WithFields(logrus.Fields{"funcPtr": &fn})

	assert.Equal(t, `can not add field "func"`, getErr(t, eWithFunc))
	assert.Equal(t, `can not add field "funcPtr"`, getErr(t, eWithFuncPtr))

	eWithFunc = eWithFunc.WithField("not_a_func", "it is a string")
	eWithFuncPtr = eWithFuncPtr.WithField("not_a_func", "it is a string")

	assert.Equal(t, `can not add field "func"`, getErr(t, eWithFunc))
	assert.Equal(t, `can not add field "funcPtr"`, getErr(t, eWithFuncPtr))

	eWithFunc = eWithFunc.WithTime(time.Now())
	eWithFuncPtr = eWithFuncPtr.WithTime(time.Now())

	assert.Equal(t, `can not add field "func"`, getErr(t, eWithFunc))
	assert.Equal(t, `can not add field "funcPtr"`, getErr(t, eWithFuncPtr))
}

func getErr(t *testing.T, e *logrus.Entry) string {
	t.Helper()

	out, err := e.String()
	require.NoError(t, err)

	var m map[string]any
	require.NoError(t, json.Unmarshal([]byte(out), &m))

	got, _ := m[logrus.FieldKeyLogrusError].(string)
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

type noopHook struct{}

func (noopHook) Levels() []logrus.Level   { return logrus.AllLevels }
func (noopHook) Fire(*logrus.Entry) error { return nil }

type nopBufferPool struct{}

func (nopBufferPool) Get() *bytes.Buffer { return new(bytes.Buffer) }
func (nopBufferPool) Put(*bytes.Buffer)  {}

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
