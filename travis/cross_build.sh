#!/bin/bash

if [[ "$TRAVIS_GO_VERSION" =~ ^1\.12\. ]] && [[ "$TRAVIS_OS_NAME" == "linux" ]]; then
    $(go env GOPATH)/bin/gox -build-lib
fi
