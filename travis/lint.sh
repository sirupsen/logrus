#!/bin/bash

if [[ "$TRAVIS_GO_VERSION" =~ ^1\.15\. ]] && [[ "$TRAVIS_OS_NAME" == "linux" ]] && [[ "$GO111MODULE" == "on" ]]; then
    $(go env GOPATH)/bin/golangci-lint run ./...
else
    echo "linter has not been run"
fi
