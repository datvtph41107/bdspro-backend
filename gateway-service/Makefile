SHELL := /bin/sh
BINARY_NAME := gateway-service
RUN_ARGS := http

include ../shared/code/development/service.mk

race: test-race

docker-build:
	docker build -f Dockerfile -t "$${QHPRO_IMAGE_NAMESPACE:-bdspro}/gateway-service:$${QHPRO_IMAGE_TAG:-local}" ..

smoke:
	go test -count=1 \
		./config \
		./internal/authroutes \
		./internal/grpcgateway \
		./internal/grpcmetadata \
		./internal/httpauth \
		./internal/httperror \
		./internal/httpmiddleware \
		./internal/httpresponse \
		./internal/tqdmultipart \
		./internal/upstream \
		./internal/v1proxy \
		./cmd/http

help:
	@printf '%s\n' 'fmt vet test test-race build server dev docker-build smoke clean check'

.PHONY: help race docker-build smoke
