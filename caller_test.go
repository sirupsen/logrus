package logrus

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryCaller(t *testing.T) {
	logger := New()
	entry := NewEntry(logger)
	logger.Level = DebugLevel

	entry.Debug()
	assert.Equal(t, entry.Data["caller"], mycaller())

	entry.Debugf("")
	assert.Equal(t, entry.Data["caller"], mycaller())

	entry.Debugln()
	assert.Equal(t, entry.Data["caller"], mycaller())

	e1 := entry.WithError(nil)
	e1.Debug()
	assert.Equal(t, e1.Data["caller"], mycaller())

	e1 = entry.WithField("first", 1)
	e1.Debug()
	assert.Equal(t, e1.Data["caller"], mycaller())

	e1 = entry.WithFields(Fields{"first": 1})
	e1.Debug()
	assert.Equal(t, e1.Data["caller"], mycaller())
}

func mycaller() (str string) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		str = "???: ?"
	} else {
		str = fmt.Sprint(filepath.Base(file), ":", line-1)
	}
	return
}
