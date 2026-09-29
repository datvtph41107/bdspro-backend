# BDSPro Service - Complete Context

## 📋 Service Overview

**Service Name**: BDSPro Service (Real Estate Management)
**Port**: 50053 (gRPC)
**Primary Role**: Real Estate Property & Asset Management
**Business Domain**: Real Estate Industry
**Status**: ✅ Production Ready

**BDSPro** = **B**ất **Đ**ộng **S**ản **Pro** (Vietnamese for "Real Estate Professional")

---

## 🎯 Core Responsibilities

### 1. Post Management (Tin đăng)
- **Post CRUD**: Create, update, delete property listings
- **Post Types**: Normal, Feed, Share
- **Post Status**: Active, Expired, Hidden
- **Post Media**: Images, videos, documents
- **Post Visibility**: Public, Private, Organization
- **Post Expiration**: Time-based expiration management
- **VIP Posts**: Premium listing features

### 2. Product Management (Sản phẩm BĐS)
- **Product CRUD**: Create, update, delete products
- **Product Types**: Land, House, Apartment, Commercial
- **Product Attributes**: Custom attributes per product type
- **Product Media**: Multi-media support
- **Product History**: Track all changes
- **Product Pricing**: Price management and history
- **Product Sharing**: Share with organization/team
- **Product Access Control**: Private field management

### 3. Asset Management (Tài sản)
- **Asset CRUD**: Comprehensive asset management
- **Asset Types**: Real estate assets
- **Asset Cost**: Track acquisition and operating costs
- **Asset Income**: Rental and other income tracking
- **Asset Legal**: Legal documentation
- **Asset Exploitation**: Usage and exploitation tracking
- **Asset Split/Merge**: Asset division and combination
- **Asset Sharing**: Organization-level asset sharing

### 4. Project Management (Dự án)
- **Project CRUD**: Real estate project management
- **Project Builds**: Building/tower management
- **Project Inventory**: Apartment/unit tracking
- **Project Progress**: Construction progress
- **Developer Info**: Developer details

### 5. Property Type Management
- **Type Hierarchy**: Land, House, Apartment, etc.
- **Type Attributes**: Custom attributes per type
- **Apartment Attributes**: Specific apartment features

### 6. Region & Location
- **Region Management**: Provinces, districts, wards
- **Location Data**: Geocoding integration
- **Address Parsing**: Smart address processing

### 7. Document Management
- **Document Types**: Contract, certificate, etc.
- **Cost Documents**: Expense documentation
- **Income Documents**: Revenue documentation

---

## 📁 Service Architecture

### Domain Layer (`internal/domain/`)

**Core Entities**:
```
- post.go                  # Property listing/post
- post_media.go           # Post images/videos
- post_user.go            # User-level posts
- post_organization.go    # Organization-level posts
- product.go              # Real estate product
- product_media.go        # Product images/videos
- product_user.go         # User-owned products
- product_organization.go # Organization products
- product_price.go        # Price history
- product_history.go      # Change history
- product_private.go      # Private field access
- asset.go                # Real estate asset
- asset_user.go           # User assets
- asset_organization.go   # Organization assets
- asset_share.go          # Shared assets
- project.go              # Real estate project
- project_build.go        # Building in project
- project_inventory.go    # Apartment units
```

**Support Entities**:
```
- property_type.go        # Land, House, Apartment, etc.
- amenity.go              # Facilities (pool, gym, etc.)
- apartment.go            # Apartment details
- apartment_attribute.go  # Apartment features
- attribute_product.go    # Custom product attributes
- attribute_value.go      # Attribute values
- developer.go            # Property developer
- region.go               # Geographic regions
- doc_type.go             # Document types
- cost_type.go            # Cost categories
```

**Asset Management**:
```
- asset_cost.go           # Asset costs
- asset_income_type.go    # Income types
- asset_legal.go          # Legal documents
- asset_exploitation.go   # Usage tracking
- cost_document.go        # Cost documentation
- income_document.go      # Income documentation
```

**Transaction**:
```
- transaction.go          # Property transactions
- deposite.go             # Deposit tracking
- customer.go             # Customer info
- payment_method.go       # Payment methods
```

---

## 🏠 Product Hierarchy

```
Product (Sản phẩm)
├── Type: Land (Đất)
├── Type: House (Nhà)
├── Type: Apartment (Căn hộ)
│   ├── Apartment Details
│   ├── Apartment Attributes
│   └── Building/Project Link
├── Type: Commercial (Thương mại)
└── Type: Industrial (Công nghiệp)

Each Product has:
├── Media (Images/Videos)
├── Attributes (Custom fields)
├── Price History
├── Location (Region)
├── Ownership (User/Organization)
└── Access Control (Private fields)
```

---

## 📰 Post vs Product Relationship

### Post (Tin đăng)
- **Purpose**: Advertise/list a property
- **Lifetime**: Temporary (can expire)
- **Visibility**: Can be hidden/shown
- **Link**: Links to a Product
- **Use Case**: Marketing, advertising

### Product (Sản phẩm)
- **Purpose**: Actual property data
- **Lifetime**: Permanent (soft delete only)
- **Ownership**: User or Organization
- **Use Case**: Property management

```
Post ──► Product ──► Asset
(Ad)     (Property)   (Real Asset)
```

**Flow**:
1. Create Product (property details)
2. Create Post (advertise the product)
3. Post can expire, but Product remains
4. Optionally link to Asset (real asset management)

---

## 🗄️ Database Schema Highlights

### product Table
```sql
CREATE TABLE product (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE,           -- Auto: SP000001
    name VARCHAR(500),
    description TEXT,
    property_type_id BIGINT,           -- Type: land, house, apartment
    
    -- Location
    province_id BIGINT,
    district_id BIGINT,
    ward_id BIGINT,
    street VARCHAR(255),
    address TEXT,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    
    -- Size & Specs
    area DECIMAL(10,2),                -- m²
    width DECIMAL(8,2),                -- m
    length DECIMAL(8,2),               -- m
    floors INT,
    bedrooms INT,
    bathrooms INT,
    
    -- Price
    price DECIMAL(18,2),
    price_per_m2 DECIMAL(18,2),
    price_type INT,                    -- 1:fixed, 2:negotiable
    
    -- Status
    status INT DEFAULT 1,              -- 1:active, 2:sold, 3:rented
    transaction_type INT,              -- 1:sell, 2:rent, 3:both
    
    -- Ownership
    owner_id BIGINT,                   -- Profile ID
    owner_type INT,                    -- 1:user, 2:organization
    visibility INT DEFAULT 1,          -- 1:public, 2:private, 3:organization
    
    -- Legal
    legal_status VARCHAR(100),
    certificate_type VARCHAR(100),
    
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_product_type ON product(property_type_id);
CREATE INDEX idx_product_owner ON product(owner_id);
CREATE INDEX idx_product_location ON product(province_id, district_id, ward_id);
CREATE INDEX idx_product_price ON product(price);
CREATE INDEX idx_product_status ON product(status);
```

### post Table
```sql
CREATE TABLE post (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,        -- Links to product
    title VARCHAR(500),
    content TEXT,
    
    -- Pricing (can differ from product)
    post_price DECIMAL(18,2),
    price_type INT,                    -- 1:fixed, 2:negotiable
    
    -- Post specific
    transaction_type INT,              -- 1:sell, 2:rent
    type VARCHAR(20),                  -- 1:normal, 2:feed, 3:share
    package_visible INT DEFAULT 1,     -- 1:normal, 2:VIP
    
    -- Timing
    expired_at TIMESTAMP,
    num_date INT,                      -- Number of days active
    
    -- Status
    status INT DEFAULT 1,
    visibility INT DEFAULT 1,
    hidden BOOLEAN DEFAULT FALSE,
    
    -- Stats
    num_like BIGINT DEFAULT 0,
    num_comment BIGINT DEFAULT 0,
    num_view BIGINT DEFAULT 0,
    
    -- Ownership
    owner_id BIGINT,
    owner_type INT,                    -- 1:user, 2:organization
    
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    FOREIGN KEY (product_id) REFERENCES product(id)
);

CREATE INDEX idx_post_product ON post(product_id);
CREATE INDEX idx_post_expired ON post(expired_at);
CREATE INDEX idx_post_status ON post(status);
CREATE INDEX idx_post_owner ON post(owner_id);
```

### asset Table
```sql
CREATE TABLE asset (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE,           -- Auto: AS000001
    product_id BIGINT,                 -- Optional link to product
    name VARCHAR(500),
    description TEXT,
    
    -- Acquisition
    acquisition_date DATE,
    acquisition_cost DECIMAL(18,2),
    acquisition_source VARCHAR(255),
    
    -- Value
    current_value DECIMAL(18,2),
    market_value DECIMAL(18,2),
    depreciation DECIMAL(18,2),
    
    -- Status
    status INT DEFAULT 1,              -- 1:active, 2:sold, 3:rented, 4:archived
    exploitation_status INT,           -- How it's being used
    
    -- Location (same as product)
    province_id BIGINT,
    district_id BIGINT,
    ward_id BIGINT,
    address TEXT,
    
    -- Ownership
    owner_id BIGINT,
    owner_type INT,
    
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    FOREIGN KEY (product_id) REFERENCES product(id)
);
```

### asset_cost Table
```sql
CREATE TABLE asset_cost (
    id BIGSERIAL PRIMARY KEY,
    asset_id BIGINT NOT NULL,
    cost_type_id BIGINT,
    amount DECIMAL(18,2),
    description TEXT,
    cost_date DATE,
    payment_method VARCHAR(100),
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (asset_id) REFERENCES asset(id)
);
```

### project Table
```sql
CREATE TABLE project (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE,
    name VARCHAR(500),
    description TEXT,
    
    -- Developer
    developer_id BIGINT,
    developer_name VARCHAR(255),
    
    -- Location
    province_id BIGINT,
    district_id BIGINT,
    ward_id BIGINT,
    address TEXT,
    
    -- Project info
    total_area DECIMAL(18,2),
    total_units INT,
    total_buildings INT,
    start_date DATE,
    expected_completion DATE,
    actual_completion DATE,
    
    -- Status
    status INT DEFAULT 1,              -- 1:planning, 2:construction, 3:completed
    
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (developer_id) REFERENCES developer(id)
);
```

---

## 🔍 Search & Filter Features

### Product Search
```go
type ProductSearchDTO struct {
    // Location
    ProvinceId   uint64   `json:"provinceId"`
    DistrictId   uint64   `json:"districtId"`
    WardId       uint64   `json:"wardId"`
    
    // Property type
    PropertyTypeIds []uint64 `json:"propertyTypeIds"`
    
    // Price range
    MinPrice     float64  `json:"minPrice"`
    MaxPrice     float64  `json:"maxPrice"`
    
    // Size range
    MinArea      float64  `json:"minArea"`
    MaxArea      float64  `json:"maxArea"`
    
    // Bedrooms/Bathrooms
    MinBedrooms  int      `json:"minBedrooms"`
    MinBathrooms int      `json:"minBathrooms"`
    
    // Transaction type
    TransactionType int   `json:"transactionType"` // 1:sell, 2:rent
    
    // Status
    Status       []int    `json:"status"`
    
    // Ownership
    OwnerId      uint64   `json:"ownerId"`
    OwnerType    int      `json:"ownerType"`
    
    // Pagination
    _dto.Pagable
}
```

### Advanced Search Query
```go
func (r *ProductRepo) Search(ctx context.Context, search *ProductSearchDTO) ([]Product, int64, error) {
    db := r.db.GetDB()
    
    query := db.Model(&Product{}).Where("deleted_at IS NULL")
    
    // Location filters
    if search.ProvinceId > 0 {
        query = query.Where("province_id = ?", search.ProvinceId)
    }
    if search.DistrictId > 0 {
        query = query.Where("district_id = ?", search.DistrictId)
    }
    
    // Property type
    if len(search.PropertyTypeIds) > 0 {
        query = query.Where("property_type_id IN ?", search.PropertyTypeIds)
    }
    
    // Price range
    if search.MinPrice > 0 {
        query = query.Where("price >= ?", search.MinPrice)
    }
    if search.MaxPrice > 0 {
        query = query.Where("price <= ?", search.MaxPrice)
    }
    
    // Area range
    if search.MinArea > 0 {
        query = query.Where("area >= ?", search.MinArea)
    }
    if search.MaxArea > 0 {
        query = query.Where("area <= ?", search.MaxArea)
    }
    
    // Bedrooms/Bathrooms
    if search.MinBedrooms > 0 {
        query = query.Where("bedrooms >= ?", search.MinBedrooms)
    }
    if search.MinBathrooms > 0 {
        query = query.Where("bathrooms >= ?", search.MinBathrooms)
    }
    
    // Transaction type
    if search.TransactionType > 0 {
        query = query.Where("transaction_type = ? OR transaction_type = 3", search.TransactionType)
    }
    
    // Status
    if len(search.Status) > 0 {
        query = query.Where("status IN ?", search.Status)
    }
    
    // Count total
    var total int64
    query.Count(&total)
    
    // Get results
    var products []Product
    err := query.
        Offset(search.GetOffset()).
        Limit(search.GetLimit()).
        Order("created_at DESC").
        Find(&products).Error
    
    return products, total, err
}
```

---

## 💰 Price Management

### Price History Tracking
```go
type ProductPrice struct {
    ID        uint64
    ProductId uint64
    Price     float64
    PriceType int      // 1:fixed, 2:negotiable
    Currency  string   // VND, USD
    ValidFrom time.Time
    ValidTo   *time.Time
    Reason    string   // Price change reason
    CreatedBy uint64
    CreatedAt time.Time
}

func (u *ProductUsecase) UpdatePrice(
    ctx context.Context,
    productId uint64,
    newPrice float64,
    reason string,
) error {
    product, _ := u.repo.GetByID(ctx, productId)
    
    // Save old price to history
    oldPrice := &ProductPrice{
        ProductId: productId,
        Price:     product.Price,
        PriceType: product.PriceType,
        ValidFrom: product.UpdatedAt,
        ValidTo:   timeNow(),
        CreatedBy: getProfileId(ctx),
    }
    u.priceRepo.Create(ctx, oldPrice)
    
    // Update product price
    product.Price = newPrice
    u.repo.Update(ctx, productId, product)
    
    // Create new price record
    newPriceRecord := &ProductPrice{
        ProductId: productId,
        Price:     newPrice,
        PriceType: product.PriceType,
        ValidFrom: timeNow(),
        Reason:    reason,
        CreatedBy: getProfileId(ctx),
    }
    u.priceRepo.Create(ctx, newPriceRecord)
    
    return nil
}
```

---

## 🏗️ Asset Management Features

### Asset Cost Tracking
```go
type AssetCost struct {
    ID              uint64
    AssetId         uint64
    CostTypeId      uint64  // Operating, Maintenance, Tax, etc.
    Amount          float64
    Description     string
    CostDate        time.Time
    PaymentMethod   string
    InvoiceNumber   string
    Status          int     // 1:pending, 2:paid, 3:cancelled
}

func (u *AssetUsecase) AddCost(
    ctx context.Context,
    assetId uint64,
    cost *AssetCost,
) error {
    // Validate asset ownership
    if err := u.validateOwnership(ctx, assetId); err != nil {
        return err
    }
    
    // Create cost record
    cost.AssetId = assetId
    if err := u.costRepo.Create(ctx, cost); err != nil {
        return err
    }
    
    // Update asset total cost
    asset, _ := u.repo.GetByID(ctx, assetId)
    asset.TotalCost += cost.Amount
    u.repo.Update(ctx, assetId, asset)
    
    // Create cost document if needed
    if cost.InvoiceNumber != "" {
        doc := &CostDocument{
            AssetCostId: cost.ID,
            DocumentType: "invoice",
            DocumentNumber: cost.InvoiceNumber,
        }
        u.docRepo.Create(ctx, doc)
    }
    
    return nil
}
```

### Asset Income Tracking
```go
type AssetIncome struct {
    ID              uint64
    AssetId         uint64
    IncomeTypeId    uint64  // Rental, Sale, etc.
    Amount          float64
    Description     string
    IncomeDate      time.Time
    PaymentMethod   string
    Status          int
}

func (u *AssetUsecase) GetROI(ctx context.Context, assetId uint64) (float64, error) {
    asset, _ := u.repo.GetByID(ctx, assetId)
    
    // Total income
    incomes, _ := u.incomeRepo.GetByAsset(ctx, assetId)
    totalIncome := 0.0
    for _, income := range incomes {
        totalIncome += income.Amount
    }
    
    // Total cost
    costs, _ := u.costRepo.GetByAsset(ctx, assetId)
    totalCost := asset.AcquisitionCost
    for _, cost := range costs {
        totalCost += cost.Amount
    }
    
    // ROI = (Total Income - Total Cost) / Total Cost * 100
    roi := (totalIncome - totalCost) / totalCost * 100
    
    return roi, nil
}
```

---

## 📊 Dashboard & Analytics

### Dashboard Metrics
```go
type DashboardDTO struct {
    TotalProduct    uint32 `json:"totalProduct"`
    TotalPost       uint32 `json:"totalPost"`
    TotalAsset      uint32 `json:"totalAsset"`
    TotalProject    uint32 `json:"totalProject"`
    TotalNewsfeed   uint32 `json:"totalNewsfeed"`
    TotalTransaction uint32 `json:"totalTransaction"`
    TotalClicked    uint32 `json:"totalClicked"`
    TotalViewHome   uint32 `json:"totalViewHome"`
}

func (u *DashboardUsecase) GetStats(ctx context.Context) (*DashboardDTO, error) {
    profileId := _utils.GetProfileIdWithContext(ctx)
    orgId := _utils.GetOrganizationIdFromContext(ctx)
    
    stats := &DashboardDTO{}
    
    // Count products
    stats.TotalProduct, _ = u.productRepo.CountByOwner(ctx, profileId, orgId)
    
    // Count posts
    stats.TotalPost, _ = u.postRepo.CountByOwner(ctx, profileId, orgId)
    
    // Count assets
    stats.TotalAsset, _ = u.assetRepo.CountByOwner(ctx, profileId, orgId)
    
    // Count projects
    stats.TotalProject, _ = u.projectRepo.CountByOwner(ctx, profileId, orgId)
    
    return stats, nil
}
```

---

## 🔄 Business Workflows

### 1. Create & Publish Post Flow
```
1. Create Product
   POST /v2/bdspro/v2/product
   - Enter property details
   - Upload images
   - Set location
   
2. Create Post (Listing)
   POST /v2/bdspro/v2/post
   - Link to product
   - Set title & description
   - Set price (can differ from product)
   - Set expiration (e.g., 30 days)
   - Choose package (normal/VIP)
   
3. Post goes live
   - Visible in search results
   - Shown on news feed
   - Send notifications
   
4. Post expires
   - Auto-hide after expiration
   - Can renew: POST /v2/bdspro/v2/post/new-expired
   
5. Product remains
   - Can create new post anytime
   - Product data persists
```

### 2. Asset Lifecycle
```
1. Create Asset
   POST /v2/bdspro/v2/asset
   - Basic info & acquisition cost
   - Link to product (optional)
   
2. Track Costs
   POST /v2/bdspro/v2/asset/cost
   - Operating costs
   - Maintenance
   - Taxes
   
3. Track Income
   POST /v2/bdspro/v2/asset/income
   - Rental income
   - Other revenue
   
4. Monitor Performance
   GET /v2/bdspro/v2/asset/roi
   - Calculate ROI
   - View profit/loss
   
5. Asset Disposition
   - Mark as sold
   - Transfer ownership
   - Archive
```

---

## 📡 External Service Integration

### 1. Organization Service
```go
// Validate organization ownership
orgClient.ValidateMembership(ctx, profileId, orgId)

// Get organization details
org, _ := orgClient.GetOrganization(ctx, orgId)
```

### 2. Social Service
```go
// When creating post, also create news feed
newsFeed := &NewsFeed{
    Title:   post.Title,
    Content: post.Content,
    PostId:  post.ID,
}
socialClient.CreateNewsFeed(ctx, newsFeed)
```

### 3. Notification Service
```go
// Notify when post expires
noti := &NotificationDTO{
    Title:   "Tin đăng hết hạn",
    Message: fmt.Sprintf("Tin '%s' đã hết hạn", post.Title),
    UserId:  post.OwnerId,
}
notificationClient.SendNotification(ctx, noti)
```

---

## 🚀 API Endpoints Summary

### Post APIs
- `POST /v2/bdspro/v2/post` - Create post
- `PUT /v2/bdspro/v2/post/{id}` - Update post
- `DELETE /v2/bdspro/v2/post/{id}` - Delete post
- `GET /v2/bdspro/v2/post/detail/{id}` - Get post detail
- `GET /v2/bdspro/v2/post/personal` - My posts
- `GET /v2/bdspro/v2/post/global` - Public posts
- `POST /v2/bdspro/v2/post/new-expired` - Renew expired post

### Product APIs
- Similar CRUD pattern for products

### Asset APIs
- Similar CRUD pattern for assets
- Additional cost/income tracking APIs

---

## 📊 Key Metrics

### Database Tables: 50+
### Domain Entities: 45+
### Usecases: 37
### Repositories: 40+
### Jobs: 3 (Dashboard stats)

---

**Last Updated**: October 15, 2025
**Service Version**: v2
**Status**: ✅ Production Ready

