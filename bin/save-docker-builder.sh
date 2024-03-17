#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"/..

. ./bin/set-env.sh

docker save battleship-builder
