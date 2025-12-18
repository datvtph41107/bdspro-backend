#!/bin/bash

# Script to generate protobuf files for bdspro-service

PROTOBUF_DIR="shared/protobuf"
SCHEMA_DIR="$PROTOBUF_DIR/schema"

echo "=== Starting BDSPro Service Protobuf Generation ==="

cd "$SCHEMA_DIR" || exit 1

buf generate --path bdspro

echo "=== Protobuf Generation Complete ==="
echo "Generated files:"
echo "  - types/bdspro/*.pb.go"
echo "  - types/bdspro/*.pb.gw.go"
echo "  - docs/bdspro/*.swagger.json"

cd - >/dev/null || exit 1
