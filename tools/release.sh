#!/bin/sh

set -e

if [ -z $CI ]; then
  exit 1
fi

version=$(go run . --version)

# if git tag -l | grep $version; then
#   exit
# fi

git config --global user.name "$GIT_USER"
git config --global user.email "$GIT_EMAIL"

git tag $version
# git push --tags

# curl -sL https://git.io/goreleaser | bash
