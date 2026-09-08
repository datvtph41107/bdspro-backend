SHELL := /bin/sh
BINARY_NAME := organization-service
RUN_ARGS := grpc
DB_URL ?= postgres://postgres:postgres@localhost:5432/organization?sslmode=disable
MIGRATION_DIR := migrate
DEV_DEPENDENCIES := postgres

include ../shared/code/development/service.mk
include ../shared/code/development/deps.mk
include ../shared/code/development/migration.mk

buf:
	$(MAKE) -C ../shared/code buf-organization

wire:
	$(MAKE) -C ../shared/code wire organization

generate: buf wire

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/organization-service:$${QHPRO_IMAGE_TAG:-local}" ..

help:
	@printf '%s\n' \
	  'fmt vet test test-race build server dev check' \
	  'deps-up deps-down migrateup migratedown1 migrate-version new_migration' \
	  'buf wire generate docker-build clean'

.PHONY: help deps-up deps-down buf wire generate docker-build
