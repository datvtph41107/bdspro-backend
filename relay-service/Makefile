SHELL := /bin/sh
BINARY_NAME := relay-service

include ../tools/development/service.mk

dev:
	@$(call with_env,air -c .air.toml)

help:
	@printf '%s\n' 'fmt vet test test-race build server dev clean check'

.PHONY: dev help
