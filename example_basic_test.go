package logrus_test

import (
	"os"

	"github.com/barryz/logrus"
)

func Example_basic() {
	var log = logrus.New()
	log.Formatter = new(logrus.JSONFormatter)
	log.Formatter = new(logrus.TextFormatter)                     //default
	log.Formatter.(*logrus.TextFormatter).DisableColors = true    // remove colors
	log.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove timestamp from test output
	log.Level = logrus.DebugLevel
	log.Out = os.Stdout

	// file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	// if err == nil {
	// 	log.Out = file
	// } else {
	// 	log.Info("Failed to log to file, using default stderr")
	// }

	defer func() {
		err := recover()
		if err != nil {
			entry := err.(*logrus.Entry)
			log.WithFields(logrus.Fields{
				"omg":         true,
				"err_animal":  entry.Data["animal"],
				"err_size":    entry.Data["size"],
				"err_level":   entry.Level,
				"err_message": entry.Message,
				"number":      100,
			}).Error("The ice breaks!") // or use Fatal() to force the process to exit with a nonzero code
		}
	}()

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
		"animal": "orca",
		"size":   9009,
	}).Panic("It's over 9000!")

	// Output:
	// level=debug msg="Started observing beach" caller="example_basic_test.go:43" animal=walrus number=8
	// level=info msg="A group of walrus emerges from the ocean" animal=walrus size=10
	// level=warning msg="The group's number increased tremendously!" caller="example_basic_test.go:53" number=122 omg=true
	// level=debug msg="Temperature changes" caller="example_basic_test.go:57" temperature=-4
	// level=panic msg="It's over 9000!" caller="example_basic_test.go:62" animal=orca size=9009
	// level=error msg="The ice breaks!" caller="example_basic_test.go:36" err_animal=orca err_level=panic err_message="It's over 9000!" err_size=9009 number=100 omg=true
}

func Example_exported_basic() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})

	defer func() {
		err := recover()
		if err != nil {
			entry := err.(*logrus.Entry)
			logrus.WithFields(logrus.Fields{
				"omg":         true,
				"err_animal":  entry.Data["animal"],
				"err_size":    entry.Data["size"],
				"err_level":   entry.Level,
				"err_message": entry.Message,
				"number":      100,
			}).Error("The ice breaks!") // or use Fatal() to force the process to exit with a nonzero code
		}
	}()

	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"number": 8,
	}).Debug("Started observing beach")

	logrus.WithFields(logrus.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	logrus.WithFields(logrus.Fields{
		"temperature": -4,
	}).Debug("Temperature changes")

	logrus.WithFields(logrus.Fields{
		"animal": "orca",
		"size":   9009,
	}).Panic("It's over 9000!")

	// Output:
	// level=info msg="A group of walrus emerges from the ocean" animal=walrus size=10
	// level=debug msg="Started observing beach" caller="example_basic_test.go:104" animal=walrus number=8
	// level=warning msg="The group's number increased tremendously!" caller="example_basic_test.go:109" number=122 omg=true
	// level=debug msg="Temperature changes" caller="example_basic_test.go:113" temperature=-4
	// level=panic msg="It's over 9000!" caller="example_basic_test.go:118" animal=orca size=9009
	// level=error msg="The ice breaks!" caller="example_basic_test.go:92" err_animal=orca err_level=panic err_message="It's over 9000!" err_size=9009 number=100 omg=true

}
