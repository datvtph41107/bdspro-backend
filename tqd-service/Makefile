SHELL := /bin/sh
BINARY_NAME := tqd-service
RUN_ARGS := -server=grpc
DB_URL ?= postgres://qhpro:qhpro@localhost:5432/qhpro_tqd?sslmode=disable
MIGRATION_DIR := database/migrations
DEV_DEPENDENCIES := postgres redis

include ../shared/code/development/service.mk
include ../shared/code/development/deps.mk
include ../shared/code/development/migration.mk

buf:
	$(MAKE) -C ../shared/code buf-tqd

wire:
	$(MAKE) -C ../shared/code wire tqd

generate: buf wire

pmtiles:
	@$(call with_env,go run . -server=pmtiles)

run-grpc: server
run-pmtiles: pmtiles

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/tqd-service:$${QHPRO_IMAGE_TAG:-local}" ..

acceptance: check

help:
	@printf '%s\n' \
	  'fmt vet test test-race build server dev check' \
	  'deps-up deps-down migrate rollback migration name=<schema_change> migration-version' \
	  'buf wire generate pmtiles docker-build clean'

.PHONY: help buf wire generate pmtiles run-grpc run-pmtiles docker-build acceptance
