SHELL := /bin/sh
include ../tools/development/go.mk

DB_URL ?= postgres://postgres:postgres@localhost:5435/payment_service?sslmode=disable
MIGRATION_DIR := migrate
ENV_FILE ?= .env
ROOT_ENV_FILE ?= ../.env
BIN := bin/payment-service
GO_TEST_FLAGS ?=

define with_env
set -a; if [ -f "$(ROOT_ENV_FILE)" ]; then . "$(abspath $(ROOT_ENV_FILE))"; fi; if [ -f "$(ENV_FILE)" ]; then . "$(abspath $(ENV_FILE))"; fi; set +a; QHPRO_EXECUTION_MODE=host $(1)
endef

help:
	@printf '%s\n' 'postgres rabbitmq migrateup migratedown1 new_migration' \
		'buf test vet build server dev docker-build clean check'

postgres:
	$(MAKE) -C .. dev-dependencies SERVICE=payment

rabbitmq:
	$(MAKE) -C .. dev-dependencies SERVICE=payment

deps-up:
	$(MAKE) -C .. dev-dependencies SERVICE=payment

deps-down:
	$(MAKE) -C .. dev-dependencies-down SERVICE=payment

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
	$(MAKE) -C ../shared/code buf-payment

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

server:
	@$(call with_env,go run . grpc)

dev:
	@$(call with_env,air -c .air.toml)

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/payment-service:$${QHPRO_IMAGE_TAG:-local}" ..

compose-up:
	$(MAKE) -C .. dev-up

compose-down:
	$(MAKE) -C .. dev-down

clean:
	rm -rf bin .tmp

check: vet test build

migrate-up: migrateup
migrate-down: migratedown1
migrate-create: new_migration
acceptance-static: check

.PHONY: help postgres rabbitmq deps-up deps-down migrateup migrateup1 migratedown \
	migratedown1 migrate-version new_migration migrate-up migrate-down migrate-create \
	buf fmt vet test test-race build server dev docker-build compose-up compose-down \
	clean check acceptance-static
