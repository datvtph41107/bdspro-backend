#!/usr/bin/env bash
set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/env.sh"
cd "$REPORT_LAB_ROOT/user-service"

exec env \
  ENV_RUNTIME= \
  DATABASE_DSN="$QHPRO_LAB_USER_DSN" \
  REDIS_HOST="$QHPRO_LAB_REDIS_ADDR" \
  REDIS_PASS="$QHPRO_LAB_REDIS_PASSWORD" \
  JWT_KEY_GENERATE="$QHPRO_LAB_JWT_SECRET" \
  QHPRO_SERVICE_ID=user-service \
  QHPRO_INTERNAL_METADATA_SECRET="$QHPRO_LAB_TRUST_SECRET" \
  QHPRO_TRUSTED_METADATA_MODE=enforce \
  go run . grpc
