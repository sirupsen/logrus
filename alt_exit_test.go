package logrus_test

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestRegister(t *testing.T) {
	if reexecTest(t, "register", func(t *testing.T) {
		outfile := os.Args[len(os.Args)-1]
		logrus.RegisterExitHandler(func() { appendLine(outfile, "first") })
		logrus.RegisterExitHandler(func() { appendLine(outfile, "second") })
		logrus.Exit(23)
	}) {
		return
	}

	outfile := filepath.Join(t.TempDir(), "out.txt")
	cmd := reexecCommand(t, "register", outfile)
	out, err := cmd.CombinedOutput()
	var ee *exec.ExitError
	if !errors.As(err, &ee) {
		t.Fatalf("expected *exec.ExitError, got %T: %v (out=%s)", err, err, out)
	}
	if ee.ExitCode() != 23 {
		t.Fatalf("expected exit 23, got %d (out=%s)", ee.ExitCode(), out)
	}
	want := "first\nsecond\n"
	assertFileContent(t, outfile, want, out)
}

func TestDefer(t *testing.T) {
	if reexecTest(t, "defer", func(t *testing.T) {
		outfile := os.Args[len(os.Args)-1]
		logrus.DeferExitHandler(func() { appendLine(outfile, "first") })
		logrus.DeferExitHandler(func() { appendLine(outfile, "second") })
		logrus.Exit(23)
	}) {
		return
	}

	outfile := filepath.Join(t.TempDir(), "out.txt")
	cmd := reexecCommand(t, "defer", outfile)
	out, err := cmd.CombinedOutput()
	var ee *exec.ExitError
	if !errors.As(err, &ee) {
		t.Fatalf("expected *exec.ExitError, got %T: %v (out=%s)", err, err, out)
	}
	if ee.ExitCode() != 23 {
		t.Fatalf("expected exit 23, got %d (out=%s)", ee.ExitCode(), out)
	}
	want := "second\nfirst\n"
	assertFileContent(t, outfile, want, out)
}

func TestHandler(t *testing.T) {
	const payload = "payload"
	if reexecTest(t, "handler", func(t *testing.T) {
		outfile := os.Args[len(os.Args)-1]
		logrus.RegisterExitHandler(func() {
			_ = os.WriteFile(outfile, []byte(payload), 0o666)
		})
		logrus.RegisterExitHandler(func() { panic("bad handler") })

		logrus.Exit(23)
	}) {
		return
	}

	outfile := filepath.Join(t.TempDir(), "outfile.out")
	cmd := reexecCommand(t, "handler", outfile)
	out, err := cmd.CombinedOutput()
	var ee *exec.ExitError
	if !errors.As(err, &ee) {
		t.Fatalf("expected *exec.ExitError, got %T: %v (out=%s)", err, err, out)
	}
	if ee.ExitCode() != 23 {
		t.Fatalf("expected exit 23, got %d (out=%s)", ee.ExitCode(), out)
	}

	want := payload
	assertFileContent(t, outfile, want, out)
}

func appendLine(path, s string) {
	b, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		_, _ = os.Stderr.WriteString("appendLine: read " + path + ": " + err.Error() + "\n")
		os.Exit(1)
	}
	b = append(b, []byte(s+"\n")...)
	if err := os.WriteFile(path, b, 0o666); err != nil {
		_, _ = os.Stderr.WriteString("appendLine: write " + path + ": " + err.Error() + "\n")
		os.Exit(1)
	}
}

func assertFileContent(t *testing.T, path, want string, childOut []byte) {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("can't read output file: %v (child out=%s)", err, childOut)
	}

	if got := string(b); got != want {
		t.Fatalf("unexpected file content: got %q, want %q (child out=%s)", got, want, childOut)
	}
}

const tokenPrefix = "reexectest-"

// argv0Token computes a short deterministic token for (t.Name(), name).
func argv0Token(t *testing.T, name string) string {
	sum := sha256.Sum256([]byte(t.Name() + "\x00" + name))
	return tokenPrefix + hex.EncodeToString(sum[:8]) // 16 hex chars
}

// reexecTest runs fn if this process is the child (argv0 == token).
// Returns true in the child (caller should return).
func reexecTest(t *testing.T, name string, f func(t *testing.T)) bool {
	t.Helper()

	if os.Args[0] != argv0Token(t, name) {
		return false
	}

	// Scrub the "-test.run=<pattern>" that was injected by reexecCommand.
	origArgs := os.Args
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test.run=") {
		os.Args = append(os.Args[:1], os.Args[2:]...)
		defer func() { os.Args = origArgs }()
	}

	f(t)
	return true
}

// reexecCommand builds a command that execs the current test binary (exe)
// with argv0 set to token and "-test.run=<pattern>" so the child runs only
// this test/subtest. extraArgs are appended after that; the parent can pass
// the outfile as extra arg, etc.
func reexecCommand(t *testing.T, name string, args ...string) *exec.Cmd {
	t.Helper()

	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable(): %v", err)
	}

	argv0 := argv0Token(t, name)
	pattern := "^" + regexp.QuoteMeta(t.Name()) + "$"

	cmd := exec.Command(exe)
	cmd.Path = exe
	cmd.Args = append([]string{argv0, "-test.run=" + pattern}, args...)
	return cmd
}
