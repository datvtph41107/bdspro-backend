SHELL := /bin/sh
BINARY_NAME := assistant-service

include ../../tools/development/service.mk

wire:
	$(MAKE) -C ../../tools/development wire assistant

buf:
	$(MAKE) -C ../../tools/development buf-assistant

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/assistant-service:$${QHPRO_IMAGE_TAG:-local}" ..

help:
	@printf '%s\n' 'fmt vet test test-race build server wire buf docker-build clean check'

.PHONY: help wire buf docker-build
