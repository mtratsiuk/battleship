#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"/..

DOCKER_BUILDKIT=1 docker build \
  -t battleship-builder \
  -f Dockerfile \
  .
