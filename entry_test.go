package logrus_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	logger.Out = &bytes.Buffer{}
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
	logger.Out = &bytes.Buffer{}
	entry := logrus.NewEntry(logger)

	assert.Equal(ctx, entry.WithContext(ctx).Context)
}

func TestEntryWithContextCopiesData(t *testing.T) {
	assert := assert.New(t)

	// Initialize a parent Entry object with a key/value set in its Data map
	logger := logrus.New()
	logger.Out = &bytes.Buffer{}
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
	logger.Out = &bytes.Buffer{}
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
	logger.Out = &bytes.Buffer{}
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
	logger.Out = &bytes.Buffer{}
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
	logger.Out = &bytes.Buffer{}
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
	logger.Out = &bytes.Buffer{}
	logger.Level = logrus.InfoLevel
	logger.Hooks.Add(&panickyHook{})

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
	logger := logrus.New()
	buffer := &bytes.Buffer{}
	logger.Out = buffer
	logger.SetLevel(logrus.InfoLevel)
	entry := logrus.NewEntry(logger)

	entry.Logf(logrus.DebugLevel, "%s", "debug")
	assert.NotContains(t, buffer.String(), "debug")

	entry.Logf(logrus.WarnLevel, "%s", "warn")
	assert.Contains(t, buffer.String(), "warn")
}

func TestEntryReportCallerRace(t *testing.T) {
	logger := logrus.New()
	entry := logrus.NewEntry(logger)

	// logging before SetReportCaller has the highest chance of causing a race condition
	// to be detected, but doing it twice just to increase the likelihood of detecting the race
	go func() {
		entry.Info("should not race")
	}()
	go func() {
		logger.SetReportCaller(true)
	}()
	go func() {
		entry.Info("should not race")
	}()
}

func TestEntryFormatterRace(t *testing.T) {
	logger := logrus.New()
	entry := logrus.NewEntry(logger)

	// logging before SetReportCaller has the highest chance of causing a race condition
	// to be detected, but doing it twice just to increase the likelihood of detecting the race
	go func() {
		entry.Info("should not race")
	}()
	go func() {
		logger.SetFormatter(&logrus.TextFormatter{})
	}()
	go func() {
		entry.Info("should not race")
	}()
}
