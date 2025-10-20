# Quick Reference Guide - BDSPro Microservices

**Last Updated**: October 18, 2025  
**Purpose**: Fast lookup reference for development

---

## 🚀 TẠO API MỚI - 7 BƯỚC

```go
1. Proto Definition: shared/protobuf/schema/<service>/<name>.proto
2. Generate: make buf-<service>
3. Handler: infra/handler/<name>_handler.go + Swagger + @bind
4. Usecase: internal/usecase/<name>_usecase.go
5. Interface: internal/interface/repo/<name>_repo.go
6. Repository: infra/postgre/<name>_postgres.go (embed _provider.CrudRepo[T])
7. Mapper: infra/mapper/<name>_mapper.go
```

**QUAN TRỌNG**:
- ❌ **KHÔNG** dùng `gin.Context` → ✅ Dùng `context.Context`
- ❌ **KHÔNG** viết test/wire → Sẽ gen sau
- ✅ **BẮT BUỘC**: `repo.Init(repo, db)` trong constructor
- ✅ Mapper dùng: `_utils.FormatTimeToString`, `_utils.ParseStringToTime`

---

## 👤 LẤY THÔNG TIN USER

```go
import _utils "common/utils"

// Lấy profileId
profileID := _utils.GetProfileIdWithContext(ctx)

// Lấy organizationId
orgID := _utils.GetOrganizationIdFromContext(ctx)
```

---

## 📦 IMPORT ALIASES (Chuẩn)

```go
import (
    _dto "common/domain/dto"
    _models "common/domain/entity"
    _enum "common/domain/enum"
    _err "common/domain/err"
    _crud "common/domain/crud"
    _utils "common/utils"
    _provider "common/provider"
    _db "common/db"
)
```

---

## 🗄️ REPOSITORY PATTERN

### ✅ Cách ĐÚNG (Infra Layer)

```go
// infra/postgre/product_postgres.go
type ProductRepo struct {
    _provider.CrudRepo[Product]  // ✅ EMBED base repo
}

func NewProductRepo(db *_db.TransactionRepo) _repo.IProductRepo {
    repo := &ProductRepo{}
    repo.Init(repo, db)  // ✅ BẮT BUỘC gọi Init
    return repo
}

// Override khi cần custom logic
func (r *ProductRepo) BeforeSave(ctx context.Context, id *uint64, entity *Product) error {
    // Validation logic
    return nil
}

func (r *ProductRepo) AfterSave(ctx context.Context, id *uint64, entity *Product) error {
    // Cache update, audit log, etc.
    return nil
}
```

### ❌ Cách SAI (Không dùng)

```go
// ❌ KHÔNG implement từ đầu
type ProductRepo struct {
    db *gorm.DB
}

func (r *ProductRepo) Create(entity *Product) error {
    // Manual implementation - KHÔNG làm thế này!
}
```

---

## 🎯 COMMON DTOs

### Pagable (Pagination)

```go
// Request DTO
type ProductListRequest struct {
    _dto.Pagable  // ✅ Embed
    Name string
}

// Usage
offset := req.GetOffset()  // (Page - 1) * Size
limit := req.GetLimit()    // Min(Size, 100)

// Response
type ProductListResponse struct {
    Data  []Product `json:"data"`
    Total int64     `json:"total"`
}
```

### ErrorDTO

```go
// Return từ usecase
func (u *ProductUsecase) Create(product *Product) (*Product, *_err.ErrorDTO) {
    if err := validate(product); err != nil {
        return nil, &_err.ErrorDTO{
            Code: 400,
            Message: "Thông tin sản phẩm không hợp lệ",
        }
    }
    return created, nil
}
```

### HistoryDTO

```go
// Log history action
historyDTO := &_dto.HistoryDTO{
    Title:      "Tạo sản phẩm",
    Note:       []string{"Thêm sản phẩm mới"},
    TargetId:   productID,
    TargetType: _enum.ETargetProduct,
    ActionType: _enum.EHistoryCreateProduct,
    OwnerID:    profileID,
    OwnerOf:    _enum.EOwnerOfMember,
}
```

---

## 🔢 CODE GENERATION

```go
import _usecase "common/domain/usecase"

// Product code: SP000001, SP000002...
code := _usecase.GetProductCode(lastProductID)

// Asset code: AS000001, AS000002...
code := _usecase.GetAssetCode(lastAssetID)

// Post code: PO000001, PO000002...
code := _usecase.GetPostCode(lastPostID)
```

---

## 🌐 SERVICES PORTS & RESPONSIBILITIES

| Service | Port | Responsibility |
|---------|------|----------------|
| **Gateway** | 8080 | API Gateway, JWT auth, routing |
| **Auth** | 50062 | Authentication, OAuth, RBAC |
| **User** | 50051 | Profile, ratings, reviews |
| **Organization** | 50052 | Org hierarchy, teams, deals |
| **BDSPro** | 50053 | Property, product, asset, posts |
| **CRM** | 50054 | Contact, lead, pipeline |
| **Payment** | 50055 | Wallet, transactions |
| **Notification** | 50056 | FCM, email, SMS, in-app |
| **Chat** | 50057 | Real-time messaging |
| **Social** | 50058 | News feed, like, comment |
| **TQD** | 50065 | POI, amenities, directories |
| **Assistant** | 50061 | AI DeepSeek, text analysis |
| **Membership** | 50063 | Subscription plans, limits |
| **Relay** | 50064 | WebSocket broadcast |
| **File** | 8083 | File upload/storage |

---

## 🔐 AUTHENTICATION FLOW

```
1. User Login → Auth Service
2. Auth generates JWT
3. Client stores token
4. Request with Authorization: Bearer <token>
5. Gateway validates JWT
6. Gateway injects profileId, organizationId vào gRPC metadata
7. Service extracts từ context
```

**JWT Claims**:
- `profileId`: User ID
- `organizationId`: Current org ID
- `role`: User role
- `permissions`: User permissions

---

## 🔄 INTER-SERVICE COMMUNICATION

### Cách gọi Service khác

```go
// 1. Định nghĩa interface trong internal/interface/provider
type AssistantProvider interface {
    AnalyzeProductText(ctx context.Context, text string) (*Result, error)
}

// 2. Implement client trong infra/client
// @bind: internal/interface/provider.AssistantProvider
type AssistantClient struct {
    grpcClient assistantpb.AssistantServiceClient
}

// 3. Inject vào usecase
type ProductUsecase struct {
    assistantClient provider.AssistantProvider
}

// 4. Sử dụng
result := u.assistantClient.AnalyzeProductText(ctx, text)
```

---

## 🗂️ FOLDER STRUCTURE

```
<service>/
├── cmd/           # Entry points (grpc.go, http.go)
├── config/        # Config files (*.yml, properties.go)
├── env/           # Environment loader
├── infra/
│   ├── client/    # gRPC clients to other services
│   ├── handler/   # gRPC/HTTP handlers (implement proto)
│   ├── mapper/    # Entity ↔ DTO conversion
│   ├── postgre/   # Repository implementations (embed CrudRepo)
│   ├── provider/  # External service providers
│   └── validator/ # Request validation
├── internal/
│   ├── domain/    # Entities (embed BaseEntity)
│   ├── dto/       # Data Transfer Objects
│   ├── enums/     # Enumerations
│   ├── interface/
│   │   ├── provider/ # Provider interfaces
│   │   └── repo/     # Repository interfaces
│   ├── usecase/   # Business logic
│   └── utils/     # Utility functions
├── initial/       # Service initialization
├── wire/          # Dependency injection (wire.go, wire_gen.go)
└── main.go        # Main entry
```

---

## 📊 SPECIAL SERVICES QUICK REF

### Payment Service - Transaction Types

```go
const (
    DEPOSIT       = "DEPOSIT"       // Nạp tiền
    WITHDRAW      = "WITHDRAW"      // Rút tiền
    PAYMENT       = "PAYMENT"       // Thanh toán
    REFUND        = "REFUND"        // Hoàn tiền
    TRANSFER_IN   = "TRANSFER_IN"   // Nhận chuyển khoản
    TRANSFER_OUT  = "TRANSFER_OUT"  // Chuyển khoản đi
    COMMISSION    = "COMMISSION"    // Hoa hồng
)

// Transaction Code Format
code := fmt.Sprintf("TRF%d%d", time.Now().UnixNano(), walletID)
code := fmt.Sprintf("DEP%d%d", time.Now().UnixNano(), walletID)
code := fmt.Sprintf("PAY%d%d", time.Now().UnixNano(), walletID)

// Row-Level Locking
wallet := repo.GetWalletByIdForUpdate(tx, walletID)  // SELECT FOR UPDATE
```

### Assistant Service - AI Functions

```go
// Analyze real estate text
POST /v2/assistant/analyze/product
{
  "content": "Nhà 3 tầng 80m2 Cầu Giấy giá 5 tỷ"
}

// Chat
POST /v2/assistant/chat
{
  "message": "Giá nhà Cầu Giấy như thế nào?",
  "history": [...]
}

// Generate content
POST /v2/assistant/generate
{
  "template": "product_description",
  "variables": {...}
}
```

### Membership Service - Plan Limits

```go
type Plan struct {
    CreateProduct      int  // Max products
    LimitPost          int  // Max posts
    CreateProductAI    int  // AI-generated products
    FriendLimit        int  // Max friends
    CustomerLimit      int  // Max customers
    JoinOrg            int  // Max orgs to join
    CreateOrganization int  // Can create org (0/1)
    OrgMemberLimit     int  // Max members in org
    OrgProductLimit    int  // Max products in org
}

// Unlimited = -1
```

### Notification Service - Types

```go
const (
    FriendRequest    = 10   // Lời mời kết bạn
    FriendResponse   = 20   // Phản hồi kết bạn
    Follow           = 30   // Theo dõi
    CommentNewsFeed  = 40   // Comment bài viết
    LikeNewsFeed     = 50   // Like bài viết
    LikeComment      = 60   // Like comment
    ShareNewsFeed    = 70   // Chia sẻ bài viết
    JoinGroup        = 80   // Tham gia nhóm
    JoinOrganization = 90   // Tham gia tổ chức
    Deposit          = 110  // Nạp tiền
    Withdraw         = 111  // Rút tiền
    Payment          = 112  // Thanh toán
)

// Send FCM notification
SendToToken(token, title, body, image)
SendToTopic(topic, title, body, image)
```

### Social Service - Interactions

```go
// NewsFeed Entity
type NewsFeed struct {
    Title       string
    Content     string
    Image       string
    NumLike     int
    NumComment  int
    NumShare    int
    NumView     int
    IsReel      bool
    ParentID    *uint64  // For comments
    Visibility  Visibility  // PUBLIC, FRIENDS, PRIVATE
    OwnerOf     OwnerOf     // USER, GROUP, ORGANIZATION
    OwnerID     *uint64
}

// Like/Unlike
POST /v2/social/like/news-feed/{id}
DELETE /v2/social/like/news-feed/{id}

// Comment (ParentID for nested replies)
POST /v2/social/comment
{
  "news_feed_id": 123,
  "content": "Nice!",
  "parent_id": 456  // Optional, for reply
}
```

### Chat Service - WebSocket

```go
// Connect
ws://localhost:8082/ws?token=<jwt>

// Redis Pub/Sub Messages
type RedisMessage struct {
    Type string  // "MESSAGE", "CREATE_ROOM", "JOIN_ROOM", etc.
    Data any
}

// Message Types
const (
    MESSAGE_TYPE                    = "MESSAGE"
    CREATE_ROOM_TYPE               = "CREATE_ROOM"
    JOIN_ROOM_TYPE                 = "JOIN_ROOM"
    LEAVE_ROOM_TYPE                = "LEAVE_ROOM"
    RECALL_MESSAGE_TYPE            = "RECALL_MESSAGE"
    EDIT_MESSAGE_TYPE              = "EDIT_MESSAGE"
    PIN_MESSAGE_TYPE               = "PIN_MESSAGE"
    REACT_MESSAGE_TYPE             = "REACT_MESSAGE"
    DELETE_VIOLATION_MESSAGE_TYPE  = "DELETE_VIOLATION_MESSAGE"
)
```

### TQD Service - Points of Interest

```go
type POI struct {
    Code           string    // Unique code
    CategoryID     uint64
    Name           string
    Address        string
    Lat, Lng       float64   // Coordinates
    GeoJSON        string    // Geographic data
    Phone          string
    Website        string
    RatingPoint    float64   // 0-5
    RatingCount    int
    Source         string    // Data source
    Confidence     float64   // Data quality (0-1)
    LastVerifiedAt *time.Time
    NextReviewAt   *time.Time
}
```

### Marketing Service - Campaigns

```go
// Campaign Types
const (
    CampaignTypeVIP       = "VIP"
    CampaignTypeMetaAds   = "META_ADS"
    CampaignTypeGoogleAds = "GOOGLE_ADS"
    CampaignTypePR        = "PR"
)

// Campaign Status
const (
    CampaignStatusDraft     = "DRAFT"
    CampaignStatusActive    = "ACTIVE"
    CampaignStatusPaused    = "PAUSED"
    CampaignStatusCompleted = "COMPLETED"
    CampaignStatusCancelled = "CANCELLED"
)

// Create campaign → auto creates wallet
// Topup budget → payment transaction
// Track views, clicks, conversions
```

### Transaction Service - Deal Workflow

```go
// Transaction Actions
const (
    EActionDeposite  = "DEPOSITE"   // Đặt cọc
    EActionSign      = "SIGN"       // Ký hợp đồng
    EActionPay       = "PAY"        // Thanh toán
    EActionHandOver  = "HAND_OVER"  // Bàn giao
    EActionCancel    = "CANCEL"     // Hủy
)

// Transaction Status
const (
    TxStatusPending = "PENDING"
    TxStatusDone    = "DONE"
    TxStatusCancel  = "CANCEL"
)

// Flow: Create → Action → Approve → Update → HandOver
```

---

## 🧩 COMMON PATTERNS

### Transaction Pattern

```go
err := u.Transaction.WithTransaction(ctx, func(ctx context.Context) error {
    // All DB operations here use same transaction
    entity1, err := u.repo1.Create(ctx, data1)
    if err != nil {
        return err  // Auto rollback
    }
    
    entity2, err := u.repo2.Create(ctx, data2)
    if err != nil {
        return err  // Auto rollback
    }
    
    return nil  // Auto commit
})
```

### Notification Pattern

```go
// Send notification async
go func() {
    u.notificationClient.SendToUser(ctx, userID, &NotificationDTO{
        Title:   "Thông báo",
        Message: "Có hoạt động mới",
        Type:    _enum.FriendRequest,
    })
}()
```

### History Logging Pattern

```go
// Log important actions
historyDTO := &_dto.HistoryDTO{
    Title:      "Tạo sản phẩm",
    TargetId:   productID,
    TargetType: _enum.ETargetProduct,
    ActionType: _enum.EHistoryCreateProduct,
    OwnerID:    profileID,
}
u.historyClient.CreateHistory(ctx, historyDTO)
```

---

## ⚡ PERFORMANCE TIPS

### Caching

```go
// Redis caching pattern
func (u *Usecase) GetProduct(id uint64) (*Product, error) {
    // Try cache first
    cached, err := u.redis.Get(ctx, fmt.Sprintf("product:%d", id))
    if err == nil {
        return cached, nil
    }
    
    // Cache miss, get from DB
    product, err := u.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Save to cache
    u.redis.Set(ctx, fmt.Sprintf("product:%d", id), product, 1*time.Hour)
    
    return product, nil
}
```

### Batch Loading

```go
// Load multiple entities at once
products := u.productRepo.GetListByIDs(ctx, productIDs)  // Single query

// Instead of N+1 queries
for _, id := range productIDs {
    product := u.productRepo.GetByID(ctx, id)  // ❌ N queries
}
```

### Pagination

```go
// Always paginate large lists
req := &dto.ProductListRequest{
    Pagable: _dto.Pagable{
        Page: 1,
        Size: 20,  // Max 100
    },
}
products, total := u.productRepo.GetList(ctx, req)
```

---

## 🐛 COMMON ERRORS & FIXES

### Error: "failed to initialize repo"
```go
// ❌ Missing Init call
repo := &ProductRepo{}
return repo  // ERROR!

// ✅ Correct
repo := &ProductRepo{}
repo.Init(repo, db)  // Must call Init
return repo
```

### Error: "cannot use gin.Context"
```go
// ❌ Wrong context type
func (u *Usecase) Create(c *gin.Context, product *Product) error

// ✅ Correct
func (u *Usecase) Create(ctx context.Context, product *Product) error
```

### Error: "profileId not found in context"
```go
// ❌ Wrong function
profileID := ctx.Value("profileId")

// ✅ Correct
profileID := _utils.GetProfileIdWithContext(ctx)
```

### Error: "transaction deadlock"
```go
// ❌ Not using row-level locking
wallet := repo.GetByID(ctx, walletID)
wallet.Balance -= amount
repo.Update(ctx, wallet)  // Race condition!

// ✅ Correct - use FOR UPDATE
tx := db.BeginTx(ctx)
wallet := repo.GetWalletByIdForUpdate(tx, walletID)  // Locked
wallet.Balance -= amount
repo.Update(tx, wallet)
tx.Commit()
```

---

## 📝 CODE TEMPLATES

### Handler Template

```go
// @Tags Product
// @Summary Create product
// @Description Create a new product
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body pb.CreateProductRequest true "Product data"
// @Success 200 {object} pb.ProductResponse
// @Router /v2/product [post]
// @bind: internal/usecase.ProductUsecase
func (h *ProductHandler) CreateProduct(
    ctx context.Context,
    req *pb.CreateProductRequest,
) (*pb.ProductResponse, error) {
    // Validation
    if err := h.validator.ValidateCreate(req); err != nil {
        return nil, status.Error(codes.InvalidArgument, err.Error())
    }
    
    // Map DTO to Entity
    entity := h.mapper.ToEntity(req)
    
    // Call usecase
    created, errDTO := h.usecase.Create(ctx, entity)
    if errDTO != nil {
        return nil, status.Error(codes.Internal, errDTO.Message)
    }
    
    // Map Entity to Response
    return h.mapper.ToResponse(created), nil
}
```

### Usecase Template

```go
type ProductUsecase struct {
    _crud.BaseUsecase[Product, _crud.ICrudRepo[Product]]
    repo          repo.ProductRepo
    validator     validator.ProductValidator
    // other dependencies...
}

func NewProductUsecase(repo repo.ProductRepo) *ProductUsecase {
    return &ProductUsecase{
        BaseUsecase: _crud.BaseUsecase[Product, _crud.ICrudRepo[Product]]{
            Repo: repo,
        },
        repo: repo,
    }
}

func (u *ProductUsecase) Create(ctx context.Context, product *Product) (*Product, *_err.ErrorDTO) {
    // Validation
    if err := u.validator.Validate(product); err != nil {
        return nil, &_err.ErrorDTO{Code: 400, Message: err.Error()}
    }
    
    // Business logic
    product.Code = _usecase.GetProductCode(lastID)
    
    // Save
    created, err := u.repo.Create(ctx, product)
    if err != nil {
        return nil, &_err.ErrorDTO{Code: 500, Message: "Tạo sản phẩm thất bại"}
    }
    
    return created, nil
}
```

### Repository Template

```go
// infra/postgre/product_postgres.go
type ProductRepo struct {
    _provider.CrudRepo[Product]
}

// @bind: internal/interface/repo.IProductRepo
func NewProductRepo(db *_db.TransactionRepo) repo.IProductRepo {
    r := &ProductRepo{}
    r.Init(r, db)
    return r
}

func (r *ProductRepo) BeforeSave(ctx context.Context, id *uint64, entity *Product) error {
    // Validation before save
    if entity.Name == "" {
        return errors.New("name is required")
    }
    return nil
}

func (r *ProductRepo) AfterSave(ctx context.Context, id *uint64, entity *Product) error {
    // Clear cache after save
    r.cache.Delete(fmt.Sprintf("product:%d", entity.ID))
    return nil
}
```

---

## 🔍 DEBUG CHECKLIST

Khi API không hoạt động, check theo thứ tự:

1. ✅ **Proto definition** đúng chưa?
2. ✅ **Generate protobuf** (`make buf-<service>`) chưa?
3. ✅ **Handler** có `@bind` comment chưa?
4. ✅ **Repository** gọi `Init()` chưa?
5. ✅ **Context.Context** (không phải gin.Context)?
6. ✅ **ProfileID** lấy đúng cách (_utils.GetProfileIdWithContext)?
7. ✅ **Wire generated** (`make wire <service>`) chưa?
8. ✅ **Swagger comments** đầy đủ chưa?
9. ✅ **Mapper** handle time fields đúng chưa?
10. ✅ **Gateway** đăng ký route chưa?

---

## 📚 USEFUL COMMANDS

```bash
# Generate protobuf
make buf-<service>

# Generate wire
make wire <service>

# Generate swagger
make swagger-<service>

# Run service
make <service>-grpc

# Run all services
make start-all

# Build service
cd <service> && go build -o main .

# Hot reload (with Air)
cd <service> && air
```

---

**Pro Tip**: Bookmark file này và Cmd+F để tìm nhanh! 🚀

