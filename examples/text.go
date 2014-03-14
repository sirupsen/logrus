package main

import (
	"github.com/Sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.Formatter = new(logrus.JSONFormatter)

	for {
		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   10,
		}).Print("A group of walrus emerges from the ocean")

		log.WithFields(logrus.Fields{
			"omg":    true,
			"number": 122,
		}).Warn("The group's number increased tremendously!")

		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   10,
		}).Print("A giant walrus appears!")

		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   9,
		}).Print("Tremendously sized cow enters the ocean.")

		log.WithFields(logrus.Fields{
			"omg":    true,
			"number": 100,
		}).Fatal("The ice breaks!")
	}
}
