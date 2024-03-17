#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"/..


echo "Verifying go code..."

not_formatted=$(gofmt -l .)

if [ ! -z "$not_formatted" ]; then
  echo "Following go files are not formatted:"
  echo "$not_formatted"
  echo "Please run 'go fmt'"
  exit 1
fi


echo "Verifying kotlin code..."

./gradlew --no-daemon ktfmtCheck test

