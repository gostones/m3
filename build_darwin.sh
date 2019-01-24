#!/usr/bin/env bash

export GOOS=darwin
export GOARCH=amd64
export CGO_ENABLED=0

export SKIP_TEST=false

#
./build.sh