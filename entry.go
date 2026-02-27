package logrus

import (
	"bytes"
	"context"
	"fmt"
	"maps"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (

	// qualified package name, cached at first use
	logrusPackage string

	// Positions in the call stack when tracing to report the calling method.
	//
	// Start at the bottom of the stack before the package-name cache is primed.
	minimumCallerDepth = 1

	// Used for caller information initialisation
	callerInitOnce sync.Once
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4
)

// ErrorKey defines the key when adding errors using [WithError], [Logger.WithError].
var ErrorKey = "error"

// Entry is the final or intermediate Logrus logging entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Trace, Debug,
// Info, Warn, Error, Fatal or Panic is called on it. These objects can be
// reused and passed around as much as you wish to avoid field duplication.
//
//nolint:recvcheck // the methods of "Entry" use pointer receiver and non-pointer receiver.
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

	// Contains the context set by the user. Useful for hook processing etc.
	Context context.Context

	// err may contain a field formatting error
	err string
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		// Default is three fields, plus one optional.  Give a little extra room.
		Data: make(Fields, 6),
	}
}

func (entry *Entry) Dup() *Entry {
	return &Entry{
		Logger:  entry.Logger,
		Data:    maps.Clone(entry.Data),
		Time:    entry.Time,
		Context: entry.Context,
		err:     entry.err,
	}
}

// Bytes returns the bytes representation of this entry from the formatter.
func (entry *Entry) Bytes() ([]byte, error) {
	return entry.Logger.Formatter.Format(entry)
}

// String returns the string representation from the reader and ultimately the
// formatter.
func (entry *Entry) String() (string, error) {
	serialized, err := entry.Bytes()
	if err != nil {
		return "", err
	}
	str := string(serialized)
	return str, nil
}

// WithError adds an error as single field (using the key defined in [ErrorKey])
// to the Entry.
func (entry *Entry) WithError(err error) *Entry {
	return entry.WithField(ErrorKey, err)
}

// WithContext adds a context to the Entry.
func (entry *Entry) WithContext(ctx context.Context) *Entry {
	return &Entry{
		Logger:  entry.Logger,
		Data:    maps.Clone(entry.Data),
		Time:    entry.Time,
		Context: ctx,
		err:     entry.err,
	}
}

// WithField adds a single field to the Entry.
func (entry *Entry) WithField(key string, value any) *Entry {
	return entry.WithFields(Fields{key: value})
}

// WithFields adds a map of fields to the Entry.
func (entry *Entry) WithFields(fields Fields) *Entry {
	data := make(Fields, len(entry.Data)+len(fields))
	maps.Copy(data, entry.Data)
	fieldErr := entry.err
	for k, v := range fields {
		isErrField := false
		if t := reflect.TypeOf(v); t != nil {
			switch {
			case t.Kind() == reflect.Func, t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Func:
				isErrField = true
			}
		}
		if isErrField {
			tmp := fmt.Sprintf("can not add field %q", k)
			if fieldErr != "" {
				fieldErr = entry.err + ", " + tmp
			} else {
				fieldErr = tmp
			}
		} else {
			data[k] = v
		}
	}
	return &Entry{Logger: entry.Logger, Data: data, Time: entry.Time, err: fieldErr, Context: entry.Context}
}

// WithTime overrides the time of the Entry.
func (entry *Entry) WithTime(t time.Time) *Entry {
	return &Entry{
		Logger:  entry.Logger,
		Data:    maps.Clone(entry.Data),
		Time:    t,
		Context: entry.Context,
		err:     entry.err,
	}
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
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)

		// dynamic get the package name and the minimum caller depth
		for i := range maximumCallerDepth {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCaller") {
				logrusPackage = getPackageName(funcName)
				break
			}
		}

		minimumCallerDepth = knownLogrusFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

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

func (entry Entry) HasCaller() (has bool) {
	return entry.Logger != nil &&
		entry.Logger.ReportCaller &&
		entry.Caller != nil
}

func (entry *Entry) log(level Level, msg string) {
	var buffer *bytes.Buffer

	newEntry := entry.Dup()

	if newEntry.Time.IsZero() {
		newEntry.Time = time.Now()
	}

	newEntry.Level = level
	newEntry.Message = msg

	newEntry.Logger.mu.Lock()
	reportCaller := newEntry.Logger.ReportCaller
	bufPool := newEntry.getBufferPool()
	newEntry.Logger.mu.Unlock()

	if reportCaller {
		newEntry.Caller = getCaller()
	}

	newEntry.fireHooks()
	buffer = bufPool.Get()
	defer func() {
		newEntry.Buffer = nil
		buffer.Reset()
		bufPool.Put(buffer)
	}()
	buffer.Reset()
	newEntry.Buffer = buffer

	newEntry.write()

	newEntry.Buffer = nil

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if level <= PanicLevel {
		panic(newEntry)
	}
}

func (entry *Entry) getBufferPool() (pool BufferPool) {
	if entry.Logger.BufferPool != nil {
		return entry.Logger.BufferPool
	}
	return bufferPool
}

func (entry *Entry) fireHooks() {
	entry.Logger.mu.Lock()
	tmpHooks := maps.Clone(entry.Logger.Hooks)
	entry.Logger.mu.Unlock()
	if len(tmpHooks) == 0 {
		return
	}

	if err := tmpHooks.Fire(entry.Level, entry); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to fire hook:", err)
	}
}

func (entry *Entry) write() {
	// Snapshot the formatter and output under the lock to protect against
	// concurrent SetFormatter/SetOutput calls, then release the lock before
	// formatting. This avoids a deadlock when Format() triggers reentrant
	// logging (e.g., a field's MarshalJSON calls logrus). See #1448, #1440.
	entry.Logger.mu.Lock()
	formatter := entry.Logger.Formatter
	out := entry.Logger.Out
	entry.Logger.mu.Unlock()

	serialized, err := formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		return
	}

	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	if _, err := out.Write(serialized); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}
}

// Log will log a message at the level given as parameter.
// Warning: using Log at Panic or Fatal level will not respectively Panic nor Exit.
// For this behaviour Entry.Panic or Entry.Fatal should be used instead.
func (entry *Entry) Log(level Level, args ...any) {
	if entry.Logger.IsLevelEnabled(level) {
		entry.log(level, fmt.Sprint(args...))
	}
}

func (entry *Entry) Trace(args ...any) {
	entry.Log(TraceLevel, args...)
}

func (entry *Entry) Debug(args ...any) {
	entry.Log(DebugLevel, args...)
}

func (entry *Entry) Print(args ...any) {
	entry.Info(args...)
}

func (entry *Entry) Info(args ...any) {
	entry.Log(InfoLevel, args...)
}

func (entry *Entry) Warn(args ...any) {
	entry.Log(WarnLevel, args...)
}

func (entry *Entry) Warning(args ...any) {
	entry.Warn(args...)
}

func (entry *Entry) Error(args ...any) {
	entry.Log(ErrorLevel, args...)
}

func (entry *Entry) Fatal(args ...any) {
	entry.Log(FatalLevel, args...)
	entry.Logger.Exit(1)
}

func (entry *Entry) Panic(args ...any) {
	entry.Log(PanicLevel, args...)
}

// Entry Printf family functions

func (entry *Entry) Logf(level Level, format string, args ...any) {
	if entry.Logger.IsLevelEnabled(level) {
		entry.Log(level, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Tracef(format string, args ...any) {
	entry.Logf(TraceLevel, format, args...)
}

func (entry *Entry) Debugf(format string, args ...any) {
	entry.Logf(DebugLevel, format, args...)
}

func (entry *Entry) Infof(format string, args ...any) {
	entry.Logf(InfoLevel, format, args...)
}

func (entry *Entry) Printf(format string, args ...any) {
	entry.Infof(format, args...)
}

func (entry *Entry) Warnf(format string, args ...any) {
	entry.Logf(WarnLevel, format, args...)
}

func (entry *Entry) Warningf(format string, args ...any) {
	entry.Warnf(format, args...)
}

func (entry *Entry) Errorf(format string, args ...any) {
	entry.Logf(ErrorLevel, format, args...)
}

func (entry *Entry) Fatalf(format string, args ...any) {
	entry.Logf(FatalLevel, format, args...)
	entry.Logger.Exit(1)
}

func (entry *Entry) Panicf(format string, args ...any) {
	entry.Logf(PanicLevel, format, args...)
}

// Entry Println family functions

func (entry *Entry) Logln(level Level, args ...any) {
	if entry.Logger.IsLevelEnabled(level) {
		entry.Log(level, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Traceln(args ...any) {
	entry.Logln(TraceLevel, args...)
}

func (entry *Entry) Debugln(args ...any) {
	entry.Logln(DebugLevel, args...)
}

func (entry *Entry) Infoln(args ...any) {
	entry.Logln(InfoLevel, args...)
}

func (entry *Entry) Println(args ...any) {
	entry.Infoln(args...)
}

func (entry *Entry) Warnln(args ...any) {
	entry.Logln(WarnLevel, args...)
}

func (entry *Entry) Warningln(args ...any) {
	entry.Warnln(args...)
}

func (entry *Entry) Errorln(args ...any) {
	entry.Logln(ErrorLevel, args...)
}

func (entry *Entry) Fatalln(args ...any) {
	entry.Logln(FatalLevel, args...)
	entry.Logger.Exit(1)
}

func (entry *Entry) Panicln(args ...any) {
	entry.Logln(PanicLevel, args...)
}

// sprintlnn => Sprint no newline. This is to get the behavior of how
// fmt.Sprintln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func (entry *Entry) sprintlnn(args ...any) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
