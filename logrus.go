package logrus

import (
	"log"
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

// Level type
type Level uint8

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// Panic level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	Panic Level = iota
	// Fatal level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	Fatal
	// Error level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	Error
	// Warn level. Non-critical entries that deserve eyes.
	Warn
	// Info level. General operational entries about what's going on inside the
	// application.
	Info
	// Debug level. Usually only enabled when debugging. Very verbose logging.
	Debug
)

// Won't compile if StdLogger can't be realized by a log.Logger
var _ StdLogger = &log.Logger{}

// StdLogger is what your logrus-enabled library should take, that way
// it'll accept a stdlib logger and a logrus logger. There's no standard
// interface, this is the closest we get, unfortunately.
type StdLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}
