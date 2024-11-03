#!/bin/sh

set -e

for tag in latest $(git tag --points-at | sed s/^v//); do
  docker buildx build \
    --platform linux/amd64,linux/arm64,linux/riscv64 \
    --tag raviqqe/muffet:$tag \
    --push \
    .
done
