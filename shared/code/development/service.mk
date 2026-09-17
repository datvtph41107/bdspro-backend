# Canonical inner-loop contract for a Go service.
#
# The including service Makefile owns business-specific targets (wire, buf,
# migrations, special actors). This file owns only generic mechanics that are
# identical across services.
.DEFAULT_GOAL := help
include ../shared/code/development/go.mk

ENV_FILE ?= .env
ROOT_ENV_FILE ?= ../.env
REPOSITORY_ROOT ?= ..
QHPRO_LOG_ROOT := $(abspath $(REPOSITORY_ROOT)/.tmp/development/logs)
BINARY_NAME ?= service
BIN ?= bin/$(BINARY_NAME)
BUILD_PACKAGE ?= .
RUN_PACKAGE ?= .
RUN_ARGS ?=
GO_TEST_FLAGS ?=

# Native development intentionally loads root environment first, then the
# service-owned environment. Docker Compose injects its own container wiring
# and does not use this helper.
define with_env
set -a; if [ -f "$(ROOT_ENV_FILE)" ]; then . "$(abspath $(ROOT_ENV_FILE))"; fi; if [ -f "$(ENV_FILE)" ]; then . "$(abspath $(ENV_FILE))"; fi; set +a; QHPRO_EXECUTION_MODE=host QHPRO_LOG_ROOT="$(QHPRO_LOG_ROOT)" $(1)
endef

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test $(GO_TEST_FLAGS) -count=1 ./...

test-race:
	go test $(GO_TEST_FLAGS) -race -count=1 ./...

build:
	mkdir -p bin
	go build -o $(BIN) $(BUILD_PACKAGE)

server run:
	@$(call with_env,go run $(RUN_PACKAGE) $(RUN_ARGS))

dev:
	@test -f .air.toml || { echo 'This service has no .air.toml; use make server or add an owner-specific dev target.' >&2; exit 2; }
	@mkdir -p .tmp/air
	@$(call with_env,air -c .air.toml)

clean:
	rm -rf bin .tmp

check: vet test build

.PHONY: fmt vet test test-race build server run dev clean check
