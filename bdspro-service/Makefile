SHELL := /bin/sh
BINARY_NAME := bdspro-service
RUN_ARGS := grpc
DB_URL ?= postgres://postgres:postgres@localhost:5432/db_bdspro?sslmode=disable
MIGRATION_DIR := database/migrations
MIGRATION_URL_ENV := BDS_DATABASE_URL

include ../shared/code/development/service.mk
include ../shared/code/development/migration.mk

buf:
	$(MAKE) -C ../shared/code buf-bdspro

help:
	@printf '%s\n' \
	  'fmt vet test test-race build server dev check' \
	  'migrate rollback migration name=<schema_change> migration-version buf clean'

.PHONY: help buf
