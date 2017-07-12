package main

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/stack"
)

var log = logrus.New()

func init() {
	log.Formatter = new(logrus.TextFormatter) // default
	log.Hooks.Add(&stack.CodeLineHook{LogLevel: logrus.DebugLevel})
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

	// If you set FieldsLogger, you can print the log directly to the object Fields
	logrus.SetFieldsLogger(log)
	logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}.Info("A group of walrus emerges from the ocean")

	log.WithFields(logrus.Fields{
		"omg":    true,
		"number": 100,
	}).Fatal("The ice breaks!")
}
