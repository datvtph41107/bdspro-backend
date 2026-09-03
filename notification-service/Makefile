SHELL := /bin/sh
include ../tools/development/go.mk

DB_URL ?= postgres://qhpro:qhpro@localhost:55435/qhpro_notification?sslmode=disable
MIGRATION_DIR := migrate
ENV_FILE ?= .env
ROOT_ENV_FILE ?= ../.env
BIN := bin/notification-service
EVENT_WORKER_BIN := bin/payment-event-worker
DELIVERY_WORKER_BIN := bin/delivery-worker
GO_TEST_FLAGS ?=

define with_env
set -a; if [ -f "$(ROOT_ENV_FILE)" ]; then . "$(abspath $(ROOT_ENV_FILE))"; fi; if [ -f "$(ENV_FILE)" ]; then . "$(abspath $(ENV_FILE))"; fi; set +a; QHPRO_EXECUTION_MODE=host $(1)
endef

help:
	@printf '%s\n' 'postgres redis rabbitmq migrateup migratedown1 new_migration' \
		'buf wire test vet build server dev event-worker delivery-worker docker-build clean check'

postgres:
	$(MAKE) -C .. dev-dependencies SERVICE=notification

redis:
	$(MAKE) -C .. dev-dependencies SERVICE=notification

rabbitmq:
	$(MAKE) -C .. dev-dependencies SERVICE=notification

deps-up:
	$(MAKE) -C .. dev-dependencies SERVICE=notification

deps-down:
	$(MAKE) -C .. dev-dependencies-down SERVICE=notification

migrateup:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$${MIGRATION_URL:-$(DB_URL)}" -verbose up)

migrateup1:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$${MIGRATION_URL:-$(DB_URL)}" -verbose up 1)

migratedown:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$${MIGRATION_URL:-$(DB_URL)}" -verbose down)

migratedown1:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$${MIGRATION_URL:-$(DB_URL)}" -verbose down 1)

migrate-version:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$${MIGRATION_URL:-$(DB_URL)}" version)

new_migration:
	@test -n "$(name)" || { echo 'usage: make new_migration name=<schema_change>' >&2; exit 2; }
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(name)

buf:
	$(MAKE) -C ../shared/code buf-notification

wire:
	$(MAKE) -C ../shared/code wire notification

generate: buf wire

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test $(GO_TEST_FLAGS) -count=1 ./...

test-race:
	go test -race -count=1 ./...

build:
	mkdir -p bin
	go build -o $(BIN) .
	go build -o $(EVENT_WORKER_BIN) ./cmd/payment-event-worker
	go build -o $(DELIVERY_WORKER_BIN) ./cmd/delivery-worker

server:
	@$(call with_env,go run . grpc)

dev:
	@$(call with_env,air -c .air.toml)

event-worker:
	@$(call with_env,go run ./cmd/payment-event-worker)

delivery-worker:
	@$(call with_env,go run ./cmd/delivery-worker)

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/notification-service:$${QHPRO_IMAGE_TAG:-local}" ..

clean:
	rm -rf bin .tmp

check: vet test build

migrate-up: migrateup
migrate-down: migratedown1
migrate-create: new_migration
run-grpc: server
run-event-worker: event-worker
run-delivery-worker: delivery-worker
acceptance: check

.PHONY: help postgres redis rabbitmq deps-up deps-down migrateup migrateup1 \
	migratedown migratedown1 migrate-version new_migration migrate-up migrate-down \
	migrate-create buf wire generate fmt vet test test-race build server dev event-worker \
	delivery-worker docker-build clean check run-grpc run-event-worker \
	run-delivery-worker acceptance
