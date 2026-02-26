package logrus_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/sirupsen/logrus"
)

const (
	envChild   = "LOGRUS_EXIT_CHILD"
	envOutfile = "LOGRUS_EXIT_OUTFILE"
)

func TestRegister(t *testing.T) {
	if os.Getenv(envChild) == "1" {
		outfile := os.Getenv(envOutfile)

		logrus.RegisterExitHandler(func() { appendLine(t, outfile, "first") })
		logrus.RegisterExitHandler(func() { appendLine(t, outfile, "second") })

		logrus.Exit(23)
		return
	}

	outfile := filepath.Join(t.TempDir(), "out.txt")
	exitCode, out := reExecTest(t, outfile)
	if exitCode != 23 {
		t.Fatalf("expected exit code 23, got %d (out=%s)", exitCode, out)
	}

	want := "first\nsecond\n"
	assertFileContent(t, outfile, want)
}

func TestDefer(t *testing.T) {
	if os.Getenv(envChild) == "1" {
		outfile := os.Getenv(envOutfile)

		logrus.DeferExitHandler(func() { appendLine(t, outfile, "first") })
		logrus.DeferExitHandler(func() { appendLine(t, outfile, "second") })

		logrus.Exit(23)
		return
	}

	outfile := filepath.Join(t.TempDir(), "out.txt")
	exitCode, out := reExecTest(t, outfile)
	if exitCode != 23 {
		t.Fatalf("expected exit code 23, got %d (out=%s)", exitCode, out)
	}

	want := "second\nfirst\n"
	assertFileContent(t, outfile, want)
}

func TestHandler(t *testing.T) {
	const payload = "payload"
	if os.Getenv(envChild) == "1" {
		outfile := os.Getenv(envOutfile)

		logrus.RegisterExitHandler(func() {
			_ = os.WriteFile(outfile, []byte(payload), 0o666)
		})
		logrus.RegisterExitHandler(func() { panic("bad handler") })

		logrus.Exit(23)
		return
	}

	outfile := filepath.Join(t.TempDir(), "outfile.out")
	exitCode, out := reExecTest(t, outfile)
	if exitCode != 23 {
		t.Fatalf("expected exit code 23, got %d (out=%s)", exitCode, out)
	}

	want := payload
	assertFileContent(t, outfile, want)
}

// reExecTest re-executes the current test binary, running only the calling
// test in a subprocess.
//
// The child process is selected via -test.run using the current test's
// name and is signaled through environment variables (envChild and
// envOutfile). The child branch is expected to call logrus.Exit, which
// terminates the process via os.Exit.
//
// reExecTest returns the child's exit code and its combined stdout/stderr
// output. If the child exits with code 0, exitCode will be 0 and the
// captured output is returned.
func reExecTest(t *testing.T, outfile string) (int, []byte) {
	t.Helper()

	pattern := "^" + regexp.QuoteMeta(t.Name()) + "$"
	cmd := exec.Command(os.Args[0], "-test.run="+pattern)
	cmd.Env = append(os.Environ(),
		envChild+"=1",
		envOutfile+"="+outfile,
	)

	out, err := cmd.CombinedOutput()
	if err == nil {
		return 0, out
	}

	var exitErr *exec.ExitError
	ok := errors.As(err, &exitErr)
	if !ok {
		t.Fatalf("expected ExitError, got %T: %v (out=%s)", err, err, out)
	}

	return exitErr.ExitCode(), out
}

func appendLine(t *testing.T, path, s string) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		t.Logf("failed to open file %s for appending: %v", path, err)
		return
	}
	defer func() { _ = f.Close() }()

	_, err = f.WriteString(s + "\n")
	if err != nil {
		t.Logf("failed to write to file %s: %v", path, err)
		return
	}
}

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("can't read output file: %v", err)
	}

	got := string(b)
	if got != want {
		t.Fatalf("unexpected file content: got %q, want %q", got, want)
	}
}
