# Shared inner-loop contract for a Go service. The including Makefile must set
# BINARY_NAME and may set RUN_ARGS. Service-specific migration/generate targets
# remain with the service owner.
.DEFAULT_GOAL := help
include ../tools/development/go.mk

ENV_FILE ?= .env
ROOT_ENV_FILE ?= ../.env
BIN ?= bin/$(BINARY_NAME)
GO_TEST_FLAGS ?=

define with_env
set -a; if [ -f "$(ROOT_ENV_FILE)" ]; then . "$(abspath $(ROOT_ENV_FILE))"; fi; if [ -f "$(ENV_FILE)" ]; then . "$(abspath $(ENV_FILE))"; fi; set +a; QHPRO_EXECUTION_MODE=host $(1)
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
	go build -o $(BIN) .

server run:
	@$(call with_env,go run . $(RUN_ARGS))

clean:
	rm -rf bin .tmp

check: vet test build

.PHONY: fmt vet test test-race build server run clean check
