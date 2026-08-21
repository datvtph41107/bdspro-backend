SHELL := /usr/bin/env bash
ROOT_DIR := ..
SHARED_CODE := $(ROOT_DIR)/shared/code
PORT ?= 0
HOST ?=

.PHONY: help build run-grpc run-pmtiles migrate migrate-up migrate-down clean bootstrap \
        generate verify-generated test test-race vet lint check run-http run-both

help:
	@printf '%s\n' \
	  'build            Build TQD binary' \
	  'run-grpc         Build and run canonical TQD gRPC process' \
	  'run-pmtiles      Build and run PMTiles process' \
	  'migrate          Build and run application migration command' \
	  'migrate-up       Run v3 SQL migrations up with migrate CLI' \
	  'migrate-down     Run one v3 SQL migration down with migrate CLI' \
	  'bootstrap        Restore module dependencies; may update go.sum if toolchain requires it' \
	  'generate         Mutating: Buf TQD contracts, then generated Wire graph' \
	  'verify-generated Read-only current tree; regenerate in disposable Git worktree' \
	  'test             Run TQD tests' \
	  'test-race        Run TQD tests with race detector' \
	  'vet              Run go vet' \
	  'lint             Run golangci-lint' \
	  'check            Read-only architecture checks + test + vet' \
	  'clean            Remove documented TQD build artifacts'

build:
	@mkdir -p bin
	@go build -o bin/tqd-service .

run-grpc: build
	@./bin/tqd-service -server=grpc $(if $(filter-out 0,$(PORT)),-port=$(PORT),) $(if $(HOST),-host=$(HOST),)

run-pmtiles: build
	@./bin/tqd-service -server=pmtiles $(if $(filter-out 0,$(PORT)),-port=$(PORT),)

migrate: build
	@./bin/tqd-service -server=migrate

# Application HTTP is intentionally disabled. The public application boundary is Gateway HTTP -> TQD gRPC.
run-http run-both:
	@echo 'TQD application HTTP is disabled; use gateway-service HTTP -> TQD gRPC.' >&2
	@exit 2

bootstrap:
	@go mod download

# Generation order is architectural: protobuf symbols must exist before Wire scans constructors.
generate:
	@$(MAKE) -C $(SHARED_CODE) buf-tqd
	@$(MAKE) -C $(SHARED_CODE) wire tqd

verify-generated:
	@bash $(ROOT_DIR)/scripts/architecture/verify-tqd-generated.sh $(ROOT_DIR)

test:
	@go test -count=1 ./...

test-race:
	@go test -count=1 -race ./...

vet:
	@go vet ./...

lint:
	@golangci-lint run

check:
	@bash $(ROOT_DIR)/scripts/architecture/verify-engineering-operating-model.sh $(ROOT_DIR)
	@$(MAKE) test
	@$(MAKE) vet

clean:
	@rm -rf bin/

migrate-up:
	@test -n "$(TQD_DATABASE_URL)" || (echo 'TQD_DATABASE_URL is required' >&2; exit 2)
	@migrate -path infra/db/migrations/v3 -database "$(TQD_DATABASE_URL)" -verbose up

migrate-down:
	@test -n "$(TQD_DATABASE_URL)" || (echo 'TQD_DATABASE_URL is required' >&2; exit 2)
	@migrate -path infra/db/migrations/v3 -database "$(TQD_DATABASE_URL)" -verbose down 1
