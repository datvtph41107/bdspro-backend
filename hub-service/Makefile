SHELL := /bin/sh
BINARY_NAME := hub-service
RUN_ARGS := grpc
DB_URL ?= postgres://postgres:postgres@localhost:5432/hub_service?sslmode=disable
MIGRATION_DIR := database/migrations
MIGRATION_URL_FALLBACK_ENV := HUB_DATABASE_URL
DEV_DEPENDENCIES := postgres redis

include ../shared/code/development/service.mk
include ../shared/code/development/deps.mk
include ../shared/code/development/migration.mk

buf:
	$(MAKE) -C ../shared/code buf-hub

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/hub-service:$${QHPRO_IMAGE_TAG:-local}" ..

help:
	@printf '%s\n' \
	  'fmt vet test test-race build server dev check' \
	  'deps-up deps-down migrateup migratedown1 migrate-version new_migration' \
	  'buf docker-build clean'

.PHONY: help deps-up deps-down buf docker-build
