package logrus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/tobi/airbrake-go"
)

type Entry struct {
	logger *Logger
	Data   Fields
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		logger: logger,
		// Default is three fields, give a little extra room. Shouldn't hurt the
		// scale.
		Data: make(Fields, 5),
	}
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
	entry.Data["time"] = time.Now().String()
	entry.Data["level"] = level
	// TODO: Is this the best name?
	entry.Data["msg"] = msg

	reader, err := entry.Reader()
	if err != nil {
		// TODO: Panic?
		entry.logger.Panicln("Failed to marshal JSON: ", err.Error())
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
		defer entry.logger.mu.Unlock()
		_, err := io.Copy(entry.logger.Out, reader)
		// TODO: Panic?
		if err != nil {
			entry.logger.Panicln("Failed to log message: ", err.Error())
		}
	}
}

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

func (entry *Entry) airbrake(exception string) {
	err := airbrake.Notify(errors.New(exception))
	if err != nil {
		entry.logger.WithFields(Fields{
			"source":   "airbrake",
			"endpoint": airbrake.Endpoint,
		}).Infof("Failed to send exception to Airbrake")
	}
}
