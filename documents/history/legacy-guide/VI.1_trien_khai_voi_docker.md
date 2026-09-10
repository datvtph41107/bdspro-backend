# IV.1 Triển Khai với Docker

## 🎯 Tổng Quan

Dự án BDSPro Microservices sử dụng Docker để containerization và Docker Compose để orchestration. Mỗi service được đóng gói thành Docker image riêng biệt và có thể deploy độc lập hoặc cùng nhau thông qua docker-compose.

---

## 🐳 Dockerfile Structure

### **Standard Dockerfile Pattern**
```dockerfile
FROM alpine:latest

WORKDIR /

COPY main .

ENV ENV_RUNTIME=develop

EXPOSE 8080

CMD ["./main", "grpc"]
```

### **Các Loại Dockerfile**

#### **1. Basic Service Dockerfile**
```dockerfile
# organization-service/Dockerfile
FROM alpine:latest

WORKDIR /

COPY main .

ENV ENV_RUNTIME=develop

CMD ["./main", "grpc"]
```

#### **2. Gateway Service Dockerfile**
```dockerfile
# gateway-service/Dockerfile
FROM alpine:latest

WORKDIR /

COPY main .
COPY docs docs

ENV ENV_RUNTIME=develop

EXPOSE 8080

CMD ["./main", "http"]
```

#### **3. File Service Dockerfile**
```dockerfile
# file-service/Dockerfile
FROM alpine:latest

WORKDIR /

COPY main .

EXPOSE 8080

CMD ["./main"]
```

---

## 🚀 Docker Compose Configuration

### **Main Docker Compose** (`shared/code/docker-compose.yml`)

#### **Service Definitions:**
```yaml
version: '3.3'

services:
  # Gateway Service
  gateway:
    build:
      context: ./gateway-service
      dockerfile: Dockerfile
    container_name: gateway
    networks:
      - internal
    ports:
      - "8000:8000"

  # User Service (HTTP)
  user:
    build:
      context: ./user-service
      dockerfile: Dockerfile
    container_name: user
    command: ["./main", "http"]
    networks:
      - internal

  # User Service (gRPC)
  user-grpc:
    build:
      context: ./user-service
      dockerfile: Dockerfile
    container_name: user
    command: ["./main", "grpc"]
    deploy:
      resources:
        limits:
          cpus: '0.50'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 128M
    networks:
      - internal

  # Organization Service
  organization:
    build:
      context: ./organization-service
      dockerfile: Dockerfile
    container_name: organization
    volumes:
      - ${PWD}/config:/app/config
    environment:
      - CONFIG_FILE=/app/config/docker_config.yaml 
    command: ["./main", "grpc-server"]
    networks:
      - internal

  # File Service
  file:
    build:
      context: ./file-service
      dockerfile: Dockerfile
    container_name: file
    command: ["./main", "http"]
    ports:
      - "8002:8002"
    networks:
      - internal
    volumes:
      - ${PWD}/files:/files

  # Chat Service
  chat:
    build:
      context: ./chat-service
      dockerfile: Dockerfile
    container_name: chat
    volumes:
      - ${PWD}/configs:/app/configs
    environment:
      - CONFIG_FILE=/app/configs/docker_config.yaml 
    command: ["/app/main", "grpc"]
    networks:
      - internal

networks:
  internal:
    driver: bridge
```

### **Individual Service Docker Compose** (`chat-service/docker-compose.yaml`)

#### **Service với Dependencies:**
```yaml
version: "3.8"

services:
  http-server:
    build: .
    container_name: http-server
    ports:
      - "8081:8081"
    command: ["sh", "-c", "sleep 5 && /app/binary http-serve --httpPort=8081 || sleep infinity"]
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    volumes:
      - ./configs:/app/configs
    environment:
      - CONFIG_FILE=/app/configs/docker_config.yaml 

  ws-server:
    build: .
    container_name: ws-server
    ports:
      - "3001:3001"
    command: ["sh", "-c", "sleep 5 && /app/binary ws-serve --wsPort=3001 || sleep infinity"]
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    volumes:
      - ./configs:/app/configs
    environment:
      - CONFIG_FILE=/app/configs/docker_config.yaml 

  redis:
    image: redis:alpine
    container_name: redis
    ports:
      - "6379:6379"
    command: ["redis-server", "--requirepass", "password"]
    environment:
      - REDIS_PASSWORD=password

  postgres:
    image: postgres:alpine
    container_name: postgres
    restart: always
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: admin
      POSTGRES_DB: mydatabase
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U admin -d mydatabase"]
      interval: 10s
      timeout: 5s
      retries: 5
```

---

## 🔧 Docker Commands

### **1. Build và Run Services**

#### **Build tất cả services:**
```bash
# Từ thư mục shared/code
docker-compose build

# Build specific service
docker-compose build gateway
docker-compose build user
docker-compose build organization
```

#### **Start tất cả services:**
```bash
# Start in background
docker-compose up -d

# Start specific service
docker-compose up -d gateway
docker-compose up -d user
```

#### **Start với rebuild:**
```bash
# Build và start
docker-compose up --build -d

# Build specific service và start
docker-compose up --build -d gateway
```

### **2. Service Management**

#### **Stop services:**
```bash
# Stop tất cả
docker-compose down

# Stop specific service
docker-compose stop gateway
```

#### **Restart services:**
```bash
# Restart tất cả
docker-compose restart

# Restart specific service
docker-compose restart gateway
```

#### **View logs:**
```bash
# Logs tất cả services
docker-compose logs -f

# Logs specific service
docker-compose logs -f gateway
docker-compose logs -f user
```

### **3. Container Management**

#### **List containers:**
```bash
# List running containers
docker-compose ps

# List all containers
docker ps -a
```

#### **Execute commands in container:**
```bash
# Execute bash in container
docker-compose exec gateway bash
docker-compose exec user bash

# Execute specific command
docker-compose exec gateway ls -la
```

---

## 🚀 Deployment Process

### **1. Local Development**

#### **Setup môi trường:**
```bash
# 1. Navigate to shared/code
cd shared/code

# 2. Build tất cả services
docker-compose build

# 3. Start services
docker-compose up -d

# 4. Check status
docker-compose ps
```

#### **Development workflow:**
```bash
# 1. Make code changes
# 2. Rebuild specific service
docker-compose build <service-name>

# 3. Restart service
docker-compose restart <service-name>

# 4. Check logs
docker-compose logs -f <service-name>
```

### **2. Production Deployment**

#### **Deploy script** (`shared/code/deploy.sh`):
```bash
#!/bin/bash

# Deploy specific service
./deploy.sh <service-name>

# Available services:
./deploy.sh gateway
./deploy.sh user
./deploy.sh auth
./deploy.sh organization
./deploy.sh bdspro
./deploy.sh crm
./deploy.sh payment
./deploy.sh notification
./deploy.sh chat
./deploy.sh social
./deploy.sh appointment
./deploy.sh transaction

# Deploy docker-compose
./deploy.sh compose
```

#### **Deploy process:**
```bash
# 1. Build Go binary for Linux
GOOS=linux GOARCH=amd64 go build -o main .

# 2. Upload binary và Dockerfile to server
scp main ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
scp Dockerfile ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}

# 3. Upload config files (nếu cần)
scp -r config ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}

# 4. Build và start Docker container
docker-compose build --no-cache ${service}
docker-compose up -d ${service}
```

---

## 📋 Service Ports & Access

### **Service Ports:**
| Service | Container Port | Host Port | Access |
|---------|---------------|-----------|---------|
| **gateway** | 8080 | 8000 | http://localhost:8000 |
| **user** | 8081 | - | Internal |
| **user-grpc** | 50051 | - | Internal |
| **organization** | 50052 | - | Internal |
| **bdspro** | 50053 | - | Internal |
| **crm** | 50054 | - | Internal |
| **payment** | 50055 | - | Internal |
| **notification** | 50056 | - | Internal |
| **chat** | 50057 | - | Internal |
| **social** | 50058 | - | Internal |
| **file** | 8081 | 8002 | http://localhost:8002 |
| **appointment** | 50059 | - | Internal |
| **transaction** | 50060 | - | Internal |
| **http-server** | 8081 | 8106 | http://localhost:8106 |
| **swagger-chat** | 3001 | 8107 | http://localhost:8107 |
| **relay** | 8308 | 8308 | http://localhost:8308 |

### **Health Checks:**
```bash
# Gateway health
curl http://localhost:8000/health

# File service health
curl http://localhost:8002/health

# Chat HTTP server
curl http://localhost:8106/health

# Relay service
curl http://localhost:8308/health
```

---

## 🔧 Configuration Management

### **1. Environment Variables**

#### **Runtime Environment:**
```dockerfile
ENV ENV_RUNTIME=develop
ENV ENV_RUNTIME=develop
```

#### **Config File Path:**
```yaml
environment:
  - CONFIG_FILE=/app/config/docker_config.yaml
  - CONFIG_FILE=/app/configs/docker_config.yaml
```

### **2. Volume Mounts**

#### **Config Files:**
```yaml
volumes:
  - ${PWD}/config:/app/config
  - ${PWD}/configs:/app/configs
  - ./configs:/app/configs
```

#### **File Storage:**
```yaml
volumes:
  - ${PWD}/files:/files
```

### **3. Resource Limits**

#### **CPU & Memory:**
```yaml
deploy:
  resources:
    limits:
      cpus: '0.50'
      memory: 512M
    reservations:
      cpus: '0.25'
      memory: 128M
```

---

## 🛠️ Troubleshooting

### **1. Common Issues**

#### **Port Already in Use:**
```bash
# Check port usage
lsof -i :8000
lsof -i :8002

# Kill process using port
sudo kill -9 $(lsof -ti:8000)
```

#### **Container Not Starting:**
```bash
# Check container logs
docker-compose logs <service-name>

# Check container status
docker-compose ps

# Restart container
docker-compose restart <service-name>
```

#### **Build Failures:**
```bash
# Clean build
docker-compose build --no-cache

# Check Dockerfile syntax
docker build -t test .

# Check Go binary
go build -o main .
```

### **2. Debug Commands**

#### **Container Debugging:**
```bash
# Execute shell in container
docker-compose exec <service> sh

# Check container processes
docker-compose exec <service> ps aux

# Check container filesystem
docker-compose exec <service> ls -la
```

#### **Network Debugging:**
```bash
# Check network connectivity
docker-compose exec <service> ping <other-service>

# Check DNS resolution
docker-compose exec <service> nslookup <service-name>
```

---

## 📊 Monitoring & Logging

### **1. Container Monitoring**

#### **Resource Usage:**
```bash
# Container stats
docker stats

# Specific container stats
docker stats <container-name>
```

#### **Container Health:**
```bash
# Health check status
docker inspect <container-name> | grep Health

# Container events
docker events
```

### **2. Log Management**

#### **Log Rotation:**
```yaml
logging:
  driver: "json-file"
  options:
    max-size: "10m"
    max-file: "3"
```

#### **Centralized Logging:**
```bash
# Collect logs from all services
docker-compose logs --tail=100 -f > all-services.log

# Filter logs by service
docker-compose logs -f gateway | grep ERROR
```

---

## 🎯 Best Practices

### **1. Dockerfile Optimization**
- **Sử dụng Alpine Linux** cho image size nhỏ
- **Multi-stage build** cho production
- **Non-root user** cho security
- **Health checks** cho monitoring

### **2. Docker Compose Best Practices**
- **Resource limits** cho mỗi service
- **Health checks** cho dependencies
- **Volume mounts** cho config và data
- **Network isolation** với custom networks

### **3. Deployment Best Practices**
- **Blue-green deployment** cho zero downtime
- **Rolling updates** cho services
- **Backup strategies** cho data
- **Monitoring và alerting**

---

## 🎯 Tóm Tắt

| Command | Mục đích | Khi nào sử dụng |
|---------|----------|-----------------|
| `docker-compose build` | Build images | Sau khi thay đổi code |
| `docker-compose up -d` | Start services | Development/Production |
| `docker-compose down` | Stop services | Khi cần dừng hệ thống |
| `docker-compose logs -f` | View logs | Debug issues |
| `docker-compose restart` | Restart service | Sau khi thay đổi config |
| `./deploy.sh <service>` | Deploy to production | Production deployment |

**Quy tắc vàng**: **Luôn sử dụng `docker-compose` cho local development và `deploy.sh` cho production deployment!**
