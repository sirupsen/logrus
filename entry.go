package logrus

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
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

	depth int
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		// Default is three fields, give a little extra room
		Data:  make(Fields, 5),
		depth: 3,
	}
}

func newEntry(logger *Logger, d int) *Entry {
	en := NewEntry(logger)
	en.depth = d
	return en
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

// Add an error as single field (using the key defined in ErrorKey) to the Entry.
func (entry *Entry) WithError(err error) *Entry {
	en := entry.WithField(ErrorKey, err)
	en.depth = entry.depth + 1
	return en
}

// Add a single field to the Entry.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	en := entry.WithFields(Fields{key: value})
	en.depth = entry.depth + 1
	return en
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
	return &Entry{Logger: entry.Logger, Data: data, depth: entry.depth + 1}
}

// This function is not declared with a pointer value because otherwise
// race conditions will occur when using multiple goroutines
func (entry Entry) log(level Level, msg string) {
	//	caller(&msg, entry.depth, entry.Logger.showCaller)
	entry.Time = time.Now()
	entry.Level = level
	entry.Message = msg

	if err := entry.Logger.Hooks.Fire(level, &entry); err != nil {
		entry.Logger.mu.Lock()
		fmt.Fprintf(os.Stderr, "Failed to fire hook: %v\n", err)
		entry.Logger.mu.Unlock()
	}

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

func (entry *Entry) Debug(args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.log(DebugLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Print(args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(InfoLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Info(args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(InfoLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Warn(args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.log(WarnLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Warning(args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.log(WarnLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Error(args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.log(ErrorLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Fatal(args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.log(FatalLevel, fmt.Sprint(args...))
	}
	os.Exit(1)
}

func (entry *Entry) Panic(args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.log(PanicLevel, fmt.Sprint(args...))
	}
	panic(fmt.Sprint(args...))
}

// Entry Printf family functions

func (entry *Entry) Debugf(format string, args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.log(DebugLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(InfoLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Printf(format string, args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(InfoLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.log(WarnLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Warningf(format string, args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.log(WarnLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.log(ErrorLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.log(FatalLevel, fmt.Sprintf(format, args...))
	}
	os.Exit(1)
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.log(PanicLevel, fmt.Sprintf(format, args...))
	}
	panic(fmt.Sprintf(format, args...))
}

// Entry Println family functions

func (entry *Entry) Debugln(args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.log(DebugLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Infoln(args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(InfoLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Println(args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(InfoLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Warnln(args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.log(WarnLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Warningln(args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.log(WarnLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Errorln(args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.log(ErrorLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Fatalln(args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.log(FatalLevel, entry.sprintlnn(args...))
	}
	os.Exit(1)
}

func (entry *Entry) Panicln(args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.log(PanicLevel, entry.sprintlnn(args...))
	}
	panic(fmt.Sprint(args...))
}

// Sprintlnn => Sprint no newline. This is to get the behavior of how
// fmt.Sprintln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func (entry *Entry) sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}

func caller(depth int) (str string) {
	_, file, line, ok := runtime.Caller(depth + 1)
	if !ok {
		str = "???: ?"
	} else {
		str = fmt.Sprint(filepath.Base(file), ":", line)
	}
	return
}
