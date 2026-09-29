SHELL := /bin/sh
BINARY_NAME := payment-service
RUN_ARGS := grpc
DB_URL ?= postgres://postgres:postgres@localhost:5432/payment_service?sslmode=disable
MIGRATION_DIR := database/migrations
DEV_DEPENDENCIES := postgres rabbitmq

include ../../tools/development/service.mk
include ../../tools/development/deps.mk
include ../../tools/development/migration.mk

buf:
	$(MAKE) -C ../../tools/development buf-payment

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/payment-service:$${QHPRO_IMAGE_TAG:-local}" ..

compose-up:
	$(MAKE) -C .. up

compose-down:
	$(MAKE) -C .. down

acceptance-static: check

help:
	@printf '%s\n' \
	  'fmt vet test test-race build server dev check' \
	  'deps-up deps-down migrate rollback migration name=<schema_change> migration-version' \
	  'buf docker-build compose-up compose-down clean'

.PHONY: help buf docker-build compose-up compose-down acceptance-static
