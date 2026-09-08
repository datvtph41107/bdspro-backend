#!/bin/bash

# Script to generate protobuf files for utility-service

PROTOBUF_DIR="shared/protobuf"
SCHEMA_DIR="$PROTOBUF_DIR/schema/utility"

echo "=== Starting Utility Service Protobuf Generation ==="

# Change to protobuf directory
cd $PROTOBUF_DIR || exit

echo "Generating event_queue.proto..."
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=. \
  --grpc-gateway_opt=paths=source_relative \
  --grpc-gateway_opt=logtostderr=true \
  --grpc-gateway_opt=generate_unbound_methods=true \
  --openapiv2_out=docs/utility \
  --openapiv2_opt=logtostderr=true \
  --proto_path=schema \
  --proto_path=../../../googleapis \
  schema/utility/event_queue.proto

echo "Generating internal.proto..."
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=. \
  --grpc-gateway_opt=paths=source_relative \
  --grpc-gateway_opt=logtostderr=true \
  --grpc-gateway_opt=generate_unbound_methods=true \
  --openapiv2_out=docs/utility \
  --openapiv2_opt=logtostderr=true \
  --proto_path=schema \
  --proto_path=../../../googleapis \
  schema/utility/internal.proto

echo "=== Protobuf Generation Complete ==="
echo "Generated files:"
echo "  - types/utility/*.pb.go"
echo "  - types/utility/*.pb.gw.go"
echo "  - docs/utility/*.swagger.json"

cd - > /dev/null

