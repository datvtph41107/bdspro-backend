SHELL := /bin/sh
BINARY_NAME := bdspro-service
RUN_ARGS := grpc
DB_URL ?= postgres://postgres:postgres@localhost:5432/db_bdspro?sslmode=disable
MIGRATION_DIR := infra/db/migrations/v2

include ../tools/development/service.mk

buf:
	$(MAKE) -C ../shared/code buf-bdspro

migrateup:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$${BDS_DATABASE_URL:-$(DB_URL)}" -verbose up)

migratedown1:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$${BDS_DATABASE_URL:-$(DB_URL)}" -verbose down 1)

dev:
	@$(call with_env,air -c .air.toml)

help:
	@printf '%s\n' 'fmt vet test test-race build server dev buf migrateup migratedown1 clean check'

.PHONY: buf migrateup migratedown1 dev help
