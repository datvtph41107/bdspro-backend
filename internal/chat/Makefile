SHELL := /bin/sh
BINARY_NAME := chat-service
RUN_ARGS := grpc

include ../../tools/development/service.mk

buf:
	$(MAKE) -C ../../tools/development buf-chat

help:
	@printf '%s\n' 'fmt vet test test-race build server dev buf clean check'

.PHONY: help buf
