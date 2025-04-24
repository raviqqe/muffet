#!/bin/sh

set -e

(
  echo '```sh'
  go run .. --help
  echo '```'
) >src/components/Help.md
