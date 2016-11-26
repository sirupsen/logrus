package logrus

import (
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	current := len(handlers)
	RegisterExitHandler(func() {})
	if len(handlers) != current+1 {
		t.Fatalf("can't add handler")
	}
}

func TestHandler(t *testing.T) {
	gofile := "/tmp/testprog.go"
	testprog := testprogleader
	testprog = append(testprog, getPackage()...)
	testprog = append(testprog, testprogtrailer...)
	if err := ioutil.WriteFile(gofile, testprog, 0666); err != nil {
		t.Fatalf("can't create go file")
	}

	outfile := "/tmp/testprog.out"
	arg := time.Now().UTC().String()
	err := exec.Command("go", "run", gofile, outfile, arg).Run()
	if err == nil {
		t.Fatalf("completed normally, should have failed")
	}

	data, err := ioutil.ReadFile(outfile)
	if err != nil {
		t.Fatalf("can't read output file %s", outfile)
	}

	if string(data) != arg {
		t.Fatalf("bad data")
	}
}

// getPackage returns the name of the current package, which makes running this
// test in a fork simpler
func getPackage() []byte {
	pc, _, _, _ := runtime.Caller(0)
	fullFuncName := runtime.FuncForPC(pc).Name()
	idx := strings.LastIndex(fullFuncName, ".")
	return []byte(fullFuncName[:idx]) // trim off function details
}

var testprogleader = []byte(`
// Test program for atexit, gets output file and data as arguments and writes
// data to output file in atexit handler.
package main

import (
	"`)
var testprogtrailer = []byte(
	`"
	"flag"
	"fmt"
	"io/ioutil"
)

var outfile = ""
var data = ""

func handler() {
	ioutil.WriteFile(outfile, []byte(data), 0666)
}

func badHandler() {
	n := 0
	fmt.Println(1/n)
}

func main() {
	flag.Parse()
	outfile = flag.Arg(0)
	data = flag.Arg(1)

	logrus.RegisterExitHandler(handler)
	logrus.RegisterExitHandler(badHandler)
	logrus.Fatal("Bye bye")
}
`)
