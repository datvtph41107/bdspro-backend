#!/bin/bash

# Script to generate protobuf files for feedback-service
# This should be run from the project root directory

echo "Generating protobuf files for feedback-service..."

# Navigate to shared/protobuf directory
cd shared/protobuf

# Generate Go files from rate.proto
echo "Generating Go files from rate.proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    schema/feedback/rate.proto

# Generate Go files from report.proto
echo "Generating Go files from report.proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    schema/feedback/report.proto

# Generate Go files from report_reason.proto
echo "Generating Go files from report_reason.proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    schema/feedback/report_reason.proto

# Generate Go files from registry.proto
echo "Generating Go files from registry.proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    schema/feedback/registry.proto

# Generate Go files from internal.proto
echo "Generating Go files from internal.proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    schema/feedback/internal.proto

# Generate swagger files
echo "Generating swagger files..."
protoc --openapiv2_out=docs/feedback \
    --openapiv2_opt=logtostderr=true \
    schema/feedback/rate.proto

protoc --openapiv2_out=docs/feedback \
    --openapiv2_opt=logtostderr=true \
    schema/feedback/report.proto

protoc --openapiv2_out=docs/feedback \
    --openapiv2_opt=logtostderr=true \
    schema/feedback/report_reason.proto

protoc --openapiv2_out=docs/feedback \
    --openapiv2_opt=logtostderr=true \
    schema/feedback/registry.proto

protoc --openapiv2_out=docs/feedback \
    --openapiv2_opt=logtostderr=true \
    schema/feedback/internal.proto

# Generate gateway files
echo "Generating gateway files..."
protoc --grpc-gateway_out=. \
    --grpc-gateway_opt=paths=source_relative \
    --grpc-gateway_opt=generate_unbound_methods=true \
    schema/feedback/rate.proto

protoc --grpc-gateway_out=. \
    --grpc-gateway_opt=paths=source_relative \
    --grpc-gateway_opt=generate_unbound_methods=true \
    schema/feedback/report.proto

protoc --grpc-gateway_out=. \
    --grpc-gateway_opt=paths=source_relative \
    --grpc-gateway_opt=generate_unbound_methods=true \
    schema/feedback/report_reason.proto

protoc --grpc-gateway_out=. \
    --grpc-gateway_opt=paths=source_relative \
    --grpc-gateway_opt=generate_unbound_methods=true \
    schema/feedback/registry.proto

protoc --grpc-gateway_out=. \
    --grpc-gateway_opt=paths=source_relative \
    --grpc-gateway_opt=generate_unbound_methods=true \
    schema/feedback/internal.proto

echo "Protobuf generation completed for feedback-service!"
echo "Generated files:"
echo "- pb/types/feedback/rate.pb.go"
echo "- pb/types/feedback/rate_grpc.pb.go"
echo "- pb/types/feedback/report.pb.go"
echo "- pb/types/feedback/report_grpc.pb.go"
echo "- pb/types/feedback/report_reason.pb.go"
echo "- pb/types/feedback/report_reason_grpc.pb.go"
echo "- pb/types/feedback/registry.pb.go"
echo "- pb/types/feedback/registry_grpc.pb.go"
echo "- pb/types/feedback/internal.pb.go"
echo "- pb/types/feedback/internal_grpc.pb.go"
echo "- docs/feedback/rate.swagger.json"
echo "- docs/feedback/report.swagger.json"
echo "- docs/feedback/report_reason.swagger.json"
echo "- docs/feedback/registry.swagger.json"
echo "- docs/feedback/internal.swagger.json"
echo "- pb/types/feedback/rate.pb.gw.go"
echo "- pb/types/feedback/report.pb.gw.go"
echo "- pb/types/feedback/report_reason.pb.gw.go"
