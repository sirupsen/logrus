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

	addr := "Address"

	log.Debug("Sample ", addr)
	log.Debugf("Sample %s", addr)
	log.Debugln("Sample", addr)
	log.WithError(nil).Debug("Sample ", addr)
	log.WithError(nil).Debugf("Sample %s", addr)
	log.WithError(nil).Debugln("Sample", addr)
	log.WithField("first", 1).Debug("Sample ", addr)
	log.WithField("first", 1).Debugf("Sample %s", addr)
	log.WithField("first", 1).Debugln("Sample", addr)
	log.WithFields(log.Fields{"first": 1, "second": 2}).Debug("Sample ", addr)
	log.WithFields(log.Fields{"first": 1, "second": 2}).Debugf("Sample %s", addr)
	log.WithFields(log.Fields{"first": 1, "second": 2}).Debugln("Sample", addr)
	log.Info("Sample ", addr)
	log.Infof("Sample %s", addr)
	log.Infoln("Sample", addr)
	log.WithError(nil).Info("Sample ", addr)
	log.WithError(nil).Infof("Sample %s", addr)
	log.WithError(nil).Infoln("Sample", addr)
	log.WithField("first", 1).Info("Sample ", addr)
	log.WithField("first", 1).Infof("Sample %s", addr)
	log.WithField("first", 1).Infoln("Sample", addr)
	log.WithFields(log.Fields{"first": 1, "second": 2}).Info("Sample ", addr)
	log.WithFields(log.Fields{"first": 1, "second": 2}).Infof("Sample %s", addr)
	log.WithFields(log.Fields{"first": 1, "second": 2}).Infoln("Sample", addr)
	log.Warn("Sample ", addr)
	log.Warnf("Sample %s", addr)
	log.Warnln("Sample", addr)
	log.WithError(nil).Warn("Sample ", addr)
	log.WithError(nil).Warnf("Sample %s", addr)
	log.WithError(nil).Warnln("Sample", addr)
	log.WithField("first", 1).Warn("Sample ", addr)
	log.WithField("first", 1).Warnf("Sample %s", addr)
	log.WithField("first", 1).Warnln("Sample", addr)
	log.WithFields(log.Fields{"first": 1, "second": 2}).Warn("Sample ", addr)
	log.WithFields(log.Fields{"first": 1, "second": 2}).Warnf("Sample %s", addr)
	log.WithFields(log.Fields{"first": 1, "second": 2}).Warnln("Sample", addr)

	log.Print("Hello")
	log.Println("Hello")
	log.Printf("Hello %d", 10)

	logg := log.StandardLogger()
	logg.Level = log.DebugLevel
	logg.Debug("Sample ", addr)
	logg.Debugf("Sample %s", addr)
	logg.Debugln("Sample", addr)
	logg.WithError(nil).Debug("Sample ", addr)
	logg.WithError(nil).Debugf("Sample %s", addr)
	logg.WithError(nil).Debugln("Sample", addr)
	logg.WithField("first", 1).Debug("Sample ", addr)
	logg.WithField("first", 1).Debugf("Sample %s", addr)
	logg.WithField("first", 1).Debugln("Sample", addr)
	logg.WithFields(log.Fields{"first": 1, "second": 2}).Debug("Sample ", addr)
	logg.WithFields(log.Fields{"first": 1, "second": 2}).Debugf("Sample %s", addr)
	logg.WithFields(log.Fields{"first": 1, "second": 2}).Debugln("Sample", addr)

	logger := log.New()
	logger.Level = log.DebugLevel
	logger.Debug("Sample ", addr)
	logger.Debugf("Sample %s", addr)
	logger.Debugln("Sample", addr)
	logger.WithError(nil).Debug("Sample ", addr)
	logger.WithError(nil).Debugf("Sample %s", addr)
	logger.WithError(nil).Debugln("Sample", addr)
	logger.WithField("first", 1).Debug("Sample ", addr)
	logger.WithField("first", 1).Debugf("Sample %s", addr)
	logger.WithField("first", 1).Debugln("Sample", addr)
	logger.WithFields(log.Fields{"first": 1, "second": 2}).Debug("Sample ", addr)
	logger.WithFields(log.Fields{"first": 1, "second": 2}).Debugf("Sample %s", addr)
	logger.WithFields(log.Fields{"first": 1, "second": 2}).Debugln("Sample", addr)
	logger.Print("Hello")
	logger.Println("Hello")
	logger.Printf("Hello %d", 10)
}
