# BDSPro Microservices - Technology Stack

## 📋 Complete Technology Overview

**Project Type**: Microservices Architecture
**Primary Language**: Go (Golang) 1.21+
**Architecture Pattern**: Clean Architecture
**Communication**: gRPC + HTTP REST
**Deployment**: Docker + Docker Compose

---

## 🔧 Core Technologies

### Programming Language
- **Go (Golang) 1.21+**
  - Concurrency support
  - Fast compilation
  - Built-in testing framework
  - Strong typing
  - Excellent performance

---

## 🌐 Web Frameworks & Libraries

### HTTP Framework
- **Gin v1.9+**
  - Fast HTTP web framework
  - Middleware support
  - JSON validation
  - Route grouping
  - Path parameters
  - Query binding

```go
r := gin.Default()
r.POST("/api/users", handler.CreateUser)
r.GET("/api/users/:id", handler.GetUser)
```

### gRPC Framework
- **gRPC-Go v1.58+**
  - High-performance RPC framework
  - Protocol Buffers
  - Bi-directional streaming
  - Service-to-service communication
  - Built-in load balancing

### gRPC-Gateway
- **grpc-gateway v2.18+**
  - HTTP REST to gRPC translation
  - Swagger/OpenAPI generation
  - Reverse proxy for gRPC services
  - JSON ↔ Protocol Buffers conversion

```go
mux := runtime.NewServeMux()
err := userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
```

---

## 🗄️ Database & Persistence

### Primary Database
- **PostgreSQL 14+**
  - ACID compliance
  - JSON support (JSONB)
  - Full-text search
  - Advanced indexing
  - Partitioning support

### ORM
- **GORM v1.25+**
  - Full-featured ORM
  - Auto migrations
  - Association management
  - Hooks (BeforeCreate, AfterUpdate)
  - Soft deletes
  - Transaction support

```go
type User struct {
    gorm.Model
    Name  string
    Email string `gorm:"uniqueIndex"`
}

db.Create(&user)
db.Where("email = ?", email).First(&user)
```

### Caching
- **Redis 7+**
  - In-memory data store
  - Session storage
  - Caching layer
  - Pub/Sub messaging
  - Distributed locks

```go
redis.Set(ctx, "key", "value", 5*time.Minute)
val, err := redis.Get(ctx, "key").Result()
```

### Document Storage (Optional)
- **MongoDB 6+**
  - NoSQL database
  - Flexible schema
  - Document-oriented
  - Horizontal scalability

---

## 🔐 Security & Authentication

### JWT (JSON Web Tokens)
- **golang-jwt/jwt v5**
  - Token generation
  - Token validation
  - Claims management
  - Refresh tokens

```go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "userId": 123,
    "exp":    time.Now().Add(15 * time.Minute).Unix(),
})
```

### Password Hashing
- **bcrypt**
  - Secure password hashing
  - Salt generation
  - Adaptive cost factor

```go
hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
err := bcrypt.CompareHashAndPassword(hash, []byte(password))
```

### OAuth2
- **golang.org/x/oauth2**
  - OAuth2 client library
  - Google, Facebook, Zalo integration
  - Token exchange
  - Refresh token management

---

## 📡 Protocol Buffers & Code Generation

### Protocol Buffers
- **protoc v24+**
  - Interface Definition Language (IDL)
  - Efficient serialization
  - Language-neutral
  - Version compatibility

### Buf
- **buf v1.28+**
  - Modern Protobuf management
  - Linting
  - Breaking change detection
  - Code generation
  - Registry management

```yaml
# buf.yaml
version: v1
lint:
  use:
    - DEFAULT
breaking:
  use:
    - FILE
```

### Generated Code
- **protoc-gen-go v1.31+**: Go structs
- **protoc-gen-go-grpc v1.3+**: gRPC service code
- **protoc-gen-grpc-gateway v2.18+**: HTTP gateway
- **protoc-gen-openapiv2 v2.18+**: OpenAPI/Swagger specs

---

## 🔄 Dependency Injection

### Wire
- **google/wire v0.5+**
  - Compile-time dependency injection
  - Auto-generate initialization code
  - Type-safe
  - No reflection at runtime

```go
//go:build wireinject
// +build wireinject

func InitializeApp() (*App, error) {
    wire.Build(
        NewDatabase,
        NewRepository,
        NewUsecase,
        NewHandler,
        NewApp,
    )
    return nil, nil
}
```

---

## 📚 API Documentation

### Swagger/OpenAPI
- **swaggo/swag v1.16+**
  - Auto-generate Swagger docs from comments
  - Interactive API documentation
  - Code-first approach
  - Swagger UI integration

```go
// @Summary Create user
// @Description Create a new user account
// @Tags Users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User data"
// @Success 200 {object} User
// @Router /api/users [post]
func CreateUser(c *gin.Context) {
    // ...
}
```

---

## 🧪 Testing

### Testing Framework
- **Go Testing Package** (built-in)
  - Unit tests
  - Benchmark tests
  - Table-driven tests

```go
func TestUserService_Create(t *testing.T) {
    tests := []struct {
        name    string
        input   *User
        wantErr bool
    }{
        {"valid user", &User{Name: "John"}, false},
        {"empty name", &User{Name: ""}, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

### Mocking
- **gomock**
  - Mock generation
  - Interface mocking
  - Behavior verification

```go
mockCtrl := gomock.NewController(t)
mockRepo := NewMockUserRepository(mockCtrl)
mockRepo.EXPECT().GetByID(gomock.Any(), uint64(1)).Return(user, nil)
```

---

## 🐳 Containerization & Deployment

### Docker
- **Docker v24+**
  - Container runtime
  - Image building
  - Multi-stage builds
  - Layer caching

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
COPY --from=builder /app/main /main
CMD ["/main"]
```

### Docker Compose
- **Docker Compose v2.23+**
  - Multi-container orchestration
  - Service dependencies
  - Network management
  - Volume management

```yaml
version: '3.8'
services:
  gateway:
    build: ./gateway-service
    ports:
      - "8080:8080"
    depends_on:
      - user
      - auth
```

---

## 🔧 Development Tools

### Hot Reload
- **Air v1.49+**
  - Live reload during development
  - Auto-rebuild on file changes
  - Custom build commands
  - Fast feedback loop

```toml
# .air.toml
[build]
  cmd = "go build -o ./tmp/main ."
  bin = "./tmp/main"
  include_ext = ["go", "tpl", "tmpl", "html"]
```

### Configuration Management
- **Viper v1.17+**
  - Configuration file reading (YAML, JSON, TOML)
  - Environment variable support
  - Remote configuration
  - Live watching

```go
viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath("./config")
viper.ReadInConfig()

port := viper.GetInt("server.port")
```

### CLI Framework
- **Cobra v1.8+**
  - Command-line interface
  - Subcommands
  - Flags management
  - Help generation

```go
var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "My application",
}

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the server",
    Run: func(cmd *cobra.Command, args []string) {
        // Start server
    },
}

rootCmd.AddCommand(serveCmd)
```

---

## 🔗 HTTP Client & Utilities

### HTTP Client
- **resty v2.11+**
  - Simple HTTP client
  - Retry mechanism
  - Request/response middleware
  - JSON handling

```go
client := resty.New()
resp, err := client.R().
    SetHeader("Content-Type", "application/json").
    SetBody(user).
    Post("https://api.example.com/users")
```

### UUID Generation
- **google/uuid v1.4+**
  - UUID v4 generation
  - Unique identifiers
  - Thread-safe

```go
id := uuid.New().String()
```

### Time Utilities
- **go-playground/validator v10**
  - Struct validation
  - Custom validators
  - Tag-based validation

```go
type User struct {
    Email string `validate:"required,email"`
    Age   int    `validate:"gte=0,lte=130"`
}

validate := validator.New()
err := validate.Struct(user)
```

---

## 📊 Monitoring & Logging (Planned)

### Logging
- **logrus** or **zap**
  - Structured logging
  - Multiple log levels
  - JSON formatting
  - Hook support

```go
log.WithFields(log.Fields{
    "userId": 123,
    "action": "create",
}).Info("User created")
```

### Metrics (Future)
- **Prometheus**
  - Metrics collection
  - Time-series database
  - PromQL queries
  - Alerting

### Tracing (Future)
- **Jaeger** or **Zipkin**
  - Distributed tracing
  - Request flow visualization
  - Performance analysis

---

## 🗂️ Data Serialization

### JSON
- **encoding/json** (built-in)
  - JSON marshaling/unmarshaling
  - Struct tags
  - Custom marshalers

### Protocol Buffers
- **protobuf**
  - Binary serialization
  - Smaller payload size
  - Fast serialization/deserialization
  - Backward compatibility

---

## 🔌 Service Communication

### gRPC
- **Primary**: Service-to-service communication
- **Features**:
  - Binary protocol
  - HTTP/2
  - Bi-directional streaming
  - Load balancing
  - Connection pooling

### HTTP REST
- **Secondary**: External client communication
- **Features**:
  - JSON over HTTP
  - RESTful conventions
  - Stateless
  - Browser-friendly

---

## 📦 Build Tools

### Make
- **GNU Make**
  - Build automation
  - Task runner
  - Dependency management

```makefile
.PHONY: build run test

build:
	go build -o bin/app main.go

run: build
	./bin/app

test:
	go test ./...
```

### Go Modules
- **go mod**
  - Dependency management
  - Version control
  - Reproducible builds

```bash
go mod init myapp
go mod tidy
go mod download
```

---

## 🌍 Cross-Cutting Concerns

### CORS
- **gin-contrib/cors**
  - Cross-Origin Resource Sharing
  - Configurable origins
  - Credential support

```go
r.Use(cors.New(cors.Config{
    AllowOrigins: []string{"http://localhost:3000"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders: []string{"Authorization", "Content-Type"},
}))
```

### Rate Limiting (Future)
- **golang.org/x/time/rate**
  - Token bucket algorithm
  - Request throttling
  - Per-user limits

---

## 📱 Push Notifications

### Firebase Cloud Messaging (FCM)
- **firebase.google.com/go**
  - Push notifications
  - Topic-based messaging
  - Device tokens
  - Platform-agnostic

```go
message := &messaging.Message{
    Notification: &messaging.Notification{
        Title: "New Message",
        Body:  "You have a new message",
    },
    Token: deviceToken,
}
response, err := fcm.Send(ctx, message)
```

---

## 🔄 Background Jobs (Future)

### Job Queue
- **asynq** or **machinery**
  - Task queue
  - Scheduled jobs
  - Retry mechanism
  - Worker pools

---

## 📈 Development Workflow

### Version Control
- **Git**
  - Source control
  - Branching strategy
  - Pull requests

### CI/CD (Future)
- **GitHub Actions** or **GitLab CI**
  - Automated testing
  - Automated deployment
  - Build pipelines

---

## 🎯 Key Technology Decisions

### Why Go?
1. **Performance**: Compiled language, fast execution
2. **Concurrency**: Goroutines & channels
3. **Simplicity**: Easy to learn and maintain
4. **Tooling**: Excellent built-in tools
5. **Community**: Large ecosystem

### Why gRPC?
1. **Performance**: Binary protocol, HTTP/2
2. **Type Safety**: Protocol Buffers
3. **Streaming**: Bi-directional streaming
4. **Multi-Language**: Polyglot support

### Why PostgreSQL?
1. **Reliability**: ACID compliance
2. **Features**: Advanced SQL features
3. **Extensions**: PostGIS, full-text search
4. **Performance**: Excellent query optimization

### Why Clean Architecture?
1. **Maintainability**: Clear separation of concerns
2. **Testability**: Easy to mock dependencies
3. **Flexibility**: Easy to swap implementations
4. **Scalability**: Modular design

---

## 📊 Technology Matrix

| Category | Technology | Version | Purpose |
|----------|-----------|---------|---------|
| **Language** | Go | 1.21+ | Primary language |
| **HTTP Framework** | Gin | 1.9+ | REST API |
| **RPC Framework** | gRPC | 1.58+ | Service communication |
| **Database** | PostgreSQL | 14+ | Primary storage |
| **ORM** | GORM | 1.25+ | Database access |
| **Cache** | Redis | 7+ | Caching & sessions |
| **DI** | Wire | 0.5+ | Dependency injection |
| **Docs** | Swagger | 1.16+ | API documentation |
| **Protobuf** | Buf | 1.28+ | Proto management |
| **Config** | Viper | 1.17+ | Configuration |
| **CLI** | Cobra | 1.8+ | Command-line |
| **Testing** | Go Test | built-in | Unit testing |
| **Container** | Docker | 24+ | Containerization |
| **Orchestration** | Docker Compose | 2.23+ | Local development |
| **Hot Reload** | Air | 1.49+ | Development |

---

## 🚀 Performance Characteristics

### Go Performance
- **Compilation**: < 5 seconds for medium projects
- **Startup Time**: < 100ms
- **Memory**: ~10-50MB per service
- **Concurrency**: Thousands of goroutines

### gRPC Performance
- **Latency**: 1-5ms (local network)
- **Throughput**: 10,000+ RPS per service
- **Payload Size**: ~50% smaller than JSON

### Database Performance
- **Query Time**: < 10ms (with indexes)
- **Connection Pool**: 10-100 connections
- **Transaction**: ACID compliant

---

**Last Updated**: October 15, 2025
**Stack Version**: v1.0
**Status**: ✅ Production Ready

