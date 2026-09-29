SHELL := /bin/sh
BINARY_NAME := notification-service
RUN_ARGS := grpc
DB_URL ?= postgres://qhpro:qhpro@localhost:5432/qhpro_notification?sslmode=disable
MIGRATION_DIR := database/migrations
DEV_DEPENDENCIES := postgres redis rabbitmq
EVENT_WORKER_BIN := bin/payment-event-worker
DELIVERY_WORKER_BIN := bin/delivery-worker

include ../../tools/development/service.mk
include ../../tools/development/deps.mk
include ../../tools/development/migration.mk

buf:
	$(MAKE) -C ../../tools/development buf-notification

wire:
	$(MAKE) -C ../../tools/development wire notification

generate: buf wire

# Compatibility/recovery binaries are still buildable, but the canonical local
# Compose process already owns these actors and does not run standing worker
# containers for them.
build: build-compat-actors
build-compat-actors:
	mkdir -p bin
	go build -o $(EVENT_WORKER_BIN) ./cmd/payment-event-worker
	go build -o $(DELIVERY_WORKER_BIN) ./cmd/delivery-worker

event-worker:
	@$(call with_env,go run ./cmd/payment-event-worker)

delivery-worker:
	@$(call with_env,go run ./cmd/delivery-worker)

run-grpc: server
run-event-worker: event-worker
run-delivery-worker: delivery-worker

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/notification-service:$${QHPRO_IMAGE_TAG:-local}" ..

acceptance: check

help:
	@printf '%s\n' \
	  'fmt vet test test-race build server dev check' \
	  'deps-up deps-down migrate rollback migration name=<schema_change> migration-version' \
	  'buf wire generate event-worker delivery-worker docker-build clean'

.PHONY: help buf wire generate \
	build-compat-actors event-worker delivery-worker run-grpc run-event-worker \
	run-delivery-worker docker-build acceptance
