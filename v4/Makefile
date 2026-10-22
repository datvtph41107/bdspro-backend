SHELL := /bin/sh
BINARY_NAME := mdn4
BUILD_PACKAGE := ./cmd/grpc
RUN_PACKAGE := ./cmd/grpc
# Alternate entry: make http  (HTTP_ADDRESS / server.port, default :8288)

http:
	@$(call with_env,go run ./cmd/http)

.PHONY: http
DB_URL ?= postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable
MIGRATION_DIR := migrate
DEV_DEPENDENCIES := postgres redis

include ../shared/code/development/service.mk
include ../shared/code/development/deps.mk
include ../shared/code/development/migration.mk

# Override service.mk with_env: YAML owns DSN/Redis/JWT/RPC. Root .env may still
# supply process secrets (QHPRO_INTERNAL_METADATA_*), but must not clobber YAML.
define with_env
set -a; \
if [ -f "$(ROOT_ENV_FILE)" ]; then . "$(abspath $(ROOT_ENV_FILE))"; fi; \
if [ -f "$(ENV_FILE)" ]; then . "$(abspath $(ENV_FILE))"; fi; \
set +a; \
unset USER_REDIS_ADDRESS REDIS_ADDRESS REDIS_HOST REDIS_PASSWORD USER_REDIS_PASSWORD \
  DATABASE_URL USER_DATABASE_URL DATABASE_DSN JWT_KEY_GENERATE \
  RPC_USER_ADDRESS RPC_AUTH_ADDRESS RPC_ORGANIZATION_ADDRESS RPC_HUB_ADDRESS \
  RPC_BDSPRO_ADDRESS RPC_NOTIFICATION_ADDRESS RPC_PAYMENT_ADDRESS RPC_CHAT_ADDRESS \
  RPC_CRM_ADDRESS RPC_TQD_ADDRESS QHPRO_USER_OUTBOUND_MESSAGING_ENABLED 2>/dev/null || true; \
QHPRO_EXECUTION_MODE=host QHPRO_LOG_CONFIG="$${QHPRO_LOG_CONFIG:-$(abspath ../logging.yml)}" $(1)
endef

buf:
	$(MAKE) -C ../shared/code buf-user-service

# Rewrites wire/wire.go from New* constructors + @bind annotations, then runs `wire gen`.
wire:
	go run ../shared/code/script/gen_wire.go v4

generate: buf wire

# One-time empty-installation operation. Inputs are supplied by the operator;
# this target never invents or persists default credentials.
bootstrap-admin:
	@$(call with_env,go run ./cmd/bootstrap-admin)

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/mdn4:$${QHPRO_IMAGE_TAG:-local}" ..

# Ship linux binary + config to production (see ./deploy.sh).
deploy:
	./deploy.sh $(if $(target),$(target),http)

deploy-compose:
	./deploy.sh compose

compose-up:
	$(MAKE) -C .. up

compose-down:
	$(MAKE) -C .. down

compose-logs:
	$(MAKE) -C .. logs service=user-service

acceptance-static: check

help:
	@printf '%s\n' \
	  'fmt vet test test-race build server dev check' \
	  'deps-up deps-down migrateup migratedown1 migrate-version new_migration' \
	  'bootstrap-admin buf wire generate docker-build compose-up compose-down compose-logs clean' \
	  'deploy [target=http|grpc] deploy-compose'

.PHONY: help buf wire generate bootstrap-admin docker-build deploy deploy-compose \
	compose-up compose-down compose-logs acceptance-static
