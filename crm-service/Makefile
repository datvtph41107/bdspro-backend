SHELL := /bin/sh
BINARY_NAME := crm-service
RUN_ARGS := grpc
DB_URL ?= postgres://postgres:postgres@localhost:5432/db_crm?sslmode=disable
MIGRATION_DIR := infra/db/migrate_v2

include ../tools/development/service.mk

buf:
	$(MAKE) -C ../shared/code buf-crm

migrateup:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$${CRM_DATABASE_URL:-$(DB_URL)}" -verbose up)

migratedown1:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$${CRM_DATABASE_URL:-$(DB_URL)}" -verbose down 1)

help:
	@printf '%s\n' 'fmt vet test test-race build server buf migrateup migratedown1 clean check'

.PHONY: buf migrateup migratedown1 help
