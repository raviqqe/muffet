#!/bin/sh

set -e

for tag in latest $(git tag --points-at | sed s/^v//); do
  docker build -t raviqqe/muffet:$tag .
  docker push raviqqe/muffet:$tag
done
