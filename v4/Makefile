SHELL := /bin/sh
BINARY_NAME := userv4
BUILD_PACKAGE := ./cmd/grpc
RUN_PACKAGE := ./cmd/grpc
# Alternate entry: make http  (HTTP_ADDRESS / server.port, default :8288)

http:
	@$(call with_env,go run ./cmd/http)

.PHONY: http
DB_URL ?= postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable
MIGRATION_DIR := migrate
DEV_DEPENDENCIES := postgres redis

include ../shared/code/development/service.mk
include ../shared/code/development/deps.mk
include ../shared/code/development/migration.mk

buf:
	$(MAKE) -C ../shared/code buf-user-service

wire:
	(cd wire && wire gen)

generate: buf wire

# One-time empty-installation operation. Inputs are supplied by the operator;
# this target never invents or persists default credentials.
bootstrap-admin:
	@$(call with_env,go run ./cmd/bootstrap-admin)

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/userv4:$${QHPRO_IMAGE_TAG:-local}" ..

# Ship linux binary + config to production (see ./deploy.sh).
deploy:
	./deploy.sh $(if $(target),$(target),http)

deploy-compose:
	./deploy.sh compose

compose-up:
	$(MAKE) -C .. up

compose-down:
	$(MAKE) -C .. down

compose-logs:
	$(MAKE) -C .. logs service=user-service

acceptance-static: check

help:
	@printf '%s\n' \
	  'fmt vet test test-race build server dev check' \
	  'deps-up deps-down migrateup migratedown1 migrate-version new_migration' \
	  'bootstrap-admin buf wire generate docker-build compose-up compose-down compose-logs clean' \
	  'deploy [target=http|grpc] deploy-compose'

.PHONY: help buf wire generate bootstrap-admin docker-build deploy deploy-compose \
	compose-up compose-down compose-logs acceptance-static
