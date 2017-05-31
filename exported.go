package logrus

import (
	"io"
)

var (
	// std is the name of the standard logger in stdlib `log`
	std = New()
)

func StandardLogger() *Logger {
	return std
}

// SetOutput sets the standard logger output.
func SetOutput(out io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.Out = out
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter Formatter) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.Formatter = formatter
}

// SetLevel sets the standard logger level.
func SetLevel(level Level) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.setLevel(level)
}

// GetLevel returns the standard logger level.
func GetLevel() Level {
	std.mu.Lock()
	defer std.mu.Unlock()
	return std.level()
}

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook Hook) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.Hooks.Add(hook)
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *Entry {
	return std.WithField(ErrorKey, err)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *Entry {
	return std.WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields Fields) *Entry {
	return std.WithFields(fields)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	if std.level() >= DebugLevel {
		entry := std.newEntry().WithSkip(skip_5)
		entry.debug(args...)
		std.releaseEntry(entry)
	}
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	std.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	if std.level() >= InfoLevel {
		entry := std.newEntry().WithSkip(skip_5)
		entry.info(args...)
		std.releaseEntry(entry)
	}
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	if std.level() >= WarnLevel {
		entry := std.newEntry().WithSkip(skip_5)
		entry.warn(args...)
		std.releaseEntry(entry)
	}
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	if std.level() >= WarnLevel {
		entry := std.newEntry().WithSkip(skip_5)
		entry.warn(args...)
		std.releaseEntry(entry)
	}
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	if std.level() >= ErrorLevel {
		entry := std.newEntry().WithSkip(skip_5)
		entry.error(args...)
		std.releaseEntry(entry)
	}
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	if std.level() >= PanicLevel {
		entry := std.newEntry().WithSkip(skip_5)
		entry.panic(args...)
		std.releaseEntry(entry)
	}
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	if std.level() >= FatalLevel {
		entry := std.newEntry().WithSkip(skip_5)
		entry.fatal(args...)
		std.releaseEntry(entry)
	}
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	if std.level() >= DebugLevel {
		entry := std.newEntry().WithSkip(skip_6)
		entry.debugf(format, args...)
		std.releaseEntry(entry)
	}
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	std.Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	std.Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	std.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	if std.level() >= DebugLevel {
		entry := std.newEntry().WithSkip(skip_6)
		entry.debugln(args...)
		std.releaseEntry(entry)
	}
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	if std.level() >= InfoLevel {
		entry := std.newEntry().WithSkip(skip_7)
		entry.println(args...)
		std.releaseEntry(entry)
	}
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	if std.level() >= InfoLevel {
		entry := std.newEntry().WithSkip(skip_6)
		entry.infoln(args...)
		std.releaseEntry(entry)
	}
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	if std.level() >= WarnLevel {
		entry := std.newEntry().WithSkip(skip_6)
		entry.warnln(args...)
		std.releaseEntry(entry)
	}
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	if std.level() >= WarnLevel {
		entry := std.newEntry().WithSkip(skip_6)
		entry.warnln(args...)
		std.releaseEntry(entry)
	}
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	if std.level() >= ErrorLevel {
		entry := std.newEntry().WithSkip(skip_6)
		entry.errorln(args...)
		std.releaseEntry(entry)
	}
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	if std.level() >= PanicLevel {
		entry := std.newEntry().WithSkip(skip_6)
		entry.panicln(args...)
		std.releaseEntry(entry)
	}
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	if std.level() >= FatalLevel {
		entry := std.newEntry().WithSkip(skip_6)
		entry.fatalln(args...)
		std.releaseEntry(entry)
	}
}
