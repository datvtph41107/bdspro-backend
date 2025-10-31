#!/bin/bash

# Script to generate protobuf files for hub-service

PROTOBUF_DIR="shared/protobuf"
SCHEMA_DIR="$PROTOBUF_DIR/schema/hub"

echo "=== Starting Hub Service Protobuf Generation ==="

# Change to protobuf directory
cd $PROTOBUF_DIR || exit

echo "Generating location_v2.proto..."
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=. \
  --grpc-gateway_opt=paths=source_relative \
  --grpc-gateway_opt=logtostderr=true \
  --grpc-gateway_opt=generate_unbound_methods=true \
  --openapiv2_out=docs/hub \
  --openapiv2_opt=logtostderr=true \
  --proto_path=schema \
  schema/hub/location_v2.proto

echo "Generating location.proto..."
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=. \
  --grpc-gateway_opt=paths=source_relative \
  --grpc-gateway_opt=logtostderr=true \
  --grpc-gateway_opt=generate_unbound_methods=true \
  --openapiv2_out=docs/hub \
  --openapiv2_opt=logtostderr=true \
  --proto_path=schema \
  schema/hub/location.proto

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
  --openapiv2_out=docs/hub \
  --openapiv2_opt=logtostderr=true \
  --proto_path=schema \
  schema/hub/event_queue.proto

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
  --openapiv2_out=docs/hub \
  --openapiv2_opt=logtostderr=true \
  --proto_path=schema \
  schema/hub/internal.proto

echo "Generating user_guide.proto..."
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=. \
  --grpc-gateway_opt=paths=source_relative \
  --grpc-gateway_opt=logtostderr=true \
  --grpc-gateway_opt=generate_unbound_methods=true \
  --openapiv2_out=docs/hub \
  --openapiv2_opt=logtostderr=true \
  --proto_path=schema \
  schema/hub/user_guide.proto

echo "=== Protobuf Generation Complete ==="
echo "Generated files:"
echo "  - types/hub/*.pb.go"
echo "  - types/hub/*.pb.gw.go"
echo "  - docs/hub/*.swagger.json"

cd - > /dev/null

