package logrus

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type color int

const (
	nocolor color = 0
	red           = 31
	green         = 32
	yellow        = 33
	blue          = 34
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
		levelText := f.stringForLevel(entry.Level)

		levelColor := f.colorForLevel(entry.Level)

		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%04d] %-44s ", levelColor, levelText, miniTS(), entry.Msg)

		keys := make([]string, 0)
		for k, _ := range entry.Data {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := entry.Data[k]
			fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=%v", levelColor, k, v)
		}
	} else {
		f.AppendKeyValue(b, "time", entry.Time)
		f.AppendKeyValue(b, "level", f.stringForLevel(entry.Level))
		f.AppendKeyValue(b, "msg", entry.Msg)

		for key, value := range entry.Data {
			f.AppendKeyValue(b, key, value)
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

	e.Level, err = f.levelForString(results[1])
	if err != nil {
		return nil, err
	}

	e.Time, err = f.timeForOffsetString(results[2])
	if err != nil {
		return nil, err
	}

	e.Msg = results[3]

	e.Data = make(Fields)
	for key, value := range f.fieldsForString(results[4]) {
		e.Data[key] = value
	}

	return &e, nil
}

func (f *TextFormatter) levelForString(s string) (Level, error) {
	switch s {
	case "PANI":
		return Panic, nil
	case "FATA":
		return Fatal, nil
	case "ERRO":
		return Error, nil
	case "WARN":
		return Warn, nil
	case "INFO":
		return Info, nil
	case "DEBU":
		return Debug, nil
	default:
		return Info, fmt.Errorf("Could not parse level: %s", s)
	}
}

func (f *TextFormatter) stringForLevel(l Level) string {
	switch l {
	case Panic:
		return "PANI"
	case Fatal:
		return "FATA"
	case Error:
		return "ERRO"
	case Warn:
		return "WARN"
	case Info:
		return "INFO"
	case Debug:
		return "DEBU"
	default:
		return ""
	}
}

func (f *TextFormatter) colorForLevel(l Level) color {
	switch l {
	case Panic, Fatal, Error:
		return red
	case Warn:
		return yellow
	default:
		return blue
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
