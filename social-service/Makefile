SHELL := /bin/sh
BINARY_NAME := social-service
RUN_ARGS := grpc

include ../shared/code/development/service.mk

buf:
	$(MAKE) -C ../shared/code buf-social

help:
	@printf '%s\n' 'fmt vet test test-race build server dev buf clean check'

.PHONY: help buf
