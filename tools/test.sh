#!/bin/sh

set -ex

go test -covermode atomic -coverprofile coverage.txt

if [ -n "$CODECOV_TOKEN" ]; then
  curl -s https://codecov.io/bash | bash
fi
