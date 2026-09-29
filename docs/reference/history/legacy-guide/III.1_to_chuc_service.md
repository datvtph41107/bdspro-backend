# III.1 Tổ Chức Service - Clean Architecture

## 🎯 Tổng Quan

Dự án BDSPro Microservices tuân theo **Clean Architecture** với cấu trúc rõ ràng, tách biệt các layer và dễ bảo trì. Auth-service được sử dụng làm chuẩn để minh họa cấu trúc tổ chức service.

---

## 📁 Cấu Trúc Thư Mục Chuẩn

```
<service-name>-service/
├── cmd/                          # Entry points
│   ├── grpc.go                   # gRPC server entry point
│   └── http.go                   # HTTP server entry point (nếu có)
├── config/                       # Configuration files
│   ├── app.yml                   # Application configuration
│   ├── properties.go             # Configuration properties
│   └── runtime.go                # Runtime configuration
├── env/                          # Environment configuration
│   └── runtime.go                # Environment-specific runtime settings
├── infra/                        # Infrastructure layer
│   ├── client/                   # External service clients
│   │   ├── grpc_client.go
│   │   ├── http_client.go
│   │   ├── notification_client.go
│   │   ├── organization_client.go
│   │   ├── payment_client.go
│   │   └── user_client.go
│   ├── cookie/                   # Cookie providers
│   │   └── cookie_provider.go
│   ├── db/                       # Database utilities
│   │   └── migrate.go
│   ├── handler/                  # gRPC/HTTP handlers
│   │   ├── auth_handler.go
│   │   ├── admin_handler.go
│   │   ├── oauth_handler.go
│   │   ├── permission_handler.go
│   │   ├── role_handler.go
│   │   └── role_group_handler.go
│   ├── mapper/                   # Data transformation
│   │   ├── role_group_mapper.go
│   │   └── role_mapper.go
│   ├── oauth/                    # OAuth providers
│   │   ├── facebook_oauth.go
│   │   ├── google_oauth.go
│   │   └── zalo_oauth.go
│   ├── postgre/                  # PostgreSQL implementations
│   │   ├── admin_access_postgres.go
│   │   ├── auth_method_postgres.go
│   │   ├── otp_postgres.go
│   │   ├── permission_repo.go
│   │   ├── role_group_postgres.go
│   │   ├── role_postgres.go
│   │   ├── session_postgres.go
│   │   └── status_postgres.go
│   ├── properties/               # Runtime properties
│   │   └── runtime.go
│   └── redis/                    # Redis providers
│       └── redis_provider.go
├── internal/                     # Internal business logic
│   ├── domain/                   # Domain entities
│   │   ├── auth_method.go
│   │   ├── auth_id.go
│   │   ├── otp.go
│   │   ├── ...
│   ├── dto/                      # Data Transfer Objects
│   │   ├── auth_dto.go
│   │   ├── permission_dto.go
│   │   ├── profile_dto.go
│   │   ├── ...
│   ├── enums/                    # Enumerations
│   │   ├── auth_enums.go
│   │   ├── permission_type.go
│   │   ├── permission.go
│   │   └── ...
│   ├── interface/                # Interface definitions
│   │   ├── factory/              # Factory interfaces
│   │   │   └── internal.go
│   │   ├── provider/             # Provider interfaces
│   │   │   ├── auth_validator.go
│   │   │   ├── cache_provider.go
│   │   │   ├── cookie_provider.go
│   │   │   ├── notification_provider.go
│   │   │   ├── organization_provider.go
│   │   │   ├── payment_provider.go
│   │   │   └── profile_provider.go
│   │   └── repo/                 # Repository interfaces
│   │       ├── admin_access_repo.go
│   │       ├── auth_method_repo.go
│   │       ├── otp_repo.go
│   │       ├── permission_repo.go
│   │       ├── profile_repo.go
│   │       ├── role_group_repo.go
│   │       ├── role_repo.go
│   │       ├── session_repo.go
│   │       └── status_repo.go
│   ├── usecase/                  # Business logic
│   │   ├── admin_access_usecase.go
│   │   ├── admin_usecase.go
│   │   ├── auth_usecase.go
│   │   ├── otp_usecase.go
│   │   ├── permission_usecase.go
│   │   ├── role_group_usecase.go
│   │   └── role_usecase.go
│   └── utils/                    # Utility functions
│       ├── code.go
│       └── password.go
├── initial/                      # Initialization
│   └── runtime.go                # Chứa các đối tượng khởi tạo project. 
│                                 # Thường bao gồm các handler hoặc ghi log hoặc khởi tạo db
├── validator/                    # Request validation
│   └── role_group_validator.go
├── wire/                         # Dependency injection
│   ├── wire.go                   # File này sẽ tự động sinh ra khi gõ script. Không cần tạo tay
│   └── wire_gen.go               # File được generate tự động bởi Wire
├── main.go                       # Main entry point
├── go.mod                        # Go module file
├── go.sum                        # Go module checksums
└── Dockerfile                    # Docker configuration
```

---

## 🏗️ Clean Architecture Layers

### **1. Domain Layer** (`internal/domain/`)
**Mục đích**: Chứa các entity và business objects cốt lõi

```go
// internal/domain/auth_method.go
package domain

import "time"

// AuthMethodDomain tương đương với Java Entity
type AuthMethodDomain struct {
    ID          uint64     `gorm:"primaryKey;autoIncrement"`
    Provider    string     `gorm:"size:10"` // enum: email, zalo, facebook, phone, admin
    Username    string     `gorm:"size:50"`
    Password    string     `gorm:"size:255"` // Mã hóa password
    OAuthID     string     `gorm:"column:oauth_id"`
    Avatar      string     `gorm:"size:255"`
    IsSensitive bool       `gorm:"column:is_sensitive"`
    Name        string     `gorm:"size:255"`
    Email       string     `gorm:"size:255"`
    Phone       string     `gorm:"size:255"`
    CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
    DeletedAt   *time.Time `gorm:"column:deleted_at"`
    UserID      uint64     `gorm:"column:user_id"` // Không cập nhật & xóa, chỉ tạo mới
}

// TableName đặt tên bảng tương đương với Java
func (AuthMethodDomain) TableName() string {
    return "auth_method"
}
```

**Quy tắc**:
- Chỉ chứa struct và business logic cơ bản
- Không phụ thuộc vào framework hay database
- Định nghĩa TableName() cho GORM

### **2. Interface Layer** (`internal/interface/`)
**Mục đích**: Định nghĩa các interface cho dependencies

#### **Repository Interfaces** (`internal/interface/repo/`)
```go
// internal/interface/repo/auth_method_repo.go
package repo

import (
    "auth/internal/domain"
    "context"
)

type AuthMethodRepository interface {
    // Existing methods
    FindByUsername(ctx context.Context, username string) (*domain.AuthMethodDomain, error)
    FindByPhone(ctx context.Context, phone string) (*domain.AuthMethodDomain, error)
    FindByEmail(ctx context.Context, email string) (*domain.AuthMethodDomain, error)
    FindByOAuthID(ctx context.Context, oauthID, provider string) (*domain.AuthMethodDomain, error)
    Create(ctx context.Context, auth *domain.AuthMethodDomain) (*domain.AuthMethodDomain, error)
    Update(ctx context.Context, auth *domain.AuthMethodDomain) (*domain.AuthMethodDomain, error)
    Delete(ctx context.Context, id uint64) error
    SoftDelete(ctx context.Context, id uint64) error
    FindByID(ctx context.Context, id uint64) (*domain.AuthMethodDomain, error)

    // Admin management methods
    FindAdminsByProvider(c context.Context, provider string, page, size int) ([]*domain.AuthMethodDomain, int64, error)
    FindByUsernameAndProvider(c context.Context, username, provider string) (*domain.AuthMethodDomain, error)
    CountByProvider(c context.Context, provider string) (int64, error)

    // Account restoration methods
    RestoreAccount(c context.Context, id uint64) error
    FindDeletedByPhone(c context.Context, phone string) (*domain.AuthMethodDomain, error)
    FindDeletedByEmail(c context.Context, email string) (*domain.AuthMethodDomain, error)
    FindDeletedByUsername(c context.Context, username string) (*domain.AuthMethodDomain, error)
}
```

#### **Provider Interfaces** (`internal/interface/provider/`)
```go
// internal/interface/provider/notification_provider.go
package provider

import "context"

type NotificationProvider interface {
    SendOTP(ctx context.Context, phone, otp string) error
    SendEmail(ctx context.Context, email, subject, body string) error
}
```

**Quy tắc**:
- Định nghĩa interface cho tất cả external dependencies
- Không chứa implementation
- Sử dụng context.Context cho tất cả methods

### **3. Usecase Layer** (`internal/usecase/`)
**Mục đích**: Chứa business logic và orchestration

```go
// internal/usecase/auth_usecase.go
package usecase

import (
    "context"
    "auth/internal/domain"
    "auth/internal/dto"
    "auth/internal/interface/provider"
    "auth/internal/interface/repo"
)

type AuthUsecase struct {
    CookieProvider     provider.CookieProvider
    CacheProvider      provider.CacheProvider
    ProfileProvider    provider.ProfileProvider
    AuthMethodRepo     repo.AuthMethodRepository
    SessionRepo        repo.SessionRepository
    OTPRepo            repo.OTPRepository
    StatusRepo         repo.StatusRepository
    OtpUsecase         *OtpUsecase
    OrganizationClient provider.OrganizationProvider
    PaymentProvider    provider.PaymentProvider
    NotificationClient provider.NotificationProvider
    properties         dto.PropertiesDTO
}

func NewAuthService(
    cookieProvider provider.CookieProvider,
    cacheProvider provider.CacheProvider,
    profileProvider provider.ProfileProvider,
    authMethodRepo repo.AuthMethodRepository,
    sessionRepo repo.SessionRepository,
    otpRepo repo.OTPRepository,
    statusRepo repo.StatusRepository,
    otpUsecase *OtpUsecase,
    organizationClient provider.OrganizationProvider,
    paymentProvider provider.PaymentProvider,
    notificationClient provider.NotificationProvider,
) *AuthUsecase {
    return &AuthUsecase{
        CookieProvider:     cookieProvider,
        CacheProvider:      cacheProvider,
        ProfileProvider:    profileProvider,
        AuthMethodRepo:     authMethodRepo,
        SessionRepo:        sessionRepo,
        OTPRepo:            otpRepo,
        StatusRepo:         statusRepo,
        OtpUsecase:         otpUsecase,
        OrganizationClient: organizationClient,
        PaymentProvider:    paymentProvider,
        NotificationClient: notificationClient,
    }
}

// Business logic methods
func (u *AuthUsecase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
    // Business logic implementation
}
```

**Quy tắc**:
- Chứa business logic chính
- Inject dependencies qua constructor
- Sử dụng interfaces, không phụ thuộc vào implementation
- Return DTOs, không return domain entities trực tiếp

### **4. Infrastructure Layer** (`infra/`)
**Mục đích**: Triển khai các interface và external dependencies

#### **Repository Implementations** (`infra/postgre/`)
```go
// infra/postgre/auth_method_postgres.go
package postgre

import (
    "auth/internal/domain"
    "auth/internal/interface/repo"
    "context"
    "gorm.io/gorm"
)

type AuthMethodPostgres struct {
    db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) repo.AuthMethodRepository {
    return &AuthMethodPostgres{db: db}
}

func (r *AuthMethodPostgres) FindByUsername(ctx context.Context, username string) (*domain.AuthMethodDomain, error) {
    var auth domain.AuthMethodDomain
    err := r.db.WithContext(ctx).Where("username = ?", username).First(&auth).Error
    if err != nil {
        return nil, err
    }
    return &auth, nil
}

// Implement other methods...
```

#### **Handler Implementations** (`infra/handler/`)
```go
// infra/handler/auth_handler.go
package handler

import (
    "auth/internal/dto"
    "auth/internal/usecase"
    "context"
    authpb "pb/types/auth"
)

type AuthHandler struct {
    authpb.UnimplementedAuthServiceServer
    AuthUsecase        *usecase.AuthUsecase
    AdminAccessUsecase *usecase.AdminAccessUsecase
    InternalHandler    *InternalHandler
    AdminUsecase       *usecase.AdminUsecase
}

func NewAuthHandler(
    service *usecase.AuthUsecase,
    adminAccessUsecase *usecase.AdminAccessUsecase,
    internalHandler *InternalHandler,
    adminUsecase *usecase.AdminUsecase,
) *AuthHandler {
    return &AuthHandler{
        AuthUsecase:        service,
        AdminAccessUsecase: adminAccessUsecase,
        InternalHandler:    internalHandler,
        AdminUsecase:       adminUsecase,
    }
}

// @Summary Gửi OTP
// @Description Gửi OTP để đăng nhập
// @Tags Auth
// @Accept json
// @Produce json
// @Param phone body authpb.RequestOTPRequest true "Số điện thoại"
// @Success 200 {object} authpb.RequestOTPResponse
// @Failure 400 {object} authpb.RequestOTPResponse
// @Failure 500 {object} authpb.RequestOTPResponse
// @Router /otp/request [post]
func (h *AuthHandler) RequestOTP(ctx context.Context, req *authpb.RequestOTPRequest) (*authpb.RequestOTPResponse, error) {
    // Convert proto request to DTO
    dto := &dto.OtpRequest{
        Phone:    req.Phone,
        Fullname: req.Fullname,
    }
    
    // Call usecase
    result, err := h.AuthUsecase.RequestOTP(ctx, dto)
    if err != nil {
        return nil, err
    }
    
    // Convert DTO to proto response
    return &authpb.RequestOTPResponse{
        Success: result.Success,
        Message: result.Message,
    }, nil
}
```

**Quy tắc**:
- Implement các interface được định nghĩa trong internal/interface
- Chứa technical details (database, HTTP, gRPC)
- Không chứa business logic
- Sử dụng @bind comment cho Wire

### **5. DTO Layer** (`internal/dto/`)
**Mục đích**: Data Transfer Objects cho giao tiếp giữa các layer

```go
// internal/dto/auth.go
package dto

type LoginRequest struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
    Token     string `json:"token"`
    UserID    uint64 `json:"user_id"`
    ExpiresAt int64  `json:"expires_at"`
}
```

**Quy tắc**:
- Chỉ dùng để giao tiếp giữa các layer
- Không chứa business logic
- Có validation tags

### **6. Environment Configuration** (`env/`)
**Mục đích**: Cấu hình môi trường runtime

```go
// env/runtime.go
package env

import "auth/config"

type RuntimeEnv struct {
    Environment string
    Debug       bool
    LogLevel    string
}

func LoadRuntimeEnv() *RuntimeEnv {
    return &RuntimeEnv{
        Environment: config.Properties.Environment,
        Debug:       config.Properties.Debug,
        LogLevel:    config.Properties.LogLevel,
    }
}
```

**Quy tắc**:
- Chứa cấu hình môi trường runtime
- Khác biệt giữa development, staging, production
- Load từ environment variables hoặc config files

### **7. Initialization** (`initial/`)
**Mục đích**: Khởi tạo các đối tượng và dependencies

```go
// initial/runtime.go
package initial

import (
    "auth/infra/handler"
    "auth/infra/db"
    "log"
)

type InitialApp struct {
    AuthHandler        *handler.AuthHandler
    AdminHandler       *handler.AdminHandler
    OAuthHandler       *handler.OAuthHandler
    PermissionHandler  *handler.PermissionHandler
    RoleHandler        *handler.RoleHandler
    RoleGroupHandler   *handler.RoleGroupHandler
    InternalHandler    *handler.InternalHandler
}

func NewApp(
    authHandler *handler.AuthHandler,
    adminHandler *handler.AdminHandler,
    oauthHandler *handler.OAuthHandler,
    permissionHandler *handler.PermissionHandler,
    roleHandler *handler.RoleHandler,
    roleGroupHandler *handler.RoleGroupHandler,
    internalHandler *handler.InternalHandler,
) *InitialApp {
    // Khởi tạo database
    db.AutoMigrate()
    
    // Khởi tạo logging
    log.Println("Application initialized successfully")
    
    return &InitialApp{
        AuthHandler:      authHandler,
        AdminHandler:     adminHandler,
        OAuthHandler:     oauthHandler,
        PermissionHandler: permissionHandler,
        RoleHandler:      roleHandler,
        RoleGroupHandler: roleGroupHandler,
        InternalHandler:  internalHandler,
    }
}
```

**Quy tắc**:
- Chứa các đối tượng khởi tạo project
- Thường bao gồm các handler hoặc ghi log hoặc khởi tạo db
- Được sử dụng trong Wire để inject dependencies

### **8. Entry Points** (`cmd/`, `main.go`)
**Mục đích**: Khởi tạo và chạy service

```go
// main.go
package main

import (
    grpc "auth/cmd"
    "auth/internal/utils"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

// @title Auth Service API
// @version 1.0
// @description API phần logic xác thực và đăng nhập
// @BasePath /v2/auth
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
var rootCmd = &cobra.Command{
    Use:   "auth-service",
    Short: "Service xử lý authentication và authorization",
}

func init() {
    rootCmd.AddCommand(grpc.GrpcCmd)
    viper.AutomaticEnv()
}

func main() {
    go func() {
        log.Println("run go-routine")
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
```

```go
// cmd/grpc.go
package cmd

import (
    "auth/config"
    "auth/infra/db"
    "auth/wire"
    _middleware "common/middleware"
    "fmt"
    "log"
    "net"
    authpb "pb/types/auth"

    "github.com/spf13/cobra"
    "google.golang.org/grpc"
)

var GrpcCmd = &cobra.Command{
    Use:   "grpc",
    Short: "Khởi chạy gRPC server",
    Run: func(cmd *cobra.Command, args []string) {
        startGRPCServer()
    },
}

func startGRPCServer() {
    config.LoadConfig()
    
    // Start gRPC server
    port := fmt.Sprintf(":%d", config.Properties.Server.TCPPort)
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    s := grpc.NewServer(
        grpc.UnaryInterceptor(_middleware.ParseGrpcMetadataContextMiddleware),
    )

    app, _ := wire.InitializeApp()
    db.AutoMigrate()

    // Register services
    authpb.RegisterAuthServiceServer(s, app.AuthHandler)
    authpb.RegisterOAuthServiceServer(s, app.OAuthHandler)
    authpb.RegisterRoleServiceServer(s, app.RoleHandler)
    authpb.RegisterPermissionServiceServer(s, app.PermissionHandler)
    authpb.RegisterRoleGroupServiceServer(s, app.RoleGroupHandler)
    authpb.RegisterAuthInternalServiceServer(s, app.InternalHandler)

    log.Printf("gRPC server listening on %s", port)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
```

---

## 🔄 Dependency Flow

### **Dependency Direction** (Clean Architecture)
```
┌─────────────────┐
│   Entry Points  │ (main.go, cmd/)
└─────────┬───────┘
          │
┌─────────▼───────┐
│   Handlers      │ (infra/handler/)
└─────────┬───────┘
          │
┌─────────▼───────┐
│   Usecases      │ (internal/usecase/)
└─────────┬───────┘
          │
┌─────────▼───────┐
│   Interfaces    │ (internal/interface/)
└─────────┬───────┘
          │
┌─────────▼───────┐
│   Domain        │ (internal/domain/)
└─────────────────┘
```

### **Implementation Flow**
```
┌─────────────────┐
│   Infrastructure│ (infra/)
│   Implementations│
└─────────┬───────┘
          │ implements
┌─────────▼───────┐
│   Interfaces    │ (internal/interface/)
└─────────────────┘
```

---

## 🎯 Quy Tắc Tổ Chức

### **1. Naming Conventions**
- **Domain**: `AuthMethodDomain`, `UserDomain`
- **Repository**: `AuthMethodRepository`, `UserRepository`
- **Usecase**: `AuthUsecase`, `UserUsecase`
- **Handler**: `AuthHandler`, `UserHandler`
- **DTO**: `LoginRequest`, `LoginResponse`

### **2. File Organization**
- **Một file = Một struct/interface**
- **Package name = Directory name**
- **Import paths rõ ràng**

### **3. Dependency Injection**
- **Sử dụng Wire cho DI**
- **Constructor injection**
- **Interface-based dependencies**

### **4. Error Handling**
- **Sử dụng context.Context**
- **Return errors, không panic**
- **Custom error types**

---

## 🚀 Best Practices

### **1. Layer Separation**
- **Domain không phụ thuộc vào bất kỳ layer nào**
- **Usecase chỉ phụ thuộc vào Domain và Interfaces**
- **Infrastructure implement Interfaces**

### **2. Testing**
- **Mock interfaces cho unit tests**
- **Integration tests cho infrastructure**
- **Test business logic trong usecase**

### **3. Configuration**
- **External configuration files**
- **Environment-based configs**
- **Validation cho configs**

### **4. Logging & Monitoring**
- **Structured logging**
- **Health checks**
- **Metrics collection**

---

## 🎯 Tóm Tắt

| Layer | Mục đích | Phụ thuộc vào |
|-------|----------|---------------|
| **Domain** | Business entities | Không phụ thuộc gì |
| **Interface** | Contract definitions | Domain |
| **Usecase** | Business logic | Domain + Interface |
| **Infrastructure** | Technical implementations | Interface |
| **Environment** | Runtime configuration | Config |
| **Initialization** | Object initialization | Tất cả layers |
| **Entry Points** | Service startup | Tất cả layers |

**Nguyên tắc vàng**: **Dependencies chỉ được point inward** - từ ngoài vào trong, không bao giờ ngược lại!

---

## 📋 Các Thư Mục Quan Trọng

### **1. Thư Mục Tự Động Generate**
- **`wire/wire.go`**: File này sẽ tự động sinh ra khi gõ script. Không cần tạo tay
- **`wire/wire_gen.go`**: File được generate tự động bởi Wire

### **2. Thư Mục Cấu Hình**
- **`config/`**: Configuration files (app.yml, properties.go, runtime.go)
- **`env/`**: Environment-specific runtime settings

### **3. Thư Mục Khởi Tạo**
- **`initial/`**: Chứa các đối tượng khởi tạo project. Thường bao gồm các handler hoặc ghi log hoặc khởi tạo db

### **4. Thư Mục Validation**
- **`validator/`**: Request validation cho các API endpoints
