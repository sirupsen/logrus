#!/bin/bash

set -e

# Only do this for go1.12 when modules are on so that it doesn't need to be done when modules are off as well.
if [[ "$TRAVIS_GO_VERSION" =~ ^1\.13\. ]] && [[ "$TRAVIS_OS_NAME" == "linux" ]] && [[ "$GO111MODULE" == "on" ]]; then
    GO111MODULE=off go get github.com/dgsb/gox
fi

if [[ "$GO111MODULE" ==  "on" ]]; then
    go mod download
fi

if [[ "$GO111MODULE" == "off" ]]; then
    # Should contain all regular (not indirect) modules from go.mod
    go get github.com/stretchr/testify golang.org/x/sys/unix github.com/konsorten/go-windows-terminal-sequences
fi
