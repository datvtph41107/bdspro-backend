# Development Workflow & Best Practices

## 📋 Complete Development Guide

**Purpose**: Step-by-step guide for developing in BDSPro Microservices
**Target Audience**: Go developers working on the project
**Status**: ✅ Production Ready

---

## 🚀 Quick Start for New Developers

### Prerequisites
```bash
# Install Go 1.21+
brew install go

# Install Air (hot reload)
go install github.com/cosmtrek/air@latest

# Install Buf CLI (protobuf)
brew install bufbuild/buf/buf

# Install Wire (dependency injection)
go install github.com/google/wire/cmd/wire@latest

# Install Swag (swagger)
go install github.com/swaggo/swag/cmd/swag@latest
```

### Clone & Setup
```bash
# Clone repository
git clone <repository-url>
cd go-microservices

# Install dependencies for all services
cd shared/code
make init

# Generate protobuf files
make buf-all
```

---

## 🏗️ Adding a New API Endpoint

### Step 1: Define in Protobuf

**Location**: `shared/protobuf/schema/<service>/<file>.proto`

```protobuf
// shared/protobuf/schema/user/profile.proto
syntax = "proto3";

package userpb;
import "google/api/annotations.proto";

service ProfileService {
    // Add new RPC method
    rpc GetUserStats(GetUserStatsRequest) returns (GetUserStatsResponse) {
        option (google.api.http) = {
            get: "/v2/user/profile/stats"
        };
    };
}

message GetUserStatsRequest {
    // Empty or add fields
}

message GetUserStatsResponse {
    int32 totalPosts = 1;
    int32 totalProducts = 2;
    int32 totalContacts = 3;
}
```

### Step 2: Generate Protobuf Code

```bash
# Generate for specific service
cd shared/code
make buf-user

# Or generate all
make buf-all
```

**Generated Files**:
- `shared/protobuf/types/user/profile.pb.go` - Go structs
- `shared/protobuf/types/user/profile_grpc.pb.go` - gRPC service
- `shared/protobuf/types/user/profile.pb.gw.go` - HTTP gateway

### Step 3: Create DTO (if needed)

**Location**: `<service>/internal/dto/`

```go
// user-service/internal/dto/stats_dto.go
package dto

type UserStatsDTO struct {
    TotalPosts    int32 `json:"totalPosts"`
    TotalProducts int32 `json:"totalProducts"`
    TotalContacts int32 `json:"totalContacts"`
}
```

### Step 4: Create Repository Method (if needed)

**Location**: `<service>/internal/interface/repo/`

```go
// user-service/internal/interface/repo/profile_repo.go
package repo

import "context"

type IProfileRepo interface {
    _crud.ICrudRepo[Profile]
    
    // Add new method
    GetStats(ctx context.Context, profileId uint64) (*dto.UserStatsDTO, error)
}
```

### Step 5: Implement Repository

**Location**: `<service>/infra/postgre/`

```go
// user-service/infra/postgre/profile_postgres.go
package postgre

import (
    _provider "common/provider"
    "context"
)

// @bind: internal/interface/repo
type ProfileRepo struct {
    _provider.CrudRepo[Profile]
}

func NewProfileRepo(db *_db.TransactionRepo) repo.IProfileRepo {
    repo := &ProfileRepo{}
    repo.Init(repo, db)
    return repo
}

func (r *ProfileRepo) GetStats(ctx context.Context, profileId uint64) (*dto.UserStatsDTO, error) {
    stats := &dto.UserStatsDTO{}
    
    // Count posts
    r.db.GetDB().Model(&Post{}).Where("owner_id = ?", profileId).Count(&stats.TotalPosts)
    
    // Count products
    r.db.GetDB().Model(&Product{}).Where("owner_id = ?", profileId).Count(&stats.TotalProducts)
    
    // Count contacts
    r.db.GetDB().Model(&Contact{}).Where("owner_id = ?", profileId).Count(&stats.TotalContacts)
    
    return stats, nil
}
```

### Step 6: Create Usecase

**Location**: `<service>/internal/usecase/`

```go
// user-service/internal/usecase/profile_usecase.go
package usecase

import (
    _crud "common/domain/crud"
    _err "common/domain/err"
    "context"
)

type ProfileUsecase struct {
    _crud.BaseUsecase[Profile, repo.IProfileRepo]
    repo repo.IProfileRepo
}

func NewProfileUsecase(repo repo.IProfileRepo) *ProfileUsecase {
    return &ProfileUsecase{
        BaseUsecase: _crud.BaseUsecase[Profile, repo.IProfileRepo]{
            Repo: repo,
        },
        repo: repo,
    }
}

func (u *ProfileUsecase) GetStats(ctx context.Context, profileId uint64) (*dto.UserStatsDTO, *_err.ErrorDTO) {
    stats, err := u.repo.GetStats(ctx, profileId)
    if err != nil {
        return nil, &_err.ErrorDTO{
            Code:    500,
            Message: "Không thể lấy thống kê",
        }
    }
    
    return stats, nil
}
```

### Step 7: Create Handler

**Location**: `<service>/infra/handler/`

```go
// user-service/infra/handler/profile_handler.go
package handler

import (
    "context"
    userpb "pb/types/user"
)

// @bind: internal/usecase
type ProfileHandler struct {
    userpb.UnimplementedProfileServiceServer
    usecase *usecase.ProfileUsecase
}

func NewProfileHandler(usecase *usecase.ProfileUsecase) *ProfileHandler {
    return &ProfileHandler{
        usecase: usecase,
    }
}

// @Summary     Get user statistics
// @Description Get statistics about user's posts, products, and contacts
// @Tags        User
// @Accept      json
// @Produce     json
// @Success     200 {object} userpb.GetUserStatsResponse
// @Failure     500 {object} common.ErrorResponse
// @Router      /v2/user/profile/stats [get]
// @Security    BearerAuth
func (h *ProfileHandler) GetUserStats(ctx context.Context, req *userpb.GetUserStatsRequest) (*userpb.GetUserStatsResponse, error) {
    profileId := _utils.GetProfileIdWithContext(ctx)
    
    stats, errDTO := h.usecase.GetStats(ctx, profileId)
    if errDTO != nil {
        return nil, status.Error(codes.Internal, errDTO.Message)
    }
    
    return &userpb.GetUserStatsResponse{
        TotalPosts:    stats.TotalPosts,
        TotalProducts: stats.TotalProducts,
        TotalContacts: stats.TotalContacts,
    }, nil
}
```

### Step 8: Generate Wire

```bash
# In shared/code directory
make wire user

# This will generate user-service/wire/wire_gen.go
```

### Step 9: Test the API

```bash
# Start the service
make user-grpc

# Test with curl
curl -H "Authorization: Bearer <token>" \
     http://localhost:8080/v2/user/profile/stats
```

---

## 🔧 Common Development Tasks

### Task 1: Add New Field to Existing Entity

```go
// 1. Update domain entity
type Product struct {
    _models.BaseEntity
    Name        string
    Description string
    NewField    string  // Add new field
}

// 2. Update protobuf
message Product {
    uint64 id = 1;
    string name = 2;
    string description = 3;
    string newField = 4;  // Add new field
}

// 3. Generate protobuf
make buf-bdspro

// 4. Create migration (if needed)
// Add migration SQL to migrate.go or create new migration file

// 5. Update mapper (if needed)
func ToProductDTO(product *Product) *ProductDTO {
    return &ProductDTO{
        ID:          product.ID,
        Name:        product.Name,
        Description: product.Description,
        NewField:    product.NewField,  // Add mapping
    }
}
```

### Task 2: Call Another Service

```go
// 1. Define interface in internal/interface/
type INotificationProvider interface {
    SendNotification(ctx context.Context, userId uint64, message string) error
}

// 2. Implement client in infra/client/
// @bind: internal/interface
type NotificationClient struct {
    conn *grpc.ClientConn
}

func NewNotificationClient(conn *grpc.ClientConn) provider.INotificationProvider {
    return &NotificationClient{conn: conn}
}

func (c *NotificationClient) SendNotification(ctx context.Context, userId uint64, message string) error {
    client := notificationpb.NewNotificationServiceClient(c.conn)
    
    req := &notificationpb.SendNotificationRequest{
        UserId:  userId,
        Message: message,
    }
    
    _, err := client.SendNotification(ctx, req)
    return err
}

// 3. Inject in usecase
type ProductUsecase struct {
    repo           repo.IProductRepo
    notificationClient provider.INotificationProvider
}

// 4. Use in usecase
func (u *ProductUsecase) Create(ctx context.Context, product *Product) error {
    // Create product
    if err := u.repo.Create(ctx, product); err != nil {
        return err
    }
    
    // Send notification
    u.notificationClient.SendNotification(ctx, product.OwnerId, "Product created")
    
    return nil
}

// 5. Generate wire
make wire bdspro
```

### Task 3: Add Pagination to List API

```go
// 1. Create request DTO
type GetProductsRequest struct {
    _dto.Pagable
    Text   string `json:"text" form:"text"`
    Status int    `json:"status" form:"status"`
}

// 2. Update repository
func (r *ProductRepo) GetList(ctx context.Context, req *GetProductsRequest) ([]Product, int64, error) {
    db := r.db.GetDB()
    
    query := db.Model(&Product{}).Where("deleted_at IS NULL")
    
    if req.Text != "" {
        query = query.Where("name LIKE ?", "%"+req.Text+"%")
    }
    
    if req.Status > 0 {
        query = query.Where("status = ?", req.Status)
    }
    
    // Count total
    var total int64
    query.Count(&total)
    
    // Get page
    var products []Product
    err := query.
        Offset(req.GetOffset()).
        Limit(req.GetLimit()).
        Order("created_at DESC").
        Find(&products).Error
    
    return products, total, err
}

// 3. Update usecase
func (u *ProductUsecase) GetList(ctx context.Context, req *GetProductsRequest) ([]Product, int64, *_err.ErrorDTO) {
    products, total, err := u.repo.GetList(ctx, req)
    if err != nil {
        return nil, 0, &_err.ErrorDTO{
            Code:    500,
            Message: "Không thể lấy danh sách sản phẩm",
        }
    }
    
    return products, total, nil
}

// 4. Update handler to return pagination info
func (h *ProductHandler) GetProducts(ctx context.Context, req *productpb.GetProductsRequest) (*productpb.GetProductsResponse, error) {
    searchReq := &dto.GetProductsRequest{
        Pagable: _dto.Pagable{
            Page: int(req.Page),
            Size: int(req.Size),
        },
        Text:   req.Text,
        Status: int(req.Status),
    }
    
    products, total, errDTO := h.usecase.GetList(ctx, searchReq)
    if errDTO != nil {
        return nil, status.Error(codes.Internal, errDTO.Message)
    }
    
    return &productpb.GetProductsResponse{
        Data:  toProductDTOs(products),
        Total: int32(total),
        Page:  req.Page,
        Size:  req.Size,
    }, nil
}
```

---

## 🎯 Best Practices

### 1. Naming Conventions

```go
// ✅ GOOD
type Product struct { }          // PascalCase for types
var totalPrice float64           // camelCase for variables
const MaxRetries = 3             // PascalCase for constants
func GetProductById() { }        // PascalCase for functions

// ❌ BAD
type product struct { }          // lowercase
var TotalPrice float64           // PascalCase for variables
const max_retries = 3            // snake_case
func get_product_by_id() { }     // snake_case
```

### 2. Error Handling

```go
// ✅ GOOD - Always handle errors
product, err := repo.GetByID(ctx, id)
if err != nil {
    return nil, &_err.ErrorDTO{
        Code:    404,
        Message: "Không tìm thấy sản phẩm",
    }
}

// ❌ BAD - Ignoring errors
product, _ := repo.GetByID(ctx, id)
```

### 3. Context Usage

```go
// ✅ GOOD - Always pass context
func (u *Usecase) GetProduct(ctx context.Context, id uint64) (*Product, error) {
    profileId := _utils.GetProfileIdWithContext(ctx)
    return u.repo.GetByID(ctx, id)
}

// ❌ BAD - No context
func (u *Usecase) GetProduct(id uint64) (*Product, error) {
    return u.repo.GetByID(id)
}
```

### 4. Using Shared Components

```go
// ✅ GOOD - Use BaseEntity
type Product struct {
    _models.BaseEntity
    Name string
}

// ✅ GOOD - Use CrudRepo
type ProductRepo struct {
    _provider.CrudRepo[Product]
}

// ✅ GOOD - Embed Pagable
type SearchRequest struct {
    _dto.Pagable
    Text string
}

// ❌ BAD - Reinvent the wheel
type Product struct {
    ID        uint64
    CreatedAt time.Time
    UpdatedAt time.Time
    // ... duplicate code
}
```

### 5. Clean Architecture Layers

```go
// ✅ GOOD - Proper layer separation

// Domain (internal/domain)
type Product struct {
    ID   uint64
    Name string
}

// Interface (internal/interface/repo)
type IProductRepo interface {
    GetByID(ctx context.Context, id uint64) (*Product, error)
}

// Usecase (internal/usecase)
type ProductUsecase struct {
    repo IProductRepo
}

// Infrastructure (infra/postgre)
type ProductRepo struct {
    db *gorm.DB
}

// ❌ BAD - Mixed concerns
type Product struct {
    ID   uint64
    Name string
    db   *gorm.DB  // Infrastructure in domain!
}
```

---

## 🔍 Debugging

### Check Service Logs
```bash
# View service logs
tail -f user-service/tmp/*.log

# Check all services
tail -f */tmp/*.log
```

### Debug Database Queries
```go
// Enable GORM debug mode
db.Debug().Where("id = ?", id).First(&product)

// This will print SQL:
// SELECT * FROM products WHERE id = 123
```

### Test gRPC Directly
```bash
# Install grpcurl
brew install grpcurl

# List services
grpcurl -plaintext localhost:50051 list

# Call method
grpcurl -plaintext -d '{"id": 1}' localhost:50051 userpb.ProfileService/GetProfileInfo
```

---

## 🚀 Deployment Workflow

### Development
```bash
# Start all services with hot reload
make start-all

# Or start individual service
make user-grpc
```

### Testing
```bash
# Run tests for a service
cd user-service
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Build
```bash
# Build for current OS
go build -o bin/user-service main.go

# Build for Linux (for Docker)
GOOS=linux GOARCH=amd64 go build -o user-service main.go
```

### Docker
```bash
# Build Docker image
docker build -t user-service .

# Run container
docker run -p 50051:50051 user-service

# Docker Compose
docker-compose up -d
```

---

## 📊 Code Review Checklist

- [ ] Protobuf definitions updated and generated
- [ ] DTO created if needed
- [ ] Repository interface defined
- [ ] Repository implementation with @bind comment
- [ ] Usecase created with business logic
- [ ] Handler created with Swagger comments
- [ ] Wire generated successfully
- [ ] Error handling implemented
- [ ] Context passed correctly
- [ ] Permissions checked if needed
- [ ] Tests written (if applicable)
- [ ] Code formatted (`go fmt`)
- [ ] Linter passed (`go vet`)

---

## 🎯 Quick Commands Reference

```bash
# Generate protobuf
make buf-<service>

# Generate wire
make wire <service>

# Generate swagger
make swag <service>

# Start service with hot reload
make <service>-grpc

# Run tests
cd <service> && go test ./...

# Format code
go fmt ./...

# Check for issues
go vet ./...

# Build
go build -o bin/app main.go
```

---

**Last Updated**: October 15, 2025
**Guide Version**: v1.0
**Status**: ✅ Complete

