SHELL := /bin/sh
BINARY_NAME := social-service
RUN_ARGS := grpc

include ../../tools/development/service.mk

buf:
	$(MAKE) -C ../../tools/development buf-social

help:
	@printf '%s\n' 'fmt vet test test-race build server dev buf clean check'

.PHONY: help buf
