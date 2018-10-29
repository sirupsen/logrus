package logrus

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	bufferPool *sync.Pool

	// qualified package name, cached at first use
	logrusPackage string

	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth int

	// Used for caller information initialisation
	callerInitOnce sync.Once
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4
)

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	// start at the bottom of the stack before the package-name cache is primed
	minimumCallerDepth = 1
}

// ErrorKey defines the key when adding errors using WithError.
var ErrorKey = "error"

// errorSlice is an implementation of error that holds multiple errors.
// Consider switching to github.com/pkg/errors in the future.
type errorSlice []error

// Error returns the errors as a string, makeing *errorSlice an implementation
// of error
func (arrOfErrs *errorSlice) Error() string {
	strBuilder := strings.Builder{}
	for i, eachErr := range *arrOfErrs {
		if i > 0 {
			strBuilder.WriteString(", ")
		}
		strBuilder.WriteString(eachErr.Error())
	}
	return strBuilder.String()
}

// appendError appends newErr to zeroOrMoreErrors error (converts to errorSlice if needed)
func appendError(zeroOrMoreErrors error, newErr error) error {
	var multipleErrors errorSlice
	if newErr == nil {
		return zeroOrMoreErrors
	} else if zeroOrMoreErrors == nil {
		return newErr
	} else if pErrSlice, ok := zeroOrMoreErrors.(*errorSlice); ok {
		multipleErrors = append(*pErrSlice, newErr)
	} else {
		multipleErrors = errorSlice{zeroOrMoreErrors, newErr}
	}
	return &multipleErrors
}

// Entry is the final or intermediate Logrus logging entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Trace, Debug,
// Info, Warn, Error, Fatal or Panic is called on it. These objects can be
// reused and passed around as much as you wish to avoid field duplication.
type Entry struct {
	Logger *Logger

	// Contains all the fields set by the user.
	Data Fields

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry was logged at: Trace, Debug, Info, Warn, Error, Fatal or Panic
	// This field will be set on entry firing and the value will be equal to the one in Logger struct field.
	Level Level

	// Calling method, with package name
	Caller *runtime.Frame

	// Message passed to Trace, Debug, Info, Warn, Error, Fatal or Panic
	Message string

	// When formatter is called in entry.log(), a Buffer may be set to entry
	Buffer *bytes.Buffer

	// fieldErrs may contain field formatting errors
	fieldErrs error
}

// NewEntry returns a new Entry.
func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		// Default is three fields, plus one optional.  Give a little extra room.
		Data: make(Fields, 6),
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

// WithError adds an error as single field (using the key defined in ErrorKey) to the Entry.
func (entry *Entry) WithError(err error) *Entry {
	return entry.WithField(ErrorKey, err)
}

// WithField adds a single field to the Entry.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return entry.WithFields(Fields{key: value})
}

// WithFields adds a map of fields to the Entry.
func (entry *Entry) WithFields(fields Fields) *Entry {
	data := make(Fields, len(entry.Data)+len(fields))
	for k, v := range entry.Data {
		data[k] = v
	}
	fieldErrs := entry.fieldErrs
	for k, v := range fields {
		isErrField := false
		if t := reflect.TypeOf(v); t != nil {
			switch t.Kind() {
			case reflect.Func:
				isErrField = true
			case reflect.Ptr:
				isErrField = t.Elem().Kind() == reflect.Func
			}
		}
		if isErrField {
			fieldErrs = appendError(fieldErrs, fmt.Errorf("can not add field %q", k))
		} else {
			data[k] = v
		}
	}
	return &Entry{Logger: entry.Logger, Data: data, Time: entry.Time, fieldErrs: fieldErrs}
}

// WithTime overrides the time of the Entry.
func (entry *Entry) WithTime(t time.Time) *Entry {
	return &Entry{Logger: entry.Logger, Data: entry.Data, Time: t, fieldErrs: entry.fieldErrs}
}

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

// getCaller retrieves the name of the first non-logrus calling function
func getCaller() *runtime.Frame {
	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		logrusPackage = getPackageName(runtime.FuncForPC(pcs[0]).Name())

		// now that we have the cache, we can skip a minimum count of known-logrus functions
		// XXX this is dubious, the number of frames may vary store an entry in a logger interface
		minimumCallerDepth = knownLogrusFrames
	})

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != logrusPackage {
			return &f
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// HasCaller reports if the given log entry has caller data
func (entry Entry) HasCaller() (has bool) {
	return entry.Logger != nil &&
		entry.Logger.ReportCaller &&
		entry.Caller != nil
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
	if entry.Logger.ReportCaller {
		entry.Caller = getCaller()
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

// LogAtLevel logs message at a given log level.
func (entry *Entry) LogAtLevel(level Level, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(level) {
		entry.log(level, fmt.Sprint(args...))
	}
	switch level {
	case FatalLevel:
		entry.Logger.Exit(1)
	case PanicLevel:
		panic(fmt.Sprint(args...))
	}
}

// LogfAtLevel logs message at a given log level.
func (entry *Entry) LogfAtLevel(level Level, format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(level) {
		entry.LogAtLevel(level, fmt.Sprintf(format, args...))
	}
	switch level {
	case FatalLevel:
		entry.Logger.Exit(1)
	}
}

// LoglnAtLevel logs message at a given log level.
func (entry *Entry) LoglnAtLevel(level Level, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(level) {
		entry.LogAtLevel(level, entry.sprintlnn(args...))
	}
	switch level {
	case FatalLevel:
		entry.Logger.Exit(1)
	}
}

// Trace logs a trace message.
func (entry *Entry) Trace(args ...interface{}) {
	entry.LogAtLevel(TraceLevel, args...)
}

// Debug logs a debug message.
func (entry *Entry) Debug(args ...interface{}) {
	entry.LogAtLevel(DebugLevel, args...)
}

// Print logs an info message.
func (entry *Entry) Print(args ...interface{}) {
	entry.Info(args...)
}

// Info logs an info message.
func (entry *Entry) Info(args ...interface{}) {
	entry.LogAtLevel(InfoLevel, args...)
}

// Warn logs a warning message.
func (entry *Entry) Warn(args ...interface{}) {
	entry.LogAtLevel(WarnLevel, args...)
}

// Warning logs a warning message.
func (entry *Entry) Warning(args ...interface{}) {
	entry.Warn(args...)
}

// Error logs an error message.
func (entry *Entry) Error(args ...interface{}) {
	entry.LogAtLevel(ErrorLevel, args...)
}

// Fatal logs a fatal error message, then exits.
func (entry *Entry) Fatal(args ...interface{}) {
	entry.LogAtLevel(FatalLevel, args...)
}

// Panic logs a panic message, then panics the *Entry.
func (entry *Entry) Panic(args ...interface{}) {
	entry.LogAtLevel(PanicLevel, args...)
}

// Entry Printf family functions

// Tracef logs a trace message.
func (entry *Entry) Tracef(format string, args ...interface{}) {
	entry.LogfAtLevel(TraceLevel, format, args...)
}

// Debugf logs a debug message.
func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.LogfAtLevel(DebugLevel, format, args...)
}

// Infof logs an info message.
func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.LogfAtLevel(InfoLevel, format, args...)
}

// Printf logs an info message.
func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.Infof(format, args...)
}

// Warnf logs a warning message.
func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.LogfAtLevel(WarnLevel, format, args...)
}

// Warningf logs a warning message.
func (entry *Entry) Warningf(format string, args ...interface{}) {
	entry.Warnf(format, args...)
}

// Errorf logs an error message.
func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.LogfAtLevel(ErrorLevel, format, args...)
}

// Fatalf logs a fatal error message, then exits.
func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.LogfAtLevel(FatalLevel, format, args...)
}

// Panicf logs a panic message, then panics the *Entry.
func (entry *Entry) Panicf(format string, args ...interface{}) {
	entry.LogfAtLevel(PanicLevel, format, args...)
}

// Entry Println family functions

// Traceln logs a trace message.
func (entry *Entry) Traceln(args ...interface{}) {
	entry.LoglnAtLevel(TraceLevel, entry.sprintlnn(args...))
}

// Debugln logs a debug message.
func (entry *Entry) Debugln(args ...interface{}) {
	entry.LoglnAtLevel(DebugLevel, entry.sprintlnn(args...))
}

// Infoln logs an info message.
func (entry *Entry) Infoln(args ...interface{}) {
	entry.LoglnAtLevel(InfoLevel, entry.sprintlnn(args...))
}

// Println logs an info message.
func (entry *Entry) Println(args ...interface{}) {
	entry.Infoln(args...)
}

// Warnln logs a warning message.
func (entry *Entry) Warnln(args ...interface{}) {
	entry.LoglnAtLevel(WarnLevel, entry.sprintlnn(args...))
}

// Warningln logs a warning message.
func (entry *Entry) Warningln(args ...interface{}) {
	entry.Warnln(args...)
}

// Errorln logs an error message.
func (entry *Entry) Errorln(args ...interface{}) {
	entry.LoglnAtLevel(ErrorLevel, entry.sprintlnn(args...))
}

// Fatalln logs a fatal error message, then exits.
func (entry *Entry) Fatalln(args ...interface{}) {
	entry.LoglnAtLevel(FatalLevel, entry.sprintlnn(args...))
}

// Panicln logs a panic message, then panics the *Entry.
func (entry *Entry) Panicln(args ...interface{}) {
	entry.LoglnAtLevel(PanicLevel, entry.sprintlnn(args...))
}

// Sprintlnn => Sprint no newline. This is to get the behavior of how
// fmt.Sprintln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func (entry *Entry) sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
