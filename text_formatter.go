package logrus

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
)

// {color}level{end-color}[seconds] msg    key=value key=value
var colorRE *regexp.Regexp

// {color}key{end-color}=value
var kvRE *regexp.Regexp

func init() {
	baseTimestamp = time.Now()

	// \x1b is CSI (Control Sequence Introducer)
	colorRE = regexp.MustCompile(`^\x1b\[\d+m(\w+)\x1b\[0m\[(\d+)\] (.*?) *(\x1b.*)?\r?\n?$`)
	kvRE = regexp.MustCompile(`\x1b\[\d+m([^\x1b]+)\x1b\[0m=([^\x1b]+)(?:\s+|$)`)
}

func miniTS() int {
	return int(time.Since(baseTimestamp) / time.Second)
}

type TextFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool
}

func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	if f.ForceColors || IsTerminal() {
		levelText := strings.ToUpper(entry.Data["level"].(string))[0:4]

		levelColor := blue

		if entry.Data["level"] == "warning" {
			levelColor = yellow
		} else if entry.Data["level"] == "error" ||
			entry.Data["level"] == "fatal" ||
			entry.Data["level"] == "panic" {
			levelColor = red
		}

		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%04d] %-44s ", levelColor, levelText, miniTS(), entry.Data["msg"])

		keys := make([]string, 0)
		for k, _ := range entry.Data {
			if k != "level" && k != "time" && k != "msg" {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := entry.Data[k]
			fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=%v", levelColor, k, v)
		}
	} else {
		f.AppendKeyValue(b, "time", entry.Data["time"].(string))
		f.AppendKeyValue(b, "level", entry.Data["level"].(string))
		f.AppendKeyValue(b, "msg", entry.Data["msg"].(string))

		for key, value := range entry.Data {
			if key != "time" && key != "level" && key != "msg" {
				f.AppendKeyValue(b, key, value)
			}
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) AppendKeyValue(b *bytes.Buffer, key, value interface{}) {
	if _, ok := value.(string); ok {
		fmt.Fprintf(b, "%v=%q ", key, value)
	} else {
		fmt.Fprintf(b, "%v=%v ", key, value)
	}
}

func (f *TextFormatter) Unformat(buffer []byte) (*Entry, error) {
	if len(buffer) == 0 {
		return nil, errors.New("Missing input")
	}

	if buffer[0] == '\x1b' {
		return f.unformatColor(buffer)
	} else {
		return f.unformatPlain(buffer)
	}
}

func (f *TextFormatter) unformatColor(buffer []byte) (*Entry, error) {
	// whole match, level, seconds, msg, key-value pairs
	results := colorRE.FindStringSubmatch(string(buffer))
	if results == nil {
		return nil, errors.New("Cannot parse input")
	}

	var err error
	var e Entry

	e.Data = make(Fields)

	e.Data["level"], err = f.standardLevelForString(results[1])
	if err != nil {
		return nil, err
	}

	t, err := f.timeForOffsetString(results[2])
	if err != nil {
		return nil, err
	}
	e.Data["time"] = t.String()

	e.Data["msg"] = results[3]

	for key, value := range f.fieldsForString(results[4]) {
		e.Data[key] = value
	}

	return &e, nil
}

func (f *TextFormatter) standardLevelForString(s string) (string, error) {
	switch s {
	case "PANI":
		return "panic", nil
	case "FATA":
		return "fatal", nil
	case "ERRO":
		return "error", nil
	case "WARN":
		return "warning", nil
	case "INFO":
		return "info", nil
	case "DEBU":
		return "debug", nil
	default:
		return "", fmt.Errorf("Could not parse level: %s", s)
	}
}

func (f *TextFormatter) timeForOffsetString(s string) (time.Time, error) {
	sec, err := strconv.Atoi(s)
	if err != nil {
		return time.Time{}, fmt.Errorf("Could not parse time: %s", s)
	}

	var t time.Time // Zero-time is the epoch (midnight on 1/1/1)
	t.Add(time.Duration(sec) * time.Second)

	return t, nil
}

func (f *TextFormatter) fieldsForString(s string) Fields {
	var fields Fields = make(Fields)
	for _, kv := range kvRE.FindAllStringSubmatch(s, -1) {
		fields[kv[1]] = kv[2]
	}

	return fields
}

func (f *TextFormatter) unformatPlain(buffer []byte) (*Entry, error) {
	return nil, errors.New("Unimplemented")
}
