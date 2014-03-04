package logrus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/burke/ttyutils"
	"github.com/tobi/airbrake-go"
)

type Entry struct {
	logger *Logger
	Data   Fields
}

var baseTimestamp time.Time

func init() {
	baseTimestamp = time.Now()
}

func miniTS() int {
	return int(time.Since(baseTimestamp) / time.Second)
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
func (entry *Entry) Reader() (*bytes.Buffer, error) {
	var serialized []byte
	var err error

	if Environment == "production" {
		serialized, err = json.Marshal(entry.Data)
		if err != nil {
			return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
		}
		serialized = append(serialized, '\n')
	} else {
		levelText := strings.ToUpper(entry.Data["level"].(string))
		levelColor := 34
		if levelText != "INFO" {
			levelColor = 31
		}
		if ttyutils.IsTerminal(os.Stdout.Fd()) {
			serialized = append(serialized, []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m[%04d] %-45s \x1b[%dm(\x1b[0m", levelColor, levelText, miniTS(), entry.Data["msg"], levelColor))...)
		}

		// TODO: Pretty-print more by coloring when stdout is a tty
		// TODO: If this is a println, it'll do a newline and then closing quote.
		keys := make([]string, 0)
		for k, _ := range entry.Data {
			if k != "level" && k != "time" && k != "msg" {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := entry.Data[k]
			serialized = append(serialized, []byte(fmt.Sprintf("\x1b[34m%s\x1b[0m='%s' ", k, v))...)
		}

		serialized = append(serialized, []byte(fmt.Sprintf("\x1b[%dm)\x1b[0m", levelColor))...)

		serialized = append(serialized, '\n')
	}

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

func (entry *Entry) log(level string, msg string) string {
	// TODO: Is the default format output from String() the one we want?
	entry.Data["time"] = time.Now().String()
	entry.Data["level"] = level
	// TODO: Is this the best name?
	entry.Data["msg"] = msg

	reader, err := entry.Reader()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v", err)
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

	entry.logger.mu.Lock()
	defer entry.logger.mu.Unlock()

	_, err = io.Copy(entry.logger.Out, reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v", err)
	}

	return reader.String()
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
		msg := entry.log("panic", fmt.Sprint(args...))
		panic(msg)
	}
	panic(fmt.Sprint(args...))
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
