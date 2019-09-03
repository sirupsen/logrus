#!/bin/bash

if [[ "$TRAVIS_GO_VERSION" =~ ^1.\12\. ]] && [[ "$TRAVIS_OS_NAME" == "linux" ]]; then
    /tmp/gox/gox -build-lib -all
fi
