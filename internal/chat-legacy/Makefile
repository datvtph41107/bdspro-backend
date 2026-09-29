SHELL := /bin/sh
BINARY_NAME := chat-v1-service
RUN_ARGS := grpc

include ../../tools/development/service.mk

buf:
	$(MAKE) -C ../../tools/development buf-chat-v1

help:
	@printf '%s\n' 'fmt vet test test-race build server dev buf clean check'

.PHONY: help buf
