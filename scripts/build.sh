#!/usr/bin/env bash
# scripts/build.sh — Compile the application binary.
#
# Usage:
#   ./scripts/build.sh [output_path]
#
# Examples:
#   ./scripts/build.sh                 # outputs to ./build/server
#   ./scripts/build.sh ./bin/myapp     # custom output path

set -euo pipefail

OUTPUT=${1:-"./build/server"}
CMD_DIR="./cmd/api"

echo ">> Building binary..."
mkdir -p "$(dirname "$OUTPUT")"

CGO_ENABLED=0 go build \
  -ldflags="-s -w -X main.Version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" \
  -o "$OUTPUT" \
  "$CMD_DIR/..."

echo ">> Binary written to: $OUTPUT"
