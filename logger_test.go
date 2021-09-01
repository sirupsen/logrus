package logrus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldValueError(t *testing.T) {
	buf := &bytes.Buffer{}
	l := &Logger{
		Out:       buf,
		Formatter: new(JSONFormatter),
		Hooks:     make(LevelHooks),
		Level:     DebugLevel,
	}
	l.WithField("func", func() {}).Info("test")
	fmt.Println(buf.String())
	var data map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		t.Error("unexpected error", err)
	}
	_, ok := data[FieldKeyLogrusError]
	require.True(t, ok, `cannot found expected "logrus_error" field: %v`, data)
}

func TestNoFieldValueError(t *testing.T) {
	buf := &bytes.Buffer{}
	l := &Logger{
		Out:       buf,
		Formatter: new(JSONFormatter),
		Hooks:     make(LevelHooks),
		Level:     DebugLevel,
	}
	l.WithField("str", "str").Info("test")
	fmt.Println(buf.String())
	var data map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		t.Error("unexpected error", err)
	}
	_, ok := data[FieldKeyLogrusError]
	require.False(t, ok)
}

func TestWarninglnNotEqualToWarning(t *testing.T) {
	buf := &bytes.Buffer{}
	bufln := &bytes.Buffer{}

	formatter := new(TextFormatter)
	formatter.DisableTimestamp = true
	formatter.DisableLevelTruncation = true

	l := &Logger{
		Out:       buf,
		Formatter: formatter,
		Hooks:     make(LevelHooks),
		Level:     DebugLevel,
	}
	l.Warning("hello,", "world")

	l.SetOutput(bufln)
	l.Warningln("hello,", "world")

	assert.NotEqual(t, buf.String(), bufln.String(), "Warning() and Wantingln() should not be equal")
}

type testBufferPool struct {
	buffers []*bytes.Buffer
	get     int
}

func (p *testBufferPool) Get() *bytes.Buffer {
	p.get++
	return new(bytes.Buffer)
}

func (p *testBufferPool) Put(buf *bytes.Buffer) {
	p.buffers = append(p.buffers, buf)
}

func TestLogger_SetBufferPool(t *testing.T) {
	out := &bytes.Buffer{}
	l := New()
	l.SetOutput(out)

	pool := new(testBufferPool)
	l.SetBufferPool(pool)

	l.Info("test")

	assert.Equal(t, pool.get, 1, "Logger.SetBufferPool(): The BufferPool.Get() must be called")
	assert.Len(t, pool.buffers, 1, "Logger.SetBufferPool(): The BufferPool.Put() must be called")
}

func TestLogger_concurrentLock(t *testing.T) {
	SetFormatter(&LogFormatter{})
	go func() {
		for {
			func() {
				defer func() {
					if p := recover(); p != nil {
					}
				}()
				hook := AddTraceIdHook("123")
				defer RemoveTraceHook(hook)
				Infof("test why ")
			}()
		}
	}()
	go func() {
		for {
			func() {
				defer func() {
					if p := recover(); p != nil {
					}
				}()
				hook := AddTraceIdHook("1233")
				defer RemoveTraceHook(hook)
				Infof("test why 2")
			}()
		}
	}()
	time.Sleep(1 * time.Minute)
}

var traceLock = &sync.Mutex{}

func AddTraceIdHook(traceId string) Hook {
	defer traceLock.Unlock()
	traceLock.Lock()
	traceHook := newTraceIdHook(traceId)
	if StandardLogger().Hooks == nil {
		hooks := new(LevelHooks)
		StandardLogger().ReplaceHooks(*hooks)
	}
	AddHook(traceHook)
	return traceHook
}

func RemoveTraceHook(hook Hook) {
	allHooks := StandardLogger().Hooks
	func() {
		defer Unlock()
		Lock()
		for key, hooks := range allHooks {
			replaceHooks := hooks
			for index, h := range hooks {
				if h == hook {
					replaceHooks = append(hooks[:index], hooks[index:]...)
					break
				}
			}
			allHooks[key] = replaceHooks
		}
	}()

	StandardLogger().ReplaceHooks(allHooks)
}

type TraceIdHook struct {
	TraceId string
	GID     uint64
}

func newTraceIdHook(traceId string) Hook {
	return &TraceIdHook{
		TraceId: traceId,
		GID:     getGID(),
	}
}

func (t TraceIdHook) Levels() []Level {
	return AllLevels
}

func (t TraceIdHook) Fire(entry *Entry) error {
	if getGID() == t.GID {
		entry.Context = context.WithValue(context.Background(), "trace_id", t.TraceId)
	}
	return nil
}

type LogFormatter struct{}


func (s *LogFormatter) Format(entry *Entry) ([]byte, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	var file string
	var line int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		line = entry.Caller.Line
	}
	level := entry.Level.String()
	if entry.Context == nil || entry.Context.Value("trace_id") == "" {
		uuid := "NO UUID"
		entry.Context = context.WithValue(context.Background(), "trace_id", uuid)
	}
	msg := fmt.Sprintf("%-15s [%-3d] [%-5s] [%s] %s:%d %s\n", timestamp, getGID(), level, entry.Context.Value("trace_id"), file, line, entry.Message)
	return []byte(msg), nil
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
