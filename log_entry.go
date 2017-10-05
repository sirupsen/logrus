package logrus

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

type formatMode int

const (
	formatted formatMode = iota
	unformatted
	newLine
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
	Level Level

	// Message passed to Debug, Info, Warn, Error, Fatal or Panic
	Message string

	// When formatter is called in entry.log(), an Buffer may be set to entry
	Buffer *bytes.Buffer
}

// NewLogEntry creates a new log entry
func NewLogEntry(logger *LogWriter) *LogEntry {
	// Default is three fields, give a little extra room
	return NewLogEntryWithFields(logger, make(Fields, 5))
}

// NewLogEntryWithFields creates a new log entry and adds a struct of fields to the entry
func NewLogEntryWithFields(logger *LogWriter, fields Fields) *LogEntry {
	return newLogEntry(logger, fields)
}

// NewLogEntryWithField creates a new log entry and adds a field to the entry
//If you want multiple fields, use `NewLogEntryWithFields`
func NewLogEntryWithField(logger *LogWriter, key string, value interface{}) *LogEntry {
	//Do not change this to Fields{key:value}. You will end up getting more allocations
	fields := make(Fields, 1)
	fields[key] = value
	return newLogEntry(logger, fields)
}

func newLogEntry(logger *LogWriter, data Fields) *LogEntry {
	return &LogEntry{
		Logger: logger,
		Data:   data,
		Level:  logger.Level,
	}
}

func (entry *LogEntry) cloneAs(level Level) *LogEntry {
	return &LogEntry{
		Logger: entry.Logger,
		Data:   entry.Data,
		Level:  level,
	}
}

// AsLevel clones the entry into a new log entry and sets the level to the specified value.
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (entry *LogEntry) AsLevel(level Level) *LogEntry {
	return entry.cloneAs(level)
}

// AsDebug clones the entry into a new log entry and sets the level to `debug`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (entry *LogEntry) AsDebug() *LogEntry {
	return entry.AsLevel(DebugLevel)
}

// AsInfo clones the entry into a new log entry and sets the level to `info`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (entry *LogEntry) AsInfo() *LogEntry {
	return entry.AsLevel(InfoLevel)
}

// AsWarning clones the entry into a new log entry and sets the level to `warning`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (entry *LogEntry) AsWarning() *LogEntry {
	return entry.AsLevel(WarnLevel)
}

// AsError clones the entry into a new log entry and sets the level to `error`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (entry *LogEntry) AsError() *LogEntry {
	return entry.AsLevel(ErrorLevel)
}

// AsFatal clones the entry into a new log entry and sets the level to `fatal`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (entry *LogEntry) AsFatal() *LogEntry {
	return entry.AsLevel(FatalLevel)
}

// AsPanic clones the entry into a new log entry and sets the level to `panic`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (entry *LogEntry) AsPanic() *LogEntry {
	return entry.AsLevel(PanicLevel)
}

// WithField adds a field to the log entry, note that it doesn't log until you call Write.
func (entry *LogEntry) WithField(key string, value interface{}) *LogEntry {
	if entry.Level > entry.Logger.level() {
		return entry
	}
	//Do not change this to Fields{key:value}. You will end up getting more allocations
	fields := make(Fields, 1)
	fields[key] = value
	return entry.WithFields(fields)
}

// WithFields adds a struct of fields to the log entry
func (entry *LogEntry) WithFields(fields Fields) *LogEntry {
	if entry.Level > entry.Logger.level() {
		return entry
	}
	data := make(Fields, len(entry.Data)+len(fields))
	for k, v := range entry.Data {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}
	return &LogEntry{
		Data:   data,
		Level:  entry.Level,
		Logger: entry.Logger,
	}
}

// WithError adds an error as single field to the log entry
func (entry *LogEntry) WithError(err error) *LogEntry {
	return entry.WithField(ErrorKey, err)
}

func (entry *LogEntry) Writef(format string, args ...interface{}) {
	entry.write(formatted, format, args...)
}

func (entry *LogEntry) Write(args ...interface{}) {
	entry.write(unformatted, "", args...)
}

func (entry *LogEntry) Writeln(args ...interface{}) {
	entry.write(newLine, "", args...)
}

func (entry *LogEntry) write(mode formatMode, format string, args ...interface{}) {
	if entry.Logger.level() >= entry.Level {
		message := constructMessage(mode, format, args...)
		entry.log(entry.Level, message)
	}
}

func constructMessage(mode formatMode, format string, args ...interface{}) string {
	switch mode {
	case formatted:
		return fmt.Sprintf(format, args...)
	case unformatted:
		return fmt.Sprint(args...)
	case newLine:
		return sprintlnn(args...)
	}
	return fmt.Sprintf(format, args...)
}

// This function is not declared with a pointer value because otherwise
// race conditions will occur when using multiple goroutines
func (entry LogEntry) log(level Level, msg string) {
	var buffer *bytes.Buffer
	entry.Time = time.Now()
	entry.Level = level
	entry.Message = msg

	entry.Logger.mu.Lock()
	err := entry.Logger.Hooks.Fire(level, &entry)
	entry.Logger.mu.Unlock()
	if err != nil {
		entry.Logger.mu.Lock()
		fmt.Fprintf(os.Stderr, "Failed to fire hook: %v\n", err)
		entry.Logger.mu.Unlock()
	}
	buffer = bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)
	entry.Buffer = buffer
	serialized, err := entry.Logger.Formatter.FormatEntry(&entry)
	entry.Buffer = nil
	if err != nil {
		entry.Logger.mu.Lock()
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		entry.Logger.mu.Unlock()
	} else {
		entry.Logger.mu.Lock()
		_, err = entry.Logger.Out.Write(serialized)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
		}
		entry.Logger.mu.Unlock()
	}

	if level == FatalLevel {
		Exit(1)
	}

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if level <= PanicLevel {
		panic(&entry)
	}
}

// Returns the string representation from the reader and ultimately the
// formatter.
func (entry *LogEntry) String() (string, error) {
	serialized, err := entry.Logger.Formatter.FormatEntry(entry)
	if err != nil {
		return "", err
	}
	str := string(serialized)
	return str, nil
}
