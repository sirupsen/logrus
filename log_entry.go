package logrus

import (
	"bytes"
	"fmt"
	"os"
	"sync/atomic"
	"time"
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
	// This is to address the scenario's in which LogEntry has not been created either
	// by the Logger or by calling NewLogEntry method. Because Level does not have
	// an Unknown state, we won't be able to determine what level we should write the logs at
	// In that case we will log an error to 'Stderr' and move on without logging
	//1 = level has been set, 0: level has not been set (Unknown)
	levelStatus int32

	// Message passed to Debug, Info, Warn, Error, Fatal or Panic
	Message string

	// When formatter is called in entry.log(), an Buffer may be set to entry
	Buffer *bytes.Buffer
}

func NewLogEntry(logger *LogWriter) *LogEntry {
	// Default is three fields, give a little extra room
	return newLogEntry(logger, make(Fields, 5))
}

func newLogEntry(logger *LogWriter, data Fields) *LogEntry {
	entry := &LogEntry{
		Logger: logger,
		Data:   data,
	}
	entry.setLevel(logger.Level)
	return entry
}

func (entry *LogEntry) AsLevel(level Level) *LogEntry {
	entry.setLevel(level)
	return entry
}

func (entry *LogEntry) AsDebug() *LogEntry {
	return entry.AsLevel(DebugLevel)
}

func (entry *LogEntry) AsInfo() *LogEntry {
	return entry.AsLevel(InfoLevel)
}

func (entry *LogEntry) AsWarning() *LogEntry {
	return entry.AsLevel(WarnLevel)
}

func (entry *LogEntry) AsError() *LogEntry {
	return entry.AsLevel(ErrorLevel)
}

func (entry *LogEntry) AsFatal() *LogEntry {
	return entry.AsLevel(FatalLevel)
}

func (entry *LogEntry) AsPanic() *LogEntry {
	return entry.AsLevel(PanicLevel)
}

func (entry *LogEntry) WithField(key string, value interface{}) *LogEntry {
	return entry.WithFields(Fields{key: value})
}

func (entry *LogEntry) WithFields(fields Fields) *LogEntry {
	level, ok := entry.getLevel()
	if !ok || level < entry.Logger.level() {
		return entry
	}
	data := make(Fields, len(entry.Data)+len(fields))
	for k, v := range entry.Data {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}
	return newLogEntry(entry.Logger, data)
}

func (entry *LogEntry) WithError(err error) *LogEntry {
	return entry.WithField(ErrorKey, err)
}

func (entry *LogEntry) Writef(format string, args ...interface{}) {
	entry.write(fmt.Sprintf(format, args...))
}

func (entry *LogEntry) Write(args ...interface{}) {
	entry.write(fmt.Sprint(args...))
}

func (entry *LogEntry) Writeln(args ...interface{}) {
	entry.write(sprintlnn(args...))
}

func (entry *LogEntry) write(message string) {
	level, ok := entry.checkLevel()
	if !ok {
		return
	}

	loggerLevel := entry.Logger.level()
	if loggerLevel >= level {
		entry.log(level, message)
	}

	if loggerLevel != level {
		//reset the LogEntry's level to the Logger's level value
		entry.setLevel(entry.Logger.Level)
	}
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

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if level <= PanicLevel {
		panic(&entry)
	}
}

func (entry *LogEntry) checkLevel() (Level, bool) {
	level, ok := entry.getLevel()
	if !ok {
		fmt.Fprintln(os.Stderr, "Unknown logging level. Call AsLevel before using the entry")
	}
	return level, ok
}

func (entry *LogEntry) setLevel(level Level) {
	atomic.StoreUint32((*uint32)(&entry.Level), uint32(level))
	atomic.StoreInt32((*int32)(&entry.levelStatus), 1)
}

func (entry *LogEntry) getLevel() (Level, bool) {
	level := Level(atomic.LoadUint32((*uint32)(&entry.Level)))
	ok := atomic.LoadInt32((*int32)(&entry.levelStatus)) == 1
	return level, ok
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
