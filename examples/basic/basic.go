package main

import (
	"github.com/sirupsen/logrus"
	// "os"
	"fmt"
)

var log = logrus.New()
var textFormatter = new(logrus.TextFormatter)

func init() {
	log.Formatter = new(logrus.JSONFormatter)
	log.Formatter = textFormatter // default

	// file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	// if err == nil {
	// 	log.Out = file
	// } else {
	// 	log.Info("Failed to log to file, using default stderr")
	// }

	log.Level = logrus.DebugLevel
}

func main() {
	printExample("Default example")

	// set different color map
	textFormatter.ColorMap = logrus.ColorMap{
		logrus.DebugLevel: logrus.Cyan,
		logrus.InfoLevel: logrus.Gray,
		logrus.WarnLevel: logrus.Yellow,
		logrus.ErrorLevel: logrus.Magenta,
		logrus.FatalLevel: logrus.Red,
		logrus.PanicLevel: logrus.Red,
	}

	printExample("Different color example")

	printErrorExample()
}

func printExample(headline string) {
	fmt.Printf("\n# %s:\n", headline)

	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"number": 8,
	}).Debug("Started observing beach")

	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	log.WithFields(logrus.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	log.WithFields(logrus.Fields{
		"temperature": -4,
	}).Debug("Temperature changes")

	log.WithFields(logrus.Fields{
		"temperature": -20,
	}).Error("Temperature changes")


}

func printErrorExample() {
	fmt.Printf("\n# %s:\n", "Error example")
	defer func() {
		err := recover()
		if err != nil {
			log.WithFields(logrus.Fields{
				"omg":    true,
				"err":    err,
				"number": 100,
			}).Fatal("The ice breaks!")
		}
	}()

	log.WithFields(logrus.Fields{
		"animal": "orca",
		"size":   9009,
	}).Panic("It's over 9000!")
}
