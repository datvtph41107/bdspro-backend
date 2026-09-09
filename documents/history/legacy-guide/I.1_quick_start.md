# I.1 Quick Start Guide - BDSPro Microservices

## 📁 Shared/Code Directory Overview

Thư mục `shared/code` chứa các script và công cụ quản lý cho toàn bộ hệ thống microservices.

---

## 🛠️ Available Commands & Scripts

### 1. **Makefile Commands** (`shared/code/Makefile`)

#### **Service Management**
```bash
# Start all services
make start-all

# Stop all services  
make stop-all

# Build all services
make build-all
```

#### **Individual Service Commands**
```bash
# Gateway Service
make gateway-http

# User Service
make user-http
make user-grpc

# Organization Service
make organization-grpc

# BDSPro Service
make bdspro-grpc

# CRM Service
make crm-grpc

# Chat Service
make chat-http
make chat-grpc

# File Service
make file-http

# Other Services
make notification-grpc
make social-grpc
make appointment-grpc
make payment-grpc
make marketing-grpc
make transaction-grpc
make auth-grpc
make membership-grpc
make relay-grpc
```

#### **Code Generation Commands**
```bash
# Generate protobuf files
make buf-all                    # Generate all services
make buf-user                   # Generate user service
make buf-organization           # Generate organization service
make buf-bdspro                 # Generate bdspro service
make buf-crm                    # Generate crm service
make buf-social                 # Generate social service
make buf-notification           # Generate notification service
make buf-appointment            # Generate appointment service
make buf-payment                # Generate payment service
make buf-chat                   # Generate chat service
make buf-task                   # Generate task service
make buf-marketing              # Generate marketing service
make buf-transaction            # Generate transaction service
make buf-auth                   # Generate auth service
make buf-shared                 # Generate shared protobuf

# Generate wire dependencies
make wire <service-name>
# Example: make wire user

# Generate Swagger documentation
make swag <service-name>
# Example: make swag user

# Merge all Swagger docs
make merge-swagger
```

#### **Development Commands**
```bash
# Install dependencies for all services
make init

# Run specific service with Air (hot reload)
make run-dev <service-name>
# Example: make run-dev user

# Run service directly (no hot reload)
make run <service-name> <args>
# Example: make run user grpc

# Debug service information
make debug-info <service-name>
# Example: make debug-info user

# Generate protobuf for specific service
make proto <service-name>
# Example: make proto user
```

#### **Utility Commands**
```bash
# Development monitor
make dev

# Generate Swagger for file service
make swag-file
```

---

### 2. **Script.sh Commands** (`shared/code/script.sh`)

#### **Swagger Generation**
```bash
# Generate Swagger docs for specific service
./script.sh swag <service-name>
# Example: ./script.sh swag user
```

#### **Wire Generation**
```bash
# Generate wire dependencies for specific service
./script.sh wire <service-name>
# Example: ./script.sh wire user
```

#### **Social Service Special Command**
```bash
# Special command for social service (generates both swagger and wire)
./script.sh social
```

---

### 3. **Deploy.sh Commands** (`shared/code/deploy.sh`)

#### **Deployment Commands**
```bash
# Deploy specific service to production server
./deploy.sh <service-name>

# Available services for deployment:
./deploy.sh gateway
./deploy.sh user
./deploy.sh auth
./deploy.sh user-grpc
./deploy.sh file
./deploy.sh map
./deploy.sh chat
./deploy.sh payment
./deploy.sh notification
./deploy.sh crm
./deploy.sh membership
./deploy.sh bdspro
./deploy.sh bdspro-grpc
./deploy.sh chathttp
./deploy.sh relay
./deploy.sh swagger-chat
./deploy.sh social
./deploy.sh appointment
./deploy.sh org
./deploy.sh transaction

# Deploy docker-compose.yml
./deploy.sh compose
```

#### **Deployment Process**
1. Builds Go binary for Linux (GOOS=linux GOARCH=amd64)
2. Uploads binary and Dockerfile to remote server
3. Uploads service-specific config files
4. Builds and starts Docker container on server

---

### 4. **Docker Compose Services** (`shared/code/docker-compose.yml`)

#### **Available Services**
```yaml
# Core Services
gateway          # API Gateway (port 8000)
user             # User HTTP service
user-grpc        # User gRPC service
auth             # Auth gRPC service

# Business Services
bdspro           # BDSPro HTTP service
bdspro-grpc      # BDSPro gRPC service
organization     # Organization gRPC service
crm              # CRM gRPC service
payment          # Payment gRPC service
notification     # Notification gRPC service
social           # Social gRPC service
appointment      # Appointment gRPC service
transaction      # Transaction gRPC service
membership       # Membership service

# Communication Services
chat             # Chat gRPC service
http-server      # Chat HTTP server (port 8106)
swagger-chat     # Chat Swagger server (port 8107)
relay            # Relay service (port 8308)

# Utility Services
file             # File service (port 8002)
map              # Map service
```

#### **Docker Commands**
```bash
# Start all services
docker-compose up -d

# Start specific service
docker-compose up -d <service-name>

# Build and start
docker-compose up --build -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f <service-name>

# Execute command in container
docker-compose exec <service-name> bash
```

---

### 5. **Monitor Tool** (`shared/code/monitor/run-dev.go`)

#### **Service Monitoring**
```bash
# Run development monitor
make dev
# or
go run monitor/run-dev.go
```

**Features:**
- Real-time service status monitoring
- Updates every 2 seconds
- Shows running/stopped status for services
- Interactive terminal interface

**Monitored Services:**
- user-service
- order-service  
- payment-service

---

### 6. **Keep.sh Commands** (`shared/code/keep.sh`)

#### **Quick Development Commands**
```bash
# Run keep.sh for quick development setup
./keep.sh
```

**Actions:**
1. Generates protobuf files for chat service
2. Generates Swagger documentation
3. Kills process on port 8204

---

## 🚀 Quick Start Workflow

### **Option 1: Development Mode (Recommended)**
```bash
# 1. Navigate to shared/code
cd shared/code

# 2. Install dependencies
make init

# 3. Generate protobuf files
make buf-all

# 4. Start all services
make start-all

# 5. Monitor services
make dev
```

### **Option 2: Docker Mode**
```bash
# 1. Navigate to shared/code
cd shared/code

# 2. Start with Docker Compose
docker-compose up -d

# 3. Check status
docker-compose ps
```

### **Option 3: Individual Service Development**
```bash
# 1. Navigate to shared/code
cd shared/code

# 2. Start specific service with hot reload
make run-dev user

# 3. In another terminal, start another service
make run-dev organization
```

---

## 📋 Service Ports & Endpoints

| Service | HTTP Port | gRPC Port | Health Check |
|---------|-----------|-----------|--------------|
| **gateway** | 8000 | - | http://localhost:8000/health |
| **user** | 8081 | 50051 | http://localhost:8081/health |
| **organization** | - | 50052 | gRPC health check |
| **bdspro** | - | 50053 | gRPC health check |
| **crm** | - | 50054 | gRPC health check |
| **payment** | - | 50055 | gRPC health check |
| **notification** | - | 50056 | gRPC health check |
| **chat** | 8082 | 50057 | http://localhost:8082/health |
| **social** | - | 50058 | gRPC health check |
| **file** | 8002 | - | http://localhost:8002/health |
| **appointment** | - | 50059 | gRPC health check |
| **marketing** | - | 50060 | gRPC health check |
| **transaction** | - | 50061 | gRPC health check |
| **auth** | - | 50062 | gRPC health check |
| **membership** | - | 50063 | gRPC health check |
| **relay** | 8308 | - | http://localhost:8308/health |

---

## 🔧 Development Commands Reference

### **Most Used Commands**
```bash
# Service management
make start-all          # Start all services
make stop-all           # Stop all services
make run-dev <service>  # Start service with hot reload

# Code generation
make buf-all            # Generate all protobuf
make wire <service>     # Generate wire dependencies
make swag <service>     # Generate Swagger docs

# Development
make init               # Install dependencies
make dev                # Run service monitor
make build-all          # Build all services
```

### **Troubleshooting Commands**
```bash
# Check running processes
ps aux | grep air

# Kill processes on specific port
lsof -ti tcp:8080 | xargs kill -9

# Check service logs
tail -f */tmp/*.log

# Debug service info
make debug-info <service>
```

---

## ⚡ Quick Test Commands

### **Health Checks**
```bash
# Test API Gateway
curl http://localhost:8000/health

# Test User Service
curl http://localhost:8081/health

# Test Chat Service
curl http://localhost:8082/health

# Test File Service
curl http://localhost:8002/health

# Test Relay Service
curl http://localhost:8308/health
```

### **API Testing**
```bash
# Test with authentication
curl -H "Authorization: Bearer <token>" \
     http://localhost:8000/v2/org/organizations

# Test gRPC services (need grpcurl)
grpcurl -plaintext localhost:50051 list
```

---

## 🆘 Troubleshooting

### **Common Issues**

#### **Port Already in Use**
```bash
# Kill process using port
lsof -ti tcp:8080 | xargs kill -9

# Or stop all and restart
make stop-all
make start-all
```

#### **Service Not Starting**
```bash
# Check if service is running
ps aux | grep air

# Check logs
tail -f */tmp/*.log

# Debug specific service
make debug-info <service>
```

#### **Protobuf Generation Errors**
```bash
# Clean and regenerate
make buf-all

# Or generate specific service
make buf-<service>
```

#### **Dependencies Issues**
```bash
# Reinstall dependencies
make init

# Check Go modules
cd <service>-service && go mod tidy
```

---

## 📚 Next Steps

1. **Explore APIs**: Visit Swagger UI at http://localhost:8000/swagger/index.html
2. **Test Endpoints**: Use the health check commands above
3. **Start Developing**: Use `make run-dev <service>` for hot reload development
4. **Monitor Services**: Use `make dev` to monitor service status
5. **Deploy**: Use `./deploy.sh <service>` to deploy to production

---

**🎉 Congratulations!** Bạn đã setup thành công hệ thống BDSPro Microservices với shared/code tools. Hệ thống đã sẵn sàng cho development và testing.