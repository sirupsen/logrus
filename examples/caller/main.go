// main
package main

import (
	"time"

	log "github.com/omidnikta/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}
func main() {
	logg := log.StandardLogger()
	logg.Level = log.DebugLevel
	logger := log.New()
	logger.Level = log.DebugLevel
	check(logg)
	check(logger)

	log.Debug(log.Caller())
	log.Debugf(log.Caller())
	log.Debugln(log.Caller())
	log.WithError(nil).Debug(log.Caller())
	log.WithError(nil).Debugf(log.Caller())
	log.WithError(nil).Debugln(log.Caller())
	log.WithField("first", 1).Debug(log.Caller())
	log.WithField("first", 1).Debugf(log.Caller())
	log.WithField("first", 1).Debugln(log.Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Debug(log.Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Debugf(log.Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Debugln(log.Caller())
	log.Info(log.Caller())
	log.Infof(log.Caller())
	log.Infoln(log.Caller())
	log.WithError(nil).Info(log.Caller())
	log.WithError(nil).Infof(log.Caller())
	log.WithError(nil).Infoln(log.Caller())
	log.WithField("first", 1).Info(log.Caller())
	log.WithField("first", 1).Infof(log.Caller())
	log.WithField("first", 1).Infoln(log.Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Info(log.Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Infof(log.Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Infoln(log.Caller())
	log.Warn(log.Caller())
	log.Warnf(log.Caller())
	log.Warnln(log.Caller())
	log.WithError(nil).Warn(log.Caller())
	log.WithError(nil).Warnf(log.Caller())
	log.WithError(nil).Warnln(log.Caller())
	log.WithField("first", 1).Warn(log.Caller())
	log.WithField("first", 1).Warnf(log.Caller())
	log.WithField("first", 1).Warnln(log.Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Warn(log.Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Warnf(log.Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Warnln(log.Caller())
	log.Error(log.Caller())
	log.Errorf(log.Caller())
	log.Errorln(log.Caller())
	log.WithError(nil).Error(log.Caller())
	log.WithError(nil).Errorf(log.Caller())
	log.WithError(nil).Errorln(log.Caller())
	log.WithField("first", 1).Error(log.Caller())
	log.WithField("first", 1).Errorf(log.Caller())
	log.WithField("first", 1).Errorln(log.Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Error(log.Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Errorf(log.Caller())
	log.WithFields(log.Fields{"first": 1, "second": 2}).Errorln(log.Caller())

	log.Print(log.Caller())
	log.Println(log.Caller())
	log.Printf(log.Caller())

	contextLogger := log.WithFields(log.Fields{
		"common": "this is a common field",
		"other":  "I also should be logged always",
	})

	contextLogger.Info("I'll be logged with common and other field")
	contextLogger.Info("Me too")

}

func check(l *log.Logger) {
	l.Debug(log.Caller())
	l.Debugf(log.Caller())
	l.Debugln(log.Caller())
	l.WithError(nil).Debug(log.Caller())
	l.WithError(nil).Debugf(log.Caller())
	l.WithError(nil).Debugln(log.Caller())
	l.WithField("first", 1).Debug(log.Caller())
	l.WithField("first", 1).Debugf(log.Caller())
	l.WithField("first", 1).Debugln(log.Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Debug(log.Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Debugf(log.Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Debugln(log.Caller())
	l.Info(log.Caller())
	l.Infof(log.Caller())
	l.Infoln(log.Caller())
	l.WithError(nil).Info(log.Caller())
	l.WithError(nil).Infof(log.Caller())
	l.WithError(nil).Infoln(log.Caller())
	l.WithField("first", 1).Info(log.Caller())
	l.WithField("first", 1).Infof(log.Caller())
	l.WithField("first", 1).Infoln(log.Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Info(log.Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Infof(log.Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Infoln(log.Caller())
	l.Warn(log.Caller())
	l.Warnf(log.Caller())
	l.Warnln(log.Caller())
	l.WithError(nil).Warn(log.Caller())
	l.WithError(nil).Warnf(log.Caller())
	l.WithError(nil).Warnln(log.Caller())
	l.WithField("first", 1).Warn(log.Caller())
	l.WithField("first", 1).Warnf(log.Caller())
	l.WithField("first", 1).Warnln(log.Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Warn(log.Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Warnf(log.Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Warnln(log.Caller())
	l.Error(log.Caller())
	l.Errorf(log.Caller())
	l.Errorln(log.Caller())
	l.WithError(nil).Error(log.Caller())
	l.WithError(nil).Errorf(log.Caller())
	l.WithError(nil).Errorln(log.Caller())
	l.WithField("first", 1).Error(log.Caller())
	l.WithField("first", 1).Errorf(log.Caller())
	l.WithField("first", 1).Errorln(log.Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Error(log.Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Errorf(log.Caller())
	l.WithFields(log.Fields{"first": 1, "second": 2}).Errorln(log.Caller())

	l.Print(log.Caller())
	l.Println(log.Caller())
	l.Printf(log.Caller())

	entry := log.NewEntry(l)
	j := 1
	for i := 0; i < j; i++ {
		go func(en *log.Entry) {
			en.Debug(log.Caller())
			en.WithField("first", 1).Debug(log.Caller())
			en.WithFields(log.Fields{"first": 1, "second": 2}).Debug(log.Caller())
		}(entry)
	}

	time.Sleep(time.Second)
	//	for i := 0; i < 100; i++ {
	//		go entry.Debug("79")
	//		go entry.Debug("80")
	//		go entry.Debug("81")
	//	}
}
