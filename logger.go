package logrus

import (
	"io"
	"os"
	"strings"
	"sync"

	"github.com/tobi/airbrake-go"
)

type Logger struct {
	Out io.Writer
	mu  sync.Mutex
}

func New() *Logger {
	environment := strings.ToLower(os.Getenv("ENV"))
	if environment == "" {
		environment = "development"
	}

	if airbrake.Environment == "" {
		airbrake.Environment = environment
	}

	return &Logger{
		Out: os.Stdout, // Default to stdout, change it if you want.
	}
}

func (logger *Logger) WithField(key string, value interface{}) *Entry {
	entry := NewEntry(logger)
	entry.WithField(key, value)
	return entry
}

func (logger *Logger) WithFields(fields Fields) *Entry {
	entry := NewEntry(logger)
	entry.WithFields(fields)
	return entry
}

// Entry Print family functions
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

func (logger *Logger) Warningf(format string, args ...interface{}) {
	NewEntry(logger).Warningf(format, args...)
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

func (logger *Logger) Warning(args ...interface{}) {
	NewEntry(logger).Warning(args...)
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

func (logger *Logger) Warningln(args ...interface{}) {
	NewEntry(logger).Warningln(args...)
}

func (logger *Logger) Fatalln(args ...interface{}) {
	NewEntry(logger).Fatalln(args...)
}

func (logger *Logger) Panicln(args ...interface{}) {
	NewEntry(logger).Panicln(args...)
}
