#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"/..

echo "Building proto..."
./bin/proto.sh


echo "Building go code..."
go build -o ./battleship-cli/dist ./battleship-cli/battleship_cli.go


echo "Building kotlin code..."
./gradlew --no-daemon build

