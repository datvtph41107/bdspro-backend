#!/usr/bin/env bash
# Wrapper chạy load test public POIs.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

go run ./scripts/load_test_public_pois.go "$@"
