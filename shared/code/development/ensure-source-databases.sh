#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"

# This helper is intentionally best-effort. Service Makefiles can still point
# at an external/native PostgreSQL through ENV. When canonical Compose postgres
# is running, create the two non-default source-health databases without
# requiring a destructive volume reset.
if ! command -v docker >/dev/null 2>&1; then
  exit 0
fi
if ! docker compose --env-file .env -f compose.yaml ps --status running --services 2>/dev/null | grep -qx postgres; then
  exit 0
fi

for db in db_bdspro db_crm; do
  if ! docker compose --env-file .env -f compose.yaml exec -T postgres \
      psql -U postgres -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='${db}'" \
      | grep -q 1; then
    echo "creating local source-health database: $db"
    docker compose --env-file .env -f compose.yaml exec -T postgres \
      psql -U postgres -d postgres -v ON_ERROR_STOP=1 -c "CREATE DATABASE ${db} OWNER postgres"
  fi
done
