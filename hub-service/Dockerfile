# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
COPY ../shared /shared

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o utility-service .

# Stage 2: Create minimal runtime image
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/utility-service .
COPY --from=builder /app/config ./config

# Expose ports
EXPOSE 8030 50030

# Set environment
ENV ENV_RUNTIME=product

# Run the binary
CMD ["./utility-service", "grpc"]

