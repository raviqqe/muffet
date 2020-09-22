#!/bin/sh

set -e

for tag in latest $(git describe --tags); do
  docker build -t raviqqe/muffet:$tag .
  docker push raviqqe/muffet:$tag
done
