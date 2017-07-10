package logrus

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	Nocolor = 0
	Black   = 30
	Red     = 31
	Green   = 32
	Yellow  = 33
	Blue    = 34
	Magenta = 35
	Cyan    = 36
	Gray    = 37

	BrightBlack   = 40
	BrightRed     = 41
	BrightGreen   = 42
	BrightYellow  = 43
	BrightBlue    = 44
	BrightMagenta = 45
	BrightCyan    = 46
	BrightGray    = 47
)

var (
	baseTimestamp time.Time
)

func init() {
	baseTimestamp = time.Now()
}

type TextFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Customize the color used for each level.  Keys are logging level constants
	// defined in this package; values are ANSI color codes.
	LevelColors map[Level]int

	// Set the number of characters of the level text that are printed.
	LevelTextLength int

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool

	// QuoteCharacter can be set to the override the default quoting character "
	// with something else. For example: ', or `.
	QuoteCharacter string

	// Whether the logger's out is to a terminal
	isTerminal bool

	sync.Once
}

func (f *TextFormatter) init(entry *Entry) {
	if len(f.QuoteCharacter) == 0 {
		f.QuoteCharacter = "\""
	}
	if entry.Logger != nil {
		f.isTerminal = IsTerminal(entry.Logger.Out)
	}

	// default level text length
	if f.LevelTextLength == 0 {
		f.LevelTextLength = 4
	}

	// replace missing level colors with defaults
	defaultLevelColors := map[Level]int{
		DebugLevel: Gray,
		WarnLevel:  Yellow,
		ErrorLevel: Red,
		FatalLevel: Red,
		PanicLevel: Red,
		InfoLevel:  Blue,
	}
	levelColors := make(map[Level]int)
	for level, _ := range defaultLevelColors {
		if color, exists := f.LevelColors[level]; !exists {
			levelColors[level] = defaultLevelColors[level]
		} else {
			levelColors[level] = color
		}
	}
	f.LevelColors = levelColors
}

func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if !f.DisableSorting {
		sort.Strings(keys)
	}
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	prefixFieldClashes(entry.Data)

	f.Do(func() { f.init(entry) })

	isColored := (f.ForceColors || f.isTerminal) && !f.DisableColors

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = DefaultTimestampFormat
	}
	if isColored {
		f.printColored(b, entry, keys, timestampFormat)
	} else {
		if !f.DisableTimestamp {
			f.appendKeyValue(b, "time", entry.Time.Format(timestampFormat))
		}
		f.appendKeyValue(b, "level", entry.Level.String())
		if entry.Message != "" {
			f.appendKeyValue(b, "msg", entry.Message)
		}
		for _, key := range keys {
			f.appendKeyValue(b, key, entry.Data[key])
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *Entry, keys []string, timestampFormat string) {
	levelColor := f.LevelColors[entry.Level]

	levelText := strings.ToUpper(entry.Level.String())
	for len(levelText) < f.LevelTextLength {
		levelText += " "
	}
	levelText = levelText[:f.LevelTextLength]

	if f.DisableTimestamp {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m %-44s ", levelColor, levelText, entry.Message)
	} else if !f.FullTimestamp {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%04d] %-44s ", levelColor, levelText, int(entry.Time.Sub(baseTimestamp)/time.Second), entry.Message)
	} else {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s] %-44s ", levelColor, levelText, entry.Time.Format(timestampFormat), entry.Message)
	}
	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColor, k)
		f.appendValue(b, v)
	}
}

func (f *TextFormatter) needsQuoting(text string) bool {
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.') {
			return true
		}
	}
	return false
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {

	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
	b.WriteByte(' ')
}

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	switch value := value.(type) {
	case string:
		if !f.needsQuoting(value) {
			b.WriteString(value)
		} else {
			b.WriteString(f.quoteString(value))
		}
	case error:
		errmsg := value.Error()
		if !f.needsQuoting(errmsg) {
			b.WriteString(errmsg)
		} else {
			b.WriteString(f.quoteString(errmsg))
		}
	default:
		fmt.Fprint(b, value)
	}
}

func (f *TextFormatter) quoteString(v string) string {
	escapedQuote := fmt.Sprintf("\\%s", f.QuoteCharacter)
	escapedValue := strings.Replace(v, f.QuoteCharacter, escapedQuote, -1)

	return fmt.Sprintf("%s%v%s", f.QuoteCharacter, escapedValue, f.QuoteCharacter)
}
