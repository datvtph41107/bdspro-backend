#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIRECTION="${1:-up}"
: "${DATABASE_URL:?DATABASE_URL is required}"

case "$DIRECTION" in
  up)
    psql "$DATABASE_URL" -v ON_ERROR_STOP=1 \
      -f "$ROOT/migrate/010_subscription_read_cutover.up.sql"
    psql "$DATABASE_URL" -v ON_ERROR_STOP=1 \
      -f "$ROOT/migrate/010_subscription_read_cutover.smoke.sql"
    ;;
  down)
    : "${ALLOW_DESTRUCTIVE_ROLLBACK:?set ALLOW_DESTRUCTIVE_ROLLBACK=1}"
    [[ "$ALLOW_DESTRUCTIVE_ROLLBACK" == "1" ]] || {
      echo "FAIL: ALLOW_DESTRUCTIVE_ROLLBACK must equal 1" >&2
      exit 1
    }
    psql "$DATABASE_URL" -v ON_ERROR_STOP=1 \
      -c "SET qhpro.allow_subscription_cutover_v2_drop = 'on';" \
      -f "$ROOT/migrate/010_subscription_read_cutover.down.sql"
    ;;
  *) echo "FAIL: use up or down" >&2; exit 1 ;;
esac
