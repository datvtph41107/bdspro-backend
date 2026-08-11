#!/usr/bin/env bash

REPORT_LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPORT_LAB_ROOT="$(cd "$REPORT_LAB_DIR/../.." && pwd)"

export REPORT_LAB_DIR REPORT_LAB_ROOT
export QHPRO_LAB_POSTGRES_USER='qhpro_lab'
export QHPRO_LAB_POSTGRES_PASSWORD='qhpro_local_only'
export QHPRO_LAB_POSTGRES_PORT='55432'
export QHPRO_LAB_REDIS_PASSWORD='qhpro_local_only'
export QHPRO_LAB_REDIS_ADDR='127.0.0.1:56379'
export QHPRO_LAB_TRUST_SECRET='qhpro-local-transport-trust-only'
export QHPRO_LAB_JWT_SECRET='qhpro-local-jwt-only'
export QHPRO_LAB_PROFILE_ID="${QHPRO_LAB_PROFILE_ID:-42}"

export QHPRO_LAB_USER_DSN="host=127.0.0.1 user=$QHPRO_LAB_POSTGRES_USER password=$QHPRO_LAB_POSTGRES_PASSWORD dbname=qhpro_user_lab port=$QHPRO_LAB_POSTGRES_PORT sslmode=disable"
export QHPRO_LAB_TQD_DSN="host=127.0.0.1 user=$QHPRO_LAB_POSTGRES_USER password=$QHPRO_LAB_POSTGRES_PASSWORD dbname=qhpro_tqd_lab port=$QHPRO_LAB_POSTGRES_PORT sslmode=disable"
export QHPRO_LAB_FILE_DSN="host=127.0.0.1 user=$QHPRO_LAB_POSTGRES_USER password=$QHPRO_LAB_POSTGRES_PASSWORD dbname=qhpro_file_lab port=$QHPRO_LAB_POSTGRES_PORT sslmode=disable"
