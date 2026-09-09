# IV.1 Gateway Service - API Gateway

## 🎯 Tổng Quan

Gateway Service là điểm vào duy nhất (Single Entry Point) của hệ thống BDSPro Microservices. Nó đóng vai trò là API Gateway, xử lý routing, authentication, load balancing, và cung cấp giao diện thống nhất cho tất cả các microservices.

---

## 🏗️ Kiến Trúc Gateway

### **Vai Trò Chính:**
- **API Gateway**: Điểm vào duy nhất cho tất cả client requests
- **Request Routing**: Định tuyến requests đến các microservices phù hợp
- **Authentication & Authorization**: Xử lý JWT authentication
- **Protocol Translation**: Chuyển đổi HTTP REST ↔ gRPC
- **Load Balancing**: Phân tải requests đến các service instances
- **Documentation**: Cung cấp Swagger documentation cho tất cả APIs

### **Hai Chế Độ Hoạt Động:**
1. **HTTP Mode**: REST API Gateway với gRPC Gateway
2. **gRPC Mode**: Pure gRPC Gateway

---

## 📁 Cấu Trúc Gateway Service

```
gateway-service/
├── cmd/                          # Entry points
│   ├── grpc/                     # gRPC mode
│   │   ├── main.go              # gRPC server entry
│   │   └── handler.go           # gRPC handlers
│   └── http/                    # HTTP mode
│       ├── main.go              # HTTP server entry
│       ├── response.go          # Response handlers
│       └── swagger.go           # Swagger handlers
├── config/                       # Configuration
│   ├── config.go                # Config loader
│   └── config.yml               # Service configuration
├── docs/                         # Swagger documentation
│   ├── <service>/               # Per-service swagger docs
│   └── merged_swagger.json      # Combined swagger
├── pkg/                          # Core packages
│   ├── gateway.go               # Gateway core
│   ├── proxy.go                 # HTTP proxy logic
│   ├── service.go               # Service registry
│   └── eureka.go                # Service discovery
├── service/                      # Service registrations
│   ├── user.go                  # User service registration
│   ├── organization.go          # Organization service registration
│   ├── auth.go                  # Auth service registration
│   └── ...                      # Other services
├── main.go                       # Main entry point
└── Dockerfile                    # Container configuration
```

---

## ⚙️ Configuration

### **Service Configuration** (`config/config.yml`):
```yaml
server:
  port: 8000          # HTTP Gateway port
  tcp_port: 8800      # gRPC Gateway port

service:
  domain: localhost   # Service domain
  port:
    gateway: 8000     # Gateway service port
    user: 8201        # User service port
    bdspro: 8202      # BDSPro service port
    file: 8203        # File service port
    notification: 8204 # Notification service port
    payment: 8205     # Payment service port
    crm: 8206         # CRM service port
    organization: 8207 # Organization service port
    chat: 8208        # Chat service port
    social: 8209      # Social service port
    map: 8210         # Map service port
    search: 8211      # Search service port
    membership: 8212  # Membership service port
    appointment: 8213 # Appointment service port
    marketing: 8214   # Marketing service port
    transaction: 8215 # Transaction service port
    auth: 8216        # Auth service port

jwt:
  key-generate: key-generate-jwt-token-authentication
  tokenPrefix: Bearer
  tokenExpirationAfterDays: 14
  accessExpAfterMinutes: 15      # Access token: 15 phút
  refreshExpAfterMinutes: 129600 # Refresh token: 3 tháng
  authorizationHeader: Authorization
```

---

## 🚀 Service Registration

### **gRPC Service Registration Pattern:**
```go
func RegisterUserService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
    port := viper.GetString("service.port.user")
    url := "localhost"
    if env != "local" {
        url = "user"  // Docker container name
    }
    endpoint := fmt.Sprintf("%s:%s", url, port)
    
    // Register multiple gRPC services
    _ = userpb.RegisterProfileServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
    _ = userpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
    _ = userpb.RegisterRateServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
    _ = userpb.RegisterReportServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
```

### **Registered Services:**
| Service | gRPC Services | Port | Container Name |
|---------|---------------|------|----------------|
| **User** | Profile, Auth, Rate, Report | 8201 | user_grpc |
| **Organization** | Organization, Branch, Group, Deal | 8207 | organization |
| **Auth** | Auth, OAuth, Token | 8216 | auth |
| **BDSPro** | Post, Product, Market | 8202 | bdspro_grpc |
| **CRM** | Contact, Lead, Customer | 8206 | crm |
| **Payment** | Payment, Wallet, Bank | 8205 | payment |
| **Notification** | Notification, Push | 8204 | notification |
| **Chat** | Chat, Message, Room | 8208 | chat |
| **Social** | NewsFeed, Comment, Like | 8209 | social |
| **Appointment** | Appointment, Schedule | 8213 | appointment |
| **Transaction** | Transaction, Deal | 8215 | transaction |
| **Marketing** | Campaign, Lead | 8214 | marketing |

---

## 🔐 Authentication & Authorization

### **JWT Middleware:**
```go
// JWT Authentication middleware
r.Any("/v2/*any", _jwt.JWTAuthMiddleware(publicRoutes), gin.WrapH(mux))
r.Any("/v1/:service/*path", _jwt.JWTAuthMiddleware(publicRoutes), forwarding)
```

### **Public Routes (Không cần authentication):**
```go
var publicRoutes = []string{
    // BDS Pro public routes
    "/v1/bdspro/v1/user/post",
    "/v1/bdspro/v2/user/post",
    "/v1/bdspro/v2/list/property-type",
    
    // User public routes
    "/v1/user/rate/list",
    "/v1/user/oauth/zalo/login",
    "/v1/user/oauth/facebook/login",
    "/v2/user/profile/info/",
    
    // Auth public routes
    "/v2/auth/otp",
    "/v2/auth/token/refresh",
    "/v2/auth/oauth",
    
    // File public routes
    "/v1/file/upload",
    "/v1/file/load",
    
    // Other public routes
    "/v1/map/locations/nearby",
    "/v1/membership/plan/list",
    "/v1/payment/bank/all",
    "/v2/social/news-feed/public",
}
```

### **JWT Token Flow:**
1. **Client** gửi request với `Authorization: Bearer <token>`
2. **Gateway** extract JWT từ header
3. **Gateway** validate JWT token
4. **Gateway** inject user info vào gRPC metadata
5. **Microservice** nhận user info từ metadata

---

## 🌐 API Routing

### **1. v2 API Routes (gRPC Gateway)**
```go
// gRPC Gateway routes
r.Any("/v2/*any", _jwt.JWTAuthMiddleware(publicRoutes), gin.WrapH(mux))
```

**URL Pattern:**
```
/v2/<service>/<method>
```

**Examples:**
```
POST /v2/user/profile/create
GET  /v2/organization/branch/list
PUT  /v2/bdspro/post/update
POST /v2/payment/wallet/create
```

### **2. v1 API Routes (HTTP Proxy)**
```go
// Legacy HTTP proxy routes
r.Any("/v1/:service/*path", _jwt.JWTAuthMiddleware(publicRoutes), forwarding)
```

**URL Pattern:**
```
/v1/<service>/<path>
```

**Examples:**
```
POST /v1/user/profile/create
GET  /v1/organization/branch/list
PUT  /v1/bdspro/post/update
POST /v1/payment/wallet/create
```

### **3. Service Port Mapping (v1 Proxy):**
```go
var servicePorts = map[string]int{
    "user":         8001,
    "file":         8002,
    "config":       8003,
    "notification": 8004,
    "payment":      8005,
    "crm":          8006,
    "membership":   8007,
    "search":       8008,
    "map":          8101,
    "bdspro":       8102,
    "tho247":       8103,
    "chat":         8104,
}
```

---

## 📚 Swagger Documentation

### **Swagger Endpoints:**
```go
// Merged swagger (tất cả services)
r.GET("/swagger/merged/*any", mergedSwaggerHandler)

// Individual service swagger
r.GET("/swagger/:service/*any", individualSwaggerHandler)

// Static swagger docs
r.Static("/swagger-doc", "docs")
```

### **Swagger URLs:**
| Endpoint | Description |
|----------|-------------|
| `/swagger/merged/` | Tất cả APIs của tất cả services |
| `/swagger/user/` | APIs của User service |
| `/swagger/organization/` | APIs của Organization service |
| `/swagger/bdspro/` | APIs của BDSPro service |
| `/swagger/payment/` | APIs của Payment service |
| `/swagger/notification/` | APIs của Notification service |
| `/swagger/crm/` | APIs của CRM service |
| `/swagger/chat/` | APIs của Chat service |
| `/swagger/social/` | APIs của Social service |
| `/swagger/appointment/` | APIs của Appointment service |
| `/swagger/transaction/` | APIs của Transaction service |
| `/swagger/auth/` | APIs của Auth service |

---

## 🔄 Request Flow

### **1. HTTP Request Flow:**
```
Client Request
    ↓
Gateway (Port 8000)
    ↓
JWT Authentication
    ↓
Route to Service
    ↓
gRPC Call to Microservice
    ↓
Response back to Client
```

### **2. gRPC Gateway Flow:**
```
HTTP Request → gRPC Gateway → gRPC Service → Response
```

### **3. HTTP Proxy Flow:**
```
HTTP Request → HTTP Proxy → HTTP Service → Response
```

---

## 🌍 CORS Configuration

### **CORS Settings:**
```go
r.Use(cors.New(cors.Config{
    AllowOrigins: []string{
        "http://localhost:3000",
        "http://localhost:3001", 
        "http://localhost:8000",
        "http://14.225.210.29:8000",
        "http://14.225.205.222:8000",
        "http://103.145.63.185:3107",
    },
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Lang"},
    ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}))
```

---

## 🚀 Running Gateway

### **1. HTTP Mode (Recommended):**
```bash
# Run HTTP Gateway
./main http

# Or with Docker
docker run -p 8000:8000 gateway-service ./main http
```

### **2. gRPC Mode:**
```bash
# Run gRPC Gateway
./main grpc

# Or with Docker
docker run -p 8800:8800 gateway-service ./main grpc
```

### **3. Environment Variables:**
```bash
# Development
ENV_RUNTIME=development

# Production
ENV_RUNTIME=release
```

---

## 🔧 Development Commands

### **Build Gateway:**
```bash
# Build binary
go build -o main .

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o main .
```

### **Generate Swagger:**
```bash
# Generate swagger docs
make swag gateway

# Merge all swagger docs
make merge-swagger
```

### **Run with Air (Hot Reload):**
```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

---

## 🐳 Docker Deployment

### **Dockerfile:**
```dockerfile
FROM alpine:latest

WORKDIR /

COPY main .
COPY docs docs

ENV ENV_RUNTIME=develop

EXPOSE 8080

CMD ["./main", "http"]
```

### **Docker Compose:**
```yaml
gateway:
  build:
    context: ./gateway-service
    dockerfile: Dockerfile
  container_name: gateway
  networks:
    - internal
  ports:
    - "8000:8000"
```

---

## 📊 Monitoring & Health Checks

### **Health Check Endpoints:**
```bash
# Gateway health
curl http://localhost:8000/health

# Service status
curl http://localhost:8000/status
```

### **Logs:**
```bash
# View gateway logs
docker logs gateway -f

# Check service connectivity
docker exec gateway ping user_grpc
docker exec gateway ping organization
```

---

## 🎯 Best Practices

### **1. Gateway Design:**
- **Single Entry Point**: Tất cả requests đi qua Gateway
- **Authentication**: JWT validation tại Gateway
- **Rate Limiting**: Implement rate limiting
- **Circuit Breaker**: Implement circuit breaker pattern
- **Load Balancing**: Distribute load across service instances

### **2. Security:**
- **HTTPS**: Sử dụng HTTPS trong production
- **CORS**: Configure CORS properly
- **JWT**: Validate JWT tokens
- **Input Validation**: Validate all inputs
- **Error Handling**: Don't expose internal errors

### **3. Performance:**
- **Caching**: Implement response caching
- **Compression**: Enable gzip compression
- **Connection Pooling**: Pool gRPC connections
- **Monitoring**: Monitor performance metrics

---

## 🎯 Tóm Tắt

| Component | Mục đích | Khi nào sử dụng |
|-----------|----------|-----------------|
| **HTTP Gateway** | REST API entry point | Client applications, web frontend |
| **gRPC Gateway** | gRPC API entry point | Service-to-service communication |
| **JWT Middleware** | Authentication & authorization | Secure API access |
| **Swagger Docs** | API documentation | Development, testing, integration |
| **CORS** | Cross-origin requests | Web applications |
| **Service Registry** | Service discovery | Dynamic service routing |

**Quy tắc vàng**: **Gateway là điểm vào duy nhất - tất cả client requests phải đi qua Gateway để đảm bảo security, authentication và routing consistency!**
