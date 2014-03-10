package main

import (
	"github.com/Sirupsen/logrus"
)

func main() {
	log := logrus.New()

	for {
		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   "10",
		}).Print("Hello WOrld!!")

		log.WithFields(logrus.Fields{
			"omg":    true,
			"number": 122,
		}).Warn("There were some omgs")

		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   "10",
		}).Print("Hello WOrld!!")

		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   "10",
		}).Print("Hello WOrld!!")

		log.WithFields(logrus.Fields{
			"omg":    true,
			"number": 122,
		}).Fatal("There were some omgs")
	}
}
