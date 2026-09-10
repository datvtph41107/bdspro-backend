# IV.2 Triển Khai Nhanh - Deploy Script

## 🎯 Tổng Quan

Script `deploy.sh` là công cụ tự động hóa quá trình triển khai các microservices lên production server. Script này thực hiện build Go binary, upload files, và deploy Docker containers một cách tự động.

---

## 🚀 Deploy Script Overview

### **Script Location:**
```bash
shared/code/deploy.sh
```

### **Server Configuration:**
```bash
REMOTE_USER="root"
REMOTE_HOST="14.225.205.222"
WORKSPACE_NAME=backend-microservice
REMOTE_PATH=/root/$WORKSPACE_NAME
```

### **Usage:**
```bash
./deploy.sh <service-name>
./deploy.sh compose
```

---

## 📋 Supported Services

### **Available Services:**
| Service | Service Path | Binary Name | Payload Build |
|---------|-------------|-------------|---------------|
| **gateway** | gateway-service | main | gateway |
| **user** | user-service | main | user |
| **auth** | auth-service | main | auth |
| **user-grpc** | user-service | main | user-grpc |
| **file** | file-service | main | file |
| **map** | map-service | main | map |
| **chat** | chat-service | main | chat |
| **payment** | payment-service | main | payment |
| **notification** | notification-service | main | notification |
| **crm** | crm-service | main | crm |
| **membership** | membership-service | main | membership |
| **bdspro** | bdspro-service | main | bdspro |
| **bdspro-grpc** | bdspro-service | main | bdspro-grpc |
| **chathttp** | chat-service | main | http-server |
| **relay** | relay-service | main | relay |
| **swagger-chat** | chat-service | main | swagger-chat |
| **social** | social-service | main | social |
| **appointment** | appointment-service | main | appointment |
| **org** | organization-service | main | organization |
| **transaction** | transaction-service | main | transaction |

---

## 🔧 Deploy Process

### **1. Build Process**
```bash
# Navigate to service directory
cd $service_path

# Build Go binary for Linux
GOOS=linux GOARCH=amd64 go build -o main .
```

### **2. File Upload Process**
```bash
# Create remote directory
ssh ${REMOTE_USER}@${REMOTE_HOST} "mkdir -p ${REMOTE_DIR}"

# Upload binary
scp ${path_jar} ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}

# Upload Dockerfile
scp Dockerfile ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/Dockerfile
```

### **3. Special Config Uploads**

#### **Notification Service:**
```bash
if [ "$1" = "notification" ]; then
  scp config/bdspro-5e599-firebase-adminsdk-fbsvc-5aff064187.json \
      ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/bdspro-5e599-firebase-adminsdk-fbsvc-5aff064187.json
fi
```

#### **Chat Service:**
```bash
if [ "$1" = "chat" ]; then
  scp -r configs ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
```

#### **Relay Service:**
```bash
if [ "$1" = "relay" ]; then
  scp -r config ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
```

#### **Organization Service:**
```bash
if [ "$1" = "org" ]; then
  scp -r config ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
```

#### **Payment Service:**
```bash
if [ "$1" = "payment" ]; then
  scp -r config ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
```

#### **Appointment Service:**
```bash
if [ "$1" = "appointment" ]; then
  scp -r config ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
```

#### **Swagger Chat:**
```bash
if [ "$1" = "swagger-chat" ]; then
  scp -r proto ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
```

#### **Gateway Service:**
```bash
if [ "$1" = "gateway" ]; then
  scp -r docs ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
```

#### **File Service:**
```bash
if [ "$1" = "file" ]; then
  scp -r docs ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
```

### **4. Docker Deployment**
```bash
# Build Docker image (no cache)
docker-compose build --no-cache ${payload_build}

# Start Docker container
docker-compose up -d ${payload_build}
```

---

## 🚀 Quick Deploy Commands

### **1. Deploy Individual Services**
```bash
# Deploy gateway service
./deploy.sh gateway

# Deploy user service
./deploy.sh user

# Deploy auth service
./deploy.sh auth

# Deploy organization service
./deploy.sh org

# Deploy bdspro service
./deploy.sh bdspro

# Deploy crm service
./deploy.sh crm

# Deploy payment service
./deploy.sh payment

# Deploy notification service
./deploy.sh notification

# Deploy chat service
./deploy.sh chat

# Deploy social service
./deploy.sh social

# Deploy appointment service
./deploy.sh appointment

# Deploy transaction service
./deploy.sh transaction
```

### **2. Deploy Special Services**
```bash
# Deploy user gRPC service
./deploy.sh user-grpc

# Deploy bdspro gRPC service
./deploy.sh bdspro-grpc

# Deploy chat HTTP server
./deploy.sh chathttp

# Deploy relay service
./deploy.sh relay

# Deploy swagger chat
./deploy.sh swagger-chat

# Deploy map service
./deploy.sh map

# Deploy membership service
./deploy.sh membership

# Deploy file service
./deploy.sh file
```

### **3. Deploy Docker Compose**
```bash
# Deploy docker-compose.yml to server
./deploy.sh compose
```

---

## ⏱️ Deployment Timing

### **Script Performance Tracking:**
```bash
start=$(date +%s)
# ... deployment process ...
end=$(date +%s)
elapsed=$((end - start))
echo "Thời gian thực thi: $elapsed giây"
```

### **Typical Deployment Times:**
| Service Type | Build Time | Upload Time | Deploy Time | Total Time |
|-------------|------------|-------------|-------------|------------|
| **Simple Service** | 10-15s | 5-10s | 15-20s | 30-45s |
| **Service with Config** | 10-15s | 10-20s | 15-20s | 35-55s |
| **Service with Docs** | 10-15s | 15-25s | 15-20s | 40-60s |
| **Docker Compose** | - | 5-10s | - | 5-10s |

---

## 🔧 Prerequisites

### **1. Local Requirements**
```bash
# Go installed
go version

# Docker installed
docker --version

# SSH access to server
ssh root@14.225.205.222
```

### **2. Server Requirements**
```bash
# Docker installed on server
docker --version

# Docker Compose installed
docker-compose --version

# SSH access configured
# SSH key authentication setup
```

### **3. Network Requirements**
```bash
# Internet connection for file upload
# SSH access to production server
# Port 22 (SSH) accessible
```

---

## 🛠️ Troubleshooting

### **1. Common Issues**

#### **SSH Connection Failed:**
```bash
# Check SSH connection
ssh root@14.225.205.222

# Check SSH key
ssh-add -l

# Test SSH with verbose
ssh -v root@14.225.205.222
```

#### **Build Failed:**
```bash
# Check Go installation
go version

# Check Go modules
go mod tidy

# Check build command
GOOS=linux GOARCH=amd64 go build -o main .
```

#### **Upload Failed:**
```bash
# Check file permissions
ls -la main
ls -la Dockerfile

# Check disk space
df -h

# Check network connectivity
ping 14.225.205.222
```

#### **Docker Deploy Failed:**
```bash
# Check Docker on server
ssh root@14.225.205.222 "docker --version"

# Check Docker Compose on server
ssh root@14.225.205.222 "docker-compose --version"

# Check docker-compose.yml exists
ssh root@14.225.205.222 "ls -la /root/backend-microservice/docker-compose.yml"
```

### **2. Debug Commands**

#### **Check Service Status:**
```bash
# Check if service is running
ssh root@14.225.205.222 "docker ps | grep <service-name>"

# Check service logs
ssh root@14.225.205.222 "docker logs <container-name>"

# Check service health
ssh root@14.225.205.222 "docker inspect <container-name> | grep Health"
```

#### **Manual Deployment:**
```bash
# Manual build
cd <service-directory>
GOOS=linux GOARCH=amd64 go build -o main .

# Manual upload
scp main root@14.225.205.222:/root/backend-microservice/<service-directory>/
scp Dockerfile root@14.225.205.222:/root/backend-microservice/<service-directory>/

# Manual deploy
ssh root@14.225.205.222 "cd /root/backend-microservice && docker-compose build --no-cache <service> && docker-compose up -d <service>"
```

---

## 📊 Monitoring Deployment

### **1. Check Deployment Status**
```bash
# Check all running containers
ssh root@14.225.205.222 "docker ps"

# Check specific service
ssh root@14.225.205.222 "docker ps | grep <service-name>"

# Check service logs
ssh root@14.225.205.222 "docker logs <container-name> -f"
```

### **2. Health Checks**
```bash
# Gateway health
curl http://14.225.205.222:8000/health

# File service health
curl http://14.225.205.222:8002/health

# Chat HTTP server
curl http://14.225.205.222:8106/health

# Relay service
curl http://14.225.205.222:8308/health
```

### **3. Resource Monitoring**
```bash
# Check container resources
ssh root@14.225.205.222 "docker stats"

# Check disk usage
ssh root@14.225.205.222 "df -h"

# Check memory usage
ssh root@14.225.205.222 "free -h"
```

---

## 🎯 Best Practices

### **1. Before Deployment**
- **Test locally** với Docker Compose
- **Check code quality** với linters
- **Run tests** trước khi deploy
- **Backup current version** nếu cần

### **2. During Deployment**
- **Monitor logs** trong quá trình deploy
- **Check service health** sau khi deploy
- **Verify functionality** với API calls
- **Monitor resources** trên server

### **3. After Deployment**
- **Check service status** và logs
- **Monitor performance** metrics
- **Verify all endpoints** hoạt động
- **Update documentation** nếu cần

### **4. Rollback Strategy**
```bash
# Stop current service
ssh root@14.225.205.222 "docker-compose stop <service>"

# Deploy previous version
./deploy.sh <service>

# Or restore from backup
ssh root@14.225.205.222 "docker-compose up -d <service>"
```

---

## 🎯 Tóm Tắt

| Command | Mục đích | Khi nào sử dụng |
|---------|----------|-----------------|
| `./deploy.sh <service>` | Deploy specific service | Production deployment |
| `./deploy.sh compose` | Deploy docker-compose | Update orchestration |
| `./deploy.sh gateway` | Deploy gateway service | API gateway updates |
| `./deploy.sh user` | Deploy user service | User management updates |
| `./deploy.sh auth` | Deploy auth service | Authentication updates |
| `./deploy.sh org` | Deploy organization service | Organization updates |

### **Deploy Process Flow:**
1. **Build** Go binary for Linux
2. **Upload** binary và Dockerfile
3. **Upload** config files (nếu cần)
4. **Build** Docker image
5. **Start** Docker container
6. **Monitor** deployment status

**Quy tắc vàng**: **Luôn test locally trước khi deploy và monitor logs sau khi deploy!**
