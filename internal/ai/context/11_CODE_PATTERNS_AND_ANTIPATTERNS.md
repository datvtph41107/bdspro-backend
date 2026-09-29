# Code Patterns & Anti-Patterns

## 📋 Overview

**Purpose**: Identify good practices and common mistakes in the codebase
**Format**: Pattern → Example → Explanation
**Status**: ✅ Complete

---

## ✅ Good Patterns (DO THIS)

### Pattern 1: Use BaseEntity for All Domain Models

**Why**: Provides consistent ID, timestamps, soft delete, audit trail

```go
// ✅ GOOD
type Product struct {
    _models.BaseEntity              // Includes ID, CreatedAt, UpdatedAt, DeletedAt, CreatedBy, UpdatedBy
    Name        string
    Price       float64
}

// ❌ BAD
type Product struct {
    ID          uint64
    Name        string
    Price       float64
    CreatedAt   time.Time
    UpdatedAt   time.Time
    // Missing: DeletedAt, CreatedBy, UpdatedBy
}
```

**Benefits**:
- Automatic timestamps
- Soft delete support
- Audit trail tracking
- Consistent across all services

---

### Pattern 2: Use CrudRepo Provider in Infrastructure

**Why**: Reuse battle-tested CRUD implementation with transaction support

```go
// ✅ GOOD - Infra layer
type ProductRepo struct {
    _provider.CrudRepo[Product]
}

func NewProductRepo(db *_db.TransactionRepo) repo.IProductRepo {
    repo := &ProductRepo{}
    repo.Init(repo, db)
    return repo
}

// Override only when needed
func (r *ProductRepo) BeforeSave(c context.Context, id *uint64, entity *Product) error {
    // Custom validation
    if entity.Name == "" {
        return errors.New("name required")
    }
    return nil
}

// ❌ BAD - Implementing from scratch
type ProductRepo struct {
    db *gorm.DB
}

func (r *ProductRepo) Create(ctx context.Context, product *Product) error {
    // Manually implementing transactions, error handling, etc.
    // Missing BeforeSave/AfterSave hooks
    // No soft delete handling
    return r.db.Create(product).Error
}
```

**Benefits**:
- Transaction support built-in
- BeforeSave/AfterSave hooks
- Consistent error handling
- Soft delete handling
- Pagination support

---

### Pattern 3: Embed Pagable in List DTOs

**Why**: Consistent pagination across all list APIs

```go
// ✅ GOOD
type GetProductsRequest struct {
    _dto.Pagable                    // Includes Page, Size, Sort with built-in logic
    Text   string `json:"text" form:"text"`
    Status int    `json:"status" form:"status"`
}

// Usage
func (r *ProductRepo) GetList(ctx context.Context, req *GetProductsRequest) ([]Product, int64, error) {
    offset := req.GetOffset()       // Automatically calculated
    limit := req.GetLimit()         // Auto-validated (max 100)
    // ...
}

// ❌ BAD - Manual pagination
type GetProductsRequest struct {
    Page   int    `json:"page"`
    Size   int    `json:"size"`
    Text   string `json:"text"`
}

// Need to manually validate and calculate offset
func (r *ProductRepo) GetList(ctx context.Context, req *GetProductsRequest) ([]Product, int64, error) {
    if req.Size <= 0 {
        req.Size = 20
    }
    if req.Size > 100 {
        req.Size = 100
    }
    offset := (req.Page - 1) * req.Size
    // Duplicated logic everywhere
}
```

**Benefits**:
- Default page size (20)
- Max page size enforcement (100)
- Automatic offset calculation
- Consistent across all APIs

---

### Pattern 4: Return ErrorDTO from Usecases

**Why**: Structured error responses with HTTP codes

```go
// ✅ GOOD
func (u *ProductUsecase) GetByID(ctx context.Context, id uint64) (*Product, *_err.ErrorDTO) {
    product, err := u.repo.GetByID(ctx, id)
    if err != nil {
        return nil, &_err.ErrorDTO{
            Code:    404,
            Message: "Không tìm thấy sản phẩm",
        }
    }
    return product, nil
}

// Handler can map to HTTP status easily
product, errDTO := u.usecase.GetByID(ctx, id)
if errDTO != nil {
    return ctx.JSON(errDTO.Code, errDTO)
}

// ❌ BAD - Generic errors
func (u *ProductUsecase) GetByID(ctx context.Context, id uint64) (*Product, error) {
    product, err := u.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err  // Lost context, no HTTP code
    }
    return product, nil
}

// Handler doesn't know what HTTP status to use
product, err := u.usecase.GetByID(ctx, id)
if err != nil {
    // What status code? 400? 404? 500?
    return ctx.JSON(500, err.Error())
}
```

**Benefits**:
- Clear HTTP status codes
- Consistent error messages
- Easy to translate to different protocols (gRPC, HTTP)

---

### Pattern 5: Always Use Context

**Why**: Propagate request-scoped values (user ID, org ID, tracing)

```go
// ✅ GOOD
func (u *ProductUsecase) Create(ctx context.Context, product *Product) (*Product, *_err.ErrorDTO) {
    profileId := _utils.GetProfileIdWithContext(ctx)
    orgId := _utils.GetOrganizationIdFromContext(ctx)
    
    product.OwnerId = profileId
    product.OrganizationId = orgId
    
    if err := u.repo.Create(ctx, product); err != nil {
        return nil, &_err.ErrorDTO{Code: 500, Message: "Không thể tạo sản phẩm"}
    }
    
    return product, nil
}

// ❌ BAD - No context
func (u *ProductUsecase) Create(product *Product) (*Product, error) {
    // How to get current user?
    // How to pass timeout?
    // How to do request tracing?
    
    if err := u.repo.Create(product); err != nil {
        return nil, err
    }
    return product, nil
}
```

**Benefits**:
- Access to user info
- Request timeout propagation
- Distributed tracing
- Cancellation support

---

### Pattern 6: Use @bind Comments for Wire

**Why**: Auto-discovery for dependency injection

```go
// ✅ GOOD
// @bind: internal/interface/repo
type ProductRepo struct {
    _provider.CrudRepo[Product]
}

// @bind: internal/usecase
type ProductUsecase struct {
    repo repo.IProductRepo
}

// @bind: infra/handler
type ProductHandler struct {
    usecase *usecase.ProductUsecase
}

// Wire script automatically finds and wires these

// ❌ BAD - Manual wire.go maintenance
// Need to manually add to wire.go every time:
// wire.Build(
//     NewProductRepo,
//     NewProductUsecase,
//     NewProductHandler,
// )
```

**Benefits**:
- Auto-discovery
- Less manual work
- Harder to forget

---

### Pattern 7: Interface in internal/interface, Implementation in infra

**Why**: Dependency inversion, testability

```go
// ✅ GOOD

// internal/interface/repo/product_repo.go
type IProductRepo interface {
    _crud.ICrudRepo[Product]
    GetByCode(ctx context.Context, code string) (*Product, error)
}

// infra/postgre/product_postgres.go
type ProductRepo struct {
    _provider.CrudRepo[Product]
}

func (r *ProductRepo) GetByCode(ctx context.Context, code string) (*Product, error) {
    var product Product
    err := r.db.GetDB().Where("code = ?", code).First(&product).Error
    return &product, err
}

// Usecase depends on interface, not implementation
type ProductUsecase struct {
    repo repo.IProductRepo  // Interface, not concrete type
}

// ❌ BAD - Direct dependency on implementation
type ProductUsecase struct {
    repo *postgre.ProductRepo  // Concrete type, hard to test
}
```

**Benefits**:
- Easy to mock for testing
- Can swap implementations
- Clean architecture compliance

---

### Pattern 8: Use Code Generators for Display Codes

**Why**: Consistent code format across entities

```go
// ✅ GOOD
func (u *ProductUsecase) Create(ctx context.Context, product *Product) (*Product, *_err.ErrorDTO) {
    codeGen := _usecase.CodeDataUsecase{}
    product.Code = codeGen.GetProductCode()  // SP000001
    
    if err := u.repo.Create(ctx, product); err != nil {
        return nil, &_err.ErrorDTO{Code: 500, Message: "Không thể tạo"}
    }
    return product, nil
}

// ❌ BAD - Manual code generation
func (u *ProductUsecase) Create(ctx context.Context, product *Product) (*Product, error) {
    // Different format in different places
    count, _ := u.repo.Count(ctx)
    product.Code = fmt.Sprintf("PROD%d", count+1)  // PROD1, PROD2 - inconsistent
    
    if err := u.repo.Create(ctx, product); err != nil {
        return nil, err
    }
    return product, nil
}
```

**Benefits**:
- Consistent format (SP000001, AS000001, PO000001)
- No duplicates
- Easy to change format globally

---

## ❌ Anti-Patterns (DON'T DO THIS)

### Anti-Pattern 1: Business Logic in Handlers

**Problem**: Handler should only handle HTTP/gRPC concerns

```go
// ❌ BAD
func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
    // Business logic in handler!
    if req.Price < 0 {
        return nil, errors.New("price must be positive")
    }
    
    // Database access in handler!
    product := &Product{
        Name:  req.Name,
        Price: req.Price,
    }
    if err := h.db.Create(product).Error; err != nil {
        return nil, err
    }
    
    // Notification logic in handler!
    h.notificationClient.Send(ctx, userId, "Product created")
    
    return &pb.CreateProductResponse{Id: product.ID}, nil
}

// ✅ GOOD
func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
    // Handler only maps request/response
    product := &Product{
        Name:  req.Name,
        Price: req.Price,
    }
    
    // Delegate to usecase
    created, errDTO := h.usecase.Create(ctx, product)
    if errDTO != nil {
        return nil, status.Error(codes.Internal, errDTO.Message)
    }
    
    return &pb.CreateProductResponse{Id: created.ID}, nil
}

// Usecase contains business logic
func (u *ProductUsecase) Create(ctx context.Context, product *Product) (*Product, *_err.ErrorDTO) {
    // Validation
    if product.Price < 0 {
        return nil, &_err.ErrorDTO{Code: 400, Message: "Giá phải dương"}
    }
    
    // Database
    if err := u.repo.Create(ctx, product); err != nil {
        return nil, &_err.ErrorDTO{Code: 500, Message: "Không thể tạo"}
    }
    
    // Notifications
    u.notificationClient.Send(ctx, product.OwnerId, "Sản phẩm đã được tạo")
    
    return product, nil
}
```

---

### Anti-Pattern 2: Using gin.Context in Domain/Usecase

**Problem**: Tight coupling to HTTP framework

```go
// ❌ BAD
func (u *ProductUsecase) Create(c *gin.Context, product *Product) error {
    // Coupled to Gin!
    userId := c.GetUint64("userId")
    product.OwnerId = userId
    return u.repo.Create(product)
}

// ✅ GOOD
func (u *ProductUsecase) Create(ctx context.Context, product *Product) (*Product, *_err.ErrorDTO) {
    // Framework-agnostic
    userId := _utils.GetProfileIdWithContext(ctx)
    product.OwnerId = userId
    
    if err := u.repo.Create(ctx, product); err != nil {
        return nil, &_err.ErrorDTO{Code: 500, Message: "Không thể tạo"}
    }
    return product, nil
}
```

**Why Bad**: Can't use usecase in gRPC handlers, CLI tools, or tests easily

---

### Anti-Pattern 3: Ignoring Errors

**Problem**: Silent failures are hard to debug

```go
// ❌ BAD
product, _ := u.repo.GetByID(ctx, id)  // Error ignored!
// product might be nil!
return product, nil

// ❌ BAD
_ = u.notificationClient.Send(ctx, userId, "Message")  // Notification might fail silently

// ✅ GOOD
product, err := u.repo.GetByID(ctx, id)
if err != nil {
    return nil, &_err.ErrorDTO{Code: 404, Message: "Không tìm thấy"}
}
return product, nil

// ✅ GOOD
if err := u.notificationClient.Send(ctx, userId, "Message"); err != nil {
    log.Printf("Failed to send notification: %v", err)
    // Continue or return error based on criticality
}
```

---

### Anti-Pattern 4: N+1 Queries

**Problem**: Performance nightmare

```go
// ❌ BAD - N+1 queries
products, _ := u.repo.GetAll(ctx)
for _, product := range products {
    // Query for each product!
    media, _ := u.mediaRepo.GetByProductId(ctx, product.ID)
    product.Media = media
    
    owner, _ := u.userClient.GetProfile(ctx, product.OwnerId)
    product.Owner = owner
}

// ✅ GOOD - Single query with joins/preload
products, _ := u.repo.GetAllWithRelations(ctx)
// Or use GORM Preload
db.Preload("Media").Preload("Owner").Find(&products)

// ✅ GOOD - Batch fetch
productIds := []uint64{}
for _, p := range products {
    productIds = append(productIds, p.ID)
}
mediaMap := u.mediaRepo.GetByProductIds(ctx, productIds)  // Single query
ownerMap := u.userClient.GetProfilesByIds(ctx, ownerIds)  // Single call
```

---

### Anti-Pattern 5: Not Using Pagination

**Problem**: Memory issues, slow responses

```go
// ❌ BAD - Load everything
func (u *ProductUsecase) GetAll(ctx context.Context) ([]Product, error) {
    return u.repo.GetAll(ctx)  // Might return 100,000+ records!
}

// ✅ GOOD - Always paginate
func (u *ProductUsecase) GetList(ctx context.Context, req *GetProductsRequest) ([]Product, int64, *_err.ErrorDTO) {
    products, total, err := u.repo.GetList(ctx, req)
    if err != nil {
        return nil, 0, &_err.ErrorDTO{Code: 500, Message: "Không thể lấy danh sách"}
    }
    return products, total, nil
}

// Response includes total for frontend
{
    "data": [...],
    "total": 1000,
    "page": 1,
    "size": 20
}
```

---

### Anti-Pattern 6: Hardcoding Values

**Problem**: Hard to change, hard to test

```go
// ❌ BAD
if product.Status == 1 {  // What does 1 mean?
    // ...
}

if user.Role == "ADMIN" {  // Hardcoded string
    // ...
}

// ✅ GOOD - Use enums/constants
const (
    ProductStatusActive   = 1
    ProductStatusInactive = 2
    ProductStatusDeleted  = 3
)

if product.Status == ProductStatusActive {
    // Clear intent
}

const (
    RoleAdmin  = "ADMIN"
    RoleMember = "MEMBER"
)

if user.Role == RoleAdmin {
    // Refactorable
}
```

---

### Anti-Pattern 7: Fat Models

**Problem**: Domain entities shouldn't have business logic

```go
// ❌ BAD
type Product struct {
    _models.BaseEntity
    Name  string
    Price float64
}

func (p *Product) CalculateDiscount(percentage float64) float64 {
    // Business logic in entity
    return p.Price * (1 - percentage/100)
}

func (p *Product) SendNotification(client NotificationClient) {
    // External dependency in entity
    client.Send("Product created")
}

// ✅ GOOD - Business logic in usecase
type Product struct {
    _models.BaseEntity
    Name  string
    Price float64
}

func (u *ProductUsecase) CalculateDiscount(product *Product, percentage float64) float64 {
    return product.Price * (1 - percentage/100)
}

func (u *ProductUsecase) NotifyCreation(ctx context.Context, product *Product) error {
    return u.notificationClient.Send(ctx, product.OwnerId, "Product created")
}
```

---

### Anti-Pattern 8: Not Using Transactions

**Problem**: Data inconsistency

```go
// ❌ BAD - No transaction
func (u *OrderUsecase) Create(ctx context.Context, order *Order) error {
    // Create order
    u.orderRepo.Create(ctx, order)
    
    // Deduct inventory
    u.inventoryRepo.Deduct(ctx, order.ProductId, order.Quantity)  // Might fail!
    
    // Create payment
    u.paymentRepo.Create(ctx, payment)  // If this fails, inventory already deducted!
    
    return nil
}

// ✅ GOOD - Use transaction
func (u *OrderUsecase) Create(ctx context.Context, order *Order) error {
    return u.db.WithTx(ctx, func(tx *gorm.DB) error {
        // All or nothing
        if err := u.orderRepo.Create(ctx, order); err != nil {
            return err
        }
        
        if err := u.inventoryRepo.Deduct(ctx, order.ProductId, order.Quantity); err != nil {
            return err  // Rollback order
        }
        
        if err := u.paymentRepo.Create(ctx, payment); err != nil {
            return err  // Rollback order and inventory
        }
        
        return nil  // Commit all
    })
}
```

---

## 🎯 Summary Checklist

### Before Committing Code

- [ ] All entities use `BaseEntity`
- [ ] All repositories use `CrudRepo` provider
- [ ] All list DTOs embed `Pagable`
- [ ] All usecases return `*_err.ErrorDTO`
- [ ] All functions accept `context.Context`
- [ ] All implementations have `@bind` comments
- [ ] No business logic in handlers
- [ ] No `gin.Context` in domain/usecase
- [ ] All errors are handled
- [ ] No N+1 queries
- [ ] Pagination used for lists
- [ ] No hardcoded values (use constants)
- [ ] Domain entities are anemic (no business logic)
- [ ] Transactions used for multi-step operations

---

**Last Updated**: October 15, 2025
**Version**: v1.0
**Status**: ✅ Complete

