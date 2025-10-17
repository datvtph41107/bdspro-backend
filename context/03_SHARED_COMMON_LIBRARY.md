# Shared Common Library - Complete Context

## 📋 Overview

**Location**: `/shared/common/`
**Purpose**: Reusable domain models, DTOs, enums, and base implementations
**Type**: Shared Go Module
**Status**: ✅ Production Ready

**Key Principle**: DRY (Don't Repeat Yourself) - All reusable code across microservices is centralized here.

---

## 🎯 Core Responsibilities

### 1. Base Entities & Models
- Provide common entity structures
- Audit trail support (CreatedBy, UpdatedBy)
- Soft delete mechanism
- Automatic timestamp management

### 2. Common DTOs
- Pagination DTOs
- Error DTOs
- History DTOs
- Notification DTOs

### 3. Shared Enumerations
- Context keys
- History action types
- Notification types
- Owner types

### 4. CRUD Base Implementations
- Generic repository interface
- Generic usecase interface
- Base repository provider
- Transaction support

### 5. Code Generators
- Product code generator
- Asset code generator
- Post code generator

---

## 📁 Directory Structure

```
shared/common/
├── domain/
│   ├── crud/
│   │   ├── repo.go              # ICrudRepo interface
│   │   └── usecase.go           # BaseUsecase generic implementation
│   ├── dto/
│   │   ├── pagable.go           # Pagination DTO
│   │   ├── history_dto.go       # History logging DTO
│   │   └── notification_dto.go  # Notification DTO
│   ├── entity/
│   │   ├── base_entity.go       # Base entity with common fields
│   │   └── audit.go             # Audit trail implementation
│   ├── enum/
│   │   ├── context_key_enum.go  # Context keys
│   │   ├── history_enum.go      # History action types
│   │   ├── notification_enum.go # Notification & owner types
│   │   └── target_notification_enum.go
│   ├── err/
│   │   └── dto.go               # Error DTO
│   └── usecase/
│       └── code_data_usecase.go # Code generators
├── provider/
│   └── repository.go            # Base CrudRepo implementation
├── go.mod
└── go.sum
```

---

## 1️⃣ CRUD Module (`domain/crud/`)

### ICrudRepo Interface

**File**: `crud/repo.go`

**Purpose**: Generic CRUD repository interface that all repositories should implement

```go
type ICrudRepo[T any] interface {
    Create(c context.Context, entity *T) error
    Update(c context.Context, id uint64, entity *T) error
    Delete(c context.Context, id uint64) error
    GetAll(c context.Context) ([]T, error)
    GetByID(c context.Context, id uint64) (*T, error)
    GetList(c context.Context, pagable _dto.IPagable) ([]T, int64, error)
    GetListByIDs(c context.Context, ids []uint64) ([]T, error)
    GetDetail(c context.Context, id uint64) (*T, error)
    BeforeSave(c context.Context, id *uint64, entity *T) error
    AfterSave(c context.Context, id *uint64, entity *T) error
}
```

**Key Methods**:
- `Create`: Insert new record
- `Update`: Update existing record
- `Delete`: Soft delete (set deleted_at)
- `GetByID`: Get one record by ID
- `GetDetail`: Get one record with relations
- `GetAll`: Get all records (no pagination)
- `GetList`: Get paginated list
- `GetListByIDs`: Get multiple records by IDs
- `BeforeSave`: Hook before create/update
- `AfterSave`: Hook after create/update

**Usage Example**:
```go
// In internal/interface/repo/product_repo.go
type IProductRepo interface {
    _crud.ICrudRepo[Product]
    // Add custom methods
    GetByCode(ctx context.Context, code string) (*Product, error)
}
```

### BaseUsecase Generic Implementation

**File**: `crud/usecase.go`

**Purpose**: Base usecase with standard CRUD logic, error handling

```go
type IBaseUsecase[T any] interface {
    Create(c context.Context, entity *T) (*T, error)
    Update(c context.Context, id uint64, entity *T) (*T, error)
    Delete(c context.Context, id uint64) (bool, error)
    GetByID(c context.Context, id uint64) (*T, error)
    GetAll(c context.Context) ([]T, error)
    GetDetail(c context.Context, id uint64) (*T, error)
    GetList(c context.Context, pagable _dto.IPagable) ([]T, int64, error)
}

type BaseUsecase[T any, R ICrudRepo[T]] struct {
    Repo R
}
```

**Features**:
- Auto error wrapping with ErrorDTO
- Vietnamese error messages
- Consistent error codes (404, 500, etc.)

**Usage Example**:
```go
// In internal/usecase/product_usecase.go
type ProductUsecase struct {
    _crud.BaseUsecase[Product, IProductRepo]
    // Add custom fields
}

func NewProductUsecase(repo IProductRepo) *ProductUsecase {
    return &ProductUsecase{
        BaseUsecase: _crud.BaseUsecase[Product, IProductRepo]{
            Repo: repo,
        },
    }
}

// Inherit all CRUD methods, add custom ones
func (u *ProductUsecase) GetByCode(ctx context.Context, code string) (*Product, error) {
    // Custom logic
}
```

---

## 2️⃣ DTO Module (`domain/dto/`)

### Pagable DTO

**File**: `dto/pagable.go`

**Purpose**: Standard pagination for list APIs

```go
type IPagable interface {
    GetPage() int
    GetSize() int
    GetSort() string
    GetOffset() int
    GetLimit() int
}

type Pagable struct {
    Page int    `json:"page" form:"page"`
    Size int    `json:"size" form:"size"`
    Sort string `json:"sort" form:"sort"`
}

const (
    DefaultSize = 20
    MaxSize     = 100
)

func (p *Pagable) GetOffset() int {
    if p.Page < 1 {
        p.Page = 1
    }
    return (p.Page - 1) * p.GetLimit()
}

func (p *Pagable) GetLimit() int {
    if p.Size <= 0 {
        p.Size = DefaultSize
    }
    if p.Size > MaxSize {
        p.Size = MaxSize
    }
    return p.Size
}
```

**Usage**:
```go
// In DTO
type GetProductListRequest struct {
    _dto.Pagable
    Text   string `json:"text" form:"text"`
    Status int    `json:"status" form:"status"`
}

// In handler
req := &GetProductListRequest{}
if err := ctx.ShouldBindQuery(req); err != nil {
    return err
}

// In usecase
products, total, err := usecase.GetList(ctx, req)

// Response
{
    "data": [...],
    "total": 100,
    "page": 1,
    "size": 20
}
```

### History DTO

**File**: `dto/history_dto.go`

**Purpose**: Log user actions across all services

```go
type HistoryDTO struct {
    Title      string   `json:"title"`
    Note       []string `json:"note"`
    TargetId   uint64   `json:"targetId"`
    TargetType int32    `json:"targetType"`   // Contact, Post, Product, etc.
    ActionType int32    `json:"actionType"`   // Create, Update, Delete, etc.
    OwnerID    uint64   `json:"ownerId"`
    OwnerOf    int32    `json:"ownerOf"`      // Member, Group, Organization
}
```

**Usage**:
```go
// After creating a contact
historyDTO := &_dto.HistoryDTO{
    Title:      "Tạo mới liên hệ",
    Note:       []string{"Tên: John Doe", "SĐT: 0901234567"},
    TargetId:   contact.ID,
    TargetType: int32(_enum.ETargetContact),
    ActionType: int32(_enum.EHistoryContactCreate),
    OwnerID:    profileId,
    OwnerOf:    int32(_enum.EOwnerOfMember),
}

// Send to history service
historyClient.CreateHistory(ctx, historyDTO)
```

### Notification DTO

**File**: `dto/notification_dto.go`

**Purpose**: Send notifications to users

```go
type NotificationDTO struct {
    Title      string   `json:"title"`
    Message    string   `json:"message"`
    Type       int32    `json:"type"`        // Info, Warning, Error, Success
    UserID     uint64   `json:"userId"`
    VisibleAt  string   `json:"visibleAt"`
}

type NotificationBatchDTO struct {
    Notifications []NotificationDTO `json:"notifications"`
}
```

**Usage**:
```go
notiDTO := &_dto.NotificationDTO{
    Title:   "Liên hệ mới",
    Message: "Bạn có 1 liên hệ mới từ John Doe",
    Type:    int32(_enum.ENotificationInfo),
    UserID:  profileId,
}

notificationClient.SendNotification(ctx, notiDTO)
```

---

## 3️⃣ Entity Module (`domain/entity/`)

### BaseEntity

**File**: `entity/base_entity.go`

**Purpose**: Common fields for all entities

```go
type BaseEntity struct {
    ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
    CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
    UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
    DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
    AuditBase            // Embedded audit fields
}
```

**Features**:
- Auto-increment ID
- Auto timestamps (CreatedAt, UpdatedAt)
- Soft delete support (DeletedAt)
- Embedded audit trail

**Usage**:
```go
type Product struct {
    _models.BaseEntity
    Name        string  `gorm:"size:255" json:"name"`
    Price       float64 `json:"price"`
    Description string  `gorm:"type:text" json:"description"`
}

// BaseEntity provides: ID, CreatedAt, UpdatedAt, DeletedAt, CreatedBy, UpdatedBy
```

### AuditBase

**File**: `entity/audit.go`

**Purpose**: Track who created/updated records

```go
type AuditBase struct {
    CreatedBy *uint64 `gorm:"column:created_by" json:"createdBy,omitempty"`
    UpdatedBy *uint64 `gorm:"column:updated_by" json:"updatedBy,omitempty"`
}

// GORM Hook - Auto set CreatedBy & UpdatedBy
func (a *AuditBase) BeforeCreate(tx *gorm.DB) error {
    ctx := tx.Statement.Context
    if profileId := getProfileIdFromContext(ctx); profileId > 0 {
        a.CreatedBy = &profileId
        a.UpdatedBy = &profileId
    }
    return nil
}

func (a *AuditBase) BeforeUpdate(tx *gorm.DB) error {
    ctx := tx.Statement.Context
    if profileId := getProfileIdFromContext(ctx); profileId > 0 {
        a.UpdatedBy = &profileId
    }
    return nil
}
```

**Key Features**:
- Auto-populate CreatedBy on insert
- Auto-populate UpdatedBy on update
- Get profileId from context automatically

**Usage**:
```go
// Just embed BaseEntity - AuditBase is included
type Contact struct {
    _models.BaseEntity
    Name  string
    Phone string
}

// When creating:
// CreatedBy and UpdatedBy will be auto-set from context
```

---

## 4️⃣ Enum Module (`domain/enum/`)

### Context Keys

**File**: `enum/context_key_enum.go`

```go
type ContextKey string

const (
    ProfileIDKey      ContextKey = "profileId"
    OrganizationIDKey ContextKey = "organizationId"
    AuthIDKey         ContextKey = "authId"
    SessionKey        ContextKey = "sessionKey"
    RoleKey           ContextKey = "role"
    TypeKey           ContextKey = "type"
    PlanIDKey         ContextKey = "planId"
    PlanFromKey       ContextKey = "planFrom"
)
```

**Usage**:
```go
// Set in middleware
ctx = context.WithValue(ctx, _enum.ProfileIDKey, profileId)

// Get in usecase
profileId := ctx.Value(_enum.ProfileIDKey).(uint64)

// Or use utility
profileId := _utils.GetProfileIdWithContext(ctx)
```

### History Action Types

**File**: `enum/history_enum.go`

```go
type EHistory int32

const (
    // Contact actions
    EHistoryContactCreate   EHistory = 101
    EHistoryContactUpdate   EHistory = 102
    EHistoryContactDelete   EHistory = 103
    EHistoryContactAssign   EHistory = 104
    EHistoryContactMessage  EHistory = 105
    
    // Post actions
    EHistoryPostCreate      EHistory = 201
    EHistoryPostUpdate      EHistory = 202
    EHistoryPostDelete      EHistory = 203
    
    // Product actions
    EHistoryProductCreate   EHistory = 301
    EHistoryProductUpdate   EHistory = 302
    EHistoryProductDelete   EHistory = 303
    EHistoryProductMerge    EHistory = 304
    EHistoryProductDivide   EHistory = 305
    
    // Asset actions
    EHistoryAssetCreate     EHistory = 401
    EHistoryAssetUpdate     EHistory = 402
    EHistoryAssetDelete     EHistory = 403
    // ... more actions
)

var HistoryMap = map[EHistory]string{
    EHistoryContactCreate: "Tạo liên hệ",
    EHistoryContactUpdate: "Cập nhật liên hệ",
    EHistoryContactDelete: "Xóa liên hệ",
    // ... more mappings
}
```

### Notification & Owner Types

**File**: `enum/notification_enum.go`

```go
type EOwnerOf int

const (
    EOwnerOfMember       EOwnerOf = 10  // Individual member
    EOwnerOfGroup        EOwnerOf = 20  // Group/Team
    EOwnerOfOrganization EOwnerOf = 30  // Organization
)

func (e EOwnerOf) IsValid() bool {
    return e == EOwnerOfMember || e == EOwnerOfGroup || e == EOwnerOfOrganization
}
```

### Target Types

**File**: `enum/target_notification_enum.go`

```go
type ETargetHistory int

const (
    ETargetLead       ETargetHistory = 1
    ETargetContact    ETargetHistory = 2
    ETargetPipeline   ETargetHistory = 3
    ETargetStage      ETargetHistory = 4
    ETargetPost       ETargetHistory = 5
    ETargetProduct    ETargetHistory = 6
    ETargetAsset      ETargetHistory = 7
    ETargetCampaign   ETargetHistory = 8
    // ... more targets
)

var TargetNotificationMap = map[ETargetHistory]string{
    ETargetContact:  "Liên hệ",
    ETargetPost:     "Tin đăng",
    ETargetProduct:  "Sản phẩm",
    ETargetAsset:    "Tài sản",
    // ... more mappings
}
```

---

## 5️⃣ Error Module (`domain/err/`)

**File**: `err/dto.go`

```go
type ErrorDTO struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func (e *ErrorDTO) Error() string {
    return e.Message
}
```

**Usage**:
```go
// In usecase
func (u *ProductUsecase) GetByID(ctx context.Context, id uint64) (*Product, *_err.ErrorDTO) {
    product, err := u.Repo.GetByID(ctx, id)
    if err != nil {
        return nil, &_err.ErrorDTO{
            Code:    404,
            Message: "Không tìm thấy sản phẩm",
        }
    }
    return product, nil
}

// In handler
product, errDTO := u.usecase.GetByID(ctx, id)
if errDTO != nil {
    return ctx.JSON(errDTO.Code, errDTO)
}
return ctx.JSON(200, product)
```

---

## 6️⃣ Code Generator Module (`domain/usecase/`)

**File**: `usecase/code_data_usecase.go`

**Purpose**: Generate display codes for entities

```go
type CodeDataUsecase struct{}

// Product codes: SP000001, SP000002, ...
func (u *CodeDataUsecase) GetProductCode() string {
    return fmt.Sprintf("SP%06d", generateSequence())
}

func (u *CodeDataUsecase) GetProductCodePtr() *string {
    code := u.GetProductCode()
    return &code
}

// Asset codes: AS000001, AS000002, ...
func (u *CodeDataUsecase) GetAssetCode() string {
    return fmt.Sprintf("AS%06d", generateSequence())
}

func (u *CodeDataUsecase) GetAssetCodePtr() *string {
    code := u.GetAssetCode()
    return &code
}

// Post codes: PO000001, PO000002, ...
func (u *CodeDataUsecase) GetPostCode() string {
    return fmt.Sprintf("PO%06d", generateSequence())
}

// Batch generation
func (u *CodeDataUsecase) GetProductCodes(count int) []string {
    codes := make([]string, count)
    for i := 0; i < count; i++ {
        codes[i] = u.GetProductCode()
    }
    return codes
}
```

**Usage**:
```go
// In mapper
func ToProductDTO(product *Product, codeGen *CodeDataUsecase) *ProductDTO {
    return &ProductDTO{
        ID:   product.ID,
        Name: product.Name,
        Code: codeGen.GetProductCode(), // SP000123
    }
}
```

---

## 7️⃣ Provider Module (`provider/`)

### CrudRepo Base Implementation

**File**: `provider/repository.go`

**Purpose**: Base repository implementation with transaction support

```go
type CrudRepo[T any] struct {
    db    *_db.TransactionRepo
    self  ICrudRepo[T]
    table string
}

func (r *CrudRepo[T]) Init(self ICrudRepo[T], db *_db.TransactionRepo) {
    r.self = self
    r.db = db
}

func (r *CrudRepo[T]) Create(c context.Context, entity *T) error {
    return r.db.WithTx(c, func(tx *gorm.DB) error {
        // Call BeforeSave hook
        if err := r.self.BeforeSave(c, nil, entity); err != nil {
            return err
        }
        
        // Insert into database
        if err := tx.Create(entity).Error; err != nil {
            return err
        }
        
        // Call AfterSave hook
        if err := r.self.AfterSave(c, nil, entity); err != nil {
            return err
        }
        
        return nil
    })
}

func (r *CrudRepo[T]) Update(c context.Context, id uint64, entity *T) error {
    return r.db.WithTx(c, func(tx *gorm.DB) error {
        // Call BeforeSave hook
        if err := r.self.BeforeSave(c, &id, entity); err != nil {
            return err
        }
        
        // Update in database
        if err := tx.Model(entity).Where("id = ?", id).Updates(entity).Error; err != nil {
            return err
        }
        
        // Call AfterSave hook
        if err := r.self.AfterSave(c, &id, entity); err != nil {
            return err
        }
        
        return nil
    })
}

func (r *CrudRepo[T]) Delete(c context.Context, id uint64) error {
    return r.db.WithTx(c, func(tx *gorm.DB) error {
        var entity T
        return tx.Where("id = ?", id).Delete(&entity).Error
    })
}

func (r *CrudRepo[T]) GetList(c context.Context, pagable _dto.IPagable) ([]T, int64, error) {
    var entities []T
    var total int64
    
    db := r.db.GetDB()
    
    // Count total
    if err := db.Model(new(T)).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Get page
    offset := pagable.GetOffset()
    limit := pagable.GetLimit()
    
    if err := db.Where("deleted_at IS NULL").
        Offset(offset).
        Limit(limit).
        Order("created_at DESC").
        Find(&entities).Error; err != nil {
        return nil, 0, err
    }
    
    return entities, total, nil
}

// BeforeSave & AfterSave default implementations (can be overridden)
func (r *CrudRepo[T]) BeforeSave(c context.Context, id *uint64, entity *T) error {
    return nil
}

func (r *CrudRepo[T]) AfterSave(c context.Context, id *uint64, entity *T) error {
    return nil
}
```

**Usage in Infrastructure Layer**:
```go
// In infra/postgre/product_postgres.go
type ProductRepo struct {
    _provider.CrudRepo[Product]
}

func NewProductRepo(db *_db.TransactionRepo) _repo.IProductRepo {
    repo := &ProductRepo{}
    repo.Init(repo, db)
    return repo
}

// Override hooks if needed
func (r *ProductRepo) BeforeSave(c context.Context, id *uint64, entity *Product) error {
    // Custom validation
    if entity.Name == "" {
        return errors.New("name is required")
    }
    return nil
}

func (r *ProductRepo) AfterSave(c context.Context, id *uint64, entity *Product) error {
    // Clear cache
    cacheKey := fmt.Sprintf("product:%d", entity.ID)
    redis.Del(cacheKey)
    return nil
}
```

---

## 📊 Usage Statistics

### Used by Services:
- ✅ Auth Service
- ✅ User Service
- ✅ Organization Service
- ✅ BDSPro Service
- ✅ CRM Service
- ✅ Payment Service
- ✅ Notification Service
- ✅ Social Service
- ✅ Chat Service
- ✅ All other services

### Components Count:
- **Base Entities**: 2 (BaseEntity, AuditBase)
- **DTOs**: 3 (Pagable, HistoryDTO, NotificationDTO)
- **Enums**: 4 modules
- **CRUD Interfaces**: 2 (ICrudRepo, IBaseUsecase)
- **Providers**: 1 (CrudRepo)
- **Code Generators**: 1 (CodeDataUsecase)

---

## 🎯 Best Practices

### 1. Always Embed BaseEntity
```go
// ✅ GOOD
type Product struct {
    _models.BaseEntity
    Name string
}

// ❌ BAD - Missing audit trail
type Product struct {
    ID   uint64
    Name string
}
```

### 2. Use CrudRepo in Infrastructure
```go
// ✅ GOOD - In infra layer
type ProductRepo struct {
    _provider.CrudRepo[Product]
}

// ❌ BAD - Implement from scratch
type ProductRepo struct {
    db *gorm.DB
}
```

### 3. Always Return ErrorDTO from Usecase
```go
// ✅ GOOD
func (u *ProductUsecase) GetByID(ctx context.Context, id uint64) (*Product, *_err.ErrorDTO) {
    // ...
}

// ❌ BAD
func (u *ProductUsecase) GetByID(ctx context.Context, id uint64) (*Product, error) {
    // ...
}
```

### 4. Embed Pagable in List DTOs
```go
// ✅ GOOD
type GetProductsRequest struct {
    _dto.Pagable
    Text   string
    Status int
}

// ❌ BAD
type GetProductsRequest struct {
    Page   int
    Size   int
    Text   string
    Status int
}
```

---

**Last Updated**: October 15, 2025
**Module Version**: v1.0
**Status**: ✅ Stable & Production Ready

