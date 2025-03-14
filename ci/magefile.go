//go:build mage

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func intersect(a, b []string) []string {
	sort.Strings(a)
	sort.Strings(b)

	res := make([]string, 0, func() int {
		if len(a) < len(b) {
			return len(a)
		}
		return len(b)
	}())

	for _, v := range a {
		idx := sort.SearchStrings(b, v)
		if idx < len(b) && b[idx] == v {
			res = append(res, v)
		}
	}
	return res
}

// getBuildMatrix returns the build matrix from the current version of the go compiler
func getFullBuildMatrix() (map[string][]string, error) {
	jsonData, err := sh.Output("go", "tool", "dist", "list", "-json")
	if err != nil {
		return nil, err
	}
	var data []struct {
		Goos   string
		Goarch string
	}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return nil, err
	}

	matrix := map[string][]string{}
	for _, v := range data {
		if val, ok := matrix[v.Goos]; ok {
			matrix[v.Goos] = append(val, v.Goarch)
		} else {
			matrix[v.Goos] = []string{v.Goarch}
		}
	}

	return matrix, nil
}

func getBuildMatrix() (map[string][]string, error) {
	minimalMatrix := map[string][]string{
		"linux":   []string{"amd64"},
		"darwin":  []string{"amd64", "arm64"},
		"freebsd": []string{"amd64"},
		"js":      []string{"wasm"},
		"solaris": []string{"amd64"},
		"windows": []string{"amd64", "arm64"},
	}

	fullMatrix, err := getFullBuildMatrix()
	if err != nil {
		return nil, err
	}

	for os, arches := range minimalMatrix {
		if fullV, ok := fullMatrix[os]; !ok {
			delete(minimalMatrix, os)
		} else {
			minimalMatrix[os] = intersect(arches, fullV)
		}
	}
	return minimalMatrix, nil
}

func CrossBuild() error {
	matrix, err := getBuildMatrix()
	if err != nil {
		return err
	}

	for os, arches := range matrix {
		for _, arch := range arches {
			env := map[string]string{
				"GOOS":   os,
				"GOARCH": arch,
			}
			if mg.Verbose() {
				fmt.Printf("Building for GOOS=%s GOARCH=%s\n", os, arch)
			}
			if err := sh.RunWith(env, "go", "build", "./..."); err != nil {
				return err
			}
		}
	}
	return nil
}

func Lint() error {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return fmt.Errorf("cannot retrieve GOPATH")
	}

	return sh.Run(path.Join(gopath, "bin", "golangci-lint"), "run", "./...")
}

// Run the test suite
func Test() error {
	return sh.RunWith(map[string]string{"GORACE": "halt_on_error=1"},
		"go", "test", "-race", "-v", "./...")
}
