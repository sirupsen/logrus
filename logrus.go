package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/tobi/airbrake-go"
)

// TODO: Type naming here feels awkward, but the exposed variable should be
// Level. That's more important than the type name, and libraries should be
// reaching for logrus.Level{Debug,Info,Warning,Fatal}, not defining the type
// themselves as an int.
type LevelType uint8
type Fields map[string]interface{}

const (
	LevelPanic LevelType = iota
	LevelFatal
	LevelWarning
	LevelInfo
	LevelDebug
)

var Level LevelType = LevelInfo
var Environment string = "development"

// StandardLogger is what your logrus-enabled library should take, that way
// it'll accept a stdlib logger and a logrus logger. There's no standard
// interface, this is the closest we get, unfortunately.
type StandardLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Printfln(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}

type Logger struct {
	Out io.Writer
	mu  sync.Mutex
}

type Entry struct {
	logger *Logger
	Data   Fields
}

// TODO: Other formats?
func (entry *Entry) Reader() (read *bytes.Buffer, err error) {
	var serialized []byte

	if Environment == "production" {
		serialized, err = json.Marshal(entry.Data)
	} else {
		// TODO: Pretty-print more by coloring when stdout is a tty
		serialized, err = json.MarshalIndent(entry.Data, "", "  ")
	}

	if err != nil {
		return nil, err
	}

	serialized = append(serialized, '\n')

	return bytes.NewBuffer(serialized), nil
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

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		logger: logger,
		// Default is three fields, give a little extra room. Shouldn't hurt the
		// scale.
		Data: make(Fields, 5),
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

func (entry *Entry) WithField(key string, value interface{}) *Entry {
	entry.Data[key] = value
	return entry
}

func (entry *Entry) WithFields(fields Fields) *Entry {
	for key, value := range fields {
		entry.WithField(key, value)
	}

	return entry
}

func (entry *Entry) log(level string, msg string) {
	// TODO: Is the default format output from String() the one we want?
	entry.Data["timestamp"] = time.Now().String()
	entry.Data["level"] = level
	// TODO: Is this the best name?
	entry.Data["msg"] = msg

	reader, err := entry.Reader()
	if err != nil {
		entry.logger.Panicln("Failed to marshal JSON ", err.Error())
	}

	// Send HTTP request in a goroutine in warning environment to not halt the
	// main thread. It's sent before logging due to panic.
	if level == "warning" {
		// TODO: new() should spawn an airbrake goroutine and this should send to
		// that channel. This prevent us from spawning hundreds of goroutines in a
		// hot code path generating a warning.
		go entry.airbrake(reader.String())
	} else if level == "fatal" || level == "panic" {
		entry.airbrake(reader.String())
	}

	if level == "panic" {
		panic(reader.String())
	} else {
		entry.logger.mu.Lock()
		io.Copy(entry.logger.Out, reader)
		entry.logger.mu.Unlock()
	}
}

func (entry *Entry) airbrake(exception string) {
	err := airbrake.Notify(errors.New(exception))
	if err != nil {
		entry.logger.WithFields(Fields{
			"source":   "airbrake",
			"endpoint": airbrake.Endpoint,
		}).Infof("Failed to send exception to Airbrake")
	}
}

// Entry Print family functions

func (entry *Entry) Debug(args ...interface{}) {
	if Level >= LevelDebug {
		entry.log("debug", fmt.Sprint(args...))
	}
}

func (entry *Entry) Info(args ...interface{}) {
	if Level >= LevelInfo {
		entry.log("info", fmt.Sprint(args...))
	}
}

func (entry *Entry) Print(args ...interface{}) {
	if Level >= LevelInfo {
		entry.log("info", fmt.Sprint(args...))
	}
}

func (entry *Entry) Warning(args ...interface{}) {
	if Level >= LevelWarning {
		entry.log("warning", fmt.Sprint(args...))
	}
}

func (entry *Entry) Fatal(args ...interface{}) {
	if Level >= LevelFatal {
		entry.log("fatal", fmt.Sprint(args...))
	}
	os.Exit(1)
}

func (entry *Entry) Panic(args ...interface{}) {
	if Level >= LevelPanic {
		entry.log("panic", fmt.Sprint(args...))
	}
}

// Entry Printf family functions

func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.Debug(fmt.Sprintf(format, args...))
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.Info(fmt.Sprintf(format, args...))
}

func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.Print(fmt.Sprintf(format, args...))
}

func (entry *Entry) Warningf(format string, args ...interface{}) {
	entry.Warning(fmt.Sprintf(format, args...))
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.Fatal(fmt.Sprintf(format, args...))
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	entry.Panic(fmt.Sprintf(format, args...))
}

// Entry Println family functions

func (entry *Entry) Debugln(args ...interface{}) {
	entry.Debug(fmt.Sprintln(args...))
}

func (entry *Entry) Infoln(args ...interface{}) {
	entry.Info(fmt.Sprintln(args...))
}

func (entry *Entry) Println(args ...interface{}) {
	entry.Print(fmt.Sprintln(args...))
}

func (entry *Entry) Warningln(args ...interface{}) {
	entry.Warning(fmt.Sprintln(args...))
}

func (entry *Entry) Fatalln(args ...interface{}) {
	entry.Fatal(fmt.Sprintln(args...))
}

func (entry *Entry) Panicln(args ...interface{}) {
	entry.Panic(fmt.Sprintln(args...))
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

// TODO: Print, Fatal, etc.

func main() {
	Environment = "development"
	Level = LevelDebug
	logger := New()
	logger.WithField("animal", "walrus").WithField("value", 10).Infof("OMG HELLO")
	logger.Infof("lolsup")
	logger.Debugf("why brackets?")
	logger.Debug("lolsup")
	logger.Fatalf("omg")
}
