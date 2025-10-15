#!/bin/bash

# Script để export permissions và role_groups ra file JSON

echo "🚀 Starting export process..."
echo ""

# Đọc config database từ file config hoặc environment
export DATABASE_URL="${DATABASE_URL:-host=localhost user=postgres password=postgres dbname=organization_db port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh}"

# Chuyển vào thư mục scripts
cd "$(dirname "$0")"

# Chạy script Go
echo "📦 Running Go script..."
go run export_permissions_roles.go

# Kiểm tra kết quả
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Export completed successfully!"
    echo ""
    
    # Hiển thị thông tin file
    if [ -f "permissions_roles_export.json" ]; then
        FILE_SIZE=$(du -h permissions_roles_export.json | cut -f1)
        echo "📁 File created: permissions_roles_export.json"
        echo "📊 File size: $FILE_SIZE"
        echo ""
        echo "You can view the file with:"
        echo "  cat permissions_roles_export.json | jq"
        echo "  or"
        echo "  open permissions_roles_export.json"
    fi
else
    echo ""
    echo "❌ Export failed! Please check the error messages above."
    exit 1
fi

