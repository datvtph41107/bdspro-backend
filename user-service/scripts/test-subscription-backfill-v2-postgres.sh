#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DATABASE_URL="${SUBSCRIPTION_BACKFILL_TEST_DATABASE_URL:-}"

fail() {
  printf 'FAIL: %s\n' "$*" >&2
  exit 1
}

[[ -n "$DATABASE_URL" ]] || \
  fail "SUBSCRIPTION_BACKFILL_TEST_DATABASE_URL is required"
command -v psql >/dev/null 2>&1 || fail "psql is required"

psql "$DATABASE_URL" \
  -v ON_ERROR_STOP=1 \
  -f "$ROOT/migrate/008_create_product_catalog_v2.up.sql" \
  -f "$ROOT/migrate/009_prepare_subscription_backfill.up.sql" \
  -f "$ROOT/migrate/009_prepare_subscription_backfill.smoke.sql"

printf '%s\n' 'PASS: PostgreSQL subscription backfill migration smoke'
