#!/bin/bash

# Script để import dữ liệu location v2 vào database

echo "========================================="
echo "Import Location V2 Data Script"
echo "========================================="

# Kiểm tra DATABASE_URL environment variable
if [ -z "$DATABASE_URL" ]; then
    echo "⚠️  DATABASE_URL not set. Using default connection."
    echo "To set custom connection, export DATABASE_URL environment variable:"
    echo "export DATABASE_URL='host=localhost user=postgres password=postgres dbname=bdspro_service port=5432 sslmode=disable'"
else
    echo "✅ Using DATABASE_URL: $DATABASE_URL"
fi

# Di chuyển tới thư mục scripts
cd "$(dirname "$0")"

# Chạy script import
echo ""
echo "Starting import..."
go run import_location_v2_data.go

# Kiểm tra kết quả
if [ $? -eq 0 ]; then
    echo ""
    echo "========================================="
    echo "✅ Import completed successfully!"
    echo "========================================="
else
    echo ""
    echo "========================================="
    echo "❌ Import failed!"
    echo "========================================="
    exit 1
fi

