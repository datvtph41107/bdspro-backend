#!/usr/bin/env bash
set -euo pipefail

tool_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$tool_dir/../../.." && pwd)"
cd "$repo_root"
umask 077

random_hex() {
  if command -v openssl >/dev/null 2>&1; then
    openssl rand -hex 32
  else
    od -An -N32 -tx1 /dev/urandom | tr -d ' \n'
  fi
}

install_env() {
  local path="$1" generator="$2" tmp
  if [[ -f "$path" ]]; then
    printf '[KEEP]   %s\n' "$path"
    return
  fi
  mkdir -p "$(dirname "$path")"
  tmp="$(mktemp "${path}.tmp.XXXXXX")"
  "$generator" >"$tmp"
  chmod 600 "$tmp"
  mv "$tmp" "$path"
  printf '[CREATE] %s\n' "$path"
}

install_missing_secret() {
  local path="$1" key="$2"
  if grep -q "^${key}=" "$path"; then
    return
  fi
  printf '%s=qhpro-development-%s\n' "$key" "$(random_hex)" >>"$path"
  printf '[ADD]    %s (%s)\n' "$path" "$key"
}

install_missing_value() {
  local path="$1" key="$2" value="$3"
  if grep -q "^${key}=" "$path"; then
    return
  fi
  printf '%s=%s\n' "$key" "$value" >>"$path"
  printf '[ADD]    %s (%s)\n' "$path" "$key"
}

# Canonical non-secret wiring is updated in existing ignored .env files when
# the development topology changes. Secrets remain create-once and are never
# printed or overwritten by this script.
install_value() {
  local path="$1" key="$2" value="$3" tmp
  tmp="$(mktemp "${path}.tmp.XXXXXX")"
  awk -v key="$key" -v value="$value" '
    BEGIN { found = 0 }
    index($0, key "=") == 1 { print key "=" value; found = 1; next }
    { print }
    END { if (!found) print key "=" value }
  ' "$path" >"$tmp"
  chmod 600 "$tmp"
  mv "$tmp" "$path"
}

remove_value() {
  local path="$1" key="$2" tmp
  [[ -f "$path" ]] || return 0
  grep -q "^${key}=" "$path" || return 0
  tmp="$(mktemp "${path}.tmp.XXXXXX")"
  awk -v key="$key" 'index($0, key "=") != 1 { print }' "$path" >"$tmp"
  chmod 600 "$tmp"
  mv "$tmp" "$path"
  printf '[REMOVE] %s (%s moved out of runtime ownership)\n' "$path" "$key"
}

root_env() {
  cat <<EOF
QHPRO_ENVIRONMENT=development
QHPRO_EXECUTION_MODE=host
GATEWAY_HTTP_PORT=8000
FILE_HTTP_PORT=8002
POSTGRES_PORT=5432
REDIS_PORT=6379
USER_GRPC_PORT=8201
ORGANIZATION_GRPC_PORT=8207
NOTIFICATION_GRPC_PORT=8204
PAYMENT_GRPC_PORT=8205
AUTH_GRPC_PORT=8216
TQD_GRPC_PORT=8219
HUB_GRPC_PORT=8280
ASSISTANT_GRPC_PORT=8218
RABBITMQ_PORT=5672
RABBITMQ_MANAGEMENT_PORT=15672
JWT_KEY_GENERATE=qhpro-development-$(random_hex)
QHPRO_INTERNAL_METADATA_SECRET=qhpro-development-$(random_hex)
QHPRO_TRUSTED_METADATA_MODE=enforce
QHPRO_INTERNAL_METADATA_MAX_AGE=30s
QHPRO_INTERNAL_METADATA_CLOCK_SKEW=5s
QHPRO_INTERNAL_METADATA_VERIFICATION_KEYS=
SERVICE_AUTH_KEY=qhpro-development-$(random_hex)
EOF
}

acceptance_env() {
  cat <<EOF
# Test-only identities used by integration-test. Application services do not
# load this file and production must never provision these accounts.
QHPRO_ACCEPTANCE_ADMIN_PASSWORD=qhpro-development-$(random_hex)
QHPRO_ACCEPTANCE_USER_PASSWORD=qhpro-development-$(random_hex)
EOF
}

auth_env() { printf '%s\n' 'QHPRO_SERVICE_ID=auth-service' 'RPC_USER_ADDRESS=localhost:8201'; }

user_env() {
  cat <<'EOF'
QHPRO_SERVICE_ID=user-service
DATABASE_URL=postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable
MIGRATION_URL=postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable
REDIS_ADDRESS=localhost:6379
USER_REDIS_DB=0
RPC_USER_ADDRESS=localhost:8201
RPC_AUTH_ADDRESS=localhost:8216
RPC_ORGANIZATION_ADDRESS=localhost:8207
RPC_HUB_ADDRESS=localhost:8280
RPC_BDSPRO_ADDRESS=localhost:8202
RPC_NOTIFICATION_ADDRESS=localhost:8204
RPC_PAYMENT_ADDRESS=localhost:8205
RPC_CHAT_ADDRESS=localhost:8208
RPC_CRM_ADDRESS=localhost:8206
QHPRO_USER_OUTBOUND_MESSAGING_ENABLED=false
EOF
}

organization_env() {
  cat <<'EOF'
QHPRO_SERVICE_ID=organization-service
ORGANIZATION_DATABASE_URL="host=localhost user=postgres password=postgres dbname=organization port=5432 sslmode=disable"
MIGRATION_URL=postgres://postgres:postgres@localhost:5432/organization?sslmode=disable
ORGANIZATION_GRPC_PORT=8207
ORGANIZATION_USER_GRPC_ADDRESS=localhost:8201
ORGANIZATION_BDSPRO_GRPC_ADDRESS=localhost:8202
ORGANIZATION_NOTIFICATION_GRPC_ADDRESS=localhost:8204
ORGANIZATION_PAYMENT_GRPC_ADDRESS=localhost:8205
ORGANIZATION_CHAT_GRPC_ADDRESS=localhost:8208
ORGANIZATION_TRANSACTION_GRPC_ADDRESS=localhost:8215
ORGANIZATION_AUTH_GRPC_ADDRESS=localhost:8216
EOF
}

payment_env() {
  cat <<EOF
QHPRO_SERVICE_ID=payment-service
PAYMENT_DATABASE_URL=postgres://postgres:postgres@localhost:5432/payment_service?sslmode=disable
MIGRATION_URL=postgres://postgres:postgres@localhost:5432/payment_service?sslmode=disable
PAYMENT_RABBIT_URL=amqp://guest:guest@localhost:5672/
PAYMENT_GRPC_ADDRESS=0.0.0.0:8205
QHPRO_USER_GRPC_ADDRESS=localhost:8201
QHPRO_AUTH_GRPC_ADDRESS=localhost:8216
QHPRO_NOTIFICATION_GRPC_ADDRESS=localhost:8204
SEPAY_API_KEY=qhpro-development-$(random_hex)
EOF
}

tqd_env() {
  cat <<'EOF'
QHPRO_SERVICE_ID=tqd-service
TQD_DATABASE_URL=postgres://qhpro:qhpro@127.0.0.1:5432/qhpro_tqd?sslmode=disable
MIGRATION_URL=postgres://qhpro:qhpro@127.0.0.1:5432/qhpro_tqd?sslmode=disable
TQD_REDIS_ADDR=127.0.0.1:6379
TQD_REDIS_DB=1
TQD_GRPC_PORT=8219
QHPRO_USER_GRPC_ADDR=localhost:8201
QHPRO_AUTH_GRPC_ADDR=localhost:8216
QHPRO_ASSISTANT_GRPC_ADDR=localhost:8218
QHPRO_FILE_HTTP_BASE_URL=http://localhost:8002
QHPRO_FILE_PUBLIC_BASE_URL=http://localhost:8002
QHPRO_REPORT_GENERATOR_MODE=production
TQD_CLASSIFY_ENABLED=false
EOF
}

notification_env() {
  cat <<'EOF'
QHPRO_SERVICE_ID=notification-service
NOTIFICATION_DATABASE_DSN=postgres://qhpro:qhpro@127.0.0.1:5432/qhpro_notification?sslmode=disable
MIGRATION_URL=postgres://qhpro:qhpro@127.0.0.1:5432/qhpro_notification?sslmode=disable
NOTIFICATION_REDIS_ADDR=127.0.0.1:6379
NOTIFICATION_REDIS_DB=2
QHPRO_RABBITMQ_URL=amqp://guest:guest@127.0.0.1:5672/
QHPRO_BUSINESS_EVENTS_EXCHANGE=qhpro.payment.events
NOTIFICATION_GRPC_PORT=8204
QHPRO_USER_GRPC_ADDR=localhost:8201
EOF
}

file_env() {
  cat <<EOF
QHPRO_SERVICE_ID=file-service
MIGRATION_URL=postgres://postgres:postgres@localhost:5432/file_service?sslmode=disable
DATABASE_DSN="host=localhost user=postgres password=postgres dbname=file_service port=5432 sslmode=disable"
FILE_SIGNATURE_KEY=qhpro-development-$(random_hex)
FILE_XOR_CRYPT_KEY=$(random_hex)
AUTH_GRPC_ADDRESS=localhost:8216
HUB_GRPC_ADDRESS=localhost:8280
EOF
}

gateway_env() { echo 'QHPRO_SERVICE_ID=gateway-service'; }

hub_env() {
  cat <<EOF
QHPRO_SERVICE_ID=hub-service
HUB_DATABASE_URL=postgres://postgres:postgres@localhost:5432/hub_service?sslmode=disable
MIGRATION_URL=postgres://postgres:postgres@localhost:5432/hub_service?sslmode=disable
HUB_REDIS_ADDRESS=localhost:6379
HUB_REDIS_DB=3
HUB_USER_GRPC_ADDRESS=localhost:8201
HUB_BDSPRO_GRPC_ADDRESS=localhost:8202
HUB_NOTIFICATION_GRPC_ADDRESS=localhost:8204
HUB_API_KEY=qhpro-development-$(random_hex)
CONFIG_ENCRYPTION_PASSPHRASE=qhpro-development-$(random_hex)
EOF
}

assistant_env() {
  cat <<'EOF'
QHPRO_SERVICE_ID=assistant-service
QHPRO_GRPC_PORT=8218
QHPRO_HTTP_PORT=8061
QHPRO_ASSISTANT_DB_DSN=postgres://postgres:postgres@localhost:5432/bdspro_db?sslmode=disable
QHPRO_ASSISTANT_REDIS_ADDR=localhost:6379
DEEPSEEK_API_KEY=
DEEPSEEK_BASE_URL=https://api.deepseek.com/v1
DEEPSEEK_MODEL=deepseek-chat
DEEPSEEK_TIMEOUT_SECONDS=30
DEEPSEEK_MAX_TOKENS=4096
OPENAI_API_KEY=
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_MODEL=gpt-4o-mini
OPENAI_TIMEOUT_SECONDS=45
OPENAI_MAX_TOKENS=4096
GEMINI_API_KEY=
GEMINI_BASE_URL=https://generativelanguage.googleapis.com/v1beta
GEMINI_MODEL=gemini-2.5-flash-lite
GEMINI_TIMEOUT_SECONDS=60
GEMINI_MAX_TOKENS=4096
QHPRO_AI_PROVIDER_MODE=stub
LOG_LEVEL=info
EOF
}

install_env .env root_env
remove_value .env QHPRO_ACCEPTANCE_ADMIN_PASSWORD
remove_value .env QHPRO_ACCEPTANCE_USER_PASSWORD
install_env integration-test/.env acceptance_env
install_missing_secret integration-test/.env QHPRO_ACCEPTANCE_ADMIN_PASSWORD
install_missing_secret integration-test/.env QHPRO_ACCEPTANCE_USER_PASSWORD
install_missing_value .env ASSISTANT_GRPC_PORT 8218
install_value .env POSTGRES_PORT 5432
install_value .env REDIS_PORT 6379
install_env auth-service/.env auth_env
install_env user-service/.env user_env
install_value user-service/.env DATABASE_URL 'postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable'
install_value user-service/.env MIGRATION_URL 'postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable'
install_value user-service/.env REDIS_ADDRESS localhost:6379
install_value user-service/.env USER_REDIS_DB 0
install_missing_value user-service/.env RPC_USER_ADDRESS localhost:8201
install_missing_value user-service/.env RPC_AUTH_ADDRESS localhost:8216
install_missing_value user-service/.env RPC_ORGANIZATION_ADDRESS localhost:8207
install_missing_value user-service/.env RPC_HUB_ADDRESS localhost:8280
install_missing_value user-service/.env RPC_BDSPRO_ADDRESS localhost:8202
install_missing_value user-service/.env RPC_NOTIFICATION_ADDRESS localhost:8204
install_missing_value user-service/.env RPC_PAYMENT_ADDRESS localhost:8205
install_missing_value user-service/.env RPC_CHAT_ADDRESS localhost:8208
install_missing_value user-service/.env RPC_CRM_ADDRESS localhost:8206
install_missing_value user-service/.env QHPRO_USER_OUTBOUND_MESSAGING_ENABLED false
install_env organization-service/.env organization_env
install_value organization-service/.env ORGANIZATION_DATABASE_URL '"host=localhost user=postgres password=postgres dbname=organization port=5432 sslmode=disable"'
install_value organization-service/.env MIGRATION_URL 'postgres://postgres:postgres@localhost:5432/organization?sslmode=disable'
install_missing_value organization-service/.env ORGANIZATION_GRPC_PORT 8207
install_missing_value organization-service/.env ORGANIZATION_USER_GRPC_ADDRESS localhost:8201
install_missing_value organization-service/.env ORGANIZATION_BDSPRO_GRPC_ADDRESS localhost:8202
install_missing_value organization-service/.env ORGANIZATION_NOTIFICATION_GRPC_ADDRESS localhost:8204
install_missing_value organization-service/.env ORGANIZATION_PAYMENT_GRPC_ADDRESS localhost:8205
install_missing_value organization-service/.env ORGANIZATION_CHAT_GRPC_ADDRESS localhost:8208
install_missing_value organization-service/.env ORGANIZATION_TRANSACTION_GRPC_ADDRESS localhost:8215
install_missing_value organization-service/.env ORGANIZATION_AUTH_GRPC_ADDRESS localhost:8216
install_env payment-service/.env payment_env
install_value payment-service/.env PAYMENT_DATABASE_URL 'postgres://postgres:postgres@localhost:5432/payment_service?sslmode=disable'
install_value payment-service/.env MIGRATION_URL 'postgres://postgres:postgres@localhost:5432/payment_service?sslmode=disable'
install_value payment-service/.env PAYMENT_RABBIT_URL 'amqp://guest:guest@localhost:5672/'
install_missing_value payment-service/.env PAYMENT_GRPC_ADDRESS 0.0.0.0:8205
install_missing_value payment-service/.env QHPRO_USER_GRPC_ADDRESS localhost:8201
install_missing_value payment-service/.env QHPRO_AUTH_GRPC_ADDRESS localhost:8216
install_missing_value payment-service/.env QHPRO_NOTIFICATION_GRPC_ADDRESS localhost:8204
install_env tqd-service/.env tqd_env
install_value tqd-service/.env TQD_DATABASE_URL 'postgres://qhpro:qhpro@127.0.0.1:5432/qhpro_tqd?sslmode=disable'
install_value tqd-service/.env MIGRATION_URL 'postgres://qhpro:qhpro@127.0.0.1:5432/qhpro_tqd?sslmode=disable'
install_value tqd-service/.env TQD_REDIS_ADDR '127.0.0.1:6379'
install_value tqd-service/.env TQD_REDIS_DB 1
install_missing_value tqd-service/.env TQD_GRPC_PORT 8219
install_missing_value tqd-service/.env QHPRO_USER_GRPC_ADDR localhost:8201
install_missing_value tqd-service/.env QHPRO_AUTH_GRPC_ADDR localhost:8216
install_missing_value tqd-service/.env QHPRO_ASSISTANT_GRPC_ADDR localhost:8218
install_missing_value tqd-service/.env QHPRO_FILE_HTTP_BASE_URL http://localhost:8002
install_missing_value tqd-service/.env QHPRO_FILE_PUBLIC_BASE_URL http://localhost:8002
install_missing_value tqd-service/.env QHPRO_REPORT_GENERATOR_MODE disabled
install_missing_value tqd-service/.env TQD_CLASSIFY_ENABLED false
install_env notification-service/.env notification_env
install_value notification-service/.env NOTIFICATION_DATABASE_DSN 'postgres://qhpro:qhpro@127.0.0.1:5432/qhpro_notification?sslmode=disable'
install_value notification-service/.env MIGRATION_URL 'postgres://qhpro:qhpro@127.0.0.1:5432/qhpro_notification?sslmode=disable'
install_value notification-service/.env NOTIFICATION_REDIS_ADDR '127.0.0.1:6379'
install_value notification-service/.env NOTIFICATION_REDIS_DB 2
install_value notification-service/.env QHPRO_RABBITMQ_URL 'amqp://guest:guest@127.0.0.1:5672/'
install_missing_value notification-service/.env QHPRO_BUSINESS_EVENTS_EXCHANGE qhpro.payment.events
install_missing_value notification-service/.env NOTIFICATION_GRPC_PORT 8204
install_missing_value notification-service/.env QHPRO_USER_GRPC_ADDR localhost:8201
install_env file-service/.env file_env
install_value file-service/.env MIGRATION_URL 'postgres://postgres:postgres@localhost:5432/file_service?sslmode=disable'
install_value file-service/.env DATABASE_DSN '"host=localhost user=postgres password=postgres dbname=file_service port=5432 sslmode=disable"'
install_missing_value file-service/.env AUTH_GRPC_ADDRESS localhost:8216
install_missing_value file-service/.env HUB_GRPC_ADDRESS localhost:8280
install_env gateway-service/.env gateway_env
install_env hub-service/.env hub_env
install_value hub-service/.env HUB_DATABASE_URL 'postgres://postgres:postgres@localhost:5432/hub_service?sslmode=disable'
install_value hub-service/.env MIGRATION_URL 'postgres://postgres:postgres@localhost:5432/hub_service?sslmode=disable'
install_value hub-service/.env HUB_REDIS_ADDRESS localhost:6379
install_value hub-service/.env HUB_REDIS_DB 3
install_missing_value hub-service/.env HUB_USER_GRPC_ADDRESS localhost:8201
install_missing_value hub-service/.env HUB_BDSPRO_GRPC_ADDRESS localhost:8202
install_missing_value hub-service/.env HUB_NOTIFICATION_GRPC_ADDRESS localhost:8204
install_env assistant-service/.env assistant_env

echo
echo 'Development configuration is ready.'
echo 'Existing files were preserved; no credential was printed.'
