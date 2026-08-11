#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
docker compose --file "$LAB_DIR/compose.yml" down

echo 'Lab containers stopped. Named volumes were preserved.'
