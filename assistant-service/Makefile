SHELL := /bin/sh
BINARY_NAME := assistant-service

include ../shared/code/development/service.mk

buf:
	$(MAKE) -C ../shared/code buf-assistant

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/assistant-service:$${QHPRO_IMAGE_TAG:-local}" ..

help:
	@printf '%s\n' 'fmt vet test test-race build server buf docker-build clean check'

.PHONY: help buf docker-build
