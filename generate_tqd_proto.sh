#!/bin/bash

# Script to generate protobuf files for tqd-service
# This should be run from the project root directory

echo "Generating protobuf files for tqd-service..."

# Navigate to shared/protobuf directory
cd shared/protobuf

# Generate Go files from contact_label.proto
echo "Generating Go files from contact_label.proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    schema/tqd/contact_label.proto

# Generate Go files from service_category.proto
echo "Generating Go files from service_category.proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    schema/tqd/service_category.proto

# Generate Go files from directory_supplier.proto
echo "Generating Go files from directory_supplier.proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    schema/tqd/directory_supplier.proto

# Generate swagger files
echo "Generating swagger files..."
protoc --openapiv2_out=docs/tqd \
    --openapiv2_opt=logtostderr=true \
    schema/tqd/contact_label.proto

protoc --openapiv2_out=docs/tqd \
    --openapiv2_opt=logtostderr=true \
    schema/tqd/service_category.proto

protoc --openapiv2_out=docs/tqd \
    --openapiv2_opt=logtostderr=true \
    schema/tqd/directory_supplier.proto

# Generate gateway files
echo "Generating gateway files..."
protoc --grpc-gateway_out=. \
    --grpc-gateway_opt=paths=source_relative \
    --grpc-gateway_opt=generate_unbound_methods=true \
    schema/tqd/contact_label.proto

protoc --grpc-gateway_out=. \
    --grpc-gateway_opt=paths=source_relative \
    --grpc-gateway_opt=generate_unbound_methods=true \
    schema/tqd/service_category.proto

protoc --grpc-gateway_out=. \
    --grpc-gateway_opt=paths=source_relative \
    --grpc-gateway_opt=generate_unbound_methods=true \
    schema/tqd/directory_supplier.proto

echo "Protobuf generation completed for tqd-service!"
echo "Generated files:"
echo "- pb/types/tqd/contact_label.pb.go"
echo "- pb/types/tqd/contact_label_grpc.pb.go"
echo "- pb/types/tqd/service_category.pb.go"
echo "- pb/types/tqd/service_category_grpc.pb.go"
echo "- docs/tqd/contact_label.swagger.json"
echo "- docs/tqd/service_category.swagger.json"
echo "- pb/types/tqd/contact_label.pb.gw.go"
echo "- pb/types/tqd/service_category.pb.gw.go"
