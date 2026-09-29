#!/bin/bash

# Script to generate protobuf files for hub-service

PROTOBUF_DIR="shared/protobuf"
SCHEMA_DIR="$PROTOBUF_DIR/schema"

echo "=== Starting Hub Service Protobuf Generation ==="

cd "$SCHEMA_DIR" || exit 1

buf generate --path hub

echo "=== Protobuf Generation Complete ==="
echo "Generated files:"
echo "  - types/hub/*.pb.go"
echo "  - types/hub/*.pb.gw.go"
echo "  - docs/hub/*.swagger.json"

cd - >/dev/null || exit 1

