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
		Out:       os.Stderr,
		Formatter: new(TextFormatter),
		Hooks:     make(LevelHooks),
		Level:     InfoLevel,
	}
}

// Adds a field to the log entry, note that you it doesn't log until you call
// Debug, Print, Info, Warn, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func (logger *Logger) WithField(key string, value interface{}) *Entry {
	return NewEntry(logger).WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (logger *Logger) WithFields(fields Fields) *Entry {
	return NewEntry(logger).WithFields(fields)
}

func (logger *Logger) Debugf(format string, args ...interface{}) *Entry {
	if logger.Level >= DebugLevel {
		return NewEntry(logger).Debugf(format, args...)
	}
	return nil
}

func (logger *Logger) Infof(format string, args ...interface{}) *Entry {
	if logger.Level >= InfoLevel {
		return NewEntry(logger).Infof(format, args...)
	}
	return nil
}

func (logger *Logger) Printf(format string, args ...interface{}) *Entry {
	return NewEntry(logger).Printf(format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) *Entry {
	if logger.Level >= WarnLevel {
		return NewEntry(logger).Warnf(format, args...)
	}
	return nil
}

func (logger *Logger) Warningf(format string, args ...interface{}) *Entry {
	if logger.Level >= WarnLevel {
		return NewEntry(logger).Warnf(format, args...)
	}
	return nil
}

func (logger *Logger) Errorf(format string, args ...interface{}) *Entry {
	if logger.Level >= ErrorLevel {
		return NewEntry(logger).Errorf(format, args...)
	}
	return nil
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	if logger.Level >= FatalLevel {
		NewEntry(logger).Fatalf(format, args...)
	}
	os.Exit(1)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	if logger.Level >= PanicLevel {
		NewEntry(logger).Panicf(format, args...)
	}
}

func (logger *Logger) Debug(args ...interface{}) *Entry {
	if logger.Level >= DebugLevel {
		return NewEntry(logger).Debug(args...)
	}
	return nil
}

func (logger *Logger) Info(args ...interface{}) *Entry {
	if logger.Level >= InfoLevel {
		return NewEntry(logger).Info(args...)
	}
	return nil
}

func (logger *Logger) Print(args ...interface{}) *Entry {
	return NewEntry(logger).Info(args...)
}

func (logger *Logger) Warn(args ...interface{}) *Entry {
	if logger.Level >= WarnLevel {
		return NewEntry(logger).Warn(args...)
	}
	return nil
}

func (logger *Logger) Warning(args ...interface{}) *Entry {
	if logger.Level >= WarnLevel {
		return NewEntry(logger).Warn(args...)
	}
	return nil
}

func (logger *Logger) Error(args ...interface{}) *Entry {
	if logger.Level >= ErrorLevel {
		return NewEntry(logger).Error(args...)
	}
	return nil
}

func (logger *Logger) Fatal(args ...interface{}) {
	if logger.Level >= FatalLevel {
		NewEntry(logger).Fatal(args...)
	}
	os.Exit(1)
}

func (logger *Logger) Panic(args ...interface{}) {
	if logger.Level >= PanicLevel {
		NewEntry(logger).Panic(args...)
	}
}

func (logger *Logger) Debugln(args ...interface{}) *Entry {
	if logger.Level >= DebugLevel {
		return NewEntry(logger).Debugln(args...)
	}
	return nil
}

func (logger *Logger) Infoln(args ...interface{}) *Entry {
	if logger.Level >= InfoLevel {
		return NewEntry(logger).Infoln(args...)
	}
	return nil
}

func (logger *Logger) Println(args ...interface{}) *Entry {
	return NewEntry(logger).Println(args...)
}

func (logger *Logger) Warnln(args ...interface{}) *Entry {
	if logger.Level >= WarnLevel {
		return NewEntry(logger).Warnln(args...)
	}
	return nil
}

func (logger *Logger) Warningln(args ...interface{}) *Entry {
	if logger.Level >= WarnLevel {
		return NewEntry(logger).Warnln(args...)
	}
	return nil
}

func (logger *Logger) Errorln(args ...interface{}) *Entry {
	if logger.Level >= ErrorLevel {
		return NewEntry(logger).Errorln(args...)
	}
	return nil
}

func (logger *Logger) Fatalln(args ...interface{}) {
	if logger.Level >= FatalLevel {
		NewEntry(logger).Fatalln(args...)
	}
	os.Exit(1)
}

func (logger *Logger) Panicln(args ...interface{}) {
	if logger.Level >= PanicLevel {
		NewEntry(logger).Panicln(args...)
	}
}
