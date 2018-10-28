package logrus

import (
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Logger object to log to.
type Logger struct {
	// The logs are `io.Copy`'d to this in a mutex. It's common to set this to a
	// file, or leave it default which is `os.Stderr`. You can also set this to
	// something more adventurous, such as logging to Kafka.
	Out io.Writer
	// Hooks for the logger instance. These allow firing events based on logging
	// levels and log entries. For example, to send errors to an error tracking
	// service, log to StatsD or dump the core on fatal errors.
	Hooks LevelHooks
	// All log entries pass through the formatter before logged to Out. The
	// included formatters are `TextFormatter` and `JSONFormatter` for which
	// TextFormatter is the default. In development (when a TTY is attached) it
	// logs with colors, but to a file it wouldn't. You can easily implement your
	// own that implements the `Formatter` interface, see the `README` or included
	// formatters for examples.
	Formatter Formatter

	// Flag for whether to log caller info (off by default)
	ReportCaller bool

	// The logging level the logger should log at. This is typically (and defaults
	// to) `logrus.Info`, which allows Info(), Warn(), Error() and Fatal() to be
	// logged.
	Level Level
	// Used to sync writing to the log. Locking is enabled by Default
	mu MutexWrap
	// Reusable empty entry
	entryPool sync.Pool
	// Function to exit the application, defaults to `os.Exit()`
	ExitFunc exitFunc
}

type exitFunc func(int)

// MutexWrap a disablable mutex
type MutexWrap struct {
	lock     sync.Mutex
	disabled bool
}

// Lock the mutex
func (mw *MutexWrap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

// Unlock the mutex
func (mw *MutexWrap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}

// Disable the mutex
func (mw *MutexWrap) Disable() {
	mw.disabled = true
}

// New creates a new logger. Configuration should be set by changing `Formatter`,
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
func New() *Logger {
	return &Logger{
		Out:          os.Stderr,
		Formatter:    new(TextFormatter),
		Hooks:        make(LevelHooks),
		Level:        InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
}

func (logger *Logger) newEntry() *Entry {
	entry, ok := logger.entryPool.Get().(*Entry)
	if ok {
		return entry
	}
	return NewEntry(logger)
}

func (logger *Logger) releaseEntry(entry *Entry) {
	entry.Data = map[string]interface{}{}
	logger.entryPool.Put(entry)
}

// WithField adds a field to the log entry, note that it doesn't log until you call
// Debug, Print, Info, Warn, Error, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func (logger *Logger) WithField(key string, value interface{}) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithField(key, value)
}

// WithFields adds a struct of fields to the log entry. All it does is call
// `WithField` for each `Field`.
func (logger *Logger) WithFields(fields Fields) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithFields(fields)
}

// WithError adds an error as single field to the log entry.  All it does is
// call `WithError` for the given `error`.
func (logger *Logger) WithError(err error) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithError(err)
}

// WithTime overrides the time of the log entry.
func (logger *Logger) WithTime(t time.Time) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithTime(t)
}

// LogfAtLevel logs a message at given level on the standard logger.
func (logger *Logger) LogfAtLevel(level Level, format string, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.LogfAtLevel(level, format, args...)
		logger.releaseEntry(entry)
	}
	switch level {
	case FatalLevel:
		logger.Exit(1)
	}
}

// LogAtLevel logs a message at given level on the standard logger.
func (logger *Logger) LogAtLevel(level Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.LogAtLevel(level, args...)
		logger.releaseEntry(entry)
	}
	switch level {
	case FatalLevel:
		logger.Exit(1)
	}
}

// LoglnAtLevel logs a message at given level on the standard logger.
func (logger *Logger) LoglnAtLevel(level Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.LoglnAtLevel(level, args...)
		logger.releaseEntry(entry)
	}
	switch level {
	case FatalLevel:
		logger.Exit(1)
	}
}

// Tracef logs a message at level Trace on the standard logger.
func (logger *Logger) Tracef(format string, args ...interface{}) {
	logger.LogfAtLevel(TraceLevel, format, args...)
}

// Debugf logs a message at level Debug on the standard logger.
func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.LogfAtLevel(DebugLevel, format, args...)
}

// Infof logs a message at level Info on the standard logger.
func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.LogfAtLevel(InfoLevel, format, args...)
}

// Printf logs a message at level Info on the standard logger.
func (logger *Logger) Printf(format string, args ...interface{}) {
	entry := logger.newEntry()
	entry.Printf(format, args...)
	logger.releaseEntry(entry)
}

// Warnf logs a message at level Warn on the standard logger.
func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.LogfAtLevel(WarnLevel, format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func (logger *Logger) Warningf(format string, args ...interface{}) {
	logger.LogfAtLevel(WarnLevel, format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.LogfAtLevel(ErrorLevel, format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.LogfAtLevel(FatalLevel, format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.LogfAtLevel(PanicLevel, format, args...)
}

// Trace logs a message at level Trace on the standard logger.
func (logger *Logger) Trace(args ...interface{}) {
	logger.LogAtLevel(TraceLevel, args...)
}

// Debug logs a message at level Debug on the standard logger.
func (logger *Logger) Debug(args ...interface{}) {
	logger.LogAtLevel(DebugLevel, args...)
}

// Info logs a message at level Info on the standard logger.
func (logger *Logger) Info(args ...interface{}) {
	logger.LogAtLevel(InfoLevel, args...)
}

// Print logs a message at level Info on the standard logger.
func (logger *Logger) Print(args ...interface{}) {
	entry := logger.newEntry()
	entry.Info(args...)
	logger.releaseEntry(entry)
}

// Warn logs a message at level Warn on the standard logger.
func (logger *Logger) Warn(args ...interface{}) {
	logger.LogAtLevel(WarnLevel, args...)
}

// Warning logs a message at level Warn on the standard logger.
func (logger *Logger) Warning(args ...interface{}) {
	logger.LogAtLevel(WarnLevel, args...)
}

// Error logs a message at level Error on the standard logger.
func (logger *Logger) Error(args ...interface{}) {
	logger.LogAtLevel(ErrorLevel, args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func (logger *Logger) Fatal(args ...interface{}) {
	logger.LogAtLevel(FatalLevel, args...)
}

// Panic logs a message at level Panic on the standard logger.
func (logger *Logger) Panic(args ...interface{}) {
	logger.LogAtLevel(PanicLevel, args...)
}

// Traceln logs a message at level Trace on the standard logger.
func (logger *Logger) Traceln(args ...interface{}) {
	logger.LoglnAtLevel(TraceLevel, args...)
}

// Debugln logs a message at level Debug on the standard logger.
func (logger *Logger) Debugln(args ...interface{}) {
	logger.LoglnAtLevel(DebugLevel, args...)
}

// Infoln logs a message at level Info on the standard logger.
func (logger *Logger) Infoln(args ...interface{}) {
	logger.LoglnAtLevel(InfoLevel, args...)
}

// Println logs a message at level Info on the standard logger.
func (logger *Logger) Println(args ...interface{}) {
	entry := logger.newEntry()
	entry.Println(args...)
	logger.releaseEntry(entry)
}

// Warnln logs a message at level Warn on the standard logger.
func (logger *Logger) Warnln(args ...interface{}) {
	logger.LoglnAtLevel(WarnLevel, args...)
}

// Warningln logs a message at level Warn on the standard logger.
func (logger *Logger) Warningln(args ...interface{}) {
	logger.LoglnAtLevel(WarnLevel, args...)
}

// Errorln logs a message at level Error on the standard logger.
func (logger *Logger) Errorln(args ...interface{}) {
	logger.LoglnAtLevel(ErrorLevel, args...)
}

// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func (logger *Logger) Fatalln(args ...interface{}) {
	logger.LoglnAtLevel(FatalLevel, args...)
}

// Panicln logs a message at level Panic on the standard logger.
func (logger *Logger) Panicln(args ...interface{}) {
	logger.LoglnAtLevel(PanicLevel, args...)
}

// Exit calls os.Exit (or logger.ExitFunc) after running handlers.
func (logger *Logger) Exit(code int) {
	runHandlers()
	if logger.ExitFunc == nil {
		logger.ExitFunc = os.Exit
	}
	logger.ExitFunc(code)
}

// SetNoLock allows user to disable locking the log file.
// When file is opened with appending mode, it's safe to
// write concurrently to a file (within 4k message on Linux).
// In these cases user can choose to disable the lock.
func (logger *Logger) SetNoLock() {
	logger.mu.Disable()
}

func (logger *Logger) level() Level {
	return Level(atomic.LoadUint32((*uint32)(&logger.Level)))
}

// SetLevel sets the logger level.
func (logger *Logger) SetLevel(level Level) {
	atomic.StoreUint32((*uint32)(&logger.Level), uint32(level))
}

// GetLevel returns the logger level.
func (logger *Logger) GetLevel() Level {
	return logger.level()
}

// AddHook adds a hook to the logger hooks.
func (logger *Logger) AddHook(hook Hook) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Hooks.Add(hook)
}

// IsLevelEnabled checks if the log level of the logger is greater than the level param
func (logger *Logger) IsLevelEnabled(level Level) bool {
	return logger.level() >= level
}

// SetFormatter sets the logger formatter.
func (logger *Logger) SetFormatter(formatter Formatter) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Formatter = formatter
}

// SetOutput sets the logger output.
func (logger *Logger) SetOutput(output io.Writer) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Out = output
}

// SetReportCaller enables/disables reporting of caller data
func (logger *Logger) SetReportCaller(reportCaller bool) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.ReportCaller = reportCaller
}

// ReplaceHooks replaces the logger hooks and returns the old ones
func (logger *Logger) ReplaceHooks(hooks LevelHooks) LevelHooks {
	logger.mu.Lock()
	oldHooks := logger.Hooks
	logger.Hooks = hooks
	logger.mu.Unlock()
	return oldHooks
}
