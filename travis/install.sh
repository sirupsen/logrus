#!/bin/bash

set -e

if [[ "$TRAVIS_GO_VERSION" =~ ^1\.12\. ]] && [[ "$TRAVIS_OS_NAME" == "linux" ]]; then
    GO111MODULE=off go get github.com/dgsb/gox
fi

if [[ "$GO111MODULE" ==  "on" ]]; then
    go mod download
fi

if [[ "$GO111MODULE" == "off" ]]; then
    go get github.com/stretchr/testify/assert golang.org/x/sys/unix github.com/konsorten/go-windows-terminal-sequences
fi
