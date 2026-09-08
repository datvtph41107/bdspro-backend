# 🚀 Quick Start Guide - BDSPro Microservices

## ⚡ Start Project trong 5 phút

### Prerequisites
- Go 1.21+
- Air (for hot reloading)
- Make
- Buf CLI
- yq (for config parsing)

---

## 🏃‍♂️ Quick Start với Makefile

### 1. Setup Environment
```bash
# Clone repository
git clone <repository-url>
cd bdspro-golang-microservices

# Install dependencies cho tất cả services
cd shared/code
make init
```

### 2. Generate Protobuf Files
```bash
# Generate tất cả protobuf files
make buf-all

# Hoặc generate từng service riêng lẻ
make buf-user
make buf-organization
make buf-bdspro
make buf-crm
# ... các service khác
```

### 3. Start All Services
```bash
# Start tất cả services cùng lúc
make start-all

# Kiểm tra services đang chạy
ps aux | grep air
```

### 4. Verify Installation
```bash
# Test API Gateway (port 8080)
curl http://localhost:8080/health

# Test File Service (port 8081)
curl http://localhost:8081/health
```

**✅ Done!** Tất cả services đã chạy và sẵn sàng sử dụng.

---

## 🔧 Start Individual Services

### Gateway Service (API Gateway)
```bash
# Start Gateway HTTP server
make gateway-http
```

### User Service
```bash
# Start User HTTP server
make user-http

# Start User gRPC server
make user-grpc
```

### Organization Service
```bash
# Start Organization gRPC server
make organization-grpc
```

### BDSPro Service
```bash
# Start BDSPro gRPC server
make bdspro-grpc
```

### CRM Service
```bash
# Start CRM gRPC server
make crm-grpc
```

### Chat Service
```bash
# Start Chat HTTP server
make chat-http

# Start Chat gRPC server
make chat-grpc
```

### File Service
```bash
# Start File HTTP server
make file-http
```

### Other Services
```bash
# Notification Service
make notification-grpc

# Social Service
make social-grpc

# Appointment Service
make appointment-grpc

# Payment Service
make payment-grpc

# Marketing Service
make marketing-grpc

# Transaction Service
make transaction-grpc

# Auth Service
make auth-grpc

# Membership Service
make membership-grpc

# Relay Service
make relay-grpc
```

---

## 📋 Available Services & Ports

| Service | HTTP Port | gRPC Port | Make Command |
|---------|-----------|-----------|--------------|
| **gateway-service** | 8080 | - | `make gateway-http` |
| user-service | 8081 | 50051 | `make user-http` / `make user-grpc` |
| organization-service | - | 50052 | `make organization-grpc` |
| bdspro-service | - | 50053 | `make bdspro-grpc` |
| crm-service | - | 50054 | `make crm-grpc` |
| payment-service | - | 50055 | `make payment-grpc` |
| notification-service | - | 50056 | `make notification-grpc` |
| chat-service | 8082 | 50057 | `make chat-http` / `make chat-grpc` |
| social-service | - | 50058 | `make social-grpc` |
| file-service | 8083 | - | `make file-http` |
| appointment-service | - | 50059 | `make appointment-grpc` |
| marketing-service | - | 50060 | `make marketing-grpc` |
| transaction-service | - | 50061 | `make transaction-grpc` |
| auth-service | - | 50062 | `make auth-grpc` |
| membership-service | - | 50063 | `make membership-grpc` |
| relay-service | - | 50064 | `make relay-grpc` |

---

## 🛠️ Development Commands

### Service Management
```bash
# Start all services
make start-all

# Stop all services
make stop-all

# Build all services
make build-all

# Start specific service with Air (hot reload)
make run-dev <service-name>
# Ví dụ: make run-dev user
```

### Code Generation
```bash
# Generate protobuf for specific service
make proto <service-name>
# Ví dụ: make proto user

# Generate wire dependencies
make wire <service-name>
# Ví dụ: make wire user

# Generate Swagger docs
make swag <service-name>
# Ví dụ: make swag user

# Merge all Swagger docs
make merge-swagger
```

### Individual Service Commands
```bash
# Run service directly (without Air)
make run <service-name> <args>
# Ví dụ: make run user grpc

# Debug service info
make debug-info <service-name>
# Ví dụ: make debug-info user
```

---

## 🌐 Access Points

### Web Interfaces
- **API Gateway**: http://localhost:8080
- **User Service**: http://localhost:8081
- **Chat Service**: http://localhost:8082
- **File Service**: http://localhost:8083

### API Endpoints
- **Health Check**: `GET http://localhost:8080/health`
- **API Base**: `http://localhost:8080/v2/`
- **Authentication**: JWT Bearer Token required

---

## ⚡ Quick Test Commands

### 1. Test API Gateway
```bash
# Health check
curl http://localhost:8080/health

# Test với authentication (cần token)
curl -H "Authorization: Bearer <your-token>" \
     http://localhost:8080/v2/org/organizations
```

### 2. Test gRPC Services
```bash
# Cần cài grpcurl: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# List available services
grpcurl -plaintext localhost:50051 list

# Test user service
grpcurl -plaintext localhost:50051 user.UserService/GetProfile
```

### 3. Test Individual Services
```bash
# Test User Service
curl http://localhost:8081/health

# Test Chat Service
curl http://localhost:8082/health

# Test File Service
curl http://localhost:8083/health
```

---

## 🔍 Troubleshooting

### Common Issues

#### 1. Port Already in Use
```bash
# Kill process using port
lsof -ti tcp:8080 | xargs kill -9

# Or stop all services and restart
make stop-all
make start-all
```

#### 2. Service Not Starting
```bash
# Check if service is running
ps aux | grep air

# Check logs
tail -f */tmp/*.log

# Debug specific service
make debug-info <service-name>
```

#### 3. Protobuf Generation Errors
```bash
# Clean and regenerate
make buf-all

# Or generate specific service
make buf-<service-name>
```

#### 4. Dependencies Issues
```bash
# Reinstall dependencies
make init

# Check Go modules
cd <service-name>-service && go mod tidy
```

### Health Checks
```bash
# Check all services health
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health

# Check running processes
ps aux | grep air
```

---

## 📚 Development Workflow

### 1. Start Development
```bash
# Start all services
make start-all

# Or start specific service for development
make run-dev user
```

### 2. Make Changes
- Edit code in your IDE
- Air will automatically reload the service
- Test your changes via API calls

### 3. Generate Code
```bash
# After changing .proto files
make buf-<service-name>

# After changing dependencies
make wire <service-name>

# After changing API docs
make swag <service-name>
```

### 4. Build for Production
```bash
# Build all services
make build-all

# Or build specific service
cd <service-name>-service && go build -o <service-name> main.go
```

---

## 🆘 Need Help?

### Quick Commands Reference
```bash
# Most used commands
make start-all          # Start all services
make stop-all           # Stop all services
make init               # Install dependencies
make buf-all            # Generate protobuf files
make build-all          # Build all services
make run-dev <service>  # Start specific service with hot reload
make debug-info <service> # Debug service info
```

### Check Service Status
```bash
# See running processes
ps aux | grep air

# See logs
tail -f */tmp/*.log

# Health checks
curl http://localhost:8080/health
```

### Documentation
- **System Overview**: `documents/overview.md`
- **API Documentation**: Swagger UI at http://localhost:8080/swagger/index.html
- **Architecture**: `README.md`

---

**🎉 Congratulations!** Bạn đã setup thành công hệ thống BDSPro Microservices với Makefile. Hệ thống đã sẵn sàng cho development và testing.

### Next Steps:
1. Explore APIs tại Swagger UI
2. Test endpoints với sample data
3. Start developing features
4. Check logs khi có issues
