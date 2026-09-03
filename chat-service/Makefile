SHELL := /bin/sh
BINARY_NAME := chat-service
RUN_ARGS := grpc

include ../tools/development/service.mk

buf:
	$(MAKE) -C ../shared/code buf-chat

dev:
	@$(call with_env,air -c .air.toml)

help:
	@printf '%s\n' 'fmt vet test test-race build server dev buf clean check'

.PHONY: buf dev help
