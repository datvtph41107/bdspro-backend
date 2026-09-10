SHELL := /bin/sh
include ../tools/development/go.mk

DB_URL ?= postgres://qhpro:qhpro@localhost:55434/qhpro_tqd?sslmode=disable
MIGRATION_DIR := migrate
ENV_FILE ?= .env
ROOT_ENV_FILE ?= ../.env
BIN := bin/tqd-service
GO_TEST_FLAGS ?=

define with_env
set -a; if [ -f "$(ROOT_ENV_FILE)" ]; then . "$(abspath $(ROOT_ENV_FILE))"; fi; if [ -f "$(ENV_FILE)" ]; then . "$(abspath $(ENV_FILE))"; fi; set +a; QHPRO_EXECUTION_MODE=host $(1)
endef

help:
	@printf '%s\n' 'postgres redis migrateup migratedown1 new_migration' \
		'buf wire test vet build server dev pmtiles docker-build clean check'

postgres:
	$(MAKE) -C .. dev-dependencies SERVICE=tqd

redis:
	$(MAKE) -C .. dev-dependencies SERVICE=tqd

deps-up:
	$(MAKE) -C .. dev-dependencies SERVICE=tqd

deps-down:
	$(MAKE) -C .. dev-dependencies-down SERVICE=tqd

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
	$(MAKE) -C ../shared/code buf-tqd

wire:
	$(MAKE) -C ../shared/code wire tqd

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

server:
	@$(call with_env,go run . -server=grpc)

dev:
	@mkdir -p .tmp/air
	@$(call with_env,air -c .air.toml)

pmtiles:
	@$(call with_env,go run . -server=pmtiles)

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/tqd-service:$${QHPRO_IMAGE_TAG:-local}" ..

clean:
	rm -rf bin .tmp

check: vet test build

migrate-up: migrateup
migrate-down: migratedown1
migrate-create: new_migration
run-grpc: server
run-pmtiles: pmtiles
acceptance: check

.PHONY: help postgres redis deps-up deps-down migrateup migrateup1 migratedown \
	migratedown1 migrate-version new_migration migrate-up migrate-down migrate-create \
	buf wire generate fmt vet test test-race build server dev pmtiles run-grpc run-pmtiles \
	docker-build clean check acceptance
