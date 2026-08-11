#!/usr/bin/env bash
set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/env.sh"
node "$REPORT_LAB_DIR/create-token.mjs"
printf '\n'
