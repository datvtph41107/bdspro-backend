#!/usr/bin/env bash
set -euo pipefail

MODE="${1:-}"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
USER_ROOT="$ROOT/user-service"
DATABASE_URL="${DATABASE_URL:-}"

fail() {
  printf 'FAIL: %s\n' "$*" >&2
  exit 1
}

[[ -n "$DATABASE_URL" ]] || fail "DATABASE_URL is required"
command -v psql >/dev/null 2>&1 || fail "psql is required"

case "$MODE" in
  up)
    psql "$DATABASE_URL" \
      -v ON_ERROR_STOP=1 \
      -f "$USER_ROOT/migrate/009_prepare_subscription_backfill.up.sql"
    ;;
  down)
    [[ "${ALLOW_DESTRUCTIVE_ROLLBACK:-}" == "1" ]] || \
      fail "set ALLOW_DESTRUCTIVE_ROLLBACK=1"
    psql "$DATABASE_URL" \
      -v ON_ERROR_STOP=1 \
      -c "SET qhpro.allow_subscription_backfill_v2_rollback = 'on';" \
      -f "$USER_ROOT/migrate/009_prepare_subscription_backfill.down.sql"
    ;;
  *)
    fail "usage: $0 up|down"
    ;;
esac
