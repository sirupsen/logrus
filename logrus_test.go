package logrus

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func LogAndAssertJSON(t *testing.T, log func(*Logger), assertions func(entry *Entry)) {
	var buffer bytes.Buffer

	logger := New()
	logger.Out = &buffer
	formatter := new(JSONFormatter)
	logger.Formatter = formatter

	log(logger)

	entry, err := formatter.Unformat(buffer.Bytes())
	assert.Nil(t, err)

	if assert.NotNil(t, entry) {
		assertions(entry)
	}
}

func TestPrint(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Print("test")
	}, func(entry *Entry) {
		assert.Equal(t, entry.Msg, "test")
		assert.Equal(t, entry.Level, Info)
	})
}

func TestInfo(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Info("test")
	}, func(entry *Entry) {
		assert.Equal(t, entry.Msg, "test")
		assert.Equal(t, entry.Level, Info)
	})
}

func TestWarn(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Warn("test")
	}, func(entry *Entry) {
		assert.Equal(t, entry.Msg, "test")
		assert.Equal(t, entry.Level, Warn)
	})
}

func TestInfolnShouldAddSpacesBetweenStrings(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln("test", "test")
	}, func(entry *Entry) {
		assert.Equal(t, entry.Msg, "test test")
	})
}

func TestInfolnShouldAddSpacesBetweenStringAndNonstring(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln("test", 10)
	}, func(entry *Entry) {
		assert.Equal(t, entry.Msg, "test 10")
	})
}

func TestInfolnShouldAddSpacesBetweenTwoNonStrings(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln(10, 10)
	}, func(entry *Entry) {
		assert.Equal(t, entry.Msg, "10 10")
	})
}

func TestInfoShouldAddSpacesBetweenTwoNonStrings(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln(10, 10)
	}, func(entry *Entry) {
		assert.Equal(t, entry.Msg, "10 10")
	})
}

func TestInfoShouldNotAddSpacesBetweenStringAndNonstring(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Info("test", 10)
	}, func(entry *Entry) {
		assert.Equal(t, entry.Msg, "test10")
	})
}

func TestInfoShouldNotAddSpacesBetweenStrings(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Info("test", "test")
	}, func(entry *Entry) {
		assert.Equal(t, entry.Msg, "testtest")
	})
}

func TestWithFieldsShouldAllowAssignments(t *testing.T) {
	var buffer bytes.Buffer
	var fields Fields

	logger := New()
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)

	localLog := logger.WithFields(Fields{
		"key1": "value1",
	})

	localLog.WithField("key2", "value2").Info("test")
	err := json.Unmarshal(buffer.Bytes(), &fields)
	assert.Nil(t, err)

	assert.Equal(t, "value2", fields["key2"])
	assert.Equal(t, "value1", fields["key1"])

	buffer = bytes.Buffer{}
	fields = Fields{}
	localLog.Info("test")
	err = json.Unmarshal(buffer.Bytes(), &fields)
	assert.Nil(t, err)

	_, ok := fields["key2"]
	assert.Equal(t, false, ok)
	assert.Equal(t, "value1", fields["key1"])
}
