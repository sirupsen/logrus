// main
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	log "github.com/omidnikta/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	formatter := &log.TextFormatter{DisableColors: true}
	//	formatter := &log.JSONFormatter{}
	f, err := os.OpenFile("caller.txt", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	log.SetOutput(f)

	logg := log.StandardLogger()
	logg.Level = log.DebugLevel
	logg.Formatter = formatter
	logg.ShowCaller(true)
	logger := log.New()
	logger.ShowCaller(true)
	logger.Level = log.DebugLevel
	logger.Formatter = formatter
	check(logg)
	check(logger)

	log.StandardLogger().ShowCaller(true)
	log.SetFormatter(formatter)
	log.Debug(Caller())
	log.Debugf(Caller())
	log.Debugln(Caller())
	log.WithError(nil).Debug(Caller())
	log.WithError(nil).Debugf(Caller())
	log.WithError(nil).Debugln(Caller())
	log.WithField("caller", 1).Debug(Caller())
	log.WithField("first", 1).Debugf(Caller())
	log.WithField("first", 1).Debugln(Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Debug(Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Debugf(Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Debugln(Caller())
	log.Info(Caller())
	log.Infof(Caller())
	log.Infoln(Caller())
	log.WithError(nil).Info(Caller())
	log.WithError(nil).Infof(Caller())
	log.WithError(nil).Infoln(Caller())
	log.WithField("first", 1).Info(Caller())
	log.WithField("first", 1).Infof(Caller())
	log.WithField("first", 1).Infoln(Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Info(Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Infof(Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Infoln(Caller())
	log.Warn(Caller())
	log.Warnf(Caller())
	log.Warnln(Caller())
	log.WithError(nil).Warn(Caller())
	log.WithError(nil).Warnf(Caller())
	log.WithError(nil).Warnln(Caller())
	log.WithField("first", 1).Warn(Caller())
	log.WithField("first", 1).Warnf(Caller())
	log.WithField("first", 1).Warnln(Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Warn(Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Warnf(Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Warnln(Caller())
	log.Error(Caller())
	log.Errorf(Caller())
	log.Errorln(Caller())
	log.WithError(nil).Error(Caller())
	log.WithError(nil).Errorf(Caller())
	log.WithError(nil).Errorln(Caller())
	log.WithField("first", 1).Error(Caller())
	log.WithField("first", 1).Errorf(Caller())
	log.WithField("first", 1).Errorln(Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Error(Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Errorf(Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Errorln(Caller())

	log.Print(Caller())
	log.Println(Caller())
	log.Printf(Caller())

	contextLogger := log.WithFields(log.Fields{
		"common": "this is a common field",
		"other":  "I also should be logged always",
	})

	contextLogger.Info(Caller())
	contextLogger.Info(Caller())

}

func check(l *log.Logger) {
	l.Debug(Caller())
	l.Debugf(Caller())
	l.Debugln(Caller())
	l.WithError(nil).Debug(Caller())
	l.WithError(nil).Debugf(Caller())
	l.WithError(nil).Debugln(Caller())
	l.WithField("first", 1).Debug(Caller())
	l.WithField("first", 1).Debugf(Caller())
	l.WithField("first", 1).Debugln(Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Debug(Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Debugf(Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Debugln(Caller())
	l.Info(Caller())
	l.Infof(Caller())
	l.Infoln(Caller())
	l.WithError(nil).Info(Caller())
	l.WithError(nil).Infof(Caller())
	l.WithError(nil).Infoln(Caller())
	l.WithField("first", 1).Info(Caller())
	l.WithField("first", 1).Infof(Caller())
	l.WithField("first", 1).Infoln(Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Info(Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Infof(Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Infoln(Caller())
	l.Warn(Caller())
	l.Warnf(Caller())
	l.Warnln(Caller())
	l.WithError(nil).Warn(Caller())
	l.WithError(nil).Warnf(Caller())
	l.WithError(nil).Warnln(Caller())
	l.WithField("first", 1).Warn(Caller())
	l.WithField("first", 1).Warnf(Caller())
	l.WithField("first", 1).Warnln(Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Warn(Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Warnf(Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Warnln(Caller())
	l.Error(Caller())
	l.Errorf(Caller())
	l.Errorln(Caller())
	l.WithError(nil).Error(Caller())
	l.WithError(nil).Errorf(Caller())
	l.WithError(nil).Errorln(Caller())
	l.WithField("first", 1).Error(Caller())
	l.WithField("first", 1).Errorf(Caller())
	l.WithField("first", 1).Errorln(Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Error(Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Errorf(Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Errorln(Caller())

	l.Print(Caller())
	l.Println(Caller())
	l.Printf(Caller())

	entry := log.NewEntry(l)
	j := 1
	for i := 0; i < j; i++ {
		go func(en *log.Entry) {
			en.Debug(Caller())
			en.WithField("first", 1).Debug(Caller())
			en.WithFields(log.Fields{"first": 1, "second": 2}).Debug(Caller())
		}(entry)
	}
}

func Caller() (str string) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		str = "???: ?"
	} else {
		str = fmt.Sprint(filepath.Base(file), ":", line)
	}
	return
}
