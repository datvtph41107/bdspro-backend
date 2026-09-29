SHELL := /bin/sh
BINARY_NAME := file-service
DB_URL ?= postgres://postgres:postgres@localhost:5432/file_service?sslmode=disable
MIGRATION_DIR := database/migrations
DEV_DEPENDENCIES := postgres

include ../../tools/development/service.mk
include ../../tools/development/deps.mk
include ../../tools/development/migration.mk

buf:
	$(MAKE) -C ../../tools/development buf-file

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/file-service:$${QHPRO_IMAGE_TAG:-local}" ..

compose-up:
	$(MAKE) -C .. up

compose-down:
	$(MAKE) -C .. down

acceptance-static: check

help:
	@printf '%s\n' \
	  'fmt vet test test-race build server dev check' \
	  'deps-up deps-down migrate rollback migration name=<schema_change> migration-version' \
	  'buf docker-build compose-up compose-down clean'

.PHONY: help buf docker-build compose-up compose-down acceptance-static
