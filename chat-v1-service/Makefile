SHELL := /bin/sh
BINARY_NAME := chat-v1-service
RUN_ARGS := grpc

include ../tools/development/service.mk

buf:
	$(MAKE) -C ../shared/code buf-chat-v1

dev:
	@$(call with_env,air -c .air.toml)

help:
	@printf '%s\n' 'fmt vet test test-race build server dev buf clean check'

.PHONY: buf dev help
