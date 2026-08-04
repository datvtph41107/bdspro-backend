#!/usr/bin/env bash
set -euo pipefail

DIRECTION="${1:-up}"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

: "${DATABASE_URL:?DATABASE_URL is required}"

case "$DIRECTION" in
  up)
    migration="$ROOT/migrate/008_create_product_catalog_v2.up.sql"
    ;;
  down)
    if [[ "${ALLOW_DESTRUCTIVE_ROLLBACK:-}" != "1" ]]; then
      printf '%s\n' \
        'Refusing destructive rollback.' \
        'Set ALLOW_DESTRUCTIVE_ROLLBACK=1 after read rollback and audit export.' >&2
      exit 1
    fi
    migration="$ROOT/migrate/008_create_product_catalog_v2.down.sql"
    export PGOPTIONS="${PGOPTIONS:-} -c qhpro.allow_catalog_v2_drop=on"
    ;;
  *)
    printf 'Usage: %s [up|down]\n' "$0" >&2
    exit 2
    ;;
esac

command -v psql >/dev/null 2>&1 || {
  printf 'psql is required\n' >&2
  exit 1
}

psql "$DATABASE_URL" \
  --set=ON_ERROR_STOP=1 \
  --file="$migration"
