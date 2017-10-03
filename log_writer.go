package logrus

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

type LogWriter struct {
	// The logs are `io.Copy`'d to this in a mutex. It's common to set this to a
	// file, or leave it default which is `os.Stderr`. You can also set this to
	// something more adventorous, such as logging to Kafka.
	Out io.Writer
	// Hooks for the logger instance. These allow firing events based on logging
	// levels and log entries. For example, to send errors to an error tracking
	// service, log to StatsD or dump the core on fatal errors.
	Hooks ExternalLevelHooks
	// All log entries pass through the formatter before logged to Out. The
	// included formatters are `TextFormatter` and `JSONFormatter` for which
	// TextFormatter is the default. In development (when a TTY is attached) it
	// logs with colors, but to a file it wouldn't. You can easily implement your
	// own that implements the `Formatter` interface, see the `README` or included
	// formatters for examples.
	Formatter Formatter
	// The logging level the logger should log at. This is typically (and defaults
	// to) `logrus.Info`, which allows Info(), Warn(), Error() and Fatal() to be
	// logged.
	Level Level

	// Used to sync writing to the log. Locking is enabled by Default
	mu MutexWrap
	// Reusable empty entry
	entryPool sync.Pool
}

// Creates a new logger. Configuration should be set by changing `Formatter`,
// `Out` and `Hooks` directly on the default logger instance. You can also just
// instantiate your own:
//
//    var log = &Logger{
//      Out: os.Stderr,
//      Formatter: new(JSONFormatter),
//      Hooks: make(LevelHooks),
//      Level: logrus.DebugLevel,
//    }
//
// It's recommended to make this a global instance called `log`.
func NewLogger(level Level) *LogWriter {
	return &LogWriter{
		Out:       os.Stderr,
		Formatter: new(JSONFormatter),
		Hooks:     make(ExternalLevelHooks),
		Level:     level,
	}
}

func (logger *LogWriter) newEntry(fields Fields) *LogEntry {
	entry, ok := logger.entryPool.Get().(*LogEntry)
	if ok {
		return entry
	}
	if fields != nil {
		return NewLogEntryWithFields(logger, fields)
	}
	return NewLogEntry(logger)
}

func (logger *LogWriter) releaseEntry(entry *LogEntry) {
	logger.entryPool.Put(entry)
}

func (logger *LogWriter) SetLevel(level Level) {
	atomic.StoreUint32((*uint32)(&logger.Level), uint32(level))
}

func (logger *LogWriter) Entry() *LogEntry {
	return logger.newEntry(nil)
}

func (logger *LogWriter) EntryWithFields(fields Fields) *LogEntry {
	return logger.newEntry(fields)
}

func (logger *LogWriter) EntryWithField(key string, value interface{}) *LogEntry {
	fields := Fields{key: value}
	return logger.newEntry(fields)
}

func (logger *LogWriter) logf(level Level, format string, args ...interface{}) {
	logger.log(level, fmt.Sprintf(format, args...))
}

func (logger *LogWriter) logln(level Level, args ...interface{}) {
	logger.log(level, sprintlnn(args...))
}

func (logger *LogWriter) log(level Level, message string) {
	if logger.level() >= level {
		entry := logger.newEntry(nil)
		entry.log(level, message)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Debugf(format string, args ...interface{}) {
	logger.logf(DebugLevel, format, args...)
}

func (logger *LogWriter) Infof(format string, args ...interface{}) {
	logger.logf(InfoLevel, format, args...)
}

func (logger *LogWriter) Warnf(format string, args ...interface{}) {
	logger.logf(WarnLevel, format, args...)
}

func (logger *LogWriter) Warningf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func (logger *LogWriter) Errorf(format string, args ...interface{}) {
	logger.logf(ErrorLevel, format, args...)
}

func (logger *LogWriter) Fatalf(format string, args ...interface{}) {
	logger.logf(FatalLevel, format, args...)
	Exit(1)
}

func (logger *LogWriter) Panicf(format string, args ...interface{}) {
	logger.logf(PanicLevel, format, args...)
	panic(fmt.Sprint(args...))
}

func (logger *LogWriter) Debug(args ...interface{}) {
	logger.log(DebugLevel, fmt.Sprint(args...))
}

func (logger *LogWriter) Info(args ...interface{}) {
	logger.log(InfoLevel, fmt.Sprint(args...))
}

func (logger *LogWriter) Warn(args ...interface{}) {
	logger.log(WarnLevel, fmt.Sprint(args...))
}

func (logger *LogWriter) Warning(args ...interface{}) {
	logger.Warn(args...)
}

func (logger *LogWriter) Error(args ...interface{}) {
	logger.log(ErrorLevel, fmt.Sprint(args...))
}

func (logger *LogWriter) Fatal(args ...interface{}) {
	logger.log(FatalLevel, fmt.Sprint(args...))
	Exit(1)
}

func (logger *LogWriter) Panic(args ...interface{}) {
	msg := fmt.Sprint(args...)
	logger.log(PanicLevel, msg)
	panic(msg)
}

func (logger *LogWriter) Debugln(args ...interface{}) {
	logger.logln(DebugLevel, args...)
}

func (logger *LogWriter) Infoln(args ...interface{}) {
	logger.logln(InfoLevel, args...)
}

func (logger *LogWriter) Warnln(args ...interface{}) {
	logger.logln(WarnLevel, args...)
}

func (logger *LogWriter) Warningln(args ...interface{}) {
	logger.logln(WarnLevel, args...)
}

func (logger *LogWriter) Errorln(args ...interface{}) {
	logger.logln(ErrorLevel, args...)
}

func (logger *LogWriter) Fatalln(args ...interface{}) {
	logger.logln(FatalLevel, args...)
	Exit(1)
}

func (logger *LogWriter) Panicln(args ...interface{}) {
	msg := fmt.Sprint(args...)
	logger.logln(PanicLevel, msg)
	panic(msg)
}

//When file is opened with appending mode, it's safe to
//write concurrently to a file (within 4k message on Linux).
//In these cases user can choose to disable the lock.
func (logger *LogWriter) SetNoLock() {
	logger.mu.Disable()
}

func (logger *LogWriter) level() Level {
	return Level(atomic.LoadUint32((*uint32)(&logger.Level)))
}

func (logger *LogWriter) AddHook(hook ExternalHook) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Hooks.Add(hook)
}
