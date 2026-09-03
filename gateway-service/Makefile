SHELL := /bin/sh
include ../tools/development/go.mk

ENV_FILE ?= .env
ROOT_ENV_FILE ?= ../.env
BIN := bin/gateway-service
GO_TEST_FLAGS ?=

define with_env
set -a; if [ -f "$(ROOT_ENV_FILE)" ]; then . "$(abspath $(ROOT_ENV_FILE))"; fi; if [ -f "$(ENV_FILE)" ]; then . "$(abspath $(ENV_FILE))"; fi; set +a; QHPRO_EXECUTION_MODE=host $(1)
endef

.PHONY: help fmt run server dev test test-race race vet build docker-build clean check smoke

help:
	@printf '%s\n' 'run dev test vet build docker-build clean check smoke'

fmt:
	go fmt ./...

run server:
	@$(call with_env,go run . http)

dev:
	@$(call with_env,air -c .air.toml)

test:
	go test $(GO_TEST_FLAGS) -count=1 ./...

test-race race:
	go test $(GO_TEST_FLAGS) -count=1 -race ./...

vet:
	go vet ./...

build:
	mkdir -p bin
	go build -o $(BIN) .

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/gateway-service:$${QHPRO_IMAGE_TAG:-local}" ..

clean:
	rm -rf bin .tmp

check: vet test build

smoke:
	go test -count=1 \
		./config \
		./internal/authroutes \
		./internal/grpcgateway \
		./internal/grpcmetadata \
		./internal/httpauth \
		./internal/httperror \
		./internal/httpmiddleware \
		./internal/httpresponse \
		./internal/tqdmultipart \
		./internal/upstream \
		./internal/v1proxy \
		./cmd/http
