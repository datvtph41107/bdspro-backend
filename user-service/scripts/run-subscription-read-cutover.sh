#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ACTION="${1:-prepare}"
shift || true

case "$ACTION" in
  prepare|activate|rollback-plan|rollback) ;;
  *) echo "FAIL: unsupported action $ACTION" >&2; exit 1 ;;
esac

cd "$ROOT"
exec go run ./cmd/catalog-subscription-cutover --action "$ACTION" "$@"
