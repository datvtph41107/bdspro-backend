#!/bin/bash

# Script to generate protobuf files for tqd-service
# This should be run from the project root directory

echo "Generating protobuf files for tqd-service with buf..."

# Navigate to shared/protobuf directory
cd shared/protobuf

buf generate --path schema/tqd

echo "Protobuf generation completed for tqd-service!"
echo "Generated files are in pb/types/tqd/ and docs/tqd/"
