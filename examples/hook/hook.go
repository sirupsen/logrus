// +build ignore
// Do NOT include the above line in your code. This is a build constraint used
// to prevent import loops in the code whilst go get'ting it.
// Read more about build constraints in golang here:
// https://golang.org/pkg/go/build/#hdr-Build_Constraints

package main

import (
	"github.com/sirupsen/logrus"
	airbrake "gopkg.in/gemnasium/logrus-airbrake-hook.v2"
)

var log = logrus.New()

func init() {
	log.Formatter = new(logrus.TextFormatter) // default
	log.Hooks.Add(airbrake.NewHook(123, "xyz", "development"))
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
