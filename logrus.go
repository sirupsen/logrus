package logrus

import ()

// TODO: Type naming here feels awkward, but the exposed variable should be
// Level. That's more important than the type name, and libraries should be
// reaching for logrus.Level{Debug,Info,Warning,Fatal}, not defining the type
// themselves as an int.
type LevelType uint8
type Fields map[string]interface{}

const (
	LevelPanic LevelType = iota
	LevelFatal
	LevelWarning
	LevelInfo
	LevelDebug
)

var Level LevelType = LevelInfo
var Environment string = "development"

// StandardLogger is what your logrus-enabled library should take, that way
// it'll accept a stdlib logger and a logrus logger. There's no standard
// interface, this is the closest we get, unfortunately.
type StandardLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Printfln(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}
