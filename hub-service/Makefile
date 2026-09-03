SHELL := /bin/sh
include ../tools/development/go.mk

DB_URL ?= postgres://postgres:postgres@localhost:5438/hub_service?sslmode=disable
MIGRATION_DIR := migrate
ENV_FILE ?= .env
ROOT_ENV_FILE ?= ../.env
BIN := bin/hub-service
GO_TEST_FLAGS ?=

define with_env
set -a; if [ -f "$(ROOT_ENV_FILE)" ]; then . "$(abspath $(ROOT_ENV_FILE))"; fi; if [ -f "$(ENV_FILE)" ]; then . "$(abspath $(ENV_FILE))"; fi; set +a; QHPRO_EXECUTION_MODE=host $(1)
endef

help:
	@printf '%s\n' 'migrateup migratedown1 new_migration' \
		'buf test vet build server dev docker-build clean check'

migrateup:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$${MIGRATION_URL:-$${HUB_DATABASE_URL:-$(DB_URL)}}" -verbose up)

migrateup1:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$${MIGRATION_URL:-$${HUB_DATABASE_URL:-$(DB_URL)}}" -verbose up 1)

migratedown:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$${MIGRATION_URL:-$${HUB_DATABASE_URL:-$(DB_URL)}}" -verbose down)

migratedown1:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$${MIGRATION_URL:-$${HUB_DATABASE_URL:-$(DB_URL)}}" -verbose down 1)

migrate-version:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$${MIGRATION_URL:-$${HUB_DATABASE_URL:-$(DB_URL)}}" version)

new_migration:
	@test -n "$(name)" || { echo 'usage: make new_migration name=<schema_change>' >&2; exit 2; }
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(name)

buf:
	$(MAKE) -C ../shared/code buf-hub

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
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/hub-service:$${QHPRO_IMAGE_TAG:-local}" ..

clean:
	rm -rf bin .tmp

check: vet test build

migrate-up: migrateup
migrate-down: migratedown1
migrate-create: new_migration

.PHONY: help migrateup migrateup1 migratedown migratedown1 migrate-version \
	new_migration migrate-up migrate-down migrate-create buf fmt vet test test-race \
	build server dev docker-build clean check
