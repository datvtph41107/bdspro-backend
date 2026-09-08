#!/bin/bash

# Script to merge all swagger JSON files and add Bearer authentication
# Usage: ./merge_swagger.sh

set -e

echo "🚀 Starting Swagger merge process..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Change to the protobuf directory
echo "$(dirname "$0")"
cd "$(dirname "$0")/../../../gateway-service"

# Check if docs directory exists
if [ ! -d "docs" ]; then
    echo "❌ docs directory not found. Please make sure you're in the protobuf directory."
    exit 1
fi

# Run the Go script
echo "📁 Processing swagger files in docs directory..."
go run scripts/merge_swagger.go

# Check if the output file was created
if [ -f "merged_swagger.json" ]; then
    echo "✅ Successfully created merged_swagger.json"
    echo "📊 File size: $(du -h merged_swagger.json | cut -f1)"
    echo "🔗 You can now use this file with Swagger UI or other documentation tools"
else
    echo "❌ Failed to create merged_swagger.json"
    exit 1
fi

echo "🎉 Merge process completed successfully!" 