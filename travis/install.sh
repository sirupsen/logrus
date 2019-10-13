#!/bin/bash

set -e

if [[ "$TRAVIS_GO_VERSION" =~ ^1\.12\. ]] && [[ "$TRAVIS_OS_NAME" == "linux" ]]; then
    git clone https://github.com/dgsb/gox.git /tmp/gox
    pushd /tmp/gox
    git checkout new_master
    go build ./
    popd
fi

if [[ "$GO111MODULE" ==  "on" ]]; then
    go mod download
fi

if [[ "$GO111MODULE" == "off" ]]; then
    go get github.com/stretchr/testify/assert golang.org/x/sys/unix github.com/konsorten/go-windows-terminal-sequences
    go get github.com/hashicorp/go-version github.com/hashicorp/go-version
fi
