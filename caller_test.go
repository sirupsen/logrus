package logrus

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Output struct {
	caller string `json:"caller"`
	msg    string `json:"msg"`
}

func TestEntryTEXTFormatterCaller(t *testing.T) {
	fileName := "caller.txt"
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	assert.NoError(t, err)

	logger := New()
	entry := NewEntry(logger)
	logger.showCaller = true
	logger.Level = DebugLevel
	logger.Formatter = &TextFormatter{DisableColors: true}
	logger.Out = f

	entry.Debug(caller(1))
	entry.Debugf(caller(1))
	entry.Debugln(caller(1))
	e1 := entry.WithError(nil)
	e1.Debug(caller(1))
	e1 = entry.WithField("first", 1)
	e1.Debug(caller(1))
	e1 = entry.WithFields(Fields{"first": 1})
	e1.Debug(caller(1))

	f.Close()

	f, err = os.Open(fileName)
	defer f.Close()
	scaner := bufio.NewScanner(f)
	for scaner.Scan() {
		ss := scaner.Text()
		fields := strings.Fields(ss)
		msg := ""
		call := ""
		for _, field := range fields {
			if len(field) >= 7 && field[:7] == "caller=" {
				msg = field[7:]
			} else if len(field) >= 4 && field[:4] == "msg=" {
				call = field[4:]
			}
		}
		assert.Equal(t, msg, call)
	}

}

func TestEntryJSONFormatterCaller(t *testing.T) {
	fileName := "caller.json"
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	assert.NoError(t, err)

	logger := New()
	entry := NewEntry(logger)
	logger.showCaller = true
	logger.Level = DebugLevel
	logger.Formatter = &JSONFormatter{}
	logger.Out = f

	entry.Debug(caller(1))
	entry.Debugf(caller(1))
	entry.Debugln(caller(1))
	e1 := entry.WithError(nil)
	e1.Debug(caller(1))
	e1 = entry.WithField("first", 1)
	e1.Debug(caller(1))
	e1 = entry.WithFields(Fields{"first": 1})
	e1.Debug(caller(1))

	f.Close()

	f, err = os.Open(fileName)
	defer f.Close()
	scaner := bufio.NewScanner(f)
	for scaner.Scan() {
		ss := scaner.Text()
		var output Output
		err = json.Unmarshal([]byte(ss), &output)
		assert.NoError(t, err)
		assert.Equal(t, output.caller, output.msg)
	}
}
