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
	// something more adventurous, such as logging to Kafka.
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

// Creates a new log writer. Configuration should be set by changing `Formatter` (default TextFormatter),
// `Out` (default os.Stderr) and `Hooks` directly on the default logger instance.
// It's recommended to make this a global instance called `log`.
func NewLogger(level Level) *LogWriter {
	return &LogWriter{
		Out: os.Stderr,
		Formatter: &TextFormatter{
			DisableSorting: true,
		},
		Hooks: make(ExternalLevelHooks),
		Level: level,
	}
}

// SetLevel sets the log writer's log level
func (logger *LogWriter) SetLevel(level Level) {
	atomic.StoreUint32((*uint32)(&logger.Level), uint32(level))
}

// WithFields adds a struct of fields to the log entry
func (logger *LogWriter) WithFields(fields Fields) *LogEntry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithFields(fields)
}

// WithField adds a field to the log entry, note that it doesn't log until you call Write.
// It only creates a log entry. If you want multiple fields, use `WithFields`.
func (logger *LogWriter) WithField(key string, value interface{}) *LogEntry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithField(key, value)
}

func (logger *LogWriter) releaseEntry(entry *LogEntry) {
	logger.entryPool.Put(entry)
}

func (logger *LogWriter) log(level Level, mode formatMode, format string, args ...interface{}) {
	if logger.level() >= level {
		entry := logger.newEntry()
		message := constructMessage(mode, format, args...)
		entry.log(level, message)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) newEntry() *LogEntry {
	entry, ok := logger.entryPool.Get().(*LogEntry)
	if ok {
		return entry
	}
	return NewLogEntry(logger)
}

func (logger *LogWriter) Debugf(format string, args ...interface{}) {
	logger.log(DebugLevel, formatted, format, args...)
}

func (logger *LogWriter) Infof(format string, args ...interface{}) {
	logger.log(InfoLevel, formatted, format, args...)
}

func (logger *LogWriter) Warningf(format string, args ...interface{}) {
	logger.log(WarnLevel, formatted, format, args...)
}

func (logger *LogWriter) Errorf(format string, args ...interface{}) {
	logger.log(ErrorLevel, formatted, format, args...)
}

func (logger *LogWriter) Fatalf(format string, args ...interface{}) {
	logger.log(FatalLevel, formatted, format, args...)
	Exit(1)
}

func (logger *LogWriter) Panicf(format string, args ...interface{}) {
	logger.log(PanicLevel, formatted, format, args...)
	panic(fmt.Sprint(args...))
}

func (logger *LogWriter) Debug(args ...interface{}) {
	logger.log(DebugLevel, unformatted, "", args...)
}

func (logger *LogWriter) Info(args ...interface{}) {
	logger.log(InfoLevel, unformatted, "", args...)
}

func (logger *LogWriter) Warning(args ...interface{}) {
	logger.log(WarnLevel, unformatted, "", args...)
}

func (logger *LogWriter) Error(args ...interface{}) {
	logger.log(ErrorLevel, unformatted, "", args...)
}

func (logger *LogWriter) Fatal(args ...interface{}) {
	logger.log(FatalLevel, unformatted, "", args...)
	Exit(1)
}

func (logger *LogWriter) Panic(args ...interface{}) {
	msg := fmt.Sprint(args...)
	logger.log(PanicLevel, unformatted, "", args...)
	panic(msg)
}

func (logger *LogWriter) Debugln(args ...interface{}) {
	logger.log(DebugLevel, newLine, "", args...)
}

func (logger *LogWriter) Infoln(args ...interface{}) {
	logger.log(InfoLevel, newLine, "", args...)
}

func (logger *LogWriter) Warningln(args ...interface{}) {
	logger.log(WarnLevel, newLine, "", args...)
}

func (logger *LogWriter) Errorln(args ...interface{}) {
	logger.log(ErrorLevel, newLine, "", args...)
}

func (logger *LogWriter) Fatalln(args ...interface{}) {
	logger.log(FatalLevel, newLine, "", args...)
	Exit(1)
}

func (logger *LogWriter) Panicln(args ...interface{}) {
	msg := fmt.Sprint(args...)
	logger.log(PanicLevel, newLine, "", msg)
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
