IMAGE_NAME=go_chat_service

SHELL := /bin/bash
PROJECT_NAME = go-chat-service
BUILD_FOLDER = ./dist
BINARY_PATH = $(BUILD_FOLDER)/$(PROJECT_NAME)

.PHONY: setup
setup:
	@if [[ ! -f ./setup.sh ]]; then \
		echo "Setup is already complete. You can delete this setup make target."; \
	else \
		chmod 755 ./setup.sh && ./setup.sh && rm ./setup.sh; \
	fi

.PHONY: build
build:
	echo "Building $(PROJECT_NAME)..."
	go build -o $(BINARY_PATH)

.PHONY: generate-mock
generate-mock:
	mockery

.PHONY: http-serve
http-serve: build
	echo "Starting $(PROJECT_NAME) server..."
	$(BINARY_PATH) http-serve --httpPort=8081

.PHONY: grpc-serve
grpc-serve: build
	echo "Starting $(PROJECT_NAME) server..."
	$(BINARY_PATH) grpc-serve --grpcPort=50053

.PHONY: ws-serve
ws-serve: build
	echo "Starting $(PROJECT_NAME) server..."
	$(BINARY_PATH) ws-serve

.PHONY: prepare-swagger
prepare-swagger: build
	echo "Prepare swagger"
	$(BINARY_PATH) prepare-swagger
	

.PHONY: swagger-serve
swagger-serve: build prepare-swagger
	echo "Starting $(PROJECT_NAME) server..."
	$(BINARY_PATH) swagger-serve

.PHONY: docker-build
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t $(IMAGE_NAME) .


.PHONY: docker-run
docker-run:
	@echo "🚀 Running Docker container..."
	docker run --rm -p 8080:8080 --name $(IMAGE_NAME) $(IMAGE_NAME)


.PHONY: up
up:
	@echo "🐳 Starting Docker Compose..."
	docker-compose up --build -d


.PHONY: down
down:
	@echo "🛑 Stopping Docker Compose..."
	docker-compose down


.PHONY: clean
clean:
	@echo "🧹 Cleaning up Docker images..."
	docker rmi -f $(IMAGE_NAME)