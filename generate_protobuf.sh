#!/bin/bash

echo "Generating protobuf files..."

# Generate protobuf for notification service
cd protobuf
buf generate

echo "Protobuf generation completed!"

# Build notification service
echo "Building notification service..."
cd ../notification-service
go mod tidy
go build -o notification-service .

# Build marketing service  
echo "Building marketing service..."
cd ../marketing-service
go mod tidy
go build -o marketing-service .

echo "Build completed!" 