---
tags:
- sentence-transformers
- sentence-similarity
- feature-extraction
- dense
- generated_from_trainer
- dataset_size:97
- loss:CosineSimilarityLoss
base_model: sentence-transformers/all-MiniLM-L6-v2
widget:
- source_sentence: Status Pattern trong bdspro-service cho Product, Asset, Post?
  sentences:
  - "Dùng `WithTransaction` pattern để wrap operations:\n\n**Interface**:\n```go\n\
    type ITransaction interface {\n    WithTransaction(ctx context.Context, fn func(ctx\
    \ context.Context) error) error\n}\n```\n\n**Implementation**:\n```go\nfunc (t\
    \ *TransactionGorm) WithTransaction(ctx context.Context, fn func(ctx context.Context)\
    \ error) error {\n    return t.DB.Transaction(func(tx *gorm.DB) error {\n    \
    \    // Preserve context values\n        profileId := _utils.GetProfileIdWithContext(ctx)\n\
    \        newCtx := context.WithValue(ctx, \"tx\", tx)\n        newCtx = context.WithValue(newCtx,\
    \ \"profileId\", profileId)\n        \n        return fn(newCtx)  // Auto rollback\
    \ nếu error\n    })\n}\n```\n\n**Usecase**:\n```go\nfunc (uc *OrderUsecase) CreateOrder(ctx\
    \ context.Context, order *domain.Order) error {\n    return uc.transaction.WithTransaction(ctx,\
    \ func(ctx context.Context) error {\n        // 1. Create order\n        if err\
    \ := uc.orderRepo.Create(ctx, order); err != nil {\n            return err  //\
    \ Rollback\n        }\n        \n        // 2. Create order items\n        for\
    \ _, item := range order.Items {\n            if err := uc.orderItemRepo.Create(ctx,\
    \ item); err != nil {\n                return err  // Rollback ALL\n         \
    \   }\n        }\n        \n        // 3. Update inventory\n        if err :=\
    \ uc.inventoryRepo.Decrease(ctx, order.ProductID); err != nil {\n            return\
    \ err  // Rollback ALL\n        }\n        \n        return nil  // Commit tất\
    \ cả\n    })\n}\n```\n\n**Auto behavior**:\n- ✅ Return nil → Commit\n- ❌ Return\
    \ error → Rollback ALL\n- Context có \"tx\" key cho nested repos"
  - "Query Assets với Legal Docs và thông tin location:\n\n```go\npackage postgre\n\
    \nimport (\n    \"context\"\n    \"bdspro/internal/domain\"\n)\n\nfunc (r *AssetRepo)\
    \ GetByOwnerWithLegal(ctx context.Context, ownerID uint64) ([]domain.Asset, error)\
    \ {\n    var assets []domain.Asset\n    \n    err := r.DB.WithContext(ctx).\n\
    \        Preload(\"LegalItems\").\n        Preload(\"PropertyType\").\n      \
    \  Preload(\"Province\").\n        Preload(\"District\").\n        Preload(\"\
    Ward\").\n        Preload(\"AssetExploitations\").\n        Where(\"owner_id =\
    \ ? AND deleted_at IS NULL\", ownerID).\n        Order(\"created_at DESC\").\n\
    \        Find(&assets).Error\n    \n    return assets, err\n}\n```\n\n**Key points**:\n\
    - Preload \"LegalItems\" để lấy danh sách giấy tờ pháp lý\n- Preload \"AssetExploitations\"\
    \ để lấy hợp đồng khai thác\n- Preload location (Province, District, Ward)\n-\
    \ Preload PropertyType\n- Filter by ownerID và soft delete\n\n**Relation**: Asset\
    \ (1) --- (N) AssetLegal [AssetID]"
  - "Các entities chính có Status Patterns riêng:\n\n**Product Status**:\n```go\n\
    SaleStatus  enums.EProductSaleStatus\n  - EProductNotSold (10): Chưa bán\n  -\
    \ EProductSelling (20): Đang bán\n  - EProductSold    (30): Đã bán\n\nRentStatus\
    \  enums.EProductRentStatus\n  - EProductNotRent (10): Chưa cho thuê\n  - EProductRenting\
    \ (20): Đang cho thuê\n  - EProductRented  (30): Đã cho thuê\n```\n\n**Asset Status**:\n\
    ```go\nRentStatus  enums.EAssetStatus\n  - EAssetStatusOwning  (10): Đang sở hữu\n\
    \  - EAssetStatusRenting (20): Đang cho thuê\n  - EAssetStatusSold    (30): Đã\
    \ bán\n  - EAssetStatusSelling (40): Đang bán\n```\n\n**Post Status**:\n```go\n\
    Status  enums.EPostStatus\n  - EPostActive   (10): Đang hoạt động\n  - EPostPending\
    \  (20): Chờ duyệt\n  - EPostApproved (30): Đã duyệt\n  - EPostRejected (40):\
    \ Từ chối\n  - EPostHidden   (50): Ẩn\n\nVisibleStatus  enums.EPostVisibleStatus\
    \ (computed)\n  - EPostTransactionSelling (100): Đang bán\n  - EPostTransactionSold\
    \    (110): Đã bán\n  - EPostTransactionRenting (200): Đang cho thuê\n  - EPostTransactionRented\
    \  (210): Đã cho thuê\n```\n\n**Pattern**: Status tracking lifecycle và business\
    \ state của entity."
- source_sentence: Sản phẩm (Product) trong BDSPro là gì?
  sentences:
  - "**Product Consistency Rules** cần validate:\n\n**Rule 1**: Nếu `Product.AssetId`\
    \ != null, Asset phải tồn tại và chưa bị xóa.\n```go\nif product.AssetId != nil\
    \ {\n    var asset Asset\n    db.First(&asset, product.AssetId)\n    if asset.ID\
    \ == 0 || asset.DeletedAt != nil {\n        return errors.New(\"asset not found\
    \ or deleted\")\n    }\n}\n```\n\n**Rule 2**: Nếu `Product.LastPriceID` != nil,\
    \ ProductPrice phải tồn tại.\n```go\nif product.LastPriceID != nil {\n    var\
    \ price ProductPrice\n    db.First(&price, product.LastPriceID)\n    if price.ID\
    \ == 0 {\n        return errors.New(\"invalid price reference\")\n    }\n}\n```\n\
    \n**Rule 3**: HouseInfo chỉ có khi PropertyType là \"Nhà riêng\" hoặc \"Biệt thự\"\
    .\n\n**Rule 4**: Apartment chỉ có khi PropertyType là \"Căn hộ/Chung cư\".\n\n\
    **Rule 5**: `SaleStatus` và `RentStatus` không thể cùng lúc = \"Sold\" và \"Rented\"\
    .\n```go\nif product.SaleStatus == enums.EProductSold && product.RentStatus ==\
    \ enums.EProductRented {\n    return errors.New(\"product cannot be both sold\
    \ and rented\")\n}\n```\n\n**Implement**: Trong BeforeSave hook hoặc Validator."
  - "**EOwnerOf** định nghĩa loại owner (người sở hữu) của entity.\n\n**Enum values**:\n\
    ```go\nEOwnerOfMember      = 10  // Thành viên (cá nhân trong org)\nEOwnerOfGroup\
    \       = 20  // Nhóm\nEOwnerOfOrgnization = 30  // Tổ chức\nEOwnerOfUser    \
    \    = 40  // Người dùng (ngoài org)\n```\n\n**Chi tiết**:\n\n**Member (10)**:\n\
    - Cá nhân trong Organization sở hữu\n- OwnerID = ProfileID của member\n- Dùng\
    \ cho BĐS của agent cá nhân\n\n**Group (20)**:\n- Nhóm trong Organization sở hữu\n\
    - OwnerID = GroupID\n- Dùng cho BĐS của team/phòng ban\n\n**Organization (30)**:\n\
    - Tổ chức sở hữu\n- OwnerID = OrganizationID\n- Dùng cho BĐS chung của công ty\n\
    \n**User (40)**:\n- Người dùng ngoài org (guest user)\n- Ít dùng trong context\
    \ org\n\n**Sử dụng trong entities**:\n```go\ntype Product struct {\n    OwnerID\
    \   uint64         // ID người sở hữu\n    OwnerOf   enums.EOwnerOf // Loại owner\n\
    }\n```\n\n**Pattern**: Tất cả entities chính đều có OwnerID + OwnerOf để tracking\
    \ ownership."
  - "**Product** trong BDSPro là **sản phẩm bất động sản** (listing) - những BĐS được\
    \ đăng bán hoặc cho thuê.\n\n**Định nghĩa**: Product đại diện cho bảng `products`\
    \ trong database, chứa thông tin chi tiết về một căn nhà, đất, hoặc bất động sản\
    \ khác mà công ty/môi giới đang quản lý.\n\n**Phân biệt với Asset**: \n- **Product**:\
    \ BĐS để bán/cho thuê, có thể là của người khác (môi giới)\n- **Asset**: BĐS mà\
    \ công ty/cá nhân đã mua và sở hữu thực tế\n\n**Entity location**: `bdspro-service/internal/domain/product.go`\n\
    \n**Table name**: `products`\n\n**Embed**: `_models.BaseEntity` (có ID, CreatedAt,\
    \ UpdatedAt, DeletedAt, CreatedBy, UpdatedBy)"
- source_sentence: CrudRepo cung cấp những method gì?
  sentences:
  - "Dùng `Preload` để eager load related entities, tránh N+1:\n\n**Entity với relationships**:\n\
    ```go\ntype Transaction struct {\n    _models.BaseEntity\n    UserID    uint64\n\
    \    ProductID uint64\n    \n    User    *User     `gorm:\"foreignKey:UserID\"\
    `\n    Product *Product  `gorm:\"foreignKey:ProductID\"`\n    Actions []Action\
    \  `gorm:\"foreignKey:TransactionID\"`\n}\n```\n\n**Repository với Preload**:\n\
    ```go\nfunc (r *TransactionRepo) GetByID(c context.Context, id uint64) (*domain.Transaction,\
    \ error) {\n    var tx domain.Transaction\n    err := r.GetDB(c).\n        Preload(\"\
    User\").\n        Preload(\"Product\").\n        Preload(\"Actions\", func(db\
    \ *gorm.DB) *gorm.DB {\n            return db.Order(\"created_at DESC\")\n   \
    \     }).\n        Where(\"id = ? AND deleted_at IS NULL\", id).\n        First(&tx).Error\n\
    \    return &tx, err\n}\n```\n\n**Nested preload**:\n```go\nfunc (r *TransactionRepo)\
    \ GetDetailWithNested(c context.Context, id uint64) (*domain.Transaction, error)\
    \ {\n    var tx domain.Transaction\n    err := r.GetDB(c).\n        Preload(\"\
    User.Profile\").\n        Preload(\"Product.Category\").\n        Preload(\"Actions.CreatedByUser\"\
    ).\n        Where(\"id = ?\", id).\n        First(&tx).Error\n    return &tx,\
    \ err\n}\n```\n\n**Quy tắc**: Preload để tránh N+1, định nghĩa relationships với\
    \ gorm tags, preload selective, nested với dot notation, filter trong callback,\
    \ KHÔNG dùng Join trừ khi filter"
  - '**ProductPrivate** là bảng chứa **thông tin nội bộ/bí mật** của sản phẩm, không
    hiển thị ra ngoài.


    **Entity**: `bdspro-service/internal/domain/product_private.go`


    **Table**: `product_private`


    **Relation**: Product (1) --- (0..1) ProductPrivate [ProductId]


    **Fields**:

    1. `ProductId` (uint64): Link tới Product

    2. `ImportPrice` (*float64): Giá nhập/giá chủ nhà muốn net - CHỈ NỘI BỘ

    3. `OperatingCost` (*float64): Chi phí vận hành (sửa chữa, marketing...)

    4. `TargetProfit` (*float64): Lợi nhuận mục tiêu

    5. `InternalNote` (string): Ghi chú nội bộ (thông tin nhạy cảm)

    6. `PrivateDocs` (string): Danh sách tài liệu riêng tư (lưu dạng CSV)


    **Methods**:

    ```go

    // Set documents

    private.SetPrivateDocs([]string{"doc1.pdf", "doc2.pdf"})

    // => PrivateDocs = "doc1.pdf,doc2.pdf"


    // Get documents

    docs := private.GetPrivateDocs()

    // => []string{"doc1.pdf", "doc2.pdf"}

    ```


    **Use Cases**:

    - Lưu giá thực tế mà chủ nhà chấp nhận (ImportPrice)

    - Tính toán margin: `margin = SalePrice - ImportPrice - OperatingCost`

    - Kiểm tra xem có đạt TargetProfit không

    - Ghi chú riêng về khách hàng, đàm phán, điều kiện đặc biệt

    - Lưu link tài liệu nhạy cảm (hợp đồng, giấy tờ pháp lý...)


    **Security**:

    - CHỈ hiển thị cho member có quyền

    - Không trả về trong public API

    - Không xuất hiện trong Post/Listing công khai


    **Access trong Product**:

    ```go

    product.PrivateData.ImportPrice // => 5.2 tỷ

    product.PrivateData.InternalNote // => "Chủ nhà cần bán gấp"

    ```'
  - "`CrudRepo[T]` từ `common/provider` cung cấp full CRUD operations:\n\n**Basic\
    \ CRUD**:\n- `Create(ctx, entity)`: Tạo mới entity\n- `Update(ctx, id, entity)`:\
    \ Cập nhật entity\n- `Delete(ctx, id)`: Soft delete entity\n\n**Query methods**:\n\
    - `GetByID(ctx, id)`: Lấy 1 entity theo ID\n- `GetDetail(ctx, id)`: Lấy entity\
    \ với relationships\n- `GetAll(ctx)`: Lấy tất cả entities\n- `GetList(ctx, pagable)`:\
    \ Lấy danh sách có phân trang (trả về data + total)\n- `GetListByIDs(ctx, ids)`:\
    \ Lấy nhiều entities theo list IDs\n\n**Hooks**:\n- `BeforeSave(ctx, id, entity)`:\
    \ Hook trước khi save\n- `AfterSave(ctx, id, entity)`: Hook sau khi save\n\n**Features\
    \ tự động**:\n- Transaction support\n- Soft delete filter (deleted_at IS NULL)\n\
    - Audit trail (CreatedBy/UpdatedBy)\n- Error handling\n\n**Usage trong repo**:\n\
    ```go\ntype ProductRepo struct {\n    _provider.CrudRepo[domain.Product]\n}\n\n\
    func NewProductRepo(db *_db.TransactionRepo) _repo.IProductRepo {\n    repo :=\
    \ &ProductRepo{}\n    repo.Init(repo, db)  // Khởi tạo CrudRepo\n    return repo\n\
    }\n```\n\n**Override khi cần custom logic**:\n```go\nfunc (r *ProductRepo) BeforeSave(ctx\
    \ context.Context, id *uint64, entity *domain.Product) error {\n    // Custom\
    \ validation\n    return nil\n}\n```"
- source_sentence: Product có những trạng thái bán và cho thuê nào? Enum values?
  sentences:
  - "**Soft Delete** và **Archive** Product có mục đích khác nhau:\n\n**1. Soft Delete\
    \ (xóa mềm)**:\n```go\nfunc (uc *ProductUsecase) Delete(ctx context.Context, productID\
    \ uint64) error {\n    product, err := uc.productRepo.GetByID(ctx, productID)\n\
    \    if err != nil {\n        return err\n    }\n    \n    // Check permission\n\
    \    profileID := _utils.GetProfileIdWithContext(ctx)\n    if product.OwnerID\
    \ != profileID {\n        return errors.New(\"no permission to delete\")\n   \
    \ }\n    \n    // Check if product has active transactions\n    if product.SaleStatus\
    \ == enums.EProductSelling || product.RentStatus == enums.EProductRenting {\n\
    \        return errors.New(\"cannot delete product with active listings\")\n \
    \   }\n    \n    // Soft delete (set deleted_at = NOW())\n    err = uc.productRepo.Delete(ctx,\
    \ productID)\n    if err != nil {\n        return err\n    }\n    \n    // Hide\
    \ all posts\n    posts, _ := uc.postRepo.GetByProduct(ctx, productID)\n    for\
    \ _, post := range posts {\n        uc.postRepo.Delete(ctx, post.ID) // Soft delete\
    \ posts too\n    }\n    \n    // Create history\n    history := &RecordHistory{\n\
    \        RecordID:   &productID,\n        RecordType: enums.ERecordTypeProduct,\n\
    \        ActionType: enums.EHistoryActionDelete,\n        Content:    []string{fmt.Sprintf(\"\
    Product deleted: %s\", product.Name)},\n    }\n    uc.historyRepo.Create(ctx,\
    \ history)\n    \n    return nil\n}\n```\n\n**2. Archive (lưu trữ)**:\n```go\n\
    func (uc *ProductUsecase) Archive(ctx context.Context, productID uint64) error\
    \ {\n    product, err := uc.productRepo.GetByID(ctx, productID)\n    if err !=\
    \ nil {\n        return err\n    }\n    \n    // Check if already archived\n \
    \   if product.Archived {\n        return errors.New(\"product already archived\"\
    )\n    }\n    \n    // Set archived flag\n    product.Archived = true\n    \n\
    \    // Auto hide if still selling\n    if product.SaleStatus == enums.EProductSelling\
    \ {\n        product.SaleStatus = enums.EProductNotSold\n    }\n    if product.RentStatus\
    \ == enums.EProductRenting {\n        product.RentStatus = enums.EProductNotRent\n\
    \    }\n    \n    err = uc.productRepo.Update(ctx, product)\n    if err != nil\
    \ {\n        return err\n    }\n    \n    // Hide all active posts\n    posts,\
    \ _ := uc.postRepo.GetByProduct(ctx, productID)\n    for _, post := range posts\
    \ {\n        if post.Status == enums.EPostActive {\n            post.Hidden =\
    \ true\n            uc.postRepo.Update(ctx, &post)\n        }\n    }\n    \n \
    \   return nil\n}\n```\n\n**3. Restore (khôi phục)**:\n```go\n// Restore from\
    \ soft delete\nfunc (uc *ProductUsecase) Restore(ctx context.Context, productID\
    \ uint64) error {\n    // Query with Unscoped to get deleted records\n    var\
    \ product Product\n    err := uc.db.Unscoped().Where(\"id = ?\", productID).First(&product).Error\n\
    \    if err != nil {\n        return err\n    }\n    \n    if product.DeletedAt\
    \ == nil {\n        return errors.New(\"product is not deleted\")\n    }\n   \
    \ \n    // Restore (set deleted_at = NULL)\n    err = uc.db.Unscoped().Model(&product).Update(\"\
    deleted_at\", nil).Error\n    if err != nil {\n        return err\n    }\n   \
    \ \n    // Create history\n    history := &RecordHistory{\n        RecordID: \
    \  &productID,\n        RecordType: enums.ERecordTypeProduct,\n        ActionType:\
    \ enums.EHistoryActionUpdate,\n        Content:    []string{\"Product restored\
    \ from deletion\"},\n    }\n    uc.historyRepo.Create(ctx, history)\n    \n  \
    \  return nil\n}\n\n// Unarchive\nfunc (uc *ProductUsecase) Unarchive(ctx context.Context,\
    \ productID uint64) error {\n    product, err := uc.productRepo.GetByID(ctx, productID)\n\
    \    if err != nil {\n        return err\n    }\n    \n    if !product.Archived\
    \ {\n        return errors.New(\"product is not archived\")\n    }\n    \n   \
    \ product.Archived = false\n    return uc.productRepo.Update(ctx, product)\n}\n\
    ```\n\n**Sự khác biệt**:\n\n| Feature | Soft Delete | Archive |\n|---------|-------------|----------|\n\
    | **Mục đích** | Xóa (có thể khôi phục) | Lưu trữ (không dùng nữa) |\n| **Field**\
    \ | deleted_at != NULL | archived = true |\n| **Query** | Bị filter tự động |\
    \ Cần filter thủ công |\n| **Hiển thị** | Không hiển thị (hidden) | Không hiển\
    \ thị (except archive view) |\n| **Restore** | Restore() | Unarchive() |\n| **Use\
    \ case** | Xóa nhầm, clean up | Sản phẩm cũ, hết hạn |\n\n**Query patterns**:\n\
    ```go\n// Active products (không deleted, không archived)\ndb.Where(\"deleted_at\
    \ IS NULL AND archived = false\")\n\n// Archived products (không deleted, nhưng\
    \ archived)\ndb.Where(\"deleted_at IS NULL AND archived = true\")\n\n// Deleted\
    \ products (soft deleted)\ndb.Unscoped().Where(\"deleted_at IS NOT NULL\")\n\n\
    // All products including deleted and archived\ndb.Unscoped()\n```\n\n**Best practices**:\n\
    1. **Soft delete**: Cho phép khôi phục trong 30 ngày, sau đó hard delete (cronjob)\n\
    2. **Archive**: Giữ vĩnh viễn nhưng không hiển thị trong list thường\n3. **Permission**:\
    \ Chỉ owner hoặc admin có thể delete/archive\n4. **Validation**: Không delete/archive\
    \ product đang có giao dịch active"
  - "Product có 2 bộ trạng thái riêng biệt cho Bán và Cho thuê.\n\n**EProductSaleStatus\
    \ - Trạng thái bán**:\n```go\nEProductNotSold = 10  // Chưa bán\nEProductSelling\
    \ = 20  // Đang bán\nEProductSold    = 30  // Đã bán\n```\n\n**EProductRentStatus\
    \ - Trạng thái cho thuê**:\n```go\nEProductNotRent = 10  // Chưa cho thuê\nEProductRenting\
    \ = 20  // Đang cho thuê\nEProductRented  = 30  // Đã cho thuê\n```\n\n**Sử dụng\
    \ trong Product**:\n```go\ntype Product struct {\n    // Sale\n    SaleStatus\
    \        enums.EProductSaleStatus\n    SaleVisibility    enums.EVisibility\n \
    \   SaleTransactionID *uint64\n    \n    // Rent\n    RentStatus        enums.EProductRentStatus\n\
    \    RentVisibility    enums.EVisibility\n    RentTransactionID *uint64\n    \n\
    \    // Transaction Type\n    TransactionType   uint  // 10=Bán, 20=Thuê (gorm\
    \ default:10)\n}\n```\n\n**Validation**: SaleStatus và RentStatus không thể cùng\
    \ lúc = \"Sold\" và \"Rented\".\n\n**Logic**: 1 Product có thể vừa bán vừa cho\
    \ thuê (đồng thời rao 2 loại), nhưng chỉ được completed 1 trong 2."
  - "**Transfer Ownership** Product sang user/group/organization khác:\n\n**1. Transfer\
    \ usecase**:\n```go\nfunc (uc *ProductUsecase) TransferOwnership(ctx context.Context,\
    \ productID uint64, newOwnerID uint64, newOwnerOf enums.EOwnerOf) error {\n  \
    \  return uc.db.Transaction(func(tx *gorm.DB) error {\n        // 1. Get product\n\
    \        product, err := uc.productRepo.GetByID(ctx, productID)\n        if err\
    \ != nil {\n            return err\n        }\n        \n        // 2. Check current\
    \ owner permission\n        profileID := _utils.GetProfileIdWithContext(ctx)\n\
    \        if product.OwnerID != profileID {\n            // Check if admin\n  \
    \          isAdmin := uc.checkAdminPermission(profileID)\n            if !isAdmin\
    \ {\n                return errors.New(\"no permission to transfer\")\n      \
    \      }\n        }\n        \n        // 3. Validate new owner exists\n     \
    \   err = uc.validateOwner(ctx, newOwnerID, newOwnerOf)\n        if err != nil\
    \ {\n            return err\n        }\n        \n        // 4. Store old owner\n\
    \        oldOwnerID := product.OwnerID\n        oldOwnerOf := product.OwnerOf\n\
    \        \n        // 5. Update product ownership\n        product.OwnerID = newOwnerID\n\
    \        product.OwnerOf = newOwnerOf\n        \n        err = uc.productRepo.Update(ctx,\
    \ product)\n        if err != nil {\n            return err\n        }\n     \
    \   \n        // 6. Update ProductUser/ProductOrganization\n        // Remove\
    \ old owner access\n        if oldOwnerOf == enums.EOwnerOfMember {\n        \
    \    tx.Where(\"product_id = ? AND profile_id = ?\", productID, oldOwnerID).\n\
    \                Delete(&ProductUser{})\n        } else if oldOwnerOf == enums.EOwnerOfOrgnization\
    \ {\n            tx.Where(\"product_id = ? AND organization_id = ?\", productID,\
    \ oldOwnerID).\n                Delete(&ProductOrganization{})\n        }\n  \
    \      \n        // Add new owner access\n        if newOwnerOf == enums.EOwnerOfMember\
    \ {\n            productUser := &ProductUser{\n                ProductID: productID,\n\
    \                ProfileID: newOwnerID,\n                IsOwner:   true,\n  \
    \              RoleID:    4, // Owner role\n            }\n            tx.Create(productUser)\n\
    \        } else if newOwnerOf == enums.EOwnerOfOrgnization {\n            productOrg\
    \ := &ProductOrganization{\n                ProductID:      productID,\n     \
    \           OrganizationID: newOwnerID,\n                IsOwner:        true,\n\
    \                RoleID:         4,\n            }\n            tx.Create(productOrg)\n\
    \        }\n        \n        // 7. Create history\n        history := &RecordHistory{\n\
    \            RecordID:   &productID,\n            RecordType: enums.ERecordTypeProduct,\n\
    \            ActionType: enums.EHistoryActionUpdate,\n            Content: []string{\n\
    \                fmt.Sprintf(\"Ownership transferred from %d (%d) to %d (%d)\"\
    ,\n                    oldOwnerID, oldOwnerOf, newOwnerID, newOwnerOf),\n    \
    \        },\n        }\n        uc.historyRepo.Create(ctx, history)\n        \n\
    \        return nil\n    })\n}\n```\n\n**2. Validate owner helper**:\n```go\n\
    func (uc *ProductUsecase) validateOwner(ctx context.Context, ownerID uint64, ownerOf\
    \ enums.EOwnerOf) error {\n    switch ownerOf {\n    case enums.EOwnerOfMember:\n\
    \        // Check member exists in organization\n        var member Member\n \
    \       err := uc.db.Where(\"id = ? AND deleted_at IS NULL\", ownerID).First(&member).Error\n\
    \        if err != nil {\n            return errors.New(\"member not found\")\n\
    \        }\n        \n    case enums.EOwnerOfGroup:\n        // Check group exists\n\
    \        var group Group\n        err := uc.db.Where(\"id = ? AND deleted_at IS\
    \ NULL\", ownerID).First(&group).Error\n        if err != nil {\n            return\
    \ errors.New(\"group not found\")\n        }\n        \n    case enums.EOwnerOfOrgnization:\n\
    \        // Check organization exists\n        var org Organization\n        err\
    \ := uc.db.Where(\"id = ? AND deleted_at IS NULL\", ownerID).First(&org).Error\n\
    \        if err != nil {\n            return errors.New(\"organization not found\"\
    )\n        }\n        \n    default:\n        return errors.New(\"invalid owner_of\
    \ type\")\n    }\n    \n    return nil\n}\n```\n\n**3. Bulk transfer (nhiều products\
    \ cùng lúc)**:\n```go\nfunc (uc *ProductUsecase) BulkTransferOwnership(ctx context.Context,\
    \ productIDs []uint64, newOwnerID uint64, newOwnerOf enums.EOwnerOf) error {\n\
    \    // Validate new owner\n    err := uc.validateOwner(ctx, newOwnerID, newOwnerOf)\n\
    \    if err != nil {\n        return err\n    }\n    \n    // Transfer each product\n\
    \    for _, productID := range productIDs {\n        err = uc.TransferOwnership(ctx,\
    \ productID, newOwnerID, newOwnerOf)\n        if err != nil {\n            //\
    \ Log error but continue\n            log.Printf(\"Failed to transfer product\
    \ %d: %v\", productID, err)\n        }\n    }\n    \n    return nil\n}\n```\n\n\
    **4. Transfer với Posts, Media, Private data**:\nKhi transfer Product, các entities\
    \ liên quan có thể:\n- **Posts**: Giữ nguyên owner (người đăng) hoặc transfer\
    \ theo\n- **Media**: Transfer theo Product (không đổi)\n- **ProductPrivate**:\
    \ Transfer theo Product\n- **ProductPrice**: Giữ nguyên history\n\n**Business\
    \ Rules**:\n1. Chỉ owner hiện tại hoặc admin có thể transfer\n2. New owner phải\
    \ tồn tại và active\n3. Transfer trong transaction để rollback nếu lỗi\n4. Update\
    \ access control tables (ProductUser/ProductOrganization)\n5. Create audit trail\
    \ (history)\n6. Notification cho new owner\n\n**API endpoint example**:\n```go\n\
    // In handler\nfunc (h *ProductHandler) TransferOwnership(ctx context.Context,\
    \ req *pb.TransferProductOwnershipRequest) (*pb.ProductResponse, error) {\n  \
    \  err := h.usecase.TransferOwnership(ctx, req.ProductId, req.NewOwnerId, enums.EOwnerOf(req.NewOwnerOf))\n\
    \    if err != nil {\n        return nil, err\n    }\n    \n    // Get updated\
    \ product\n    product, err := h.usecase.GetByID(ctx, req.ProductId)\n    if err\
    \ != nil {\n        return nil, err\n    }\n    \n    return h.mapper.DomainToPb(product),\
    \ nil\n}\n```"
- source_sentence: Cách lấy profileId và organizationId từ context trong usecase?
  sentences:
  - "**Query pattern** để filter Products với nhiều điều kiện:\n\n**Pattern 1: Filter\
    \ by Status & Visibility (Public Listings)**\n```go\n// Lấy Products đang bán\
    \ công khai\nfunc (r *ProductRepo) GetPublicForSale(ctx context.Context, page,\
    \ size int) ([]Product, int64, error) {\n    var products []Product\n    var total\
    \ int64\n    \n    query := r.DB.WithContext(ctx).\n        Where(\"sale_status\
    \ = ? AND sale_visibility = ? AND deleted_at IS NULL\",\n            enums.EProductSelling,\n\
    \            enums.EVisiblePublic).\n        Where(\"archived = ?\", false)\n\
    \    \n    // Count total\n    query.Model(&Product{}).Count(&total)\n    \n \
    \   // Get paginated data\n    err := query.\n        Preload(\"PropertyType\"\
    ).\n        Preload(\"Province\").\n        Preload(\"District\").\n        Preload(\"\
    Ward\").\n        Preload(\"Price\").\n        Preload(\"MediaList\").\n     \
    \   Order(\"created_at DESC\").\n        Offset((page - 1) * size).\n        Limit(size).\n\
    \        Find(&products).Error\n    \n    return products, total, err\n}\n```\n\
    \n**Pattern 2: Filter by Owner & Organization**\n```go\n// Lấy Products của user\
    \ trong org\nfunc (r *ProductRepo) GetByOwnerInOrg(ctx context.Context, ownerID,\
    \ orgID uint64) ([]Product, error) {\n    var products []Product\n    \n    err\
    \ := r.DB.WithContext(ctx).\n        Where(\"owner_id = ? AND owner_of = ? AND\
    \ deleted_at IS NULL\",\n            ownerID, enums.EOwnerOfMember).\n       \
    \ Where(\"organization_id = ?\", orgID). // Assuming có field này\n        Preload(\"\
    PropertyType\").\n        Preload(\"Price\").\n        Find(&products).Error\n\
    \    \n    return products, err\n}\n```\n\n**Pattern 3: Filter by Location (Province,\
    \ District, Ward)**\n```go\n// Lấy Products theo khu vực\nfunc (r *ProductRepo)\
    \ GetByLocation(ctx context.Context, provinceID, districtID, wardID *uint64) ([]Product,\
    \ error) {\n    var products []Product\n    \n    query := r.DB.WithContext(ctx).Where(\"\
    deleted_at IS NULL\")\n    \n    if provinceID != nil {\n        query = query.Where(\"\
    province_id = ?\", *provinceID)\n    }\n    if districtID != nil {\n        query\
    \ = query.Where(\"district_id = ?\", *districtID)\n    }\n    if wardID != nil\
    \ {\n        query = query.Where(\"ward_id = ?\", *wardID)\n    }\n    \n    err\
    \ := query.\n        Preload(\"Province\").\n        Preload(\"District\").\n\
    \        Preload(\"Ward\").\n        Find(&products).Error\n    \n    return products,\
    \ err\n}\n```\n\n**Pattern 4: Filter by PropertyType & Price Range**\n```go\n\
    // Lấy Products theo loại và giá\nfunc (r *ProductRepo) GetByTypeAndPrice(ctx\
    \ context.Context, propertyTypeID uint64, minPrice, maxPrice float64) ([]Product,\
    \ error) {\n    var products []Product\n    \n    err := r.DB.WithContext(ctx).\n\
    \        Joins(\"JOIN product_price ON products.last_price_id = product_price.id\"\
    ).\n        Where(\"products.property_type_id = ? AND products.deleted_at IS NULL\"\
    , propertyTypeID).\n        Where(\"product_price.sale_price BETWEEN ? AND ?\"\
    , minPrice, maxPrice).\n        Preload(\"PropertyType\").\n        Preload(\"\
    Price\").\n        Find(&products).Error\n    \n    return products, err\n}\n\
    ```\n\n**Pattern 5: Filter với Access Control (Owner + Shared)**\n```go\n// Lấy\
    \ Products user có quyền xem (own + shared)\nfunc (r *ProductRepo) GetAccessibleProducts(ctx\
    \ context.Context, userID uint64) ([]Product, error) {\n    var products []Product\n\
    \    \n    // Subquery: IDs của products được share\n    sharedQuery := r.DB.Table(\"\
    sharing_access\").\n        Select(\"domain_id\").\n        Where(\"to_id = ?\
    \ AND domain = ? AND deleted_at IS NULL\",\n            userID, enums.EDomainProduct)\n\
    \    \n    err := r.DB.WithContext(ctx).\n        Where(\"(owner_id = ? AND owner_of\
    \ = ?) OR id IN (?)\",\n            userID, enums.EOwnerOfMember, sharedQuery).\n\
    \        Where(\"deleted_at IS NULL\").\n        Preload(\"PropertyType\").\n\
    \        Preload(\"Price\").\n        Find(&products).Error\n    \n    return\
    \ products, err\n}\n```\n\n**Pattern 6: Search by Name/Address (Full-text)**\n\
    ```go\n// Tìm kiếm Products theo tên hoặc địa chỉ\nfunc (r *ProductRepo) Search(ctx\
    \ context.Context, keyword string) ([]Product, error) {\n    var products []Product\n\
    \    \n    searchPattern := \"%\" + keyword + \"%\"\n    \n    err := r.DB.WithContext(ctx).\n\
    \        Where(\"(name ILIKE ? OR address ILIKE ?) AND deleted_at IS NULL\",\n\
    \            searchPattern, searchPattern).\n        Preload(\"PropertyType\"\
    ).\n        Preload(\"Province\").\n        Preload(\"District\").\n        Limit(50).\
    \ // Limit results\n        Find(&products).Error\n    \n    return products,\
    \ err\n}\n```\n\n**Best Practices**:\n1. **Luôn check** `deleted_at IS NULL` (soft\
    \ delete)\n2. **Use Preload** cho relations thay vì N+1 queries\n3. **Apply pagination**\
    \ cho list endpoints (Offset + Limit)\n4. **Use indexes** trên filter columns:\
    \ status, visibility, owner_id, province_id\n5. **WithContext** để support timeout/cancellation\n\
    6. **Validate enums** trước khi query\n\n**Performance tip**: Dùng `ProductMarket`\
    \ view (denormalized) cho public listings thay vì join nhiều bảng."
  - "Project có quy ước import alias chuẩn:\n\n**Common libraries (underscore prefix)**:\n\
    ```go\nimport (\n    _utils \"common/utils\"\n    _models \"common/models\"\n\
    \    _dto \"common/domain/dto\"\n    _enum \"common/domain/enum\"\n    _provider\
    \ \"common/provider\"\n    _db \"common/db\"\n)\n```\n\n**Internal packages**:\n\
    ```go\nimport (\n    \"your-service/internal/domain\"\n    _repo \"your-service/internal/interface/repo\"\
    \n)\n```\n\n**Protobuf packages**:\n```go\nimport (\n    pb_bdspro \"pb/types/bdspro\"\
    \n    pb_auth \"pb/types/auth\"\n    shared_enum \"pb/enums\"\n)\n```\n\n**Lý\
    \ do**:\n- Tránh conflict package cùng tên\n- Phân biệt common vs service-specific\n\
    - Consistent naming\n- Dễ đọc và maintain"
  - "Để lấy thông tin user từ context, LUÔN dùng các hàm trong `common/utils`:\n\n\
    **Lấy profileId**:\n```go\nimport _utils \"common/utils\"\n\nfunc (uc *YourUsecase)\
    \ DoSomething(ctx context.Context) error {\n    profileID := _utils.GetProfileIdWithContext(ctx)\n\
    \    // profileID là uint64\n}\n```\n\n**Lấy organizationId**:\n```go\nfunc (uc\
    \ *YourUsecase) DoSomething(ctx context.Context) error {\n    organizationID :=\
    \ _utils.GetOrganizationIdFromContext(ctx)\n    // organizationID là uint64\n\
    }\n```\n\n**Quy tắc**:\n- KHÔNG tự parse từ JWT hay context\n- LUÔN import utils\
    \ với alias `_utils`\n- Chỉ gọi ở layer usecase, handler\n- Context tự động có\
    \ từ middleware gateway"
pipeline_tag: sentence-similarity
library_name: sentence-transformers
---

# SentenceTransformer based on sentence-transformers/all-MiniLM-L6-v2

This is a [sentence-transformers](https://www.SBERT.net) model finetuned from [sentence-transformers/all-MiniLM-L6-v2](https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2). It maps sentences & paragraphs to a 384-dimensional dense vector space and can be used for semantic textual similarity, semantic search, paraphrase mining, text classification, clustering, and more.

## Model Details

### Model Description
- **Model Type:** Sentence Transformer
- **Base model:** [sentence-transformers/all-MiniLM-L6-v2](https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2) <!-- at revision c9745ed1d9f207416be6d2e6f8de32d1f16199bf -->
- **Maximum Sequence Length:** 256 tokens
- **Output Dimensionality:** 384 dimensions
- **Similarity Function:** Cosine Similarity
<!-- - **Training Dataset:** Unknown -->
<!-- - **Language:** Unknown -->
<!-- - **License:** Unknown -->

### Model Sources

- **Documentation:** [Sentence Transformers Documentation](https://sbert.net)
- **Repository:** [Sentence Transformers on GitHub](https://github.com/UKPLab/sentence-transformers)
- **Hugging Face:** [Sentence Transformers on Hugging Face](https://huggingface.co/models?library=sentence-transformers)

### Full Model Architecture

```
SentenceTransformer(
  (0): Transformer({'max_seq_length': 256, 'do_lower_case': False, 'architecture': 'BertModel'})
  (1): Pooling({'word_embedding_dimension': 384, 'pooling_mode_cls_token': False, 'pooling_mode_mean_tokens': True, 'pooling_mode_max_tokens': False, 'pooling_mode_mean_sqrt_len_tokens': False, 'pooling_mode_weightedmean_tokens': False, 'pooling_mode_lasttoken': False, 'include_prompt': True})
  (2): Normalize()
)
```

## Usage

### Direct Usage (Sentence Transformers)

First install the Sentence Transformers library:

```bash
pip install -U sentence-transformers
```

Then you can load this model and run inference.
```python
from sentence_transformers import SentenceTransformer

# Download from the 🤗 Hub
model = SentenceTransformer("sentence_transformers_model_id")
# Run inference
sentences = [
    'Cách lấy profileId và organizationId từ context trong usecase?',
    'Để lấy thông tin user từ context, LUÔN dùng các hàm trong `common/utils`:\n\n**Lấy profileId**:\n```go\nimport _utils "common/utils"\n\nfunc (uc *YourUsecase) DoSomething(ctx context.Context) error {\n    profileID := _utils.GetProfileIdWithContext(ctx)\n    // profileID là uint64\n}\n```\n\n**Lấy organizationId**:\n```go\nfunc (uc *YourUsecase) DoSomething(ctx context.Context) error {\n    organizationID := _utils.GetOrganizationIdFromContext(ctx)\n    // organizationID là uint64\n}\n```\n\n**Quy tắc**:\n- KHÔNG tự parse từ JWT hay context\n- LUÔN import utils với alias `_utils`\n- Chỉ gọi ở layer usecase, handler\n- Context tự động có từ middleware gateway',
    'Project có quy ước import alias chuẩn:\n\n**Common libraries (underscore prefix)**:\n```go\nimport (\n    _utils "common/utils"\n    _models "common/models"\n    _dto "common/domain/dto"\n    _enum "common/domain/enum"\n    _provider "common/provider"\n    _db "common/db"\n)\n```\n\n**Internal packages**:\n```go\nimport (\n    "your-service/internal/domain"\n    _repo "your-service/internal/interface/repo"\n)\n```\n\n**Protobuf packages**:\n```go\nimport (\n    pb_bdspro "pb/types/bdspro"\n    pb_auth "pb/types/auth"\n    shared_enum "pb/enums"\n)\n```\n\n**Lý do**:\n- Tránh conflict package cùng tên\n- Phân biệt common vs service-specific\n- Consistent naming\n- Dễ đọc và maintain',
]
embeddings = model.encode(sentences)
print(embeddings.shape)
# [3, 384]

# Get the similarity scores for the embeddings
similarities = model.similarity(embeddings, embeddings)
print(similarities)
# tensor([[1.0000, 0.9729, 0.9543],
#         [0.9729, 1.0000, 0.9791],
#         [0.9543, 0.9791, 1.0000]])
```

<!--
### Direct Usage (Transformers)

<details><summary>Click to see the direct usage in Transformers</summary>

</details>
-->

<!--
### Downstream Usage (Sentence Transformers)

You can finetune this model on your own dataset.

<details><summary>Click to expand</summary>

</details>
-->

<!--
### Out-of-Scope Use

*List how the model may foreseeably be misused and address what users ought not to do with the model.*
-->

<!--
## Bias, Risks and Limitations

*What are the known or foreseeable issues stemming from this model? You could also flag here known failure cases or weaknesses of the model.*
-->

<!--
### Recommendations

*What are recommendations with respect to the foreseeable issues? For example, filtering explicit content.*
-->

## Training Details

### Training Dataset

#### Unnamed Dataset

* Size: 97 training samples
* Columns: <code>sentence_0</code>, <code>sentence_1</code>, and <code>label</code>
* Approximate statistics based on the first 97 samples:
  |         | sentence_0                                                                        | sentence_1                                                                            | label                                                         |
  |:--------|:----------------------------------------------------------------------------------|:--------------------------------------------------------------------------------------|:--------------------------------------------------------------|
  | type    | string                                                                            | string                                                                                | float                                                         |
  | details | <ul><li>min: 9 tokens</li><li>mean: 17.77 tokens</li><li>max: 31 tokens</li></ul> | <ul><li>min: 215 tokens</li><li>mean: 254.34 tokens</li><li>max: 256 tokens</li></ul> | <ul><li>min: 1.0</li><li>mean: 1.0</li><li>max: 1.0</li></ul> |
* Samples:
  | sentence_0                                                         | sentence_1                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    | label            |
  |:-------------------------------------------------------------------|:------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:-----------------|
  | <code>ProductPrice là gì? Quan hệ với Product ra sao?</code>       | <code>**ProductPrice** là bảng lưu **lịch sử giá** của sản phẩm.<br><br>**Why separate table?**<br>- Product có thể thay đổi giá nhiều lần<br>- Cần tracking lịch sử giá để phân tích<br>- Giá bán và giá thuê khác nhau<br>- Hỗ trợ nhiều loại tiền tệ<br><br>**Relation**:<br>- Product (1) --- (N) ProductPrice<br>- Product.LastPriceID → ProductPrice.ID (giá hiện tại)<br><br>**Typical fields** (cần xem chi tiết trong code):<br>```go<br>type ProductPrice struct {<br>    ID          uint64<br>    ProductID   uint64<br>    SalePrice   *float64  // Giá bán<br>    RentPrice   *float64  // Giá thuê/tháng<br>    Currency    string    // VND, USD<br>    Unit        string    // tỷ, triệu, m²<br>    ValidFrom   time.Time<br>    ValidTo     *time.Time<br>}<br>```<br><br>**Usage**:<br>```go<br>// Get current price<br>product.Price.SalePrice  // => 5.2 (tỷ)<br>product.Price.RentPrice  // => 15 (triệu/tháng)<br><br>// Query price history<br>prices := GetProductPriceHistory(productId)<br>```<br><br>**Update giá**:<br>1. Tạo ProductPrice mới với giá mới<br>2. Update `Product.LastPriceID` về record mới<br>3. Giữ nguyên price cũ...</code> | <code>1.0</code> |
  | <code>Data Consistency Rules cho Product là gì?</code>             | <code>**Product Consistency Rules** cần validate:<br><br>**Rule 1**: Nếu `Product.AssetId` != null, Asset phải tồn tại và chưa bị xóa.<br>```go<br>if product.AssetId != nil {<br>    var asset Asset<br>    db.First(&asset, product.AssetId)<br>    if asset.ID == 0 \|\| asset.DeletedAt != nil {<br>        return errors.New("asset not found or deleted")<br>    }<br>}<br>```<br><br>**Rule 2**: Nếu `Product.LastPriceID` != nil, ProductPrice phải tồn tại.<br>```go<br>if product.LastPriceID != nil {<br>    var price ProductPrice<br>    db.First(&price, product.LastPriceID)<br>    if price.ID == 0 {<br>        return errors.New("invalid price reference")<br>    }<br>}<br>```<br><br>**Rule 3**: HouseInfo chỉ có khi PropertyType là "Nhà riêng" hoặc "Biệt thự".<br><br>**Rule 4**: Apartment chỉ có khi PropertyType là "Căn hộ/Chung cư".<br><br>**Rule 5**: `SaleStatus` và `RentStatus` không thể cùng lúc = "Sold" và "Rented".<br>```go<br>if product.SaleStatus == enums.EProductSold && product.RentStatus == enums.EProductRented {<br>    return errors.New("product cannot be both sold and rented")<br>}<br>```<br>...</code>              | <code>1.0</code> |
  | <code>Product có những computed fields nào? Dùng để làm gì?</code> | <code>Product có các **computed fields** (không lưu trong DB, tính toán runtime):<br><br>**Location Names**:<br>```go<br>ProvinceName   string `gorm:"-"` // Tên tỉnh (computed)<br>DistrictName   string `gorm:"-"` // Tên quận (computed)<br>WardName       string `gorm:"-"` // Tên phường (computed)<br>```<br><br>**Cách tính toán**:<br>- Load từ relations: Province, District, Ward<br>- Mapper/DTO layer set giá trị từ Region.Name<br><br>**Lợi ích**:<br>- Giảm số lần join trong query<br>- Response DTO có sẵn tên địa chỉ, không cần client join<br>- Frontend hiển thị trực tiếp<br><br>**Ví dụ sử dụng**:<br>```go<br>// Trong mapper<br>func ToProductDTO(product *Product) *ProductDTO {<br>    dto := &ProductDTO{<br>        ID:   product.ID,<br>        Name: product.Name,<br>        // ... other fields<br>    }<br>    <br>    if product.Province != nil {<br>        dto.ProvinceName = product.Province.Name<br>    }<br>    if product.District != nil {<br>        dto.DistrictName = product.District.Name<br>    }<br>    if product.Ward != nil {<br>        dto.WardName = product.Ward.Name<br>    }<br>    <br>    return d...</code>    | <code>1.0</code> |
* Loss: [<code>CosineSimilarityLoss</code>](https://sbert.net/docs/package_reference/sentence_transformer/losses.html#cosinesimilarityloss) with these parameters:
  ```json
  {
      "loss_fct": "torch.nn.modules.loss.MSELoss"
  }
  ```

### Training Hyperparameters
#### Non-Default Hyperparameters

- `per_device_train_batch_size`: 16
- `per_device_eval_batch_size`: 16
- `num_train_epochs`: 16
- `multi_dataset_batch_sampler`: round_robin

#### All Hyperparameters
<details><summary>Click to expand</summary>

- `overwrite_output_dir`: False
- `do_predict`: False
- `eval_strategy`: no
- `prediction_loss_only`: True
- `per_device_train_batch_size`: 16
- `per_device_eval_batch_size`: 16
- `per_gpu_train_batch_size`: None
- `per_gpu_eval_batch_size`: None
- `gradient_accumulation_steps`: 1
- `eval_accumulation_steps`: None
- `torch_empty_cache_steps`: None
- `learning_rate`: 5e-05
- `weight_decay`: 0.0
- `adam_beta1`: 0.9
- `adam_beta2`: 0.999
- `adam_epsilon`: 1e-08
- `max_grad_norm`: 1
- `num_train_epochs`: 16
- `max_steps`: -1
- `lr_scheduler_type`: linear
- `lr_scheduler_kwargs`: {}
- `warmup_ratio`: 0.0
- `warmup_steps`: 0
- `log_level`: passive
- `log_level_replica`: warning
- `log_on_each_node`: True
- `logging_nan_inf_filter`: True
- `save_safetensors`: True
- `save_on_each_node`: False
- `save_only_model`: False
- `restore_callback_states_from_checkpoint`: False
- `no_cuda`: False
- `use_cpu`: False
- `use_mps_device`: False
- `seed`: 42
- `data_seed`: None
- `jit_mode_eval`: False
- `bf16`: False
- `fp16`: False
- `fp16_opt_level`: O1
- `half_precision_backend`: auto
- `bf16_full_eval`: False
- `fp16_full_eval`: False
- `tf32`: None
- `local_rank`: 0
- `ddp_backend`: None
- `tpu_num_cores`: None
- `tpu_metrics_debug`: False
- `debug`: []
- `dataloader_drop_last`: False
- `dataloader_num_workers`: 0
- `dataloader_prefetch_factor`: None
- `past_index`: -1
- `disable_tqdm`: False
- `remove_unused_columns`: True
- `label_names`: None
- `load_best_model_at_end`: False
- `ignore_data_skip`: False
- `fsdp`: []
- `fsdp_min_num_params`: 0
- `fsdp_config`: {'min_num_params': 0, 'xla': False, 'xla_fsdp_v2': False, 'xla_fsdp_grad_ckpt': False}
- `fsdp_transformer_layer_cls_to_wrap`: None
- `accelerator_config`: {'split_batches': False, 'dispatch_batches': None, 'even_batches': True, 'use_seedable_sampler': True, 'non_blocking': False, 'gradient_accumulation_kwargs': None}
- `parallelism_config`: None
- `deepspeed`: None
- `label_smoothing_factor`: 0.0
- `optim`: adamw_torch_fused
- `optim_args`: None
- `adafactor`: False
- `group_by_length`: False
- `length_column_name`: length
- `project`: huggingface
- `trackio_space_id`: trackio
- `ddp_find_unused_parameters`: None
- `ddp_bucket_cap_mb`: None
- `ddp_broadcast_buffers`: False
- `dataloader_pin_memory`: True
- `dataloader_persistent_workers`: False
- `skip_memory_metrics`: True
- `use_legacy_prediction_loop`: False
- `push_to_hub`: False
- `resume_from_checkpoint`: None
- `hub_model_id`: None
- `hub_strategy`: every_save
- `hub_private_repo`: None
- `hub_always_push`: False
- `hub_revision`: None
- `gradient_checkpointing`: False
- `gradient_checkpointing_kwargs`: None
- `include_inputs_for_metrics`: False
- `include_for_metrics`: []
- `eval_do_concat_batches`: True
- `fp16_backend`: auto
- `push_to_hub_model_id`: None
- `push_to_hub_organization`: None
- `mp_parameters`: 
- `auto_find_batch_size`: False
- `full_determinism`: False
- `torchdynamo`: None
- `ray_scope`: last
- `ddp_timeout`: 1800
- `torch_compile`: False
- `torch_compile_backend`: None
- `torch_compile_mode`: None
- `include_tokens_per_second`: False
- `include_num_input_tokens_seen`: no
- `neftune_noise_alpha`: None
- `optim_target_modules`: None
- `batch_eval_metrics`: False
- `eval_on_start`: False
- `use_liger_kernel`: False
- `liger_kernel_config`: None
- `eval_use_gather_object`: False
- `average_tokens_across_devices`: True
- `prompts`: None
- `batch_sampler`: batch_sampler
- `multi_dataset_batch_sampler`: round_robin
- `router_mapping`: {}
- `learning_rate_mapping`: {}

</details>

### Framework Versions
- Python: 3.13.5
- Sentence Transformers: 5.1.1
- Transformers: 4.57.1
- PyTorch: 2.9.0
- Accelerate: 1.10.1
- Datasets: 4.2.0
- Tokenizers: 0.22.1

## Citation

### BibTeX

#### Sentence Transformers
```bibtex
@inproceedings{reimers-2019-sentence-bert,
    title = "Sentence-BERT: Sentence Embeddings using Siamese BERT-Networks",
    author = "Reimers, Nils and Gurevych, Iryna",
    booktitle = "Proceedings of the 2019 Conference on Empirical Methods in Natural Language Processing",
    month = "11",
    year = "2019",
    publisher = "Association for Computational Linguistics",
    url = "https://arxiv.org/abs/1908.10084",
}
```

<!--
## Glossary

*Clearly define terms in order to be accessible across audiences.*
-->

<!--
## Model Card Authors

*Lists the people who create the model card, providing recognition and accountability for the detailed work that goes into its construction.*
-->

<!--
## Model Card Contact

*Provides a way for people who have updates to the Model Card, suggestions, or questions, to contact the Model Card authors.*
-->