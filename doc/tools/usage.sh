#!/bin/sh

set -e

(
  echo '```text'
  go run .. --help
  echo '```'
) >src/components/Help.md
