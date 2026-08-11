#!/usr/bin/env bash
set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/env.sh"
cd "$REPORT_LAB_ROOT/file-service"

exec env \
  ENV_RUNTIME= \
  DATABASE_DSN="$QHPRO_LAB_FILE_DSN" \
  JWT_KEY_GENERATE="$QHPRO_LAB_JWT_SECRET" \
  QHPRO_SERVICE_ID=file-service \
  QHPRO_INTERNAL_METADATA_SECRET="$QHPRO_LAB_TRUST_SECRET" \
  QHPRO_TRUSTED_METADATA_MODE=enforce \
  QHPRO_FILE_GRPC_ENABLED=true \
  QHPRO_FILE_GRPC_PORT=8203 \
  go run .
