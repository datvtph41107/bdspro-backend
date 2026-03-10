# TQD Service Makefile

.PHONY: build run-http run-grpc run-both clean help test lint deps

# Build the main binary
build:
	@echo "Building TQD Service..."
	@go build -o bin/tqd-service .

# Run HTTP server
run-http: build
	@echo "Starting HTTP server..."
	@./bin/tqd-service http

# Run gRPC server
run-grpc: build
	@echo "Starting gRPC server..."
	@./bin/tqd-service grpc

# Run both servers
run-both: build
	@echo "Starting both HTTP and gRPC servers..."
	@./bin/tqd-service both

# Run HTTP server in background
run-http-bg: build
	@echo "Starting HTTP server in background..."
	@./bin/tqd-service http &

# Run gRPC server in background
run-grpc-bg: build
	@echo "Starting gRPC server in background..."
	@./bin/tqd-service grpc &

# Run both servers in background
run-all-bg: build
	@echo "Starting both servers in background..."
	@./bin/tqd-service both &

# Run HTTP server on custom port
run-http-port: build
	@echo "Starting HTTP server on custom port..."
	@./bin/tqd-service http --port $(PORT)

# Run gRPC server on custom port
run-grpc-port: build
	@echo "Starting gRPC server on custom port..."
	@./bin/tqd-service grpc --port $(PORT)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Run linting
lint:
	@echo "Running linter..."
	@golangci-lint run

# Show help
help:
	@echo "Available commands:"
	@echo "  build          Build the main binary"
	@echo "  run-http       Build and run HTTP server"
	@echo "  run-grpc       Build and run gRPC server"
	@echo "  run-both       Build and run both servers"
	@echo "  run-http-bg    Run HTTP server in background"
	@echo "  run-grpc-bg    Run gRPC server in background"
	@echo "  run-all-bg     Run both servers in background"
	@echo "  run-http-port  Run HTTP server on custom port (set PORT=xxxx)"
	@echo "  run-grpc-port  Run gRPC server on custom port (set PORT=xxxx)"
	@echo "  clean          Clean build artifacts"
	@echo "  deps           Install dependencies"
	@echo "  test           Run tests"
	@echo "  lint           Run linter"
	@echo "  help           Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make run-http-port PORT=8080"
	@echo "  make run-grpc-port PORT=9090"

migrateup:
	migrate -path infra/db/migrations/v2 -database "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable" -verbose up

migratedown:
	migrate -path infra/db/migrations/v2 -database "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable" -verbose down

