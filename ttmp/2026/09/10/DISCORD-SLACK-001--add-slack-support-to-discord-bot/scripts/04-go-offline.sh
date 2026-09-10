#!/bin/sh
# Reproducible offline checks on this workstation, without toolchain downloads.
export GOROOT=/home/manuel/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.4.linux-amd64
export GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/slack-go-cache GOWORK=off
exec "$GOROOT/bin/go" "$@"
