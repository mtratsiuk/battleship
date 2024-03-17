#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"/..

if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi
