package logrus

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test that logger.Fatal() exits even if log level set to set Panic.
func TestFatal(t *testing.T) {
	assert := assert.New(t)

	exiter := new(exiter)
	logger := New()
	logger.Out = &bytes.Buffer{}
	logger.Level = PanicLevel
	logger.Exit = exiter.Exit

	logger.Fatal("kaboom")

	assert.Equal(1, exiter.exitCode)
}

// Test that logger.Fatalf() exits even if log level set to set Panic.
func TestFatalf(t *testing.T) {
	assert := assert.New(t)

	exiter := new(exiter)
	logger := New()
	logger.Out = &bytes.Buffer{}
	logger.Level = PanicLevel
	logger.Exit = exiter.Exit

	logger.Fatalf("kaboom")

	assert.Equal(1, exiter.exitCode)
}

// Test that logger.Fatalln() exits even if log level set to set Panic.
func TestFatalln(t *testing.T) {
	assert := assert.New(t)

	exiter := new(exiter)
	logger := New()
	logger.Out = &bytes.Buffer{}
	logger.Level = PanicLevel
	logger.Exit = exiter.Exit

	logger.Fatalln("kaboom")

	assert.Equal(1, exiter.exitCode)
}
