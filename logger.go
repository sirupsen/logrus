package logrus

import (
	"io"
	"os"
	"sync"
)

type Logger struct {
	Out       io.Writer
	Hooks     levelHooks
	Formatter Formatter
	mu        sync.Mutex
}

func New() *Logger {
	return &Logger{
		Out:       os.Stdout, // Default to stdout, change it if you want.
		Formatter: new(TextFormatter),
		Hooks:     make(levelHooks),
	}
}

func (logger *Logger) WithField(key string, value interface{}) *Entry {
	return NewEntry(logger).WithField(key, value)
}

func (logger *Logger) WithFields(fields Fields) *Entry {
	return NewEntry(logger).WithFields(fields)
}

// Logger Printf family functions

func (logger *Logger) Debugf(format string, args ...interface{}) {
	NewEntry(logger).Debugf(format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	NewEntry(logger).Infof(format, args...)
}

func (logger *Logger) Printf(format string, args ...interface{}) {
	NewEntry(logger).Printf(format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	NewEntry(logger).Warnf(format, args...)
}

func (logger *Logger) Warningf(format string, args ...interface{}) {
	NewEntry(logger).Warnf(format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	NewEntry(logger).Errorf(format, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	NewEntry(logger).Fatalf(format, args...)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	NewEntry(logger).Panicf(format, args...)
}

// Logger Print family functions

func (logger *Logger) Debug(args ...interface{}) {
	NewEntry(logger).Debug(args...)
}

func (logger *Logger) Info(args ...interface{}) {
	NewEntry(logger).Info(args...)
}

func (logger *Logger) Print(args ...interface{}) {
	NewEntry(logger).Print(args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	NewEntry(logger).Warn(args...)
}

func (logger *Logger) Warning(args ...interface{}) {
	NewEntry(logger).Warn(args...)
}

func (logger *Logger) Error(args ...interface{}) {
	NewEntry(logger).Error(args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	NewEntry(logger).Fatal(args...)
}

func (logger *Logger) Panic(args ...interface{}) {
	NewEntry(logger).Panic(args...)
}

// Logger Println family functions

func (logger *Logger) Debugln(args ...interface{}) {
	NewEntry(logger).Debugln(args...)
}

func (logger *Logger) Infoln(args ...interface{}) {
	NewEntry(logger).Infoln(args...)
}

func (logger *Logger) Println(args ...interface{}) {
	NewEntry(logger).Println(args...)
}

func (logger *Logger) Warnln(args ...interface{}) {
	NewEntry(logger).Warnln(args...)
}

func (logger *Logger) Warningln(args ...interface{}) {
	NewEntry(logger).Warnln(args...)
}

func (logger *Logger) Errorln(args ...interface{}) {
	NewEntry(logger).Errorln(args...)
}

func (logger *Logger) Fatalln(args ...interface{}) {
	NewEntry(logger).Fatalln(args...)
}

func (logger *Logger) Panicln(args ...interface{}) {
	NewEntry(logger).Panicln(args...)
}
