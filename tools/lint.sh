#!/bin/sh

set -ex

go run github.com/golangci/golangci-lint/cmd/golangci-lint run
