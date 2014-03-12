package main

import (
	"github.com/Sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.Formatter = logrus.NewZalgoFormatterrrrrr()
	for {
		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   "10",
		}).Print("To invoke the hive-mind representing chaos.")

		log.WithFields(logrus.Fields{
			"omg":    true,
			"number": 122,
		}).Warn("Invoking the feeling of chaos.")

		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   "10",
		}).Print("With out order.")

		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   "9",
		}).Error("The Nezperdian hive-mind of chaos. Zalgo.")

		log.WithFields(logrus.Fields{
			"omg":    true,
			"number": 100,
		}).Warn("He who Waits Behind The Wall.")

		log.Fatal("ZALGO !")
	}
}
