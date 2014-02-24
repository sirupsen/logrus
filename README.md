# Logrus

Logrus is a simple, opinionated logging package for Go. Features include:

* **Level logging**. Logrus has the levels: Debug, Info, Warning and Fatal.
* **Exceptions**. Warnings will log as an exception along with logging it to
  out, without quitting. Fatal will do the same, but call `os.Exit(1)` after
  emitting the exception.
* **JSON**. Logrus currently logs as JSON by default.

The API is completely compatible with the Go standard lib logger, with only the
features above added.

## Motivation

The motivation for this library came out of a pattern seen in Go applications me
and others have been writing with functions such as:

```go
func reportFatalError(err error) {
  airbrake.Notify(err)
  log.Fatal(err)
}

func reportWarning(err error) {
  airbrake.Notify(err)
}
```

JSON logging is excellent for parsing logs for analysis and troubleshooting.
It's supported natively by log aggregators such as logstash and Splunk. Logging
JSON with logrus with the `WithFields` and `WithField` API in logrus forces you
to think about what context to log, to provide valuable troubleshoot information
later.

## Example

```go
import (
  "github.com/Sirupsen/logrus"
)

var logger logrus.New()
func main() {
  logger.WithFields(Fields{
      "animal":   "walrus",
      "location": "New York Aquarium",
      "weather":  "rain",
      "name":     "Wally",
      "event":    "escape",
      }).Info("Walrus has escaped the aquarium! Action required!")
  // {
  //   "level": "info",
  //   "animal": "walrus",
  //   "location": "New York Aquarium",
  //   "weather":"rain",
  //   "name": "Wally",
  //   "event":"escape",
  //   "msg": "Walrus has escaped the aquarium! Action required!")
  //   "time": "2014-02-23 19:57:35.862271048 -0500 EST"
  // }

  logger.WithField("source", "kafka").Infof("Connection to Kafka failed with %s", "some error")
  // {
  //   "level": "info",
  //   "source": "kafka",
  //   "msg": "Connection to Kafka failed with some error",
  //   "time": "2014-02-23 19:57:35.862271048 -0500 EST"
  // }
}
```

Using `Warning` and `Fatal` to log to `airbrake` requires setting
`airbrake.Endpoint` and `airbrake.ApiKey`. See
[tobi/airbrake-go](https://github.com/tobi/airbrake-go).
