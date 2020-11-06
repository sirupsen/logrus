#!/bin/bash

set -e

# Install golanci 1.32.2
if [[ "$TRAVIS_GO_VERSION" =~ ^1\.15\. ]]; then
    curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(go env GOPATH)/bin v1.32.2
fi

# Only do this for go1.12 when modules are on so that it doesn't need to be done when modules are off as well.
if [[ "$TRAVIS_GO_VERSION" =~ ^1\.13\. ]] && [[ "$TRAVIS_OS_NAME" == "linux" ]] && [[ "$GO111MODULE" == "on" ]]; then
    GO111MODULE=off go get github.com/dgsb/gox
fi

if [[ "$GO111MODULE" ==  "on" ]]; then
    go mod download
fi
