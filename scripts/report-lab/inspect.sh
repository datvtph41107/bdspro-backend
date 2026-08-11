#!/usr/bin/env bash
set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/env.sh"

docker compose --file "$REPORT_LAB_DIR/compose.yml" exec --no-TTY postgres \
  psql --username qhpro_lab --dbname qhpro_tqd_lab --command \
  "SELECT id,user_id,status,format,client_request_id,job_id,pdf_url,file_size FROM user_reported ORDER BY id; SELECT operation,operation_id,idempotency_key,amount,reservation_id FROM quota_usage_events ORDER BY id; SELECT id,report_id,status,attempts,last_error FROM report_jobs ORDER BY created_at;"

docker compose --file "$REPORT_LAB_DIR/compose.yml" exec --no-TTY redis \
  redis-cli --no-auth-warning --pass qhpro_local_only --scan --pattern 'quota:*'
