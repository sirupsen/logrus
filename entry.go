package logrus

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

// Defines the key when adding errors using WithError.
var ErrorKey = "error"

// An entry is the final or intermediate Logrus logging entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Debug, Info,
// Warn, Error, Fatal or Panic is called on it. These objects can be reused and
// passed around as much as you wish to avoid field duplication.
type Entry struct {
	Logger *Logger

	// Contains all the fields set by the user.
	Data Fields

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry was logged at: Debug, Info, Warn, Error, Fatal or Panic
	Level Level

	// Message passed to Debug, Info, Warn, Error, Fatal or Panic
	Message string

	// A receive-only notification channel for knowing when all hooks have
	// finished firing.
	HooksDone <-chan struct{}
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		// Default is three fields, give a little extra room
		Data: make(Fields, 5),
	}
}

// Returns a reader for the entry, which is a proxy to the formatter.
func (entry *Entry) Reader() (*bytes.Buffer, error) {
	serialized, err := entry.Logger.Formatter.Format(entry)
	return bytes.NewBuffer(serialized), err
}

// Returns the string representation from the reader and ultimately the
// formatter.
func (entry *Entry) String() (string, error) {
	reader, err := entry.Reader()
	if err != nil {
		return "", err
	}

	return reader.String(), err
}

// ToError returns the field value of ErrorKey (nil)
func (entry *Entry) ToError() error {
	if err, ok := entry.Data[ErrorKey].(error); ok {
		return err
	}
	return nil
}

// Add an error as single field (using the key defined in ErrorKey) to the Entry.
func (entry *Entry) WithError(err error) *Entry {
	return entry.WithField(ErrorKey, err)
}

// Add a single field to the Entry.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return entry.WithFields(Fields{key: value})
}

// Add a map of fields to the Entry.
func (entry *Entry) WithFields(fields Fields) *Entry {
	data := Fields{}
	for k, v := range entry.Data {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}
	return &Entry{Logger: entry.Logger, Data: data}
}

// This function is not declared with a pointer value because otherwise
// race conditions will occur when using multiple goroutines.
func (entry Entry) log(level Level, msg string, hooksDone chan struct{}) {
	entry.Time = time.Now()
	entry.Level = level
	entry.Message = msg

	entry.Logger.Hooks.Fire(level, &entry, hooksDone)

	reader, err := entry.Reader()
	if err != nil {
		entry.Logger.mu.Lock()
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		entry.Logger.mu.Unlock()
	}

	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()

	_, err = io.Copy(entry.Logger.Out, reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if level <= PanicLevel {
		panic(&entry)
	}
}

func (entry *Entry) Debug(args ...interface{}) *Entry {
	if entry.Logger.Level >= DebugLevel {
		hooksDone := make(chan struct{}, 1)
		entry.HooksDone = hooksDone
		entry.log(DebugLevel, fmt.Sprint(args...), hooksDone)
	}
	return entry
}

func (entry *Entry) Print(args ...interface{}) *Entry {
	return entry.Info(args...)
}

func (entry *Entry) Info(args ...interface{}) *Entry {
	if entry.Logger.Level >= InfoLevel {
		hooksDone := make(chan struct{}, 1)
		entry.HooksDone = hooksDone
		entry.log(InfoLevel, fmt.Sprint(args...), hooksDone)
	}
	return entry
}

func (entry *Entry) Warn(args ...interface{}) *Entry {
	if entry.Logger.Level >= WarnLevel {
		hooksDone := make(chan struct{}, 1)
		entry.HooksDone = hooksDone
		entry.log(WarnLevel, fmt.Sprint(args...), hooksDone)
	}
	return entry
}

func (entry *Entry) Warning(args ...interface{}) *Entry {
	return entry.Warn(args...)
}

func (entry *Entry) Error(args ...interface{}) *Entry {
	if entry.Logger.Level >= ErrorLevel {
		hooksDone := make(chan struct{}, 1)
		entry.HooksDone = hooksDone
		entry.log(ErrorLevel, fmt.Sprint(args...), hooksDone)
	}
	return entry
}

func (entry *Entry) Fatal(args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		hooksDone := make(chan struct{}, 1)
		entry.HooksDone = hooksDone
		entry.log(FatalLevel, fmt.Sprint(args...), hooksDone)
	}
	os.Exit(1)
}

func (entry *Entry) Panic(args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		hooksDone := make(chan struct{}, 1)
		entry.HooksDone = hooksDone
		entry.log(PanicLevel, fmt.Sprint(args...), hooksDone)
	}
	panic(fmt.Sprint(args...))
}

// Entry Printf family functions

func (entry *Entry) Debugf(format string, args ...interface{}) *Entry {
	if entry.Logger.Level >= DebugLevel {
		entry.Debug(fmt.Sprintf(format, args...))
	}
	return entry
}

func (entry *Entry) Infof(format string, args ...interface{}) *Entry {
	if entry.Logger.Level >= InfoLevel {
		entry.Info(fmt.Sprintf(format, args...))
	}
	return entry
}

func (entry *Entry) Printf(format string, args ...interface{}) *Entry {
	return entry.Infof(format, args...)
}

func (entry *Entry) Warnf(format string, args ...interface{}) *Entry {
	if entry.Logger.Level >= WarnLevel {
		entry.Warn(fmt.Sprintf(format, args...))
	}
	return entry
}

func (entry *Entry) Warningf(format string, args ...interface{}) *Entry {
	return entry.Warnf(format, args...)
}

func (entry *Entry) Errorf(format string, args ...interface{}) *Entry {
	if entry.Logger.Level >= ErrorLevel {
		entry.Error(fmt.Sprintf(format, args...))
	}
	return entry
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.Fatal(fmt.Sprintf(format, args...))
	}
	os.Exit(1)
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.Panic(fmt.Sprintf(format, args...))
	}
}

// Entry Println family functions

func (entry *Entry) Debugln(args ...interface{}) *Entry {
	if entry.Logger.Level >= DebugLevel {
		entry.Debug(entry.sprintlnn(args...))
	}
	return entry
}

func (entry *Entry) Infoln(args ...interface{}) *Entry {
	if entry.Logger.Level >= InfoLevel {
		entry.Info(entry.sprintlnn(args...))
	}
	return entry
}

func (entry *Entry) Println(args ...interface{}) *Entry {
	return entry.Infoln(args...)
}

func (entry *Entry) Warnln(args ...interface{}) *Entry {
	if entry.Logger.Level >= WarnLevel {
		entry.Warn(entry.sprintlnn(args...))
	}
	return entry
}

func (entry *Entry) Warningln(args ...interface{}) *Entry {
	return entry.Warnln(args...)
}

func (entry *Entry) Errorln(args ...interface{}) *Entry {
	if entry.Logger.Level >= ErrorLevel {
		entry.Error(entry.sprintlnn(args...))
	}
	return entry
}

func (entry *Entry) Fatalln(args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.Fatal(entry.sprintlnn(args...))
	}
	os.Exit(1)
}

func (entry *Entry) Panicln(args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.Panic(entry.sprintlnn(args...))
	}
}

// Sprintlnn => Sprint no newline. This is to get the behavior of how
// fmt.Sprintln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func (entry *Entry) sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
