#!/bin/sh

set -e

image=raviqqe/muffet
version=$(git tag --points-at | sed s/^v//)

docker buildx build \
  --platform linux/386,linux/amd64,linux/arm,linux/arm64 \
  --tag $image:latest \
  ${version:+--tag $image:$version} \
  "$@" \
  .
