#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"/..

docker run \
  --rm \
  --mount "type=bind,src=$(pwd),dst=/battleship" \
  -w /battleship \
  battleship-builder \
  "$@"
