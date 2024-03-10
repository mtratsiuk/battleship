#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"/..

rm -rf gen
cd ./battleship-proto
buf generate
