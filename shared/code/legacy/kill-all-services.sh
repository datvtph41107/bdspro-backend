#!/bin/bash

# Script để kill tất cả các service của BDS Pro
# Dựa trên cấu hình port từ gateway-service/config/config.yml

echo "🔄 Đang dừng tất cả các service BDS Pro..."

# Danh sách các port từ config.yml
PORTS=(
    8000  # gateway
    8800  # gateway tcp
    8201  # user
    8202  # bdspro
    8203  # file
    8204  # notification
    8205  # payment
    8206  # crm
    8207  # organization
    8208  # chat
    8209  # social
    8210  # map
    8211  # search
    8212  # membership
    8213  # appointment
    8214  # marketing
    8215  # transaction
    8216  # auth
)

# Tên các service tương ứng
SERVICES=(
    "gateway"
    "gateway-tcp"
    "user"
    "bdspro"
    "file"
    "notification"
    "payment"
    "crm"
    "organization"
    "chat"
    "social"
    "map"
    "search"
    "membership"
    "appointment"
    "marketing"
    "transaction"
    "auth"
)

# Hàm kill process theo port
kill_by_port() {
    local port=$1
    local service_name=$2
    
    # Tìm process đang chạy trên port
    local pid=$(lsof -ti:$port 2>/dev/null)
    
    if [ ! -z "$pid" ]; then
        echo "🔴 Dừng service $service_name (port $port, PID: $pid)"
        kill -9 kill -9 $pid
        sleep 1
        
        # Kiểm tra lại xem process đã dừng chưa
        local check_pid=$(lsof -ti:$port 2>/dev/null)
        if [ ! -z "$check_pid" ]; then
            echo "⚠️  Service $service_name vẫn đang chạy, thử kill lần nữa..."
            kill -9 $check_pid
        else
            echo "✅ Service $service_name đã dừng thành công"
        fi
    else
        echo "ℹ️  Service $service_name (port $port) không đang chạy"
    fi
}

# Kill tất cả các service
for i in "${!PORTS[@]}"; do
    kill_by_port "${PORTS[$i]}" "${SERVICES[$i]}"
done

echo ""
echo "🧹 Dọn dẹp các process còn sót lại..."

# Kill các process có tên chứa tên service
SERVICE_NAMES=(
    "gateway-service"
    "user-service"
    "bdspro-service"
    "file-service"
    "notification-service"
    "payment-service"
    "crm-service"
    "organization-service"
    "chat-service"
    "social-service"
    "map-service"
    "search-service"
    "membership-service"
    "appointment-service"
    "marketing-service"
    "transaction-service"
    "auth-service"
)

for service in "${SERVICE_NAMES[@]}"; do
    local pids=$(pgrep -f "$service" 2>/dev/null)
    if [ ! -z "$pids" ]; then
        echo "🔴 Dừng process $service (PIDs: $pids)"
        echo $pids | xargs kill -9 2>/dev/null
    fi
done

echo ""
echo "📊 Kiểm tra các port còn đang được sử dụng:"
echo "----------------------------------------"

# Kiểm tra lại các port
for i in "${!PORTS[@]}"; do
    local port="${PORTS[$i]}"
    local service="${SERVICES[$i]}"
    local pid=$(lsof -ti:$port 2>/dev/null)
    
    if [ ! -z "$pid" ]; then
        echo "❌ Port $port ($service) vẫn đang được sử dụng bởi PID: $pid"
    else
        echo "✅ Port $port ($service) đã được giải phóng"
    fi
done

echo ""
echo "🎉 Hoàn thành! Tất cả các service BDS Pro đã được dừng."
echo ""
echo "💡 Để khởi động lại các service, sử dụng:"
echo "   - make run-all (nếu có Makefile)"
echo "   - Hoặc chạy từng service riêng lẻ"
