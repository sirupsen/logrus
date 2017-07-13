package main

import (
	"github.com/sirupsen/logrus"
	testHook "github.com/sirupsen/logrus/hooks/test"
)

var log = logrus.New()

func init() {
	log.Formatter = new(logrus.TextFormatter) // default
	_, hook := testHook.NewNullLogger()
	log.Hooks.Add(hook)
}

func main() {
	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	log.WithFields(logrus.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	log.WithFields(logrus.Fields{
		"omg":    true,
		"number": 100,
	}).Fatal("The ice breaks!")
}
