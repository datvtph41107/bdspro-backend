SHELL := /bin/sh
BINARY_NAME := chat-service
RUN_ARGS := grpc

include ../shared/code/development/service.mk

buf:
	$(MAKE) -C ../shared/code buf-chat

help:
	@printf '%s\n' 'fmt vet test test-race build server dev buf clean check'

.PHONY: help buf
