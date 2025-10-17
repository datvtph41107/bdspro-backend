# Data Access & Transaction Patterns

## 📋 Overview

**Purpose**: Deep dive into data access patterns and transaction management
**Key Components**: TransactionRepo, CrudRepo, GORM patterns
**Status**: ✅ Production Implementation

---

## 🗄️ Transaction Management

### TransactionRepo - Core Transaction Support

**File**: `shared/common/db/db_transaction.go`

```go
type ITransactionRepo interface {
    Begin() *gorm.DB
    Commit() *gorm.DB
    Rollback() *gorm.DB
    WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type TransactionRepo struct {
    db *gorm.DB
}

// Automatic transaction with rollback on error
func (t *TransactionRepo) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
    return t.db.Transaction(func(tx *gorm.DB) error {
        // Inject transaction into context
        txCtx := context.WithValue(ctx, "tx", tx)
        
        // Execute function
        return fn(txCtx)
    })
}

// Get DB from context (transaction or normal)
func (t *TransactionRepo) GetDB(ctx context.Context) *gorm.DB {
    // If context has transaction, use it
    if tx, ok := ctx.Value("tx").(*gorm.DB); ok && tx != nil {
        return tx.WithContext(ctx)
    }
    
    // Otherwise use normal DB connection
    return t.db.WithContext(ctx)
}
```

**Key Features**:
- ✅ Automatic transaction management
- ✅ Rollback on error
- ✅ Commit on success
- ✅ Context-based transaction passing
- ✅ Nested transaction support

---

## 🏗️ CrudRepo Provider Implementation

### File: `shared/common/provider/repository.go`

**Complete Implementation**:

```go
type CrudRepo[T any] struct {
    *_db.TransactionRepo
    implements _crud.ICrudRepo[T]
}

func (r *CrudRepo[T]) Init(repo _crud.ICrudRepo[T], db *_db.TransactionRepo) {
    r.implements = repo
    r.TransactionRepo = db
}

// Create with transaction and hooks
func (r *CrudRepo[T]) Create(c context.Context, entity *T) error {
    // Call BeforeSave hook
    if err := r.implements.BeforeSave(c, nil, entity); err != nil {
        return err
    }
    
    // Execute in transaction
    err := r.TransactionRepo.WithTransaction(c, func(c context.Context) error {
        // Insert into database
        exec := r.GetDB(c).Create(entity)
        if exec.Error != nil {
            return exec.Error
        }
        
        // Call AfterSave hook
        if err := r.implements.AfterSave(c, nil, entity); err != nil {
            return err
        }
        
        return nil
    })
    
    return err
}

// Update with transaction and hooks
func (r *CrudRepo[T]) Update(c context.Context, id uint64, entity *T) error {
    // Call BeforeSave hook
    if err := r.implements.BeforeSave(c, &id, entity); err != nil {
        return err
    }
    
    // Execute in transaction
    err := r.TransactionRepo.WithTransaction(c, func(c context.Context) error {
        // Update in database
        err := r.GetDB(c).
            Model(entity).
            Where("id = ? AND deleted_at IS NULL", id).
            Updates(entity).Error
        if err != nil {
            return err
        }
        
        // Call AfterSave hook
        return r.implements.AfterSave(c, &id, entity)
    })
    
    return err
}

// Soft delete
func (r *CrudRepo[T]) Delete(c context.Context, id uint64) error {
    return r.GetDB(c).
        Model(new(T)).
        Where("id = ? AND deleted_at IS NULL", id).
        Update("deleted_at", time.Now()).
        Error
}

// Get with soft delete filter
func (r *CrudRepo[T]) GetByID(c context.Context, id uint64) (*T, error) {
    var entity T
    err := r.GetDB(c).
        Where("id = ? AND deleted_at IS NULL", id).
        First(&entity).
        Error
    if err != nil {
        return nil, err
    }
    return &entity, nil
}

// Paginated list with soft delete filter
func (r *CrudRepo[T]) GetList(c context.Context, pagable _dto.IPagable) ([]T, int64, error) {
    var entities []T
    var total int64
    
    db := r.GetDB(c).Where("deleted_at IS NULL")
    
    // Count total
    err := db.Model(new(T)).Count(&total).Error
    if err != nil {
        return nil, 0, err
    }
    
    // Get page
    err = db.
        Offset(pagable.GetOffset()).
        Limit(pagable.GetLimit()).
        Find(&entities).Error
    if err != nil {
        return nil, 0, err
    }
    
    return entities, total, nil
}

// Get multiple by IDs
func (r *CrudRepo[T]) GetListByIDs(c context.Context, ids []uint64) ([]T, error) {
    var entities []T
    err := r.GetDB(c).
        Where("id IN (?) AND deleted_at IS NULL", ids).
        Find(&entities).Error
    return entities, err
}

// Hooks (can be overridden)
func (r *CrudRepo[T]) BeforeSave(c context.Context, id *uint64, entity *T) error {
    return nil
}

func (r *CrudRepo[T]) AfterSave(c context.Context, id *uint64, entity *T) error {
    return nil
}
```

---

## 🔄 Hook Pattern for Custom Logic

### BeforeSave Hook - Validation

```go
// In infra/postgre/product_postgres.go
type ProductRepo struct {
    _provider.CrudRepo[Product]
}

// Override BeforeSave for validation
func (r *ProductRepo) BeforeSave(c context.Context, id *uint64, entity *Product) error {
    // Validate required fields
    if entity.Name == "" {
        return errors.New("tên sản phẩm không được để trống")
    }
    
    if entity.Area <= 0 {
        return errors.New("diện tích phải lớn hơn 0")
    }
    
    if entity.Price < 0 {
        return errors.New("giá không thể âm")
    }
    
    // Check duplicates
    if id == nil {  // Creating new
        existing, _ := r.GetByCode(c, entity.Code)
        if existing != nil {
            return errors.New("mã sản phẩm đã tồn tại")
        }
    }
    
    // Auto-generate code if empty
    if entity.Code == "" {
        codeGen := _usecase.CodeDataUsecase{}
        entity.Code = codeGen.GetProductCode()
    }
    
    return nil
}
```

### AfterSave Hook - Side Effects

```go
func (r *ProductRepo) AfterSave(c context.Context, id *uint64, entity *Product) error {
    // Clear cache
    cacheKey := fmt.Sprintf("product:%d", entity.ID)
    redis.Del(c, cacheKey)
    
    // Update search index
    if err := r.searchService.IndexProduct(c, entity); err != nil {
        log.Printf("Failed to index product: %v", err)
        // Don't fail the transaction
    }
    
    // Send notification (async, non-blocking)
    go func() {
        r.notificationClient.Send(context.Background(), entity.OwnerId, "Sản phẩm đã được cập nhật")
    }()
    
    return nil
}
```

---

## 💾 Common Query Patterns

### Pattern 1: Find with Relations (Preload)

```go
func (r *ProductRepo) GetDetailWithRelations(ctx context.Context, id uint64) (*Product, error) {
    var product Product
    
    err := r.GetDB(ctx).
        Preload("Media").                    // Load product_media
        Preload("Owner").                    // Load owner profile
        Preload("PropertyType").             // Load property type
        Preload("Province").                 // Load province
        Preload("District").                 // Load district
        Preload("Ward").                     // Load ward
        Preload("Amenities").                // Load amenities (many-to-many)
        Where("id = ? AND deleted_at IS NULL", id).
        First(&product).Error
    
    if err != nil {
        return nil, err
    }
    
    return &product, nil
}
```

### Pattern 2: Complex Search with Dynamic Filters

```go
func (r *ProductRepo) Search(ctx context.Context, search *ProductSearchDTO) ([]Product, int64, error) {
    db := r.GetDB(ctx).Model(&Product{})
    
    // Base filter - exclude deleted
    query := db.Where("deleted_at IS NULL")
    
    // Text search
    if search.Text != "" {
        query = query.Where("name LIKE ? OR description LIKE ?", 
            "%"+search.Text+"%", "%"+search.Text+"%")
    }
    
    // Location filters
    if search.ProvinceId > 0 {
        query = query.Where("province_id = ?", search.ProvinceId)
    }
    if search.DistrictId > 0 {
        query = query.Where("district_id = ?", search.DistrictId)
    }
    if search.WardId > 0 {
        query = query.Where("ward_id = ?", search.WardId)
    }
    
    // Property type
    if search.PropertyTypeId > 0 {
        query = query.Where("property_type_id = ?", search.PropertyTypeId)
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
    
    // Transaction type
    if search.TransactionType > 0 {
        query = query.Where("transaction_type = ? OR transaction_type = 3", search.TransactionType)
    }
    
    // Status filter (multiple)
    if len(search.Status) > 0 {
        query = query.Where("status IN ?", search.Status)
    }
    
    // Ownership
    if search.OwnerId > 0 {
        query = query.Where("owner_id = ? AND owner_type = ?", search.OwnerId, search.OwnerType)
    }
    
    // Count total
    var total int64
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Get page
    var products []Product
    err := query.
        Offset(search.GetOffset()).
        Limit(search.GetLimit()).
        Order("created_at DESC").
        Find(&products).Error
    
    return products, total, err
}
```

### Pattern 3: Aggregate Queries

```go
func (r *ProductRepo) GetStatsByOwner(ctx context.Context, ownerId uint64, ownerType int) (*ProductStats, error) {
    stats := &ProductStats{}
    
    db := r.GetDB(ctx).
        Model(&Product{}).
        Where("owner_id = ? AND owner_type = ? AND deleted_at IS NULL", ownerId, ownerType)
    
    // Count total products
    db.Count(&stats.TotalProducts)
    
    // Count by status
    db.Where("status = ?", StatusActive).Count(&stats.ActiveProducts)
    db.Where("status = ?", StatusSold).Count(&stats.SoldProducts)
    
    // Sum total value
    db.Select("COALESCE(SUM(price), 0)").Row().Scan(&stats.TotalValue)
    
    // Average price
    db.Select("COALESCE(AVG(price), 0)").Row().Scan(&stats.AveragePrice)
    
    // Count by transaction type
    db.Where("transaction_type = ?", TransactionSale).Count(&stats.ForSale)
    db.Where("transaction_type = ?", TransactionRent).Count(&stats.ForRent)
    
    return stats, nil
}
```

### Pattern 4: Bulk Operations

```go
func (r *ContactRepo) BulkCreate(ctx context.Context, contacts []Contact) ([]Contact, error) {
    // Use transaction for bulk insert
    err := r.TransactionRepo.WithTransaction(ctx, func(ctx context.Context) error {
        // Batch insert (more efficient than loop)
        if err := r.GetDB(ctx).Create(&contacts).Error; err != nil {
            return err
        }
        
        return nil
    })
    
    if err != nil {
        return nil, err
    }
    
    return contacts, nil
}

// Alternative: Batch insert with upsert
func (r *ContactRepo) BulkUpsert(ctx context.Context, contacts []Contact) error {
    return r.GetDB(ctx).
        Clauses(clause.OnConflict{
            Columns:   []clause.Column{{Name: "phone"}},
            DoUpdates: clause.AssignmentColumns([]string{"full_name", "email", "updated_at"}),
        }).
        Create(&contacts).Error
}
```

---

## 🎯 Real Implementation Examples

### Example 1: Product Repository

**File**: `bdspro-service/infra/postgres/product_postgre.go`

```go
type ProductRepo struct {
    _provider.CrudRepo[Product]  // Embed base CRUD
    mediaRepo MediaRepo           // Related repositories
}

func NewProductRepo(db *_db.TransactionRepo, mediaRepo MediaRepo) repo.IProductRepo {
    repo := &ProductRepo{
        mediaRepo: mediaRepo,
    }
    repo.Init(repo, db)
    return repo
}

// Custom method - not in base CRUD
func (r *ProductRepo) GetByCode(ctx context.Context, code string) (*Product, error) {
    var product Product
    err := r.GetDB(ctx).
        Where("code = ? AND deleted_at IS NULL", code).
        First(&product).Error
    return &product, err
}

// Override BeforeSave for validation
func (r *ProductRepo) BeforeSave(c context.Context, id *uint64, entity *Product) error {
    // Validation
    if entity.Name == "" {
        return errors.New("name is required")
    }
    
    // Auto-generate code
    if entity.Code == "" {
        codeGen := _usecase.CodeDataUsecase{}
        entity.Code = codeGen.GetProductCode()  // SP000001
    }
    
    // Calculate price per m2
    if entity.Area > 0 && entity.Price > 0 {
        entity.PricePerM2 = entity.Price / entity.Area
    }
    
    return nil
}

// Override AfterSave for related operations
func (r *ProductRepo) AfterSave(c context.Context, id *uint64, entity *Product) error {
    // If creating (id is nil before insert, but entity.ID is set after)
    if id == nil {
        // Handle media
        if len(entity.MediaItems) > 0 {
            for i, media := range entity.MediaItems {
                media.ProductId = entity.ID
                media.SortOrder = i + 1
            }
            r.mediaRepo.BulkCreate(c, entity.MediaItems)
        }
    }
    
    // Clear cache
    redis.Del(fmt.Sprintf("product:%d", entity.ID))
    
    return nil
}
```

---

## 🔍 Advanced Query Patterns

### Pattern 1: Subqueries

```go
func (r *PostRepo) GetWithLatestPrice(ctx context.Context) ([]Post, error) {
    var posts []Post
    
    // Subquery to get latest price for each product
    subQuery := r.GetDB(ctx).
        Model(&ProductPrice{}).
        Select("product_id, MAX(created_at) as max_created").
        Group("product_id")
    
    // Main query with join
    err := r.GetDB(ctx).
        Model(&Post{}).
        Joins("LEFT JOIN product ON post.product_id = product.id").
        Joins("LEFT JOIN product_price ON product.id = product_price.product_id AND product_price.created_at IN (?)", subQuery).
        Where("post.deleted_at IS NULL").
        Find(&posts).Error
    
    return posts, err
}
```

### Pattern 2: Aggregation with GROUP BY

```go
func (r *AssetRepo) GetCostSummaryByType(ctx context.Context, assetId uint64) ([]CostSummary, error) {
    var summaries []CostSummary
    
    err := r.GetDB(ctx).
        Model(&AssetCost{}).
        Select("cost_type_id, SUM(amount) as total_amount, COUNT(*) as count").
        Where("asset_id = ? AND deleted_at IS NULL", assetId).
        Group("cost_type_id").
        Scan(&summaries).Error
    
    return summaries, err
}
```

### Pattern 3: Raw SQL for Complex Queries

```go
func (r *DealRepo) GetCommissionReport(ctx context.Context, startDate, endDate time.Time) ([]CommissionReport, error) {
    var report []CommissionReport
    
    sql := `
        SELECT 
            d.id as deal_id,
            d.name as deal_name,
            dm.user_id,
            u.full_name,
            SUM(dm.commission_amount) as total_commission,
            COUNT(d.id) as deal_count
        FROM deal d
        JOIN deal_member dm ON d.id = dm.deal_id
        JOIN profile u ON dm.user_id = u.id
        WHERE d.created_at BETWEEN ? AND ?
          AND d.deleted_at IS NULL
        GROUP BY d.id, dm.user_id, u.full_name
        ORDER BY total_commission DESC
    `
    
    err := r.GetDB(ctx).Raw(sql, startDate, endDate).Scan(&report).Error
    return report, err
}
```

---

## 🔒 Transaction Use Cases

### Use Case 1: Create Product with Media

```go
func (u *ProductUsecase) CreateProductWithMedia(
    ctx context.Context,
    product *Product,
    mediaItems []ProductMedia,
) (*Product, error) {
    // All operations in single transaction
    err := u.repo.(*ProductRepo).TransactionRepo.WithTransaction(ctx, func(ctx context.Context) error {
        // 1. Create product
        if err := u.repo.Create(ctx, product); err != nil {
            return err
        }
        
        // 2. Create media
        for i := range mediaItems {
            mediaItems[i].ProductId = product.ID
            mediaItems[i].SortOrder = i + 1
        }
        if err := u.mediaRepo.BulkCreate(ctx, mediaItems); err != nil {
            return err  // Rollback product creation
        }
        
        // 3. Update product media count
        product.MediaCount = len(mediaItems)
        if err := u.repo.Update(ctx, product.ID, product); err != nil {
            return err  // Rollback all
        }
        
        return nil  // Commit all
    })
    
    return product, err
}
```

### Use Case 2: Transfer Asset Ownership

```go
func (u *AssetUsecase) TransferOwnership(
    ctx context.Context,
    assetId uint64,
    newOwnerId uint64,
    newOwnerType int,
) error {
    return u.repo.(*AssetRepo).TransactionRepo.WithTransaction(ctx, func(ctx context.Context) error {
        // 1. Get asset
        asset, err := u.repo.GetByID(ctx, assetId)
        if err != nil {
            return err
        }
        
        // 2. Validate current ownership
        currentOwnerId := _utils.GetProfileIdWithContext(ctx)
        if asset.OwnerId != currentOwnerId {
            return errors.New("not owner")
        }
        
        // 3. Create transfer record
        transfer := &AssetTransfer{
            AssetId:      assetId,
            FromOwnerId:  asset.OwnerId,
            FromOwnerType: asset.OwnerType,
            ToOwnerId:    newOwnerId,
            ToOwnerType:  newOwnerType,
            TransferDate: time.Now(),
        }
        if err := u.transferRepo.Create(ctx, transfer); err != nil {
            return err
        }
        
        // 4. Update asset ownership
        asset.OwnerId = newOwnerId
        asset.OwnerType = newOwnerType
        if err := u.repo.Update(ctx, assetId, asset); err != nil {
            return err
        }
        
        // 5. Log history
        u.historyRepo.Create(ctx, &History{
            AssetId:    assetId,
            ActionType: ActionTransfer,
            OldValue:   fmt.Sprintf("%d", asset.OwnerId),
            NewValue:   fmt.Sprintf("%d", newOwnerId),
        })
        
        return nil
    })
}
```

### Use Case 3: Split Product into Children

```go
func (u *ProductUsecase) SplitProduct(
    ctx context.Context,
    parentId uint64,
    children []Product,
) error {
    return u.repo.(*ProductRepo).TransactionRepo.WithTransaction(ctx, func(ctx context.Context) error {
        // 1. Get parent product
        parent, err := u.repo.GetByID(ctx, parentId)
        if err != nil {
            return err
        }
        
        // 2. Validate total area
        totalChildArea := 0.0
        for _, child := range children {
            totalChildArea += child.Area
        }
        if totalChildArea > parent.Area {
            return errors.New("tổng diện tích con vượt quá diện tích cha")
        }
        
        // 3. Create child products
        for i := range children {
            children[i].ParentId = &parentId
            children[i].OwnerId = parent.OwnerId
            children[i].OwnerType = parent.OwnerType
            
            if err := u.repo.Create(ctx, &children[i]); err != nil {
                return err
            }
        }
        
        // 4. Update parent status
        parent.HasChildren = true
        parent.AvailableArea = parent.Area - totalChildArea
        if err := u.repo.Update(ctx, parentId, parent); err != nil {
            return err
        }
        
        // 5. Log split history
        u.historyRepo.Create(ctx, &ProductHistory{
            ProductId:  parentId,
            ActionType: ActionSplit,
            ChildCount: len(children),
        })
        
        return nil
    })
}
```

---

## 🚀 Performance Optimizations

### Optimization 1: Batch Loading

```go
// ❌ BAD - N+1 queries
func (r *PostRepo) GetAllWithProducts(ctx context.Context) ([]Post, error) {
    posts, _ := r.GetAll(ctx)
    
    for i := range posts {
        // N queries!
        product, _ := r.productRepo.GetByID(ctx, posts[i].ProductId)
        posts[i].Product = product
    }
    
    return posts, nil
}

// ✅ GOOD - 2 queries total
func (r *PostRepo) GetAllWithProducts(ctx context.Context) ([]Post, error) {
    var posts []Post
    
    // Use GORM Preload - automatic join
    err := r.GetDB(ctx).
        Preload("Product").
        Where("deleted_at IS NULL").
        Find(&posts).Error
    
    return posts, err
}

// ✅ GOOD - Manual batch
func (r *PostRepo) GetAllWithProducts(ctx context.Context) ([]Post, error) {
    // Query 1: Get posts
    posts, _ := r.GetAll(ctx)
    
    // Extract product IDs
    productIds := make([]uint64, len(posts))
    for i, post := range posts {
        productIds[i] = post.ProductId
    }
    
    // Query 2: Batch get products
    products, _ := r.productRepo.GetByIDs(ctx, productIds)
    
    // Create map for O(1) lookup
    productMap := make(map[uint64]*Product)
    for i := range products {
        productMap[products[i].ID] = &products[i]
    }
    
    // Attach products
    for i := range posts {
        posts[i].Product = productMap[posts[i].ProductId]
    }
    
    return posts, nil
}
```

### Optimization 2: Select Only Needed Fields

```go
// ❌ BAD - Select all fields (expensive for large tables)
func (r *ProductRepo) GetList(ctx context.Context) ([]Product, error) {
    var products []Product
    r.GetDB(ctx).Find(&products)  // SELECT * FROM products
    return products, nil
}

// ✅ GOOD - Select only needed fields
func (r *ProductRepo) GetListSummary(ctx context.Context) ([]ProductSummary, error) {
    var summaries []ProductSummary
    
    err := r.GetDB(ctx).
        Model(&Product{}).
        Select("id", "code", "name", "price", "area", "status", "created_at").
        Where("deleted_at IS NULL").
        Scan(&summaries).Error
    
    return summaries, err
}
```

### Optimization 3: Use Indexes

```sql
-- Create composite index for common queries
CREATE INDEX idx_product_search ON product(
    province_id, 
    district_id, 
    property_type_id, 
    status, 
    deleted_at
);

-- Query will use this index
SELECT * FROM product 
WHERE province_id = 1 
  AND district_id = 2 
  AND property_type_id = 3 
  AND status = 1 
  AND deleted_at IS NULL;
```

---

## 📊 Query Performance Monitoring

### Enable GORM Debug Mode

```go
// In development
if env == "development" {
    db = db.Debug()  // Prints all SQL queries
}

// Selective debugging
db.Debug().Where("id = ?", id).First(&product)
// Output: [0.002s] SELECT * FROM products WHERE id = 123
```

### Query Logging

```go
func (r *ProductRepo) GetByID(ctx context.Context, id uint64) (*Product, error) {
    start := time.Now()
    
    var product Product
    err := r.GetDB(ctx).Where("id = ?", id).First(&product).Error
    
    duration := time.Since(start)
    if duration > 50*time.Millisecond {
        log.Printf("SLOW QUERY: GetByID took %v", duration)
    }
    
    return &product, err
}
```

---

## 🎯 Best Practices Summary

### DO ✅

1. **Use TransactionRepo** for multi-step operations
2. **Use CrudRepo** as base for all repositories
3. **Use BeforeSave** for validation
4. **Use AfterSave** for side effects (cache clear, notifications)
5. **Use Preload** to avoid N+1 queries
6. **Use context.Context** throughout
7. **Filter deleted_at IS NULL** explicitly
8. **Use batch operations** for multiple records
9. **Use indexes** for frequently queried fields
10. **Monitor slow queries** in production

### DON'T ❌

1. **Don't** implement Create/Update manually (use CrudRepo)
2. **Don't** forget to check deleted_at
3. **Don't** use N+1 queries (batch instead)
4. **Don't** forget to handle errors
5. **Don't** ignore transaction errors
6. **Don't** put business logic in repositories
7. **Don't** forget to add indexes
8. **Don't** select all fields when few are needed

---

**Last Updated**: October 15, 2025
**Version**: v1.0
**Status**: ✅ Production Patterns

