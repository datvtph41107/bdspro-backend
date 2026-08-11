#!/usr/bin/env bash
set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/env.sh"
cd "$REPORT_LAB_ROOT/gateway-service"

exec env \
  ENV_RUNTIME= \
  JWT_KEY_GENERATE="$QHPRO_LAB_JWT_SECRET" \
  QHPRO_SERVICE_ID=gateway-service \
  QHPRO_INTERNAL_METADATA_SECRET="$QHPRO_LAB_TRUST_SECRET" \
  QHPRO_TRUSTED_METADATA_MODE=enforce \
  QHPRO_REPORT_TARGET_ENABLED=true \
  go run . http
