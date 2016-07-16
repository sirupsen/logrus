package logrus

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Output struct {
	caller string `json:"caller"`
	msg    string `json:"msg"`
}

func TestEntryTEXTFormatterCaller(t *testing.T) {

	sample := []byte{}
	logger := New()
	logger.showCaller = true
	logger.Level = DebugLevel
	logger.Formatter = &TextFormatter{DisableColors: true}
	logger.Out = bytes.NewBuffer(sample)

	check(logger)
	scanner := bufio.NewScanner(bytes.NewReader(sample))
	for scanner.Scan() {
		ss := scanner.Text()
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
	sample := []byte{}
	logger := New()
	logger.showCaller = true
	logger.Level = DebugLevel
	logger.Formatter = &JSONFormatter{}
	logger.Out = bytes.NewBuffer(sample)

	check(logger)

	scanner := bufio.NewScanner(bytes.NewReader(sample))
	for scanner.Scan() {
		ss := scanner.Text()
		var output Output
		err := json.Unmarshal([]byte(ss), &output)
		assert.NoError(t, err)
		assert.Equal(t, output.caller, output.msg)
	}
}

func check(l *Logger) {
	l.Debug(caller(1))
	l.Debugf(caller(1))
	l.Debugln(caller(1))
	l.WithError(nil).Debug(caller(1))
	l.WithError(nil).Debugf(caller(1))
	l.WithError(nil).Debugln(caller(1))
	l.WithField("caller", 1).Debug(caller(1))
	l.WithField("first", 1).Debugf(caller(1))
	l.WithField("first", 1).Debugln(caller(1))
	l.WithFields(Fields{"caller": 1, "second": 2}).Debug(caller(1))
	l.WithFields(Fields{"first": 1, "second": 2}).Debugf(caller(1))
	l.WithFields(Fields{"first": 1, "second": 2}).Debugln(caller(1))
	l.Info(caller(1))
	l.Infof(caller(1))
	l.Infoln(caller(1))
	l.WithError(nil).Info(caller(1))
	l.WithError(nil).Infof(caller(1))
	l.WithError(nil).Infoln(caller(1))
	l.WithField("caller", 1).Info(caller(1))
	l.WithField("first", 1).Infof(caller(1))
	l.WithField("first", 1).Infoln(caller(1))
	l.WithFields(Fields{"caller": 1, "second": 2}).Info(caller(1))
	l.WithFields(Fields{"first": 1, "second": 2}).Infof(caller(1))
	l.WithFields(Fields{"first": 1, "second": 2}).Infoln(caller(1))
	l.Warn(caller(1))
	l.Warnf(caller(1))
	l.Warnln(caller(1))
	l.WithError(nil).Warn(caller(1))
	l.WithError(nil).Warnf(caller(1))
	l.WithError(nil).Warnln(caller(1))
	l.WithField("caller", 1).Warn(caller(1))
	l.WithField("first", 1).Warnf(caller(1))
	l.WithField("first", 1).Warnln(caller(1))
	l.WithFields(Fields{"caller": 1, "second": 2}).Warn(caller(1))
	l.WithFields(Fields{"first": 1, "second": 2}).Warnf(caller(1))
	l.WithFields(Fields{"first": 1, "second": 2}).Warnln(caller(1))
	l.Error(caller(1))
	l.Errorf(caller(1))
	l.Errorln(caller(1))
	l.WithError(nil).Error(caller(1))
	l.WithError(nil).Errorf(caller(1))
	l.WithError(nil).Errorln(caller(1))
	l.WithField("caller", 1).Error(caller(1))
	l.WithField("first", 1).Errorf(caller(1))
	l.WithField("first", 1).Errorln(caller(1))
	l.WithFields(Fields{"caller": 1, "second": 2}).Error(caller(1))
	l.WithFields(Fields{"first": 1, "second": 2}).Errorf(caller(1))
	l.WithFields(Fields{"first": 1, "second": 2}).Errorln(caller(1))

	l.Print(caller(1))
	l.Println(caller(1))
	l.Printf(caller(1))

	entry := NewEntry(l)
	j := 100
	var wg sync.WaitGroup
	for i := 0; i < j; i++ {
		wg.Add(1)
		go func(en *Entry, wg *sync.WaitGroup) {
			en.Debug(caller(1))
			en.WithField("caller", 1).Debug(caller(1))
			en.WithFields(Fields{"caller": 1, "second": 2}).Debug(caller(1))
			wg.Done()
		}(entry, &wg)
	}
	wg.Wait()
}
