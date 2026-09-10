SHELL := /bin/sh
BINARY_NAME := file-service
DB_URL ?= postgres://postgres:postgres@localhost:5432/file_service?sslmode=disable
MIGRATION_DIR := migrate
DEV_DEPENDENCIES := postgres

include ../shared/code/development/service.mk
include ../shared/code/development/deps.mk
include ../shared/code/development/migration.mk

buf:
	$(MAKE) -C ../shared/code buf-file

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
	  'deps-up deps-down migrateup migratedown1 migrate-version new_migration' \
	  'buf docker-build compose-up compose-down clean'

.PHONY: help buf docker-build compose-up compose-down acceptance-static
