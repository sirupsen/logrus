package logrus

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// skipFrameCount calling depth:
// 0. it's runtime.Caller self invocation.
// 1. content() invocation.
// 2. entry.log() invocation.
// 3. entry.Info/Infof/Infoln/Error() ... entry instance invocation.
// 4. logger.Info/Infof/Infoln/Error() ... logger instance invocation.
const skipFrameCount = 4

var bufferPool *sync.Pool

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
}

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
	// This field will be set on entry firing and the value will be equal to the one in Logger struct field.
	Level Level

	// Message passed to Debug, Info, Warn, Error, Fatal or Panic
	Message string

	// Caller passed to Debug, Info, Warn, Error, Fatal or Panic.
	// To indicates where the log call came from and formats it for output.
	Caller string

	// When formatter is called in entry.log(), a Buffer may be set to entry
	Buffer *bytes.Buffer
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		// Default is five fields, give a little extra room
		Data: make(Fields, 5),
	}
}

// Returns the string representation from the reader and ultimately the
// formatter.
func (entry *Entry) String() (string, error) {
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return "", err
	}
	str := string(serialized)
	return str, nil
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
	data := make(Fields, len(entry.Data)+len(fields))
	for k, v := range entry.Data {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}
	return &Entry{Logger: entry.Logger, Data: data, Time: entry.Time}
}

// Overrides the time of the Entry.
func (entry *Entry) WithTime(t time.Time) *Entry {
	return &Entry{Logger: entry.Logger, Data: entry.Data, Time: t}
}

// This function is not declared with a pointer value because otherwise
// race conditions will occur when using multiple goroutines
func (entry Entry) log(level Level, msg string) {
	var buffer *bytes.Buffer

	// Default to now, but allow users to override if they want.
	//
	// We don't have to worry about polluting future calls to Entry#log()
	// with this assignment because this function is declared with a
	// non-pointer receiver.
	if entry.Time.IsZero() {
		entry.Time = time.Now()
	}

	entry.Level = level
	entry.Message = msg

	// Info level log does not need to print caller information
	defer entry.Logger.resetCallerOffset()
	if entry.Level != InfoLevel {
		depth := entry.Logger.getCallerOffset() + skipFrameCount
		entry.Caller = context(depth)
	}

	entry.fireHooks()

	buffer = bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)
	entry.Buffer = buffer

	entry.write()

	entry.Buffer = nil

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if level <= PanicLevel {
		// In order to prevent subsequent logging getting an incorrect offset.
		// We must reset the offset before the panic.
		entry.Logger.resetCallerOffset()
		panic(&entry)
	}

}

func (entry *Entry) fireHooks() {
	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	err := entry.Logger.Hooks.Fire(entry.Level, entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fire hook: %v\n", err)
	}
}

func (entry *Entry) write() {
	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
	} else {
		_, err = entry.Logger.Out.Write(serialized)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
		}
	}
}

func (entry *Entry) Debug(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(DebugLevel) {
		entry.log(DebugLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Print(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(InfoLevel) {
		entry.log(InfoLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Info(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(InfoLevel) {
		entry.log(InfoLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Warn(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(WarnLevel) {
		entry.log(WarnLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Warning(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(WarnLevel) {
		entry.log(WarnLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Error(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(ErrorLevel) {
		entry.log(ErrorLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Fatal(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(FatalLevel) {
		entry.log(FatalLevel, fmt.Sprint(args...))
	}
	Exit(1)
}

func (entry *Entry) Panic(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(PanicLevel) {
		entry.log(PanicLevel, fmt.Sprint(args...))
	}
	panic(fmt.Sprint(args...))
}

// Entry Printf family functions

func (entry *Entry) Debugf(format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(DebugLevel) {
		entry.log(DebugLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(InfoLevel) {
		entry.log(InfoLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Printf(format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(InfoLevel) {
		entry.log(InfoLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(WarnLevel) {
		entry.log(WarnLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Warningf(format string, args ...interface{}) {
	entry.log(WarnLevel, fmt.Sprintf(format, args...))
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(ErrorLevel) {
		entry.log(ErrorLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(FatalLevel) {
		entry.log(FatalLevel, fmt.Sprintf(format, args...))
	}
	Exit(1)
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(PanicLevel) {
		entry.log(PanicLevel, fmt.Sprintf(format, args...))
	}
	panic(fmt.Sprintf(format, args...))
}

// Entry Println family functions

func (entry *Entry) Debugln(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(DebugLevel) {
		entry.log(DebugLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Infoln(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(InfoLevel) {
		entry.log(InfoLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Println(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(InfoLevel) {
		entry.log(InfoLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Warnln(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(WarnLevel) {
		entry.log(WarnLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Warningln(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(WarnLevel) {
		entry.log(WarnLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Errorln(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(ErrorLevel) {
		entry.log(ErrorLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Fatalln(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(FatalLevel) {
		entry.log(FatalLevel, entry.sprintlnn(args...))
	}
	Exit(1)
}

func (entry *Entry) Panicln(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(PanicLevel) {
		entry.log(PanicLevel, entry.sprintlnn(args...))
	}
	panic(entry.sprintlnn(args...))
}

// Sprintlnn => Sprint no newline. This is to get the behavior of how
// fmt.Sprintln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func (entry *Entry) sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}

// Captures where the log call came from and formats it for output.
func context(depth int) string {
	if _, file, line, ok := runtime.Caller(depth); ok {
		return strings.Join([]string{filepath.Base(file), strconv.Itoa(line)}, ":")
	}

	return "???:0"
}
