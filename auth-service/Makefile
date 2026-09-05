SHELL := /bin/sh
include ../tools/development/go.mk

ENV_FILE ?= .env
ROOT_ENV_FILE ?= ../.env
BIN := bin/auth-service
GO_TEST_FLAGS ?=

define with_env
set -a; if [ -f "$(ROOT_ENV_FILE)" ]; then . "$(abspath $(ROOT_ENV_FILE))"; fi; if [ -f "$(ENV_FILE)" ]; then . "$(abspath $(ENV_FILE))"; fi; set +a; QHPRO_EXECUTION_MODE=host $(1)
endef

help:
	@printf '%s\n' 'buf test vet build server dev docker-build clean check'

buf:
	$(MAKE) -C ../shared/code buf-auth

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
	@$(call with_env,go run .)

dev:
	@$(call with_env,air -c .air.toml)

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/auth-service:$${QHPRO_IMAGE_TAG:-local}" ..

clean:
	rm -rf bin .tmp

check: vet test build

.PHONY: help buf fmt vet test test-race build server dev docker-build clean check
