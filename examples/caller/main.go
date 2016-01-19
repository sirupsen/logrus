// main
package main

import (
	log "github.com/omid/logrus"
)

func init() {
	log.ShowCaller(true)
	log.SetLevel(log.DebugLevel)
}
func main() {
	logg := log.StandardLogger()
	logg.Level = log.DebugLevel
	logger := log.New()
	logger.Level = log.DebugLevel
	check(logg)
	check(logger)
}

func check(l *log.Logger) {
	addr := "Address"

	l.Debug("Sample ", addr)
	l.Debugf("Sample %s", addr)
	l.Debugln("Sample", addr)
	l.WithError(nil).Debug("Sample ", addr)
	l.WithError(nil).Debugf("Sample %s", addr)
	l.WithError(nil).Debugln("Sample", addr)
	l.WithField("first", 1).Debug("Sample ", addr)
	l.WithField("first", 1).Debugf("Sample %s", addr)
	l.WithField("first", 1).Debugln("Sample", addr)
	l.WithFields(log.Fields{"first": 1, "second": 2}).Debug("Sample ", addr)
	l.WithFields(log.Fields{"first": 1, "second": 2}).Debugf("Sample %s", addr)
	l.WithFields(log.Fields{"first": 1, "second": 2}).Debugln("Sample", addr)
	l.Info("Sample ", addr)
	l.Infof("Sample %s", addr)
	l.Infoln("Sample", addr)
	l.WithError(nil).Info("Sample ", addr)
	l.WithError(nil).Infof("Sample %s", addr)
	l.WithError(nil).Infoln("Sample", addr)
	l.WithField("first", 1).Info("Sample ", addr)
	l.WithField("first", 1).Infof("Sample %s", addr)
	l.WithField("first", 1).Infoln("Sample", addr)
	l.WithFields(log.Fields{"first": 1, "second": 2}).Info("Sample ", addr)
	l.WithFields(log.Fields{"first": 1, "second": 2}).Infof("Sample %s", addr)
	l.WithFields(log.Fields{"first": 1, "second": 2}).Infoln("Sample", addr)
	l.Warn("Sample ", addr)
	l.Warnf("Sample %s", addr)
	l.Warnln("Sample", addr)
	l.WithError(nil).Warn("Sample ", addr)
	l.WithError(nil).Warnf("Sample %s", addr)
	l.WithError(nil).Warnln("Sample", addr)
	l.WithField("first", 1).Warn("Sample ", addr)
	l.WithField("first", 1).Warnf("Sample %s", addr)
	l.WithField("first", 1).Warnln("Sample", addr)
	l.WithFields(log.Fields{"first": 1, "second": 2}).Warn("Sample ", addr)
	l.WithFields(log.Fields{"first": 1, "second": 2}).Warnf("Sample %s", addr)
	l.WithFields(log.Fields{"first": 1, "second": 2}).Warnln("Sample", addr)
	l.Error("Sample ", addr)
	l.Errorf("Sample %s", addr)
	l.Errorln("Sample", addr)
	l.WithError(nil).Error("Sample ", addr)
	l.WithError(nil).Errorf("Sample %s", addr)
	l.WithError(nil).Errorln("Sample", addr)
	l.WithField("first", 1).Error("Sample ", addr)
	l.WithField("first", 1).Errorf("Sample %s", addr)
	l.WithField("first", 1).Errorln("Sample", addr)
	l.WithFields(log.Fields{"first": 1, "second": 2}).Error("Sample ", addr)
	l.WithFields(log.Fields{"first": 1, "second": 2}).Errorf("Sample %s", addr)
	l.WithFields(log.Fields{"first": 1, "second": 2}).Errorln("Sample", addr)

	l.Print("Hello")
	l.Println("Hello")
	l.Printf("Hello %d", 10)
}
