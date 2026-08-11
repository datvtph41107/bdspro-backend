#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

command -v docker >/dev/null 2>&1 || {
  echo 'Docker is required. On WSL, enable Docker Desktop > Settings > Resources > WSL Integration.' >&2
  exit 1
}

docker info >/dev/null 2>&1 || {
  echo 'Docker CLI exists but the daemon is unavailable.' >&2
  echo 'Start Docker Desktop and enable WSL integration, then retry.' >&2
  exit 1
}

docker compose --file "$LAB_DIR/compose.yml" up --detach --wait
echo 'PostgreSQL/PostGIS: 127.0.0.1:55432'
echo 'Redis:              127.0.0.1:56379'
