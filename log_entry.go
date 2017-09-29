package logrus

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
	"sync/atomic"
)

const (
	UnknownDesiredLevel Level = iota
	PanicDesiredLevel
	FatalDesiredLevel
	ErrorDesiredLevel
	WarnDesiredLevel
	InfoDesiredLevel
	DebugDesiredLevel
)

// An entry is the final or intermediate Logrus logging entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Debug, Info,
// Warn, Error, Fatal or Panic is called on it. These objects can be reused and
// passed around as much as you wish to avoid field duplication.
type LogEntry struct {
	Logger *LogWriter

	// Contains all the fields set by the user.
	Data Fields

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry was logged at: Debug, Info, Warn, Error, Fatal or Panic
	// This field will be set on entry firing and the value will be equal to the one in Logger struct field.
	level Level

	desiredLevel Level

	// Message passed to Debug, Info, Warn, Error, Fatal or Panic
	Message string

	// When formatter is called in entry.log(), an Buffer may be set to entry
	Buffer *bytes.Buffer
}

func NewLogEntry(logger *LogWriter) *LogEntry {
	return &LogEntry{
		Logger:       logger,
		desiredLevel: Level(logger.Level + 1),
		// Default is three fields, give a little extra room
		Data: make(Fields, 5),
	}
}

func (entry *LogEntry) AsLevel(level Level) *LogEntry {
	return logger.newEntry()
}

func (entry *LogEntry) AsDebug() *LogEntry {
	return logger.AsLevel(DebugLevel)
}

func (entry *LogEntry) AsInfo() *LogEntry {
	return logger.AsLevel(InfoLevel)
}

func (entry *LogEntry) AsWarning() *LogEntry {
	return logger.AsLevel(WarnLevel)
}

func (entry *LogEntry) AsError() *LogEntry {
	return logger.AsLevel(ErrorLevel)
}

func (entry *LogEntry) AsFatal() *LogEntry {
	return logger.AsLevel(FatalLevel)
}

func (entry *LogEntry) AsPanic() *LogEntry {
	return logger.AsLevel(PanicLevel)
}

func (entry *LogEntry) WithField(key string, value interface{}) *LogEntry {
	return nil
}

func (entry *LogEntry) WithFields(fields Fields) *LogEntry {
	return nil
}

func (entry *LogEntry) WithError(err error) *LogEntry {
	return nil
}

func (entry *LogEntry) WriteF(format string, args ...interface{}) {

}

func (entry *LogEntry) setDesiredLevel(level Level) {
	atomic.StoreInt32((*uint32)(&entry.desiredLevel), int32(level))
}
