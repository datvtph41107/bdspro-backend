#!/usr/bin/env bash
set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/env.sh"

if [[ "${CONFIRM_QHPRO_REPORT_LAB_RESET:-}" != "yes" ]]; then
  echo 'This clears report/quota/file rows and the dedicated lab Redis database.' >&2
  echo 'Re-run with CONFIRM_QHPRO_REPORT_LAB_RESET=yes.' >&2
  exit 1
fi

docker compose --file "$REPORT_LAB_DIR/compose.yml" exec --no-TTY postgres \
  psql --username qhpro_lab --dbname qhpro_tqd_lab \
  --set ON_ERROR_STOP=1 \
  --command 'TRUNCATE TABLE report_jobs, quota_usage_events, user_reported RESTART IDENTITY CASCADE;'

docker compose --file "$REPORT_LAB_DIR/compose.yml" exec --no-TTY postgres \
  psql --username qhpro_lab --dbname qhpro_file_lab \
  --set ON_ERROR_STOP=1 \
  --command 'TRUNCATE TABLE file RESTART IDENTITY CASCADE;'

docker compose --file "$REPORT_LAB_DIR/compose.yml" exec --no-TTY redis \
  redis-cli --no-auth-warning --pass qhpro_local_only FLUSHDB

echo 'Dedicated report lab state was reset.'
