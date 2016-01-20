// main
package main

import (
	"time"

	log "github.com/omidnikta/logrus"
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
	log.Error("Sample ", addr)
	log.Errorf("Sample %s", addr)
	log.Errorln("Sample", addr)
	log.WithError(nil).Error("Sample ", addr)
	log.WithError(nil).Errorf("Sample %s", addr)
	log.WithError(nil).Errorln("Sample", addr)
	log.WithField("first", 1).Error("Sample ", addr)
	log.WithField("first", 1).Errorf("Sample %s", addr)
	log.WithField("first", 1).Errorln("Sample", addr)
	log.WithFields(log.Fields{"first": 1, "second": 2}).Error("Sample ", addr)
	log.WithFields(log.Fields{"first": 1, "second": 2}).Errorf("Sample %s", addr)
	log.WithFields(log.Fields{"first": 1, "second": 2}).Errorln("Sample", addr)

	log.Print("Hello")
	log.Println("Hello")
	log.Printf("Hello %d", 10)

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

	entry := log.NewEntry(l)
	j := 1
	for i := 0; i < j; i++ {
		go func(en *log.Entry) {
			en.WithField("first", 1).Debug("----- -------- 83")
			en.WithFields(log.Fields{"first": 1, "second": 2}).Debug("----- -------- 84")
		}(entry)
	}
	for i := 0; i < j; i++ {
		go func(en *log.Entry) {
			en.WithField("first", 1).Debug("----- -------- 89")
			en.WithFields(log.Fields{"first": 1, "second": 2}).Debug("----- -------- 90")
		}(entry)
	}
	time.Sleep(time.Second)
	//	for i := 0; i < 100; i++ {
	//		go entry.Debug("79")
	//		go entry.Debug("80")
	//		go entry.Debug("81")
	//	}
}
