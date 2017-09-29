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
	Hooks LevelHooks
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
func NewLogger() *LogWriter {
	return &LogWriter{
		Out:       os.Stderr,
		Formatter: new(TextFormatter),
		Hooks:     make(LevelHooks),
		Level:     InfoLevel,
	}
}

func (logger *LogWriter) newEntry() *LogEntry {
	entry, ok := logger.entryPool.Get().(*LogEntry)
	if ok {
		return entry
	}
	return NewLogEntry(logger)
}

func (logger *LogWriter) releaseEntry(entry *LogEntry) {
	logger.entryPool.Put(entry)
}

func (logger *LogWriter) AsLevel(level Level) *LogEntry {
	return logger.newEntry()
}

func (logger *LogWriter) AsDebug() *LogEntry {
	return logger.AsLevel(DebugLevel)
}

func (logger *LogWriter) AsInfo() *LogEntry {
	return logger.AsLevel(InfoLevel)
}

func (logger *LogWriter) AsWarning() *LogEntry {
	return logger.AsLevel(WarnLevel)
}

func (logger *LogWriter) AsError() *LogEntry {
	return logger.AsLevel(ErrorLevel)
}

func (logger *LogWriter) AsFatal() *LogEntry {
	return logger.AsLevel(FatalLevel)
}

func (logger *LogWriter) AsPanic() *LogEntry {
	return logger.AsLevel(PanicLevel)
}

func (logger *LogWriter) Debugf(format string, args ...interface{}) {
	if logger.level() >= DebugLevel {
		entry := logger.newEntry()
		entry.Debugf(format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Infof(format string, args ...interface{}) {
	if logger.level() >= InfoLevel {
		entry := logger.newEntry()
		entry.Infof(format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Printf(format string, args ...interface{}) {
	entry := logger.newEntry()
	entry.Printf(format, args...)
	logger.releaseEntry(entry)
}

func (logger *LogWriter) Warnf(format string, args ...interface{}) {
	if logger.level() >= WarnLevel {
		entry := logger.newEntry()
		entry.Warnf(format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Warningf(format string, args ...interface{}) {
	if logger.level() >= WarnLevel {
		entry := logger.newEntry()
		entry.Warnf(format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Errorf(format string, args ...interface{}) {
	if logger.level() >= ErrorLevel {
		entry := logger.newEntry()
		entry.Errorf(format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Fatalf(format string, args ...interface{}) {
	if logger.level() >= FatalLevel {
		entry := logger.newEntry()
		entry.Fatalf(format, args...)
		logger.releaseEntry(entry)
	}
	Exit(1)
}

func (logger *LogWriter) Panicf(format string, args ...interface{}) {
	if logger.level() >= PanicLevel {
		entry := logger.newEntry()
		entry.Panicf(format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Debug(args ...interface{}) {
	if logger.level() >= DebugLevel {
		entry := logger.newEntry()
		entry.Debug(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Info(args ...interface{}) {
	if logger.level() >= InfoLevel {
		entry := logger.newEntry()
		entry.Info(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Print(args ...interface{}) {
	entry := logger.newEntry()
	entry.Info(args...)
	logger.releaseEntry(entry)
}

func (logger *LogWriter) Warn(args ...interface{}) {
	if logger.level() >= WarnLevel {
		entry := logger.newEntry()
		entry.Warn(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Warning(args ...interface{}) {
	if logger.level() >= WarnLevel {
		entry := logger.newEntry()
		entry.Warn(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Error(args ...interface{}) {
	if logger.level() >= ErrorLevel {
		entry := logger.newEntry()
		entry.Error(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Fatal(args ...interface{}) {
	if logger.level() >= FatalLevel {
		entry := logger.newEntry()
		entry.Fatal(args...)
		logger.releaseEntry(entry)
	}
	Exit(1)
}

func (logger *LogWriter) Panic(args ...interface{}) {
	if logger.level() >= PanicLevel {
		entry := logger.newEntry()
		entry.Panic(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Debugln(args ...interface{}) {
	if logger.level() >= DebugLevel {
		entry := logger.newEntry()
		entry.Debugln(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Infoln(args ...interface{}) {
	if logger.level() >= InfoLevel {
		entry := logger.newEntry()
		entry.Infoln(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Println(args ...interface{}) {
	entry := logger.newEntry()
	entry.Println(args...)
	logger.releaseEntry(entry)
}

func (logger *LogWriter) Warnln(args ...interface{}) {
	if logger.level() >= WarnLevel {
		entry := logger.newEntry()
		entry.Warnln(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Warningln(args ...interface{}) {
	if logger.level() >= WarnLevel {
		entry := logger.newEntry()
		entry.Warnln(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Errorln(args ...interface{}) {
	if logger.level() >= ErrorLevel {
		entry := logger.newEntry()
		entry.Errorln(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) Fatalln(args ...interface{}) {
	if logger.level() >= FatalLevel {
		entry := logger.newEntry()
		entry.Fatalln(args...)
		logger.releaseEntry(entry)
	}
	Exit(1)
}

func (logger *LogWriter) Panicln(args ...interface{}) {
	if logger.level() >= PanicLevel {
		entry := logger.newEntry()
		entry.Panicln(args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogWriter) WriteF(format string, args ...interface{}) {
	if !logger.hasDesiredStatus() {
		fmt.Fprint(os.Stderr, "Unknown log level. Call SetLevel or any of AsXYZ methods before calling Write methods")
		return
	}
	desired := logger.desiredLevel()
	if desired >= logger.level() {
		entry := logger.newEntry()
		entry.log(desired, fmt.Sprintf(format, args...))
		logger.releaseEntry(entry)
	}
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

func (logger *LogWriter) desiredLevel() Level {
	return Level(atomic.LoadUint32((*uint32)(&logger.asLevel)))
}

func (logger *LogWriter) setDesiredStatus(set bool) {
	var status int32
	if set {
		status = 1
	}
	atomic.StoreInt32(&logger.asLevelStatus, status)
}

func (logger *LogWriter) hasDesiredStatus() bool {
	return atomic.LoadInt32(&logger.asLevelStatus) == 1
}

func (logger *LogWriter) SetLevel(level Level) {
	atomic.StoreUint32((*uint32)(&logger.Level), uint32(level))
	logger.setDesiredLevel(level)
}

func (logger *LogWriter) setDesiredLevel(level Level) {
	atomic.StoreUint32((*uint32)(&logger.asLevel), uint32(level))
}

func (logger *LogWriter) AddHook(hook Hook) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Hooks.Add(hook)
}
