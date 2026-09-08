SHELL := /bin/sh
BINARY_NAME := user-service
BUILD_PACKAGE := ./cmd/grpc
RUN_PACKAGE := ./cmd/grpc
DB_URL ?= postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable
MIGRATION_DIR := database/migrations
DEV_DEPENDENCIES := postgres redis

include ../shared/code/development/service.mk
include ../shared/code/development/deps.mk
include ../shared/code/development/migration.mk

buf:
	$(MAKE) -C ../shared/code buf-user-service

wire:
	$(MAKE) -C ../shared/code wire user

generate: buf wire

# Run the User database-owned seed orchestration. Mandatory reference data that is
# versioned by migrations is intentionally not duplicated here, and identities
# belong to bootstrap or integration/E2E fixture lifecycles instead.
seed:
	@$(call with_env,go run ./cmd/seed)

# One-time empty-installation operation. Inputs are supplied by the operator;
# this target never invents or persists default credentials.
bootstrap-admin:
	@$(call with_env,go run ./cmd/bootstrap-admin)

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/user-service:$${QHPRO_IMAGE_TAG:-local}" ..

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
	  'seed bootstrap-admin buf wire generate docker-build compose-up compose-down compose-logs clean'

.PHONY: help buf wire generate seed bootstrap-admin docker-build \
	compose-up compose-down compose-logs acceptance-static
