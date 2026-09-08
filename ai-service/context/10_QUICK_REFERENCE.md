# Quick Reference Guide

## ⚡ Fast Lookup for Common Tasks

**Purpose**: Quick answers to common development questions
**Format**: Question → Answer
**Status**: ✅ Complete

---

## 🚀 Getting Started

### Q: How do I start all services?
```bash
cd shared/code
make start-all
```

### Q: How do I start a specific service?
```bash
make <service>-grpc
# Examples:
make user-grpc
make auth-grpc
make bdspro-grpc
```

### Q: How do I stop all services?
```bash
make stop-all
# Or manually:
ps aux | grep air | awk '{print $2}' | xargs kill
```

---

## 📝 Adding New Features

### Q: How do I add a new API endpoint?
**Step 1**: Define in protobuf
```protobuf
// shared/protobuf/schema/user/profile.proto
rpc GetStats(GetStatsRequest) returns (GetStatsResponse) {
    option (google.api.http) = {
        get: "/v2/user/stats"
    };
}
```

**Step 2**: Generate protobuf
```bash
make buf-user
```

**Step 3**: Implement handler
```go
// user-service/infra/handler/profile_handler.go
// @Summary Get user statistics
// @Router /v2/user/stats [get]
func (h *Handler) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
    // Implementation
}
```

**Step 4**: Generate wire
```bash
make wire user
```

### Q: How do I add a new database table?
**Step 1**: Define entity
```go
// internal/domain/example.go
type Example struct {
    _models.BaseEntity
    Name        string `gorm:"size:255"`
    Description string `gorm:"type:text"`
}
```

**Step 2**: Create repository interface
```go
// internal/interface/repo/example_repo.go
type IExampleRepo interface {
    _crud.ICrudRepo[Example]
}
```

**Step 3**: Implement repository
```go
// infra/postgre/example_postgres.go
// @bind: internal/interface/repo
type ExampleRepo struct {
    _provider.CrudRepo[Example]
}

func NewExampleRepo(db *_db.TransactionRepo) repo.IExampleRepo {
    repo := &ExampleRepo{}
    repo.Init(repo, db)
    return repo
}
```

**Step 4**: Create migration (optional)
```sql
-- migrate/create_example.sql
CREATE TABLE example (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## 🔧 Common Code Patterns

### Q: How do I get current user ID?
```go
import _utils "common/utils"

profileId := _utils.GetProfileIdWithContext(ctx)
organizationId := _utils.GetOrganizationIdFromContext(ctx)
```

### Q: How do I implement pagination?
```go
// DTO
type GetListRequest struct {
    _dto.Pagable
    Text string `json:"text" form:"text"`
}

// Repository
func (r *Repo) GetList(ctx context.Context, req *GetListRequest) ([]Entity, int64, error) {
    var entities []Entity
    var total int64
    
    db := r.db.GetDB()
    query := db.Model(&Entity{}).Where("deleted_at IS NULL")
    
    if req.Text != "" {
        query = query.Where("name LIKE ?", "%"+req.Text+"%")
    }
    
    // Count total
    query.Count(&total)
    
    // Get page
    err := query.
        Offset(req.GetOffset()).
        Limit(req.GetLimit()).
        Order("created_at DESC").
        Find(&entities).Error
    
    return entities, total, err
}
```

### Q: How do I return errors properly?
```go
// In usecase
func (u *Usecase) DoSomething(ctx context.Context) (*Result, *_err.ErrorDTO) {
    result, err := u.repo.DoSomething(ctx)
    if err != nil {
        return nil, &_err.ErrorDTO{
            Code:    500,
            Message: "Không thể thực hiện",
        }
    }
    return result, nil
}

// In handler
result, errDTO := u.usecase.DoSomething(ctx)
if errDTO != nil {
    return nil, status.Error(codes.Internal, errDTO.Message)
}
return &pb.Response{Data: result}, nil
```

### Q: How do I call another service?
```go
// Define interface
type INotificationProvider interface {
    Send(ctx context.Context, userId uint64, msg string) error
}

// Implement client
// @bind: internal/interface/provider
type NotificationClient struct {
    conn *grpc.ClientConn
}

func (c *NotificationClient) Send(ctx context.Context, userId uint64, msg string) error {
    client := notificationpb.NewNotificationServiceClient(c.conn)
    _, err := client.SendNotification(ctx, &notificationpb.SendRequest{
        UserId:  userId,
        Message: msg,
    })
    return err
}

// Use in usecase
type Usecase struct {
    notiClient provider.INotificationProvider
}

func (u *Usecase) Create(ctx context.Context, data *Data) error {
    // Create something
    
    // Send notification
    u.notiClient.Send(ctx, userId, "Created successfully")
    
    return nil
}
```

---

## 📊 Database Queries

### Q: How do I query with GORM?
```go
// Find by ID
var product Product
db.First(&product, id)

// Find with condition
db.Where("status = ?", 1).First(&product)

// Find multiple
var products []Product
db.Where("price > ?", 1000).Find(&products)

// Count
var count int64
db.Model(&Product{}).Where("status = ?", 1).Count(&count)

// With relations
db.Preload("Media").Preload("Owner").First(&product, id)

// Join
db.Joins("LEFT JOIN product_media ON products.id = product_media.product_id").
   Where("products.id = ?", id).
   Find(&product)

// Pagination
db.Where("deleted_at IS NULL").
   Offset(offset).
   Limit(limit).
   Order("created_at DESC").
   Find(&products)
```

### Q: How do I use transactions?
```go
// Using TransactionRepo
return r.db.WithTx(ctx, func(tx *gorm.DB) error {
    // Create record
    if err := tx.Create(&entity).Error; err != nil {
        return err
    }
    
    // Update related record
    if err := tx.Model(&related).Update("count", gorm.Expr("count + 1")).Error; err != nil {
        return err
    }
    
    return nil
})
```

### Q: How do soft deletes work?
```go
// Soft delete (sets deleted_at)
db.Delete(&product, id)

// Query excludes soft deleted automatically
db.Where("status = ?", 1).Find(&products)  // Only non-deleted

// Include soft deleted
db.Unscoped().Where("status = ?", 1).Find(&products)

// Permanent delete
db.Unscoped().Delete(&product, id)
```

---

## 🔐 Authentication & Authorization

### Q: How do I protect an API endpoint?
```go
// In protobuf - no special annotation needed
// Gateway handles JWT validation

// In handler - get user from context
profileId := _utils.GetProfileIdWithContext(ctx)
if profileId == 0 {
    return nil, status.Error(codes.Unauthenticated, "Not authenticated")
}
```

### Q: How do I check permissions?
```go
// Define permission
const (
    PermissionProductView   = "product.view"
    PermissionProductCreate = "product.create"
    PermissionProductUpdate = "product.update"
)

// Check permission in usecase
func (u *Usecase) Create(ctx context.Context, product *Product) error {
    if err := u.checkPermission(ctx, PermissionProductCreate); err != nil {
        return err
    }
    // Create product
}

func (u *Usecase) checkPermission(ctx context.Context, permission string) error {
    roleKey := ctx.Value("role").(string)
    // Check if role has permission
    hasPermission := u.permissionRepo.Has(roleKey, permission)
    if !hasPermission {
        return errors.New("forbidden")
    }
    return nil
}
```

### Q: How do I get JWT token for testing?
```bash
# Login to get token
curl -X POST http://localhost:8080/v2/auth/otp/verify \
  -H "Content-Type: application/json" \
  -d '{"phone": "0901234567", "otp": "123456"}'

# Response includes access token
{
  "accessToken": "eyJhbGc...",
  "refreshToken": "eyJhbGc...",
  "profileId": 123
}

# Use in requests
curl -H "Authorization: Bearer eyJhbGc..." \
     http://localhost:8080/v2/user/profile/me
```

---

## 🧪 Testing

### Q: How do I test an API locally?
```bash
# Using curl
curl -X GET http://localhost:8080/v2/user/profile/me \
  -H "Authorization: Bearer <token>"

# Using grpcurl
grpcurl -plaintext \
  -d '{"id": 1}' \
  localhost:50051 \
  userpb.ProfileService/GetProfileInfo

# Using Postman
# Import: shared/protobuf/docs/<service>/swagger.json
```

### Q: How do I run unit tests?
```bash
# Run tests for a service
cd user-service
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestUserService_Create ./internal/usecase/
```

### Q: How do I write a unit test?
```go
func TestProductUsecase_Create(t *testing.T) {
    // Setup
    mockRepo := NewMockProductRepo()
    usecase := NewProductUsecase(mockRepo)
    
    // Test cases
    tests := []struct {
        name    string
        input   *Product
        wantErr bool
    }{
        {
            name:    "valid product",
            input:   &Product{Name: "Test"},
            wantErr: false,
        },
        {
            name:    "empty name",
            input:   &Product{Name: ""},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := usecase.Create(context.Background(), tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

---

## 🐛 Debugging

### Q: How do I check service logs?
```bash
# View logs for specific service
tail -f user-service/tmp/*.log

# View all services
tail -f */tmp/*.log

# Search logs
grep "error" user-service/tmp/*.log
```

### Q: How do I debug SQL queries?
```go
// Enable debug mode
db.Debug().Where("id = ?", id).First(&product)

// This prints SQL:
// [0.001s] SELECT * FROM products WHERE id = 123
```

### Q: Service won't start, what to check?
```bash
# Check if port is in use
lsof -ti tcp:50051

# Kill process on port
lsof -ti tcp:50051 | xargs kill -9

# Check config file
cat user-service/config/local.yml

# Check database connection
psql -h localhost -U postgres -d user_db
```

---

## 📦 Code Generation

### Q: When do I need to generate protobuf?
- After adding/modifying `.proto` files
- After pulling new changes that include proto changes

```bash
make buf-<service>
# Or all services
make buf-all
```

### Q: When do I need to generate wire?
- After adding new repository
- After adding new usecase
- After adding new handler
- After changing dependency injection

```bash
make wire <service>
```

### Q: How do I generate Swagger docs?
```bash
make swag <service>

# Merge all swagger
make merge-swagger
```

---

## 🔍 Finding Things

### Q: Where is X defined?
```bash
# Find files
find . -name "*product*"

# Search in code
grep -r "ProductService" .

# Search in specific directory
grep -r "GetByID" user-service/internal/
```

### Q: Which service handles X?
| Feature | Service |
|---------|---------|
| Authentication | Auth Service |
| User profiles | User Service |
| Organizations | Organization Service |
| Real estate | BDSPro Service |
| Contacts/Leads | CRM Service |
| Payments | Payment Service |
| Notifications | Notification Service |
| Chat | Chat Service |
| Social posts | Social Service |
| Files | File Service |

### Q: Where are API endpoints defined?
- **Protobuf**: `shared/protobuf/schema/<service>/`
- **Handler**: `<service>/infra/handler/`
- **Routes**: Registered in Gateway Service

---

## 🚢 Deployment

### Q: How do I build for production?
```bash
# Build binary
cd user-service
go build -o user-service main.go

# Build Docker image
docker build -t user-service:latest .

# Run with Docker Compose
docker-compose up -d
```

### Q: How do I check if services are running?
```bash
# Check processes
ps aux | grep "service"

# Check Docker containers
docker ps

# Health checks
curl http://localhost:8080/health
curl http://localhost:50051/health
```

---

## 💡 Best Practices

### Q: Naming conventions?
- **Types**: `PascalCase` (Product, UserService)
- **Variables**: `camelCase` (productId, userName)
- **Constants**: `PascalCase` (MaxRetries, DefaultSize)
- **Functions**: `PascalCase` (GetByID, CreateProduct)
- **Files**: `snake_case.go` (product_repo.go, user_service.go)

### Q: Import conventions?
```go
// Standard library first
import (
    "context"
    "fmt"
    "time"
)

// External packages
import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// Internal packages with aliases
import (
    _dto "common/domain/dto"
    _models "common/domain/entity"
    _enum "common/domain/enum"
    _err "common/domain/err"
    _utils "common/utils"
)
```

### Q: Error handling best practices?
```go
// ✅ GOOD - Always handle errors
result, err := doSomething()
if err != nil {
    return &_err.ErrorDTO{
        Code:    500,
        Message: "Không thể thực hiện",
    }
}

// ❌ BAD - Ignoring errors
result, _ := doSomething()

// ✅ GOOD - Wrap with context
if err != nil {
    return fmt.Errorf("failed to create product: %w", err)
}
```

---

## 📞 Getting Help

### Q: Where do I ask questions?
- **Code questions**: Team lead
- **Architecture**: Solution architect
- **Business logic**: Product manager
- **Database**: DBA
- **DevOps**: DevOps team

### Q: Documentation not clear?
1. Check other context files in `/context/`
2. Look at existing code examples
3. Ask team member
4. Update documentation after you learn!

---

## ⚡ Performance Tips

### Q: How to optimize queries?
```go
// Use indexes
db.Where("status = ?", 1).Find(&products)  // Ensure status is indexed

// Avoid N+1 queries
db.Preload("Media").Preload("Owner").Find(&products)

// Use pagination
db.Limit(20).Offset(0).Find(&products)

// Select specific fields
db.Select("id", "name", "price").Find(&products)
```

### Q: How to optimize API response?
- Use pagination for lists
- Return only needed fields
- Cache frequently accessed data (Redis)
- Use gRPC for service-to-service calls
- Compress large responses

---

**Last Updated**: October 15, 2025
**Version**: v1.0
**Status**: ✅ Complete

---

**Pro Tip**: Bookmark this page for instant answers during development!

