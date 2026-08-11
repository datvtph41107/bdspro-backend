#!/usr/bin/env bash
set -euo pipefail

for database_name in qhpro_user_lab qhpro_tqd_lab qhpro_file_lab; do
  if ! psql --username "$POSTGRES_USER" --dbname postgres --tuples-only --command \
    "SELECT 1 FROM pg_database WHERE datname = '$database_name'" | grep -q 1; then
    createdb --username "$POSTGRES_USER" "$database_name"
  fi
done

psql --username "$POSTGRES_USER" --dbname qhpro_user_lab \
  --set ON_ERROR_STOP=1 \
  --file /workspace/user-service/migrate/008_create_product_catalog_v2.up.sql
psql --username "$POSTGRES_USER" --dbname qhpro_user_lab \
  --set ON_ERROR_STOP=1 \
  --file /workspace/scripts/report-lab/sql/seed-user-access.sql

psql --username "$POSTGRES_USER" --dbname qhpro_tqd_lab \
  --set ON_ERROR_STOP=1 \
  --file /workspace/scripts/report-lab/sql/create-tqd-report-base.sql
psql --username "$POSTGRES_USER" --dbname qhpro_tqd_lab \
  --set ON_ERROR_STOP=1 \
  --file /workspace/tqd-service/infra/db/migrations/v3/000009_report_quota_target.up.sql
