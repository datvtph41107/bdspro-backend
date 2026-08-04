#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ACTION="${1:-}"
shift || true

fail() {
  printf 'FAIL: %s\n' "$*" >&2
  exit 1
}

[[ -n "$ACTION" ]] || fail "usage: $0 review|plan|execute|report|readiness [flags]"
[[ -n "${QHPRO_DATABASE_DSN:-}" ]] || fail "QHPRO_DATABASE_DSN is required"

case "$ACTION" in
  review)
    [[ "${QHPRO_ALLOW_SUBSCRIPTION_BACKFILL_REVIEW:-}" == "1" ]] || \
      fail "set QHPRO_ALLOW_SUBSCRIPTION_BACKFILL_REVIEW=1"
    ;;
  execute)
    [[ "${QHPRO_ALLOW_SUBSCRIPTION_BACKFILL_EXECUTE:-}" == "1" ]] || \
      fail "set QHPRO_ALLOW_SUBSCRIPTION_BACKFILL_EXECUTE=1"
    ;;
  plan|report|readiness)
    ;;
  *)
    fail "unsupported action: $ACTION"
    ;;
esac

cd "$ROOT"
exec go run ./cmd/catalog-subscription-backfill \
  --action "$ACTION" \
  "$@"
