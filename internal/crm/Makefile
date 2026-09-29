SHELL := /bin/sh
BINARY_NAME := crm-service
RUN_ARGS := grpc
DB_URL ?= postgres://postgres:postgres@localhost:5432/db_crm?sslmode=disable
MIGRATION_DIR := database/migrations
MIGRATION_URL_ENV := CRM_DATABASE_URL

include ../../tools/development/service.mk
include ../../tools/development/migration.mk

buf:
	$(MAKE) -C ../../tools/development buf-crm

payment-completed-consumer:
	@$(call with_env,go run ./cmd/payment-completed-consumer)

help:
	@printf '%s\n' \
	  'fmt vet test test-race build server dev check' \
	  'migrate rollback migration name=<schema_change> migration-version buf payment-completed-consumer clean'

.PHONY: help buf payment-completed-consumer
