#!/bin/bash

# Script để export dữ liệu từ user-service database ra file JSON

echo "🚀 Starting user-service data export..."
echo ""

# Đọc config database từ file config hoặc environment
# Default từ config/app.yml
export DATABASE_URL="${DATABASE_URL:-host=14.225.210.29 user=hieubap password=Matkhau.1 dbname=db_user port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh}"

# Chuyển vào thư mục scripts
cd "$(dirname "$0")"

# Chạy script Go
echo "📦 Running Go script..."
go run export_data.go

# Kiểm tra kết quả
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Export completed successfully!"
    echo ""
    
    # Hiển thị thông tin file
    if [ -f "user_service_export.json" ]; then
        FILE_SIZE=$(du -h user_service_export.json | cut -f1)
        echo "📁 File created: user_service_export.json"
        echo "📊 File size: $FILE_SIZE"
        echo ""
        echo "You can view the file with:"
        echo "  cat user_service_export.json | jq"
        echo "  or"
        echo "  open user_service_export.json"
    fi
else
    echo ""
    echo "❌ Export failed! Please check the error messages above."
    exit 1
fi

