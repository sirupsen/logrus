package writer_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"github.com/stretchr/testify/assert"
)

func TestDifferentLevelsGoToDifferentWriters(t *testing.T) {
	var a, b bytes.Buffer

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	})
	log.SetOutput(io.Discard) // Send all logs to nowhere by default

	log.AddHook(&writer.Hook{
		Writer: &a,
		LogLevels: []logrus.Level{
			logrus.WarnLevel,
		},
	})
	log.AddHook(&writer.Hook{ // Send info and debug logs to stdout
		Writer: &b,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
		},
	})
	log.Warn("send to a")
	log.Info("send to b")

	assert.Equal(t, "level=warning msg=\"send to a\"\n", a.String())
	assert.Equal(t, "level=info msg=\"send to b\"\n", b.String())
}
