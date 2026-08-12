package slog_test

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/sirupsen/logrus"
	lslog "github.com/sirupsen/logrus/hooks/slog"
)

// ExampleNewHandler demonstrates using slog alongside an existing Logrus logger.
func ExampleNewHandler() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})

	slog.SetDefault(slog.New(lslog.NewHandler(logger)))

	// Both slog and Logrus write through the same Logrus backend.
	slog.Info("hello from slog", "source", "slog")
	logger.WithField("source", "logrus").Info("hello from logrus")

	// Output:
	// level=info msg="hello from slog" source=slog
	// level=info msg="hello from logrus" source=logrus
}

// ExampleNewHook demonstrates forwarding existing Logrus logging to an slog
// logger, allowing applications to migrate their logging backend independently
// from code that still uses Logrus.
func ExampleNewHook() {
	var slogOutput bytes.Buffer
	slogger := slog.New(slog.NewTextHandler(&slogOutput, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return attr
		},
	}))

	logger := logrus.New()
	// Discard Logrus's own output; the hook forwards the entry to slog.
	logger.SetOutput(io.Discard)
	logger.AddHook(lslog.NewHook(slogger))

	// Log through Logrus; the hook forwards this to slog.
	logger.WithField("animal", "walrus").Info("hello from logrus")

	// Show what was emitted by the slog handler.
	fmt.Print(slogOutput.String())

	// Output:
	// level=INFO msg="hello from logrus" animal=walrus
}
