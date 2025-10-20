# Common Code Patterns - BDSPro Microservices

Tài liệu này tổng hợp các pattern code phổ biến trong BDSPro microservices.

---

## 1. Timestamp Mapping Pattern

**Câu hỏi**: Làm sao map timestamp từ protobuf request sang domain entity?

**Trả lời**: Timestamp trong protobuf request và response đều là string. Để convert:

### Từ string (protobuf) sang time.Time (domain):
```go
import _utils "common/utils"

// Trong mapper
func (m *YourMapper) PbToDomain(pb *pb.YourRequest) *domain.YourEntity {
    return &domain.YourEntity{
        StartDate: _utils.ParseStringToTime(pb.StartDate),
        EndDate: _utils.ParseStringToTime(pb.EndDate),
    }
}
```

### Từ time.Time (domain) sang string (protobuf):
```go
func (m *YourMapper) DomainToPb(entity *domain.YourEntity) *pb.YourResponse {
    return &pb.YourResponse{
        StartDate: _utils.FormatTimeToString(entity.StartDate),
        EndDate: _utils.FormatTimeToString(entity.EndDate),
        CreatedAt: _utils.FormatTimeToString(entity.CreatedAt),
    }
}
```

### Chi tiết hàm utils:
- `_utils.ParseStringToTime(string)` → `*time.Time`: Parse RFC3339 string thành time pointer
- `_utils.FormatTimeToString(*time.Time)` → `string`: Format time pointer thành RFC3339 string
- Cả 2 hàm đều handle nil safety

**Lưu ý**: LUÔN dùng `_utils` alias khi import `common/utils`

---

## 2. Context Management Pattern

**Câu hỏi**: Cách lấy profileId và organizationId từ context trong usecase?

**Trả lời**: Để lấy thông tin user từ context, LUÔN dùng các hàm trong `common/utils`:

### Lấy profileId (ID của user đang login):
```go
import _utils "common/utils"

func (uc *YourUsecase) DoSomething(ctx context.Context) error {
    profileID := _utils.GetProfileIdWithContext(ctx)
    // profileID là uint64
    // Sử dụng profileID...
}
```

### Lấy organizationId (ID của tổ chức hiện tại):
```go
import _utils "common/utils"

func (uc *YourUsecase) DoSomething(ctx context.Context) error {
    organizationID := _utils.GetOrganizationIdFromContext(ctx)
    // organizationID là uint64
    // Sử dụng organizationID...
}
```

### Quy tắc quan trọng:
- KHÔNG được tự parse thông tin từ JWT hay context
- LUÔN import utils với alias `_utils`
- Chỉ gọi các hàm này ở layer usecase, handler, hoặc service cần thông tin user
- Context được set từ middleware gateway, tự động có trong ctx parameter

### Ví dụ thực tế:
```go
func (uc *ProductUsecase) CreateProduct(ctx context.Context, dto *ProductDTO) error {
    profileID := _utils.GetProfileIdWithContext(ctx)
    orgID := _utils.GetOrganizationIdFromContext(ctx)
    
    product := &domain.Product{
        Name: dto.Name,
        CreatedBy: &profileID,
        OrganizationID: orgID,
    }
    return uc.repo.Create(ctx, product)
}
```

---

## 3. BaseEntity Pattern

**Câu hỏi**: Cách sử dụng BaseEntity khi định nghĩa domain entity?

**Trả lời**: Khi tạo domain entity mới, LUÔN embed `BaseEntity` từ `common/models`:

### Cấu trúc chuẩn:
```go
package domain

import (
    _models "common/models"
    "time"
)

type YourEntity struct {
    _models.BaseEntity  // LUÔN embed đầu tiên
    Name        string
    Description string
    Status      uint
    // ... các field khác
}

func (e *YourEntity) TableName() string {
    return "your_table_name"
}
```

### BaseEntity cung cấp:
- `ID uint64`: Primary key
- `CreatedAt time.Time`: Timestamp tạo (tự động)
- `UpdatedAt time.Time`: Timestamp cập nhật (tự động)
- `DeletedAt *time.Time`: Soft delete (GORM)
- `CreatedBy *uint64`: User tạo (tự động từ context)
- `UpdatedBy *uint64`: User cập nhật (tự động từ context)

### Lợi ích:
- GORM tự động handle CreatedAt/UpdatedAt
- Soft delete tự động (query sẽ filter deleted_at IS NULL)
- Audit trail tự động qua AuditBase callbacks
- Giảm code duplication

**Lưu ý**:
- Import với alias `_models "common/models"`
- BaseEntity phải là field đầu tiên khi embed
- LUÔN định nghĩa TableName() method

---

## 4. Repository Pattern với CrudRepo

**Câu hỏi**: Cách viết repository implementation trong infra/postgre?

**Trả lời**: Repository trong infra/postgre phải:
1. Embed `CrudRepo` từ `common/provider`
2. Gọi `Init()` trong constructor
3. Có comment `@bind` để wire generate

### Template chuẩn:
```go
package postgre

import (
    "context"
    _db "common/db"
    _provider "common/provider"
    "your-service/internal/domain"
    _repo "your-service/internal/interface/repo"
)

// @bind: your-service/internal/interface/repo.IYourEntityRepo
type YourEntityRepo struct {
    _provider.CrudRepo[domain.YourEntity]
}

func NewYourEntityRepo(db *_db.TransactionRepo) _repo.IYourEntityRepo {
    repo := &YourEntityRepo{}
    repo.Init(repo, db)  // BẮT BUỘC gọi Init
    return repo
}

// Override BeforeSave nếu cần custom validation
func (r *YourEntityRepo) BeforeSave(c context.Context, id *uint64, entity *domain.YourEntity) error {
    // Custom validation logic
    if entity.Name == "" {
        return errors.New("name is required")
    }
    return nil
}

// Custom query methods
func (r *YourEntityRepo) FindByName(c context.Context, name string) (*domain.YourEntity, error) {
    var entity domain.YourEntity
    err := r.DB.WithContext(c).
        Where("name = ? AND deleted_at IS NULL", name).
        First(&entity).Error
    return &entity, err
}
```

### CrudRepo tự động cung cấp:
- `Create`, `Update`, `Delete`
- `GetByID`, `GetDetail`, `GetAll`, `GetList`
- Transaction support tự động
- Soft delete filter
- Hook callbacks (BeforeSave/AfterSave)

---

## 5. Import Conventions

**Câu hỏi**: Quy tắc import và alias trong Go microservices?

**Trả lời**: Project có quy ước import alias chuẩn:

### Common libraries (LUÔN dùng underscore prefix):
```go
import (
    _utils "common/utils"           // Utils
    _models "common/models"         // Base entities
    _dto "common/domain/dto"        // Common DTOs
    _enum "common/domain/enum"      // Enums
    _err "common/domain/err"        // Error handling
    _crud "common/domain/crud"      // CRUD interfaces
    _provider "common/provider"     // Base repository
    _db "common/db"                 // Database
)
```

### Internal packages:
```go
import (
    "your-service/internal/domain"
    "your-service/internal/dto"
    _repo "your-service/internal/interface/repo"
    _handler "your-service/internal/interface/handler"
)
```

### Protobuf packages:
```go
import (
    pb_bdspro "pb/types/bdspro"
    pb_auth "pb/types/auth"
    pb_user "pb/types/user"
    shared_enum "pb/enums"
)
```

### Lý do dùng alias:
- Tránh conflict giữa các package cùng tên
- Dễ phân biệt common library vs service-specific code
- Consistent naming convention
- Code dễ đọc và maintain

---

## 6. Mapper Pattern

**Câu hỏi**: Cách viết mapper để convert giữa domain, DTO và protobuf?

**Trả lời**: Mapper nằm trong `infra/mapper` và handle conversion giữa các layers.

### Cấu trúc mapper chuẩn:
```go
package mapper

import (
    _utils "common/utils"
    "your-service/internal/domain"
    "your-service/internal/dto"
    pb "pb/types/your-service"
)

type YourEntityMapper struct {
    // Inject mappers khác nếu cần
}

func NewYourEntityMapper() *YourEntityMapper {
    return &YourEntityMapper{}
}

// 1. Domain → Protobuf (cho response)
func (m *YourEntityMapper) DomainToPb(entity *domain.YourEntity) *pb.YourEntityDTO {
    if entity == nil {
        return nil
    }
    
    return &pb.YourEntityDTO{
        Id: entity.ID,
        Name: entity.Name,
        Amount: entity.Amount,
        Status: uint32(entity.Status),
        // Timestamp: time.Time → string
        StartDate: _utils.FormatTimeToString(entity.StartDate),
        CreatedAt: _utils.FormatTimeToString(entity.CreatedAt),
    }
}

// 2. Protobuf → Domain (cho request)
func (m *YourEntityMapper) PbToDomain(pb *pb.YourEntityRequest) *domain.YourEntity {
    return &domain.YourEntity{
        Name: pb.Name,
        Amount: pb.Amount,
        Status: domain.StatusEnum(pb.Status),
        // Timestamp: string → time.Time
        StartDate: _utils.ParseStringToTime(pb.StartDate),
    }
}

// 3. Batch conversion
func (m *YourEntityMapper) DomainListToPb(entities []domain.YourEntity) []*pb.YourEntityDTO {
    result := make([]*pb.YourEntityDTO, len(entities))
    for i, entity := range entities {
        result[i] = m.DomainToPb(&entity)
    }
    return result
}
```

### Quy tắc mapper:
1. **Timestamp conversion**: 
   - Domain → PB: `_utils.FormatTimeToString()`
   - PB → Domain: `_utils.ParseStringToTime()`

2. **Enum conversion**:
   - Domain → PB: Cast to uint32: `uint32(entity.Status)`
   - PB → Domain: Cast to domain enum: `domain.StatusEnum(pb.Status)`

3. **Nil safety**: LUÔN check nil cho pointers

4. **NO business logic**: Mapper CHỈ convert data, KHÔNG có logic

---

## 7. Context.Context vs gin.Context

**Câu hỏi**: Khi nào dùng context.Context thay vì gin.Context?

**Trả lời**: Trong Clean Architecture, LUÔN dùng `context.Context` thay vì `gin.Context`.

### KHÔNG dùng gin.Context trong:
- Usecase layer
- Repository layer
- Domain layer
- Mapper layer
- Any internal business logic

### Handler layer: Dùng context.Context
```go
// ✅ ĐÚNG
func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
    // ctx là context.Context
    profileID := _utils.GetProfileIdWithContext(ctx)
    // ...
}
```

### Usecase layer: Dùng context.Context
```go
// ✅ ĐÚNG
func (uc *ProductUsecase) Create(ctx context.Context, dto *dto.ProductDTO) (*domain.Product, error) {
    profileID := _utils.GetProfileIdWithContext(ctx)
    organizationID := _utils.GetOrganizationIdFromContext(ctx)
    // ...
}
```

### Lý do:
1. **Framework agnostic**: Code không phụ thuộc vào Gin framework
2. **Testability**: Dễ test với mock context
3. **Reusability**: Code có thể reuse cho gRPC, CLI, worker...
4. **Clean Architecture**: Tách biệt delivery mechanism khỏi business logic
5. **Standard library**: `context.Context` là Go standard, portable

---

## 8. Soft Delete Pattern

**Câu hỏi**: Làm sao handle soft delete và query filter trong repository?

**Trả lời**: Soft delete được handle tự động qua `BaseEntity` và `CrudRepo`.

### Soft Delete tự động:

Khi entity embed `BaseEntity`:
```go
type YourEntity struct {
    _models.BaseEntity  // Có DeletedAt *time.Time
    Name string
}
```

GORM tự động:
- Filter `deleted_at IS NULL` trong queries
- Set `deleted_at = NOW()` khi delete
- KHÔNG xóa record khỏi database

### Custom query với soft delete filter:
```go
func (r *YourRepo) FindByName(c context.Context, name string) (*domain.YourEntity, error) {
    var entity domain.YourEntity
    err := r.DB.WithContext(c).
        Where("name = ? AND deleted_at IS NULL", name).  // LUÔN check deleted_at
        First(&entity).Error
    return &entity, err
}

func (r *YourRepo) GetList(c context.Context, page, size int) ([]domain.YourEntity, int64, error) {
    var entities []domain.YourEntity
    var total int64
    
    query := r.DB.Model(&domain.YourEntity{}).
        Where("deleted_at IS NULL")  // Filter deleted
    
    query.Count(&total)
    err := query.
        Offset(page * size).
        Limit(size).
        Find(&entities).Error
    
    return entities, total, err
}
```

### Quy tắc:
- LUÔN dùng `BaseEntity` để có soft delete
- LUÔN check `deleted_at IS NULL` trong custom queries
- CrudRepo methods tự động handle soft delete
- Dùng `Unscoped()` nếu cần query deleted records

---

## 9. @bind Comment for Dependency Injection

**Câu hỏi**: Cách viết comment @bind để wire generate dependency injection?

**Trả lời**: Comment `@bind` dùng để wire tool tự động generate dependency injection code.

### Cú pháp:
```go
// @bind: <interface_path>.<InterfaceName>
```

### Ví dụ trong handler:
```go
package handler

import (
    "context"
    pb "pb/types/bdspro"
    "bdspro/internal/usecase"
)

// @bind: bdspro/internal/interface/handler.IProductHandler
type ProductHandler struct {
    productUsecase usecase.IProductUsecase
}

func NewProductHandler(productUsecase usecase.IProductUsecase) *ProductHandler {
    return &ProductHandler{
        productUsecase: productUsecase,
    }
}
```

### Quy tắc:
1. Comment `@bind` phải ở NGAY TRƯỚC struct definition
2. Path phải chính xác từ service root
3. Tên interface phải match với interface được implement
4. Constructor function (New...) phải return interface type
5. Sau khi thêm @bind, chạy `wire gen` để generate code

---

---

## 10. Handler Implementation với Protobuf

**Câu hỏi**: Cách implement handler để expose protobuf API?

**Trả lời**: Handler trong `infra/handler` triển khai các method từ protobuf service definition.

### Template handler chuẩn:
```go
package handler

import (
    "context"
    _utils "common/utils"
    pb "pb/types/your-service"
    _dto "common/domain/dto"
    "your-service/internal/usecase"
    "your-service/infra/mapper"
)

// @bind: your-service/internal/interface/handler.IYourEntityHandler
type YourEntityHandler struct {
    pb.UnimplementedYourServiceServer  // BẮT BUỘC embed để implement gRPC service
    usecase usecase.IYourEntityUsecase
    mapper  *mapper.YourEntityMapper
}

func NewYourEntityHandler(
    usecase usecase.IYourEntityUsecase,
    mapper *mapper.YourEntityMapper,
) *YourEntityHandler {
    return &YourEntityHandler{
        usecase: usecase,
        mapper:  mapper,
    }
}

// @Summary Tạo mới entity
// @Description Tạo mới entity trong hệ thống
// @Tags YourService
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param req body pb.CreateYourEntityRequest true "CreateYourEntityRequest"
// @Success 200 {object} pb.YourEntityResponse "Entity created successfully"
// @Failure 400 {object} error "Bad Request"
// @Failure 500 {object} error "Internal Server Error"
// @Router /your-entity [post]
func (h *YourEntityHandler) Create(
    ctx context.Context, 
    req *pb.CreateYourEntityRequest,
) (*pb.YourEntityResponse, error) {
    // 1. Lấy thông tin user từ context
    profileID := _utils.GetProfileIdWithContext(ctx)
    organizationID := _utils.GetOrganizationIdFromContext(ctx)
    
    // 2. Convert protobuf request sang domain entity
    entity := h.mapper.PbToDomain(req)
    entity.CreatedBy = &profileID
    entity.OrganizationID = organizationID
    
    // 3. Call usecase
    result, err := h.usecase.Create(ctx, entity)
    if err != nil {
        return nil, err
    }
    
    // 4. Convert domain entity sang protobuf response
    return h.mapper.DomainToPb(result), nil
}

// @Summary Lấy danh sách entities
// @Description Lấy danh sách entities có phân trang
// @Tags YourService
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} pb.YourEntityListResponse "List retrieved successfully"
// @Router /your-entity/list [get]
func (h *YourEntityHandler) GetList(
    ctx context.Context,
    req *pb.YourEntityListRequest,
) (*pb.YourEntityListResponse, error) {
    // Parse pagination từ request
    pagable := _dto.Pagable{
        Page: uint32(req.Page),
        Size: uint32(req.Size),
    }
    
    entities, total, err := h.usecase.GetList(ctx, &pagable)
    if err != nil {
        return nil, err
    }
    
    // Convert list sang protobuf
    pbEntities := make([]*pb.YourEntityResponse, len(entities))
    for i, entity := range entities {
        pbEntities[i] = h.mapper.DomainToPb(&entity)
    }
    
    return &pb.YourEntityListResponse{
        Data:  pbEntities,
        Total: int32(total),
    }, nil
}
```

### Quy tắc quan trọng:
1. **LUÔN embed `Unimplemented...Server`** từ protobuf generated code
2. **Comment `@bind`** để wire generate dependency injection
3. **Comment swagger** (`@Summary`, `@Tags`, `@Router`...) cho documentation
4. **Dùng context.Context**, KHÔNG dùng gin.Context
5. **Lấy profileID/organizationID** từ context qua `_utils`
6. **Mapper** handle toàn bộ conversion giữa PB ↔ Domain
7. **Error** return trực tiếp, middleware sẽ handle format

---

## 11. Error Handling với gRPC Status

**Câu hỏi**: Cách định nghĩa và throw custom errors trong microservices?

**Trả lời**: Sử dụng gRPC status codes kết hợp với custom error codes.

### Định nghĩa custom errors:
```go
package custom_error

import (
    _errors "common/errors"
    "errors"
    "google.golang.org/grpc/codes"
    "google.golang.org/protobuf/protoadapt"
)

var (
    InvalidRequestCode       codes.Code = 4001
    EntityNotFoundCode       codes.Code = 4002
    DuplicateEntityCode      codes.Code = 4003
    PermissionDeniedCode     codes.Code = 4004
    InternalServerErrorCode  codes.Code = 5001
)

// Factory functions cho từng loại error
func InvalidRequest(details ...protoadapt.MessageV1) error {
    return _errors.ThrowError(
        int32(InvalidRequestCode), 
        errors.New("invalid request"), 
        details...,
    )
}

func EntityNotFound(entityType string) error {
    return _errors.ThrowError(
        int32(EntityNotFoundCode), 
        errors.New(entityType + " not found"),
    )
}

func DuplicateEntity(entityType string) error {
    return _errors.ThrowError(
        int32(DuplicateEntityCode), 
        errors.New(entityType + " already exists"),
    )
}

func PermissionDenied(message string) error {
    return _errors.ThrowError(
        int32(codes.PermissionDenied), 
        errors.New(message),
    )
}

func InternalServerError() error {
    return _errors.ThrowError(
        int32(InternalServerErrorCode), 
        errors.New("internal server error"),
    )
}
```

### Sử dụng trong usecase:
```go
func (uc *YourUsecase) Create(ctx context.Context, entity *domain.YourEntity) error {
    // Validation
    if entity.Name == "" {
        return custom_error.InvalidRequest()
    }
    
    // Check duplicate
    existing, _ := uc.repo.FindByName(ctx, entity.Name)
    if existing != nil {
        return custom_error.DuplicateEntity("entity")
    }
    
    // Permission check
    profileID := _utils.GetProfileIdWithContext(ctx)
    if !uc.HasPermission(ctx, profileID, "create") {
        return custom_error.PermissionDenied("you don't have permission to create")
    }
    
    // Business logic
    err := uc.repo.Create(ctx, entity)
    if err != nil {
        return custom_error.InternalServerError()
    }
    
    return nil
}
```

### Quy tắc:
1. **Error codes** nằm trong range: 4xxx (client errors), 5xxx (server errors)
2. **Factory pattern** để tạo errors, dễ quản lý
3. **Middleware** tự động convert error code → HTTP status code
4. **KHÔNG return raw Go errors**, luôn wrap bằng ThrowError

---

## 12. gRPC Client Pattern

**Câu hỏi**: Cách tạo gRPC client để gọi service khác?

**Trả lời**: Client nằm trong `infra/client`, kết nối tới service khác qua gRPC.

### Template client chuẩn:
```go
package client

import (
    "context"
    "fmt"
    "log"
    "time"
    _middleware "common/middleware"
    pb "pb/types/target-service"
    "your-service/env"
    "google.golang.org/grpc"
)

// @bind: your-service/internal/interface/provider.ITargetServiceClient
type TargetServiceClient struct {
    client pb.TargetServiceClient
}

func NewTargetServiceClient() *TargetServiceClient {
    clientWrapper := &TargetServiceClient{client: nil}
    
    go func() {
        // Xác định URL dựa trên environment
        url := "localhost"
        name := "target-service"
        if env.ENV_RUNTIME == "release" {
            url = name  // Trong production dùng service name
        }
        
        // Tạo gRPC connection
        conn, err := grpc.Dial(
            fmt.Sprintf("%s:8202", url),
            grpc.WithInsecure(),
            grpc.WithBackoffMaxDelay(5*time.Second),
            grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
        )
        if err != nil {
            log.Fatalf("failed to dial: %v", err)
        }
        
        ctx, _ := context.WithCancel(context.Background())
        
        // Monitor connection health
        go _middleware.MonitorConnection(ctx, name, conn)
        
        // Tạo client từ connection
        client := pb.NewTargetServiceClient(conn)
        clientWrapper.client = client
    }()
    
    return clientWrapper
}

// Wrapper methods để gọi gRPC APIs
func (c *TargetServiceClient) GetEntityByID(ctx context.Context, id uint64) (*pb.EntityResponse, error) {
    request := &pb.GetEntityRequest{
        Id: id,
    }
    
    response, err := c.client.GetEntity(ctx, request)
    if err != nil {
        return nil, err
    }
    
    return response, nil
}

func (c *TargetServiceClient) GetEntitiesByIDs(ctx context.Context, ids []uint64) ([]*pb.EntityResponse, error) {
    request := &pb.GetEntitiesRequest{
        Ids: ids,
    }
    
    response, err := c.client.GetEntities(ctx, request)
    if err != nil {
        return nil, err
    }
    
    return response.Data, nil
}

// Helper method để map data vào DTO
func (c *TargetServiceClient) EnrichDTOWithEntities(ctx context.Context, dto *YourDTO, entityIDs []uint64) error {
    entities, err := c.GetEntitiesByIDs(ctx, entityIDs)
    if err != nil {
        return err
    }
    
    dto.Entities = entities
    return nil
}
```

### Khi nào cần client?
- **Cross-service communication**: Gọi API từ service khác
- **Data enrichment**: Lấy thông tin bổ sung (user info, product details...)
- **Notification**: Gửi notification qua notification-service
- **History logging**: Log history qua history-service

### Quy tắc:
1. **Client init trong goroutine** để không block startup
2. **WithUnaryInterceptor** để pass context (JWT token, profileId...)
3. **MonitorConnection** để health check
4. **Comment @bind** để wire generate
5. **Wrapper methods** thay vì expose raw gRPC client

---

## 13. Transaction Pattern

**Câu hỏi**: Cách handle database transaction trong usecase?

**Trả lời**: Sử dụng `WithTransaction` pattern để wrap business logic trong transaction.

### Định nghĩa transaction interface:
```go
package iface

type ITransaction interface {
    WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
```

### Triển khai transaction trong infra/postgre:
```go
package postgre

import (
    "context"
    _utils "common/utils"
    "gorm.io/gorm"
)

// @bind: your-service/internal/interface.ITransaction
type TransactionGorm struct {
    DB *gorm.DB
}

func NewTransactionGorm(DB *gorm.DB) *TransactionGorm {
    return &TransactionGorm{DB: DB}
}

func (t *TransactionGorm) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
    return t.DB.Transaction(func(tx *gorm.DB) error {
        // Lấy profileId và organizationId từ context gốc
        profileId := _utils.GetProfileIdWithContext(ctx)
        organizationId := _utils.GetOrganizationIdFromContext(ctx)
        
        // Tạo context mới với tx và preserve user info
        newCtx := context.WithValue(ctx, "tx", tx)
        newCtx = context.WithValue(newCtx, "profileId", profileId)
        newCtx = context.WithValue(newCtx, "organizationId", organizationId)
        
        return fn(newCtx)
    })
}

// Helper function để lấy DB từ context
func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
    if tx, ok := ctx.Value("tx").(*gorm.DB); ok && tx != nil {
        return tx.WithContext(ctx)
    }
    return defaultDB
}
```

### Sử dụng transaction trong usecase:
```go
func (uc *TransactionUsecase) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
    // Validation logic
    if !tx.IsValid() {
        return errors.New("invalid transaction")
    }
    
    // Wrap trong transaction
    err := uc.transaction.WithTransaction(ctx, func(ctx context.Context) error {
        // 1. Create main transaction
        err := uc.transactionRepo.Create(ctx, tx)
        if err != nil {
            return err  // Auto rollback
        }
        
        // 2. Create related action
        action := &domain.Action{
            TransactionID: tx.ID,
            Type:          "create",
            Value:         tx.Amount,
        }
        err = uc.actionRepo.Create(ctx, action)
        if err != nil {
            return err  // Auto rollback
        }
        
        // 3. Update main transaction with action ID
        tx.LastActionID = &action.ID
        err = uc.transactionRepo.Update(ctx, tx)
        if err != nil {
            return err  // Auto rollback
        }
        
        // 4. Notify user
        err = uc.notificationClient.NotifyUser(ctx, tx.UserID, "Transaction created")
        if err != nil {
            return err  // Auto rollback
        }
        
        return nil  // Commit nếu không có error
    })
    
    return err
}
```

### Repository phải dùng GetDB:
```go
func (r *YourRepo) Create(ctx context.Context, entity *domain.Entity) error {
    // LUÔN dùng GetDB để lấy tx từ context nếu có
    return GetDB(ctx, r.DB).Create(entity).Error
}

func (r *YourRepo) Update(ctx context.Context, entity *domain.Entity) error {
    return GetDB(ctx, r.DB).
        Where("id = ?", entity.ID).
        Updates(entity).Error
}
```

### Quy tắc:
1. **WithTransaction** tự động rollback nếu fn return error
2. **Preserve context values** (profileId, organizationId) vào tx context
3. **GetDB helper** để repository tự động dùng transaction
4. **Nested operations** (create → update → notify) đều trong cùng transaction
5. **External calls** (notification, client) cũng rollback nếu fail

---

## 14. Validator Pattern

**Câu hỏi**: Cách viết validator để validate request DTO?

**Trả lời**: Validator nằm trong `infra/validator`, validate trước khi gọi usecase.

### Template validator chuẩn:
```go
package validator

import (
    "errors"
    pb "pb/types/your-service"
    "google.golang.org/protobuf/protoadapt"
)

type YourEntityValidator struct{}

func NewYourEntityValidator() *YourEntityValidator {
    return &YourEntityValidator{}
}

// Validate create request
func (v *YourEntityValidator) ValidateCreateRequest(req *pb.CreateYourEntityRequest) error {
    details := []protoadapt.MessageV1{}
    
    // 1. Required field validation
    if req.Name == "" {
        details = append(details, &pb.ErrorDetail{
            Field:   "name",
            Message: "Name is required",
        })
    }
    
    // 2. Length validation
    if len(req.Name) > 255 {
        details = append(details, &pb.ErrorDetail{
            Field:   "name",
            Message: "Name must not exceed 255 characters",
        })
    }
    
    // 3. Range validation
    if req.Price < 0 {
        details = append(details, &pb.ErrorDetail{
            Field:   "price",
            Message: "Price must be greater than or equal to 0",
        })
    }
    
    // 4. Enum validation
    if !isValidStatus(req.Status) {
        details = append(details, &pb.ErrorDetail{
            Field:   "status",
            Message: "Invalid status value",
        })
    }
    
    // 5. Date validation
    if req.StartDate != "" && req.EndDate != "" {
        startDate := _utils.ParseStringToTime(req.StartDate)
        endDate := _utils.ParseStringToTime(req.EndDate)
        
        if startDate != nil && endDate != nil && endDate.Before(*startDate) {
            details = append(details, &pb.ErrorDetail{
                Field:   "endDate",
                Message: "End date must be after start date",
            })
        }
    }
    
    // 6. Return error nếu có validation failures
    if len(details) > 0 {
        return custom_error.InvalidRequest(details...)
    }
    
    return nil
}

// Validate update request
func (v *YourEntityValidator) ValidateUpdateRequest(req *pb.UpdateYourEntityRequest) error {
    details := []protoadapt.MessageV1{}
    
    // ID required for update
    if req.Id == 0 {
        details = append(details, &pb.ErrorDetail{
            Field:   "id",
            Message: "ID is required",
        })
    }
    
    // Validate other fields...
    
    if len(details) > 0 {
        return custom_error.InvalidRequest(details...)
    }
    
    return nil
}

// Helper validation functions
func isValidStatus(status int32) bool {
    validStatuses := []int32{1, 2, 3, 4}
    for _, v := range validStatuses {
        if v == status {
            return true
        }
    }
    return false
}
```

### Sử dụng trong handler:
```go
func (h *YourEntityHandler) Create(ctx context.Context, req *pb.CreateYourEntityRequest) (*pb.YourEntityResponse, error) {
    // 1. Validate trước khi process
    if err := h.validator.ValidateCreateRequest(req); err != nil {
        return nil, err
    }
    
    // 2. Process business logic
    entity := h.mapper.PbToDomain(req)
    result, err := h.usecase.Create(ctx, entity)
    if err != nil {
        return nil, err
    }
    
    return h.mapper.DomainToPb(result), nil
}
```

### Quy tắc validator:
1. **Collect all errors** vào details array thay vì return ngay
2. **Field-level errors** với field name và message
3. **Return InvalidRequest** với details khi có lỗi
4. **Validator riêng** cho Create, Update, Delete operations
5. **Helper functions** cho complex validations (date, enum, regex...)

---

## 15. Pagination Response Pattern

**Câu hỏi**: Cách chuẩn hóa response có phân trang?

**Trả lời**: Response phân trang LUÔN có `data` (array) và `total` (int64).

### Protobuf definition:
```protobuf
message YourEntityListRequest {
  int32 page = 1;
  int32 size = 2;
  string sort = 3;
  // Filter fields...
}

message YourEntityListResponse {
  repeated YourEntityResponse data = 1;
  int32 total = 2;
}
```

### Handler implementation:
```go
func (h *YourEntityHandler) GetList(ctx context.Context, req *pb.YourEntityListRequest) (*pb.YourEntityListResponse, error) {
    // 1. Parse pagination
    pagable := _dto.Pagable{
        Page: uint32(req.Page),
        Size: uint32(req.Size),
    }
    
    // 2. Build filters
    filters := &dto.YourEntityFilters{
        Pagable:    pagable,
        Status:     req.Status,
        SearchTerm: req.SearchTerm,
    }
    
    // 3. Get data from usecase
    entities, total, err := h.usecase.GetList(ctx, filters)
    if err != nil {
        return nil, err
    }
    
    // 4. Map to protobuf
    pbEntities := make([]*pb.YourEntityResponse, len(entities))
    for i, entity := range entities {
        pbEntities[i] = h.mapper.DomainToPb(&entity)
    }
    
    // 5. Return with data + total
    return &pb.YourEntityListResponse{
        Data:  pbEntities,
        Total: int32(total),
    }, nil
}
```

### Repository với pagination:
```go
func (r *YourEntityRepo) GetList(c context.Context, pagable _dto.IPagable, filters *dto.Filters) ([]domain.YourEntity, int64, error) {
    var entities []domain.YourEntity
    var total int64
    
    // Build base query
    query := r.GetDB(c).Model(&domain.YourEntity{}).
        Where("deleted_at IS NULL")
    
    // Apply filters
    if filters.Status != 0 {
        query = query.Where("status = ?", filters.Status)
    }
    if filters.SearchTerm != "" {
        query = query.Where("name LIKE ?", "%"+filters.SearchTerm+"%")
    }
    
    // Get total count
    query.Count(&total)
    
    // Apply pagination và lấy data
    err := query.
        Order("created_at DESC").
        Offset(pagable.GetOffset()).
        Limit(pagable.GetLimit()).
        Find(&entities).Error
    
    return entities, total, err
}
```

### Quy tắc:
1. **LUÔN trả về `data` và `total`** trong list response
2. **Default page = 0, size = 20, max size = 100**
3. **Count total trước** khi apply pagination
4. **Apply filters → Count → Paginate → Fetch**
5. **DTO embed `_dto.Pagable`** khi cần phân trang

---

## 16. Wire Dependency Injection Pattern

**Câu hỏi**: Cách setup dependency injection với Google Wire?

**Trả lời**: Wire tự động generate code inject dependencies dựa trên `wire.go`.

### wire.go structure:
```go
//go:build wireinject
// +build wireinject

package wire

import (
    "github.com/google/wire"
    _db "common/db"
    
    "your-service/infra/client"
    "your-service/infra/handler"
    "your-service/infra/mapper"
    "your-service/infra/postgre"
    "your-service/infra/validator"
    "your-service/initial"
    "your-service/internal/interface/provider"
    "your-service/internal/interface/repo"
    "your-service/internal/usecase"
)

// Inject dependencies
var wireSet = wire.NewSet(
    // Database
    _db.NewDB,
    _db.NewTransactionRepo,
    
    // Bindings: Interface → Implementation
    wire.Bind(new(provider.IUserClient), new(*client.UserClient)),
    wire.Bind(new(provider.INotificationClient), new(*client.NotificationClient)),
    wire.Bind(new(repo.IYourEntityRepo), new(*postgre.YourEntityRepo)),
    wire.Bind(new(provider.ITransaction), new(*postgre.TransactionGorm)),
    
    // Clients
    client.NewUserClient,
    client.NewNotificationClient,
    
    // Repositories
    postgre.NewYourEntityRepo,
    postgre.NewTransactionGorm,
    
    // Usecases
    usecase.NewYourEntityUsecase,
    
    // Mappers
    mapper.NewYourEntityMapper,
    
    // Validators
    validator.NewYourEntityValidator,
    
    // Handlers
    handler.NewYourEntityHandler,
    
    // Initial
    initial.NewApp,
)

func InitializeApp() (*initial.App, error) {
    wire.Build(wireSet)
    return nil, nil
}
```

### Chạy wire để generate:
```bash
cd your-service/wire
wire
```

### wire_gen.go (auto-generated):
```go
// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
    // imports...
)

func InitializeApp() (*initial.App, error) {
    db := _db.NewDB()
    transactionRepo := _db.NewTransactionRepo(db)
    
    yourEntityRepo := postgre.NewYourEntityRepo(transactionRepo)
    yourEntityUsecase := usecase.NewYourEntityUsecase(yourEntityRepo)
    
    yourEntityMapper := mapper.NewYourEntityMapper()
    yourEntityValidator := validator.NewYourEntityValidator()
    
    yourEntityHandler := handler.NewYourEntityHandler(
        yourEntityUsecase,
        yourEntityMapper,
        yourEntityValidator,
    )
    
    app := initial.NewApp(yourEntityHandler)
    return app, nil
}
```

### Quy tắc wire:
1. **`@bind` comment** trên struct để wire detect interface binding
2. **`wire.NewSet`** định nghĩa tất cả providers
3. **`wire.Bind`** để bind interface → implementation
4. **Chạy `wire` command** sau khi thay đổi dependencies
5. **KHÔNG edit `wire_gen.go`** - file này auto-generated

---

## 17. Gateway Registration Pattern

**Câu hỏi**: Cách đăng ký handler mới vào gateway-service?

**Trả lời**: Khi tạo service mới hoặc thêm handler mới, phải đăng ký trong gateway.

### 1. Tạo file service registration trong gateway-service:
```go
// gateway-service/service/your_service.go
package service

import (
    "common"
    "context"
    "fmt"
    pb "pb/types/your-service"
    
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "github.com/spf13/viper"
    "google.golang.org/grpc"
)

func RegisterYourService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
    // Lấy port từ config
    port := viper.GetString("service.port.your-service")
    url := common.GetPrefixProtobufUrl(env, "your-service")
    
    endpoint := fmt.Sprintf("%s:%s", url, port)
    
    // Đăng ký tất cả handlers của service
    _ = pb.RegisterYourServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
    _ = pb.RegisterYourAdminServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
    // Thêm các handlers khác...
}
```

### 2. Đăng ký trong main.go của gateway:
```go
// gateway-service/cmd/http/main.go
func main() {
    // ... setup code ...
    
    ctx := context.Background()
    opts := []grpc.DialOption{
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithConnectParams(grpc.ConnectParams{
            Backoff: backoff.Config{
                BaseDelay:  5 * time.Second,
                Multiplier: 1.0,
                Jitter:     0.0,
                MaxDelay:   5 * time.Second,
            },
            MinConnectTimeout: 5 * time.Second,
        }),
    }
    
    // Đăng ký các gRPC services
    service.RegisterNotificationService(ctx, mux, env, opts)
    service.RegisterPaymentService(ctx, mux, env, opts)
    service.RegisterOrganizationService(ctx, mux, env, opts)
    service.RegisterBdsproService(ctx, mux, env, opts)
    service.RegisterCrmService(ctx, mux, env, opts)
    service.RegisterSocialService(ctx, mux, env, opts)
    service.RegisterYourService(ctx, mux, env, opts)  // THÊM service mới
    
    // ... start server ...
}
```

### 3. Cập nhật config:
```yaml
# gateway-service/config/config.yml
service:
  port:
    your-service: "8210"  # Port của service mới
```

### Quy tắc:
1. **Mỗi service có 1 file registration** trong `gateway-service/service/`
2. **Đăng ký TẤT CẢ handlers** của service trong function
3. **Thêm service vào main.go** để enable routing
4. **Config port** trong config.yml
5. **Endpoint format**: `{service-name}:{port}` (production) hoặc `localhost:{port}` (development)

---

## 18. Enum Validation Pattern

**Câu hỏi**: Cách định nghĩa và validate enum values?

**Trả lời**: Enum được định nghĩa trong `internal/enum`, có validation method.

### Định nghĩa enum:
```go
package enum

type EYourStatus int32

const (
    StatusPending   EYourStatus = 1
    StatusActive    EYourStatus = 2
    StatusInactive  EYourStatus = 3
    StatusDeleted   EYourStatus = 4
)

// Map enum → display name
var StatusMap = map[EYourStatus]string{
    StatusPending:  "Đang chờ",
    StatusActive:   "Hoạt động",
    StatusInactive: "Ngưng hoạt động",
    StatusDeleted:  "Đã xóa",
}

// Validation method
func (e EYourStatus) IsValid() bool {
    _, ok := StatusMap[e]
    return ok
}

// Get display name
func (e EYourStatus) String() string {
    if name, ok := StatusMap[e]; ok {
        return name
    }
    return "Unknown"
}

// Get all valid values
func GetAllStatuses() []EYourStatus {
    return []EYourStatus{
        StatusPending,
        StatusActive,
        StatusInactive,
        StatusDeleted,
    }
}
```

### Sử dụng enum trong domain:
```go
package domain

import (
    "your-service/internal/enum"
    _models "common/models"
)

type YourEntity struct {
    _models.BaseEntity
    Name   string
    Status enum.EYourStatus `gorm:"type:int;not null"`
}
```

### Validation trong usecase:
```go
func (uc *YourUsecase) Create(ctx context.Context, entity *domain.YourEntity) error {
    // Validate enum value
    if !entity.Status.IsValid() {
        return errors.New("invalid status value")
    }
    
    // Business logic...
    return uc.repo.Create(ctx, entity)
}
```

### Enum trong protobuf:
```protobuf
enum YourStatus {
  PENDING = 1;
  ACTIVE = 2;
  INACTIVE = 3;
  DELETED = 4;
}

message YourEntityResponse {
  uint64 id = 1;
  string name = 2;
  YourStatus status = 3;
}
```

### Mapping enum giữa domain và protobuf:
```go
// Domain → Protobuf
func (m *YourMapper) DomainToPb(entity *domain.YourEntity) *pb.YourEntityResponse {
    return &pb.YourEntityResponse{
        Id:     entity.ID,
        Name:   entity.Name,
        Status: pb.YourStatus(entity.Status),  // Cast enum
    }
}

// Protobuf → Domain
func (m *YourMapper) PbToDomain(req *pb.CreateRequest) *domain.YourEntity {
    return &domain.YourEntity{
        Name:   req.Name,
        Status: enum.EYourStatus(req.Status),  // Cast enum
    }
}
```

### Quy tắc enum:
1. **Type alias `int32`** cho enum
2. **Map enum → string** cho display name tiếng Việt
3. **`IsValid()` method** để validation
4. **`String()` method** để get display name
5. **Cast khi mapping** giữa domain ↔ protobuf

---

## 19. Preload với GORM Pattern

**Câu hỏi**: Cách load relationships (foreign keys) với GORM?

**Trả lời**: Dùng `Preload` để eager load related entities, tránh N+1 query.

### Entity với relationships:
```go
type Transaction struct {
    _models.BaseEntity
    Name        string
    Amount      float64
    Status      enum.ETransactionStatus
    
    // Foreign keys
    UserID      uint64
    ProductID   uint64
    CategoryID  *uint64  // Optional FK
    
    // Relationships (không lưu vào DB)
    User        *User        `gorm:"foreignKey:UserID"`
    Product     *Product     `gorm:"foreignKey:ProductID"`
    Category    *Category    `gorm:"foreignKey:CategoryID"`
    Actions     []Action     `gorm:"foreignKey:TransactionID"`
}

type Action struct {
    _models.BaseEntity
    TransactionID uint64
    Type          string
    Value         float64
}
```

### Repository với Preload:
```go
// Get by ID với relationships
func (r *TransactionRepo) GetByID(c context.Context, id uint64) (*domain.Transaction, error) {
    var tx domain.Transaction
    err := r.GetDB(c).
        Preload("User").              // Load User
        Preload("Product").           // Load Product
        Preload("Category").          // Load Category (nullable)
        Preload("Actions", func(db *gorm.DB) *gorm.DB {
            return db.Order("created_at DESC")  // Order actions
        }).
        Where("id = ? AND deleted_at IS NULL", id).
        First(&tx).Error
    
    if err != nil {
        return nil, err
    }
    return &tx, nil
}

// Get list với relationships
func (r *TransactionRepo) GetList(c context.Context, pagable _dto.IPagable) ([]domain.Transaction, int64, error) {
    var transactions []domain.Transaction
    var total int64
    
    query := r.GetDB(c).
        Model(&domain.Transaction{}).
        Preload("User").
        Preload("Product").
        Preload("Actions").
        Where("deleted_at IS NULL")
    
    // Count total
    query.Count(&total)
    
    // Fetch with pagination
    err := query.
        Order("created_at DESC").
        Offset(pagable.GetOffset()).
        Limit(pagable.GetLimit()).
        Find(&transactions).Error
    
    return transactions, total, err
}

// Nested preload (load relationships của relationship)
func (r *TransactionRepo) GetDetailWithNested(c context.Context, id uint64) (*domain.Transaction, error) {
    var tx domain.Transaction
    err := r.GetDB(c).
        Preload("User.Profile").                    // Load User → Profile
        Preload("Product.Category").                // Load Product → Category
        Preload("Actions.CreatedByUser").           // Load Actions → User
        Where("id = ? AND deleted_at IS NULL", id).
        First(&tx).Error
    
    return &tx, err
}
```

### Conditional Preload:
```go
func (r *TransactionRepo) GetByIDWithOptions(c context.Context, id uint64, includeUser, includeProduct bool) (*domain.Transaction, error) {
    var tx domain.Transaction
    query := r.GetDB(c)
    
    // Conditional preload
    if includeUser {
        query = query.Preload("User")
    }
    if includeProduct {
        query = query.Preload("Product")
    }
    
    err := query.
        Where("id = ? AND deleted_at IS NULL", id).
        First(&tx).Error
    
    return &tx, err
}
```

### Quy tắc Preload:
1. **Dùng Preload** để tránh N+1 query problem
2. **Định nghĩa relationships** trong struct với `gorm:"foreignKey:..."`
3. **Preload selective** - chỉ load khi cần
4. **Nested preload** với dot notation: `"User.Profile"`
5. **Filter trong Preload** với callback function
6. **KHÔNG dùng Join** trừ khi cần filter theo related table

---

## 20. Generic Response Pattern

**Câu hỏi**: Cách chuẩn hóa response structure cho consistency?

**Trả lời**: Sử dụng generic DTO để wrap response data.

### Generic Response DTO:
```go
package dto

type Response[T any] struct {
    Data  T     `json:"data"`
    Total int64 `json:"total,omitempty"`
}

// Single entity response
type SingleResponse[T any] struct {
    Data T `json:"data"`
}

// List response with pagination
type ListResponse[T any] struct {
    Data  []T   `json:"data"`
    Total int64 `json:"total"`
}

// Success response without data
type SuccessResponse struct {
    Message string `json:"message"`
}
```

### Sử dụng trong handler:
```go
// Single entity
func (h *YourHandler) GetByID(ctx context.Context, req *pb.GetRequest) (*pb.YourEntityResponse, error) {
    entity, err := h.usecase.GetByID(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    
    // Trả về entity đơn
    return h.mapper.DomainToPb(entity), nil
}

// List với pagination
func (h *YourHandler) GetList(ctx context.Context, req *pb.ListRequest) (*pb.YourEntityListResponse, error) {
    pagable := _dto.Pagable{
        Page: uint32(req.Page),
        Size: uint32(req.Size),
    }
    
    entities, total, err := h.usecase.GetList(ctx, &pagable)
    if err != nil {
        return nil, err
    }
    
    // Map entities
    pbEntities := make([]*pb.YourEntityResponse, len(entities))
    for i, entity := range entities {
        pbEntities[i] = h.mapper.DomainToPb(&entity)
    }
    
    // Trả về list với total
    return &pb.YourEntityListResponse{
        Data:  pbEntities,
        Total: int32(total),
    }, nil
}

// Success operation
func (h *YourHandler) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.SuccessResponse, error) {
    err := h.usecase.Delete(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    
    return &pb.SuccessResponse{
        Message: "Entity deleted successfully",
    }, nil
}
```

### Quy tắc response:
1. **Single entity**: Trả về entity trực tiếp hoặc wrap trong `data`
2. **List**: LUÔN trả `data` (array) và `total` (int64)
3. **Success operation**: Trả message hoặc empty response
4. **Error**: Throw error, middleware sẽ format

---

## 21. Batch Operations Pattern

**Câu hỏi**: Cách xử lý bulk operations hiệu quả?

**Trả lời**: Sử dụng batch processing với transaction và bulk insert.

### Bulk create:
```go
func (r *YourEntityRepo) BulkCreate(c context.Context, entities []domain.YourEntity) error {
    // GORM bulk insert (1 query)
    return r.GetDB(c).CreateInBatches(entities, 100).Error
}
```

### Bulk update:
```go
func (r *YourEntityRepo) BulkUpdateStatus(c context.Context, ids []uint64, status enum.EYourStatus) error {
    return r.GetDB(c).
        Model(&domain.YourEntity{}).
        Where("id IN (?) AND deleted_at IS NULL", ids).
        Update("status", status).Error
}
```

### Bulk delete (soft delete):
```go
func (r *YourEntityRepo) BulkDelete(c context.Context, ids []uint64) error {
    return r.GetDB(c).
        Model(&domain.YourEntity{}).
        Where("id IN (?) AND deleted_at IS NULL", ids).
        Update("deleted_at", time.Now()).Error
}
```

### Usecase với batch + transaction:
```go
func (uc *YourUsecase) BulkCreate(ctx context.Context, dtos []*dto.CreateDTO) error {
    // Convert DTOs to entities
    entities := make([]domain.YourEntity, len(dtos))
    for i, dto := range dtos {
        entities[i] = *uc.mapper.DTOToDomain(dto)
    }
    
    // Wrap trong transaction
    return uc.transaction.WithTransaction(ctx, func(ctx context.Context) error {
        // 1. Bulk insert
        err := uc.repo.BulkCreate(ctx, entities)
        if err != nil {
            return err
        }
        
        // 2. Create history logs
        for _, entity := range entities {
            historyDTO := &_dto.HistoryDTO{
                TargetID:   entity.ID,
                TargetType: enum.ETargetYourEntity,
                ActionType: enum.EHistoryCreate,
            }
            err = uc.historyClient.CreateHistory(ctx, historyDTO)
            if err != nil {
                return err  // Rollback all
            }
        }
        
        return nil
    })
}
```

### Quy tắc batch operations:
1. **Dùng `CreateInBatches`** với batch size = 100
2. **Wrap trong transaction** để ensure consistency
3. **Bulk update/delete** dùng `Where IN`
4. **Avoid N+1** - process in batches, not loops
5. **History/Audit logs** sau khi batch operation success

---

## Tổng kết

Các pattern này là nền tảng của Clean Architecture trong BDSPro microservices. Tuân thủ những quy tắc này sẽ giúp code:
- Dễ đọc và maintain
- Testable
- Reusable
- Consistent across services
- Following SOLID principles
- Scalable và performant

**Nguyên tắc vàng**:
1. **Clean Architecture**: Tách biệt layers rõ ràng
2. **Interface-driven**: Code against interfaces, not implementations
3. **Dependency Injection**: Dùng Wire để quản lý dependencies
4. **Error Handling**: Custom errors với gRPC status codes
5. **Transaction Safety**: Wrap critical operations trong transaction
6. **Validation**: Validate early, fail fast
7. **Consistency**: Follow patterns đã thiết lập trong project

