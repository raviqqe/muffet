#!/bin/sh

set -ex

go get \
  github.com/golangci/golangci-lint/cmd/golangci-lint \
  github.com/m3ng9i/ran

go get -d ./...
