package logrus

import (
	"io"
	"os"
	"sync"
)

type Logger struct {
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
	// logged. `logrus.Debug` is useful in
	Level Level
	// Used to sync writing to the log.
	mu sync.Mutex

	showCaller bool
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
func New() *Logger {
	return &Logger{
		Out:        os.Stderr,
		Formatter:  new(TextFormatter),
		Hooks:      make(LevelHooks),
		Level:      InfoLevel,
		showCaller: true,
	}
}

// Adds a field to the log entry, note that you it doesn't log until you call
// Debug, Print, Info, Warn, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func (logger *Logger) WithField(key string, value interface{}) *Entry {
	return newEntry(logger, 4).WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (logger *Logger) WithFields(fields Fields) *Entry {
	return newEntry(logger, 4).WithFields(fields)
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (logger *Logger) WithError(err error) *Entry {
	return newEntry(logger, 3).WithError(err)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	if logger.Level >= DebugLevel {
		newEntry(logger, 5).Debugf(format, args...)
	}
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	if logger.Level >= InfoLevel {
		newEntry(logger, 5).Infof(format, args...)
	}
}

func (logger *Logger) Printf(format string, args ...interface{}) {
	newEntry(logger, 5).Printf(format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	if logger.Level >= WarnLevel {
		newEntry(logger, 5).Warnf(format, args...)
	}
}

func (logger *Logger) Warningf(format string, args ...interface{}) {
	if logger.Level >= WarnLevel {
		newEntry(logger, 5).Warnf(format, args...)
	}
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	if logger.Level >= ErrorLevel {
		newEntry(logger, 5).Errorf(format, args...)
	}
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	if logger.Level >= FatalLevel {
		newEntry(logger, 5).Fatalf(format, args...)
	}
	os.Exit(1)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	if logger.Level >= PanicLevel {
		newEntry(logger, 5).Panicf(format, args...)
	}
}

func (logger *Logger) Debug(args ...interface{}) {
	if logger.Level >= DebugLevel {
		newEntry(logger, 5).Debug(args...)
	}
}

func (logger *Logger) Info(args ...interface{}) {
	if logger.Level >= InfoLevel {
		newEntry(logger, 5).Info(args...)
	}
}

func (logger *Logger) Print(args ...interface{}) {
	newEntry(logger, 5).Info(args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	if logger.Level >= WarnLevel {
		newEntry(logger, 5).Warn(args...)
	}
}

func (logger *Logger) Warning(args ...interface{}) {
	if logger.Level >= WarnLevel {
		newEntry(logger, 5).Warn(args...)
	}
}

func (logger *Logger) Error(args ...interface{}) {
	if logger.Level >= ErrorLevel {
		newEntry(logger, 5).Error(args...)
	}
}

func (logger *Logger) Fatal(args ...interface{}) {
	if logger.Level >= FatalLevel {
		newEntry(logger, 5).Fatal(args...)
	}
	os.Exit(1)
}

func (logger *Logger) Panic(args ...interface{}) {
	if logger.Level >= PanicLevel {
		newEntry(logger, 5).Panic(args...)
	}
}

func (logger *Logger) Debugln(args ...interface{}) {
	if logger.Level >= DebugLevel {
		newEntry(logger, 5).Debugln(args...)
	}
}

func (logger *Logger) Infoln(args ...interface{}) {
	if logger.Level >= InfoLevel {
		newEntry(logger, 5).Infoln(args...)
	}
}

func (logger *Logger) Println(args ...interface{}) {
	newEntry(logger, 5).Println(args...)
}

func (logger *Logger) Warnln(args ...interface{}) {
	if logger.Level >= WarnLevel {
		newEntry(logger, 5).Warnln(args...)
	}
}

func (logger *Logger) Warningln(args ...interface{}) {
	if logger.Level >= WarnLevel {
		newEntry(logger, 5).Warnln(args...)
	}
}

func (logger *Logger) Errorln(args ...interface{}) {
	if logger.Level >= ErrorLevel {
		newEntry(logger, 5).Errorln(args...)
	}
}

func (logger *Logger) Fatalln(args ...interface{}) {
	if logger.Level >= FatalLevel {
		newEntry(logger, 5).Fatalln(args...)
	}
	os.Exit(1)
}

func (logger *Logger) Panicln(args ...interface{}) {
	if logger.Level >= PanicLevel {
		newEntry(logger, 5).Panicln(args...)
	}
}

// SetFileLineEnabled sets the standard logger level.
func (logger *Logger) ShowCaller(b bool) {
	logger.showCaller = b
}
