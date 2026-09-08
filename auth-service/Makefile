SHELL := /bin/sh
BINARY_NAME := auth-service

include ../shared/code/development/service.mk

buf:
	$(MAKE) -C ../shared/code buf-auth

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/auth-service:$${QHPRO_IMAGE_TAG:-local}" ..

help:
	@printf '%s\n' 'fmt vet test test-race build server dev buf docker-build clean check'

.PHONY: help buf docker-build
