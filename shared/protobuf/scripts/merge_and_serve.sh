#!/bin/bash

# Script to merge swagger files and provide access to merged swagger UI
# Usage: ./merge_and_serve.sh

set -e

echo "🚀 Starting Swagger merge and serve process..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Change to the protobuf directory
cd "$(dirname "$0")/.."

# Check if docs directory exists
if [ ! -d "docs" ]; then
    echo "❌ docs directory not found. Please make sure you're in the protobuf directory."
    exit 1
fi

# Run the merge script
echo "📁 Processing swagger files..."
go run scripts/merge_swagger.go

# Check if the output file was created
if [ -f "merged_swagger.json" ]; then
    echo "✅ Successfully created merged_swagger.json"
    echo "📊 File size: $(du -h merged_swagger.json | cut -f1)"
else
    echo "❌ Failed to create merged_swagger.json"
    exit 1
fi

# Check if gateway docs directory was created
if [ -f "../gateway-service/docs/merged_swagger.json" ]; then
    echo "✅ Successfully copied to gateway-service/docs/merged_swagger.json"
else
    echo "⚠️  Warning: Could not copy to gateway service docs"
fi

echo ""
echo "🎉 Merge process completed successfully!"
echo ""
echo "📋 Next steps:"
echo "1. Start the gateway service:"
echo "   cd ../gateway-service"
echo "   go run cmd/http/main.go"
echo ""
echo "2. Access the merged Swagger UI:"
echo "   🌐 http://localhost:8080/swagger/merged/index.html"
echo ""
echo "3. Access individual service Swagger UI:"
echo "   🌐 http://localhost:8080/swagger/organization/index.html"
echo "   🌐 http://localhost:8080/swagger/user/index.html"
echo "   🌐 http://localhost:8080/swagger/bdspro/index.html"
echo "   🌐 ... (other services)"
echo ""
echo "🔗 Direct access to merged swagger JSON:"
echo "   📄 http://localhost:8080/swagger-doc/merged_swagger.json"
echo ""
echo "💡 The merged swagger includes:"
echo "   ✅ All API endpoints from all microservices"
echo "   ✅ Bearer authentication for all APIs"
echo "   ✅ All schemas/definitions"
echo "   ✅ All tags (no duplicates)" 