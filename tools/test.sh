#!/bin/sh

set -ex

go test -covermode atomic -coverprofile coverage.txt
