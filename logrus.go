package logrus

type Fields map[string]interface{}

type LevelType uint8

const (
	LevelPanic LevelType = iota
	LevelFatal
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

var Level LevelType = LevelInfo

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
