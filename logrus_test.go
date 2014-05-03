package logrus

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func LogAndAssertJSON(t *testing.T, log func(*Logger), assertions func(fields Fields)) {
	var buffer bytes.Buffer
	var fields Fields

	logger := New()
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)

	log(logger)

	err := json.Unmarshal(buffer.Bytes(), &fields)
	assert.Nil(t, err)

	assertions(fields)
}

func TestPrint(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Print("test")
	}, func(fields Fields) {
		assert.Equal(t, fields["msg"], "test")
		assert.Equal(t, fields["level"], "info")
	})
}

func TestInfo(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Info("test")
	}, func(fields Fields) {
		assert.Equal(t, fields["msg"], "test")
		assert.Equal(t, fields["level"], "info")
	})
}

func TestWarn(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Warn("test")
	}, func(fields Fields) {
		assert.Equal(t, fields["msg"], "test")
		assert.Equal(t, fields["level"], "warning")
	})
}

func TestInfolnShouldAddSpacesBetweenStrings(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln("test", "test")
	}, func(fields Fields) {
		assert.Equal(t, fields["msg"], "test test")
	})
}

func TestInfolnShouldAddSpacesBetweenStringAndNonstring(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln("test", 10)
	}, func(fields Fields) {
		assert.Equal(t, fields["msg"], "test 10")
	})
}

func TestInfolnShouldAddSpacesBetweenTwoNonStrings(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln(10, 10)
	}, func(fields Fields) {
		assert.Equal(t, fields["msg"], "10 10")
	})
}

func TestInfoShouldAddSpacesBetweenTwoNonStrings(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln(10, 10)
	}, func(fields Fields) {
		assert.Equal(t, fields["msg"], "10 10")
	})
}

func TestInfoShouldNotAddSpacesBetweenStringAndNonstring(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Info("test", 10)
	}, func(fields Fields) {
		assert.Equal(t, fields["msg"], "test10")
	})
}

func TestInfoShouldNotAddSpacesBetweenStrings(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Info("test", "test")
	}, func(fields Fields) {
		assert.Equal(t, fields["msg"], "testtest")
	})
}

type SlowString string

func (s SlowString) String() string {
	time.Sleep(time.Millisecond)
	return string(s)
}

func getLogAtLevel(l Level) *Logger {
	log := New()
	log.Level = l
	log.Out = ioutil.Discard
	return log
}

func BenchmarkLevelDisplayed(b *testing.B) {
	log := getLogAtLevel(Info)
	for i := 0; i < b.N; i++ {
		log.Info(SlowString("foo"))
	}
}

func BenchmarkLevelHidden(b *testing.B) {
	log := getLogAtLevel(Info)
	for i := 0; i < b.N; i++ {
		log.Debug(SlowString("foo"))
	}
}

func BenchmarkLevelfDisplayed(b *testing.B) {
	log := getLogAtLevel(Info)
	for i := 0; i < b.N; i++ {
		log.Infof("%s", SlowString("foo"))
	}
}

func BenchmarkLevelfHidden(b *testing.B) {
	log := getLogAtLevel(Info)
	for i := 0; i < b.N; i++ {
		log.Debugf("%s", SlowString("foo"))
	}
}

func BenchmarkLevellnDisplayed(b *testing.B) {
	log := getLogAtLevel(Info)
	for i := 0; i < b.N; i++ {
		log.Infoln(SlowString("foo"))
	}
}

func BenchmarkLevellnHidden(b *testing.B) {
	log := getLogAtLevel(Info)
	for i := 0; i < b.N; i++ {
		log.Debugln(SlowString("foo"))
	}
}
