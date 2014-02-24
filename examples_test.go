package logrus

import (
	"github.com/tobi/airbrake-go"
)

func ExampleLogger_Info() {
	logger := New()
	logger.Info("Simple logging call, compatible with the standard logger")
	// {
	//   "level": "info",
	//   "msg": "Simple logging call, compatible with the standard logger",
	//   "time": "2014-02-23 19:57:35.862271048 -0500 EST"
	// }
}

func ExampleLogger_Warning() {
	logger := New()

	airbrake.Environment = "production"
	airbrake.ApiKey = "valid"
	airbrake.Endpoint = "https://exceptions.example.com/notifer_api/v2/notices"

	// This will send an exception with Airbrake now that it has been setup.
	logger.Warning("Something failed: %s", "failure")
	// {
	//   "level": "warning",
	//   "msg": "Something failed: failure",
	//   "time": "2014-02-23 19:57:35.862271048 -0500 EST"
	// }
}

func ExampleLogger_WithField() {
	logger := New()
	logger.WithField("source", "kafka").Infof("Connection to Kafka failed with %s", "some error")
	// {
	//   "level": "info",
	//   "source": "kafka",
	//   "msg": "Connection to Kafka failed with some error",
	//   "time": "2014-02-23 19:57:35.862271048 -0500 EST"
	// }
}

func ExampleLogger_WithFields() {
	logger := New()
	logger.WithFields(Fields{
		"animal":   "walrus",
		"location": "New York Aquarium",
		"weather":  "rain",
		"name":     "Wally",
		"event":    "escape",
	}).Info("Walrus has escaped the aquarium! Action required!")
	// {
	//   "level": "info",
	// 	 "animal": "walrus",
	// 	 "location": "New York Aquarium",
	// 	 "weather":"rain",
	// 	 "name": "Wally",
	// 	 "event":"escape",
	//   "msg": "Walrus has escaped the aquarium! Action required!")
	//   "time": "2014-02-23 19:57:35.862271048 -0500 EST"
	// }
}
