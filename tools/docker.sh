#!/bin/sh

set -e

for tag in latest $(git tag --points-at | sed s/^v//); do
  docker buildx build \
    --platform linux/386,linux/amd64,linux/arm,linux/arm64,windows/amd64 \
    --tag raviqqe/muffet:$tag \
    --push \
    .
done
