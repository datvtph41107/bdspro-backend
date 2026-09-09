SHELL := /bin/sh
BINARY_NAME := assistant-service

include ../shared/code/development/service.mk

wire:
	$(MAKE) -C ../shared/code wire assistant

buf:
	$(MAKE) -C ../shared/code buf-assistant

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/assistant-service:$${QHPRO_IMAGE_TAG:-local}" ..

help:
	@printf '%s\n' 'fmt vet test test-race build server wire buf docker-build clean check'

.PHONY: help wire buf docker-build
