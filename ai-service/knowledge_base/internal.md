# BDSPro Internal Services - API Documentation

Tài liệu này mô tả chi tiết các **Internal gRPC Services** mà bdspro-service cung cấp cho các services khác trong hệ thống microservices.

---

## Tổng quan Internal Services

BDSPro service expose **7 gRPC services chính**:

1. **BdsproInternalService** - Internal APIs cho services khác
2. **ProductService** - Quản lý sản phẩm BĐS
3. **AssetService** - Quản lý tài sản
4. **PostService** - Quản lý bài đăng
5. **TransactionService** - Quản lý giao dịch
6. **BdsproPublicService** - Public APIs (market, reference data)
7. **Admin Services** - Quản trị (Product, Asset, Post, Project, Transaction)

---

## 1. BdsproInternalService - Core Internal APIs

**Package**: `bdspropb.BdsproInternalService`

**Purpose**: Cung cấp APIs đơn giản cho services khác để query thông tin cơ bản.

### 1.1. GetCountByOwner

**Mục đích**: Lấy số lượng Product/Asset/Post/Project của một owner.

**gRPC Method**:
```protobuf
rpc GetCountByOwner(GetCountByOwnerRequest) returns (GetCountByOwnerResponse);
```

**Request**:
```protobuf
message GetCountByOwnerRequest {
    sharepb.OwnerOf ownerOf = 1;  // 10=Member, 20=Group, 30=Organization
    uint64 ownerId = 2;            // ID của owner
}
```

**Response**:
```protobuf
message GetCountByOwnerResponse {
    uint32 totalProduct = 1;   // Tổng số Product
    uint32 totalAsset = 2;     // Tổng số Asset
    uint32 totalPost = 3;      // Tổng số Post
    uint32 totalProject = 4;   // Tổng số Project
}
```

**Use Cases**:
- Dashboard hiển thị số lượng items của user
- Organization service query số lượng BĐS của member
- Analytics service thu thập metrics

**Example Call** (từ Go client):
```go
client := bdspropb.NewBdsproInternalServiceClient(conn)
resp, err := client.GetCountByOwner(ctx, &bdspropb.GetCountByOwnerRequest{
    OwnerOf: sharepb.OwnerOf_MEMBER,
    OwnerId: 12345,
})
// resp.TotalProduct, resp.TotalAsset, ...
```

---

## 2. ProductService - Quản lý Sản phẩm BĐS

**Package**: `bdspropb.ProductService`

**Purpose**: CRUD và các operations phức tạp trên Product.

### 2.1. Internal Query APIs

#### 2.1.1. GetDetailByIds

**Mục đích**: Lấy thông tin chi tiết Products theo danh sách IDs (internal use).

```protobuf
rpc GetDetailByIds (IdRequest) returns (ProductInternalResponse);
```

**Use Cases**: Service khác cần thông tin Products để enrich data.

---

#### 2.1.2. GetProductAttachmentByIds

**Mục đích**: Lấy attachments/media của Products theo IDs.

```protobuf
rpc GetProductAttachmentByIds (IdRequest) returns (ProductAttachmentResponse);
```

---

### 2.2. CRUD Operations

#### 2.2.1. ProductNew - Tạo Product mới

```protobuf
rpc ProductNew (ProductSaveRequest) returns (Response) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/product/new"
        body: "*"
    };
}
```

**Request** (`ProductSaveRequest`):
```protobuf
message ProductSaveRequest {
    optional uint64 parentId = 1;          // Parent Product (nếu là child)
    string name = 2;                       // Tên sản phẩm (required)
    string code = 3;                       // Mã sản phẩm
    double area = 4;                       // Diện tích m²
    string description = 5;                // Mô tả
    string note = 6;                       // Ghi chú
    optional uint64 propertyTypeId = 7;    // Loại BĐS
    optional uint64 docTypeId = 8;         // Loại giấy tờ
    repeated uint64 amenityIds = 9;        // Danh sách tiện ích
    optional uint64 projectId = 10;        // Dự án
    
    // Location
    optional uint64 provinceId = 11;       // Tỉnh
    optional uint64 districtId = 12;       // Quận
    optional uint64 wardId = 13;           // Phường
    
    int32 transactionType = 14;            // Bán/Thuê
    ProductPrice priceData = 15;           // Thông tin giá
    repeated MediaItem mediaItems = 16;    // Ảnh/video
    HouseInfo houseInfo = 17;              // Thông tin nhà
    ProductPrivate productPrivate = 18;    // Thông tin nội bộ
    string googleMapLink = 19;             // Link Google Maps
    
    uint32 sourceType = 20;                // Nguồn
    uint32 saleStatus = 21;                // Trạng thái bán
    uint32 rentStatus = 22;                // Trạng thái thuê
    
    PostSaveWithProduct post = 23;         // Tạo Post cùng lúc
    AssetCreateWithProduct asset = 24;     // Tạo Asset cùng lúc
    
    uint32 saleVisibility = 25;            // Visibility bán
    uint32 rentVisibility = 26;            // Visibility thuê
    uint32 ownerType = 27;                 // Loại owner
    uint64 id = 28;                        // ID (for update)
    uint64 ownerId = 29;                   // ID owner
}
```

**Nested Messages**:

**ProductPrice**:
```protobuf
message ProductPrice {
    uint64 id = 1;
    uint64 productId = 2;
    string currency = 3;
    optional double priceOwner = 4;        // Giá chủ nhà
    optional double salePrice = 5;         // Giá bán
    optional double saleCommission = 6;    // Hoa hồng bán
    optional double deposite = 7;          // Tiền đặt cọc
    optional double rentPrice = 8;         // Giá thuê
    optional double rentCommission = 9;    // Hoa hồng thuê
    optional uint32 rentPaymentCycle = 10; // Chu kỳ thanh toán
}
```

**HouseInfo**:
```protobuf
message HouseInfo {
    optional int32 numBedroom = 2;     // Số phòng ngủ
    optional int32 numBathroom = 3;    // Số WC
    optional int32 numFloor = 4;       // Số tầng
    optional int32 numFront = 5;       // Số mặt tiền
    optional int32 numCarPark = 6;     // Chỗ đậu xe
    optional string furniture = 7;     // Nội thất
    optional string orientation = 8;   // Hướng nhà
    optional int32 numToilet = 9;      // Số toilet
}
```

**PostSaveWithProduct** (tạo Post cùng lúc):
```protobuf
message PostSaveWithProduct {
    uint64 productId = 1;
    string title = 2;
    string content = 3;
    double postPrice = 4;
    int32 transactionType = 5;
    int32 visibility = 6;
    int32 packageVisible = 7;      // 1=Thường, 2=VIP
    int32 expiredAt = 8;
    int32 numDate = 9;
    bool hidden = 10;
    repeated MediaItem metadatas = 11;
}
```

**AssetCreateWithProduct** (tạo Asset cùng lúc):
```protobuf
message AssetCreateWithProduct {
    double purchasePrice = 1;
    int32 purchaseDate = 2;
    int32 legalStatus = 3;
    int32 rentStatus = 4;
    string description = 5;
    repeated LegalItem legalItems = 6;
}
```

---

#### 2.2.2. UpdateProduct - Cập nhật Product

```protobuf
rpc UpdateProduct (ProductSaveRequest) returns (Response) {
    option (google.api.http) = {
        put: "/v2/bdspro/v2/product/edit/{id}"
        body: "*"
    };
}
```

**Lưu ý**: Request giống ProductNew, nhưng phải có `id` field.

---

#### 2.2.3. DeleteProduct - Xóa Product

```protobuf
rpc DeleteProduct (IdRequest) returns (Response) {
    option (google.api.http) = {
        delete: "/v2/bdspro/v2/product/{id}"
    };
}
```

**Lưu ý**: Soft delete (set deleted_at).

---

### 2.3. Query & Search APIs

#### 2.3.1. ProductMe - Products của user hiện tại

```protobuf
rpc ProductMe (ProductSearch) returns (SearchResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/product/me"
    };
}
```

**ProductSearch** (filter phức tạp):
```protobuf
message ProductSearch {
    uint32 page = 1;
    uint32 size = 2;
    string sort = 3;
    string text = 4;                  // Search text
    uint64 ownerId = 5;
    uint64 ownerType = 6;
    optional uint64 parentId = 7;
    uint64 propertyTypeId = 8;
    uint64 docTypeId = 9;
    uint64 projectId = 10;
    uint64 provinceId = 11;
    uint64 districtId = 12;
    uint64 wardId = 13;
    int32 transactionType = 14;
    int32 saleStatus = 15;
    int32 rentStatus = 16;
    int32 saleVisibility = 17;
    int32 rentVisibility = 18;
    int32 sourceType = 19;
    int32 transactionPrice = 20;
    int32 salePrice = 21;
    int32 rentPrice = 22;
    int32 saleCommission = 23;
    int32 rentCommission = 24;
    int32 rentPaymentCycle = 25;
    optional bool archived = 26;
}
```

**Response**:
```protobuf
message SearchResponse {
    repeated ProductResponse data = 1;
    int64 totalElements = 2;
}
```

---

#### 2.3.2. GetDetail - Chi tiết Product

```protobuf
rpc GetDetail (IdRequest) returns (ProductResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/product/{id}/detail"
    };
}
```

**ProductResponse** (đầy đủ thông tin):
```protobuf
message ProductResponse {
    uint64 id = 1;
    optional uint64 parentId = 2;
    optional uint64 assetId = 3;
    optional uint64 propertyTypeId = 4;
    optional ItemResponse propertyType = 5;
    optional uint64 docTypeId = 6;
    optional DocType docType = 7;
    repeated AmenityItem amenities = 8;
    uint32 visibility = 9;
    string name = 10;
    string code = 11;
    optional uint64 categoryId = 12;
    double area = 13;
    string positionUrl = 14;
    string description = 15;
    uint64 organizationId = 16;
    int32 orgStatus = 17;
    optional uint64 lastPriceId = 18;
    optional ProductPrice price = 19;
    optional uint64 apartmentId = 20;
    Apartment apartment = 21;
    ApartmentAttribute apartmentAttribute = 22;
    ProjectBuild build = 23;
    Project project = 24;
    Developer developer = 25;
    int32 sourceType = 26;
    optional uint64 customerId = 27;
    optional ProductPrivate privateData = 28;
    string note = 29;
    int32 transactionType = 30;
    uint32 saleStatus = 31;
    uint32 saleVisibility = 32;
    uint32 rentStatus = 33;
    uint32 rentVisibility = 34;
    optional uint64 ownerId = 35;
    uint32 ownerType = 36;
    optional uint64 provinceId = 37;
    optional Region province = 38;
    optional uint64 districtId = 39;
    optional Region district = 40;
    optional uint64 wardId = 41;
    optional Region ward = 42;
    string address = 43;
    string googleMapLink = 44;
    repeated MediaItem mediaList = 45;
    HouseInfo houseInfo = 46;
    // ... more fields
}
```

---

#### 2.3.3. CurrentOrganizationProducts - Products của Organization

```protobuf
rpc CurrentOrganizationProducts (ProductSearch) returns (SearchResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/product/organization/current"
    };
}
```

---

#### 2.3.4. GetGroupProducts - Products của Group

```protobuf
rpc GetGroupProducts (GroupProductSearch) returns (SearchResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/product/group/{groupId}"
    };
}
```

---

### 2.4. Status Update APIs

#### 2.4.1. SaleStatus - Cập nhật trạng thái bán

```protobuf
rpc SaleStatus (StatusRequest) returns (Response) {
    option (google.api.http) = {
        put: "/v2/bdspro/v2/product/sale-status"
        body: "*"
    };
}
```

**StatusRequest**:
```protobuf
message StatusRequest {
    uint64 id = 1;
    int32 status = 2;          // New status (10=NotSold, 20=Selling, 30=Sold)
    uint64 contactId = 3;      // Contact ID (nếu sold)
    string note = 4;           // Ghi chú
    double amount = 5;         // Số tiền
    string phone = 6;          // Số điện thoại
}
```

---

#### 2.4.2. RentStatus - Cập nhật trạng thái thuê

```protobuf
rpc RentStatus (StatusRequest) returns (Response) {
    option (google.api.http) = {
        put: "/v2/bdspro/v2/product/rent-status"
        body: "*"
    };
}
```

---

#### 2.4.3. SaleVisibility & RentVisibility

```protobuf
rpc SaleVisibility (VisibilityRequest) returns (Response);
rpc RentVisibility (VisibilityRequest) returns (Response);
```

**VisibilityRequest**:
```protobuf
message VisibilityRequest {
    uint64 id = 1;
    int32 visibility = 2;  // 10=Public, 20=Internal, 30=Private
}
```

---

### 2.5. Child Product Operations (Split/Merge)

#### 2.5.1. ChildSplit - Tách Product thành nhiều Products con

```protobuf
rpc ChildSplit (SaveChildRequest) returns (Response) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/product/child/split"
        body: "*"
    };
}
```

**SaveChildRequest**:
```protobuf
message SaveChildRequest {
    uint64 parentId = 1;               // ID Product cha
    string name = 2;
    string code = 3;
    double area = 4;                   // Diện tích con
    string description = 5;
    string note = 6;
    // ... các fields khác giống ProductSaveRequest
}
```

**Flow**:
1. Tạo Product con với `parentId`
2. Tạo ProductHistory với action = Split
3. Archive Product cha (optional)

---

#### 2.5.2. ChildDevideProduct - Chia Product thành nhiều phần

```protobuf
rpc ChildDevideProduct (DevideChildRequest) returns (Response) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/product/child/devide"
        body: "*"
    };
}
```

**DevideChildRequest**:
```protobuf
message DevideChildRequest {
    uint64 parentId = 1;
    repeated SaveChildRequest childs = 2;  // Tạo nhiều child cùng lúc
}
```

---

#### 2.5.3. ChildMergeProduct - Gộp nhiều Products thành 1

```protobuf
rpc ChildMergeProduct (MergeChildRequest) returns (Response) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/product/child/merge"
        body: "*"
    };
}
```

**MergeChildRequest**:
```protobuf
message MergeChildRequest {
    uint64 parentId = 1;           // ID Product cha (result)
    repeated uint64 childIds = 2;  // IDs các Product con cần gộp
}
```

---

#### 2.5.4. ChildMergeAll - Gộp tất cả child Products

```protobuf
rpc ChildMergeAll (IdRequest) returns (Response) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/product/child/merge-all/{id}"
        body: "*"
    };
}
```

---

### 2.6. Other Operations

#### 2.6.1. ArchiveProduct - Archive/Unarchive Product

```protobuf
rpc ArchiveProduct (ArchivedRequest) returns (Response) {
    option (google.api.http) = {
        put: "/v2/bdspro/v2/product/archived"
        body: "*"
    };
}
```

---

#### 2.6.2. GetMembers - Lấy danh sách members có quyền

```protobuf
rpc GetMembers (IdRequest) returns (MemberResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/product/{id}/members"
    };
}
```

**MemberResponse**:
```protobuf
message MemberResponse {
    repeated ProductAccess data = 1;
}

message ProductAccess {
    uint64 id = 1;
    uint64 productId = 2;
    string fields = 3;           // Fields được phép xem (CSV)
    uint32 targetType = 4;       // 1=User, 2=Group, 3=Org
    uint64 targetId = 5;
    float commission = 6;        // Hoa hồng
    string targetName = 7;
    string avatar = 8;
}
```

---

#### 2.6.3. GetPosts - Lấy danh sách Posts của Product

```protobuf
rpc GetPosts (IdRequest) returns (PostResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/product/{id}/posts"
    };
}
```

---

#### 2.6.4. AreaInfo - Thông tin diện tích Product

```protobuf
rpc AreaInfo (IdRequest) returns (AreaInfoResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/product/area-info/{id}"
    };
}
```

**AreaInfoResponse**:
```protobuf
message AreaInfoResponse {
    double totalArea = 1;      // Tổng diện tích
    double avaiableArea = 2;   // Diện tích còn lại
    double splitArea = 3;      // Diện tích đã tách
}
```

---

#### 2.6.5. Suggest - AI Suggest Product info

```protobuf
rpc Suggest (SuggestRequest) returns (ProductResponse) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/product/suggest"
        body: "*"
    };
}
```

**SuggestRequest**:
```protobuf
message SuggestRequest {
    string content = 1;  // Mô tả ngắn gọn, AI parse thành Product
}
```

---

#### 2.6.6. CreateProductForUser - Tạo Product cho User khác

```protobuf
rpc CreateProductForUser (CreateProductForUserRequest) returns (CreateProductResponse) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/product/create-for-user"
        body: "*"
    };
}
```

**CreateProductForUserRequest**:
```protobuf
message CreateProductForUserRequest {
    ProductSaveRequest product = 1;
    uint64 profileId = 2;    // User nhận Product
    bool isOwner = 3;         // Là owner?
    uint64 roleId = 4;        // Role ID
}
```

---

#### 2.6.7. CreateProductForOrg - Tạo Product cho Organization

```protobuf
rpc CreateProductForOrg (CreateProductForOrgRequest) returns (CreateProductResponse) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/product/create-for-org"
        body: "*"
    };
}
```

---

## 3. AssetService - Quản lý Tài sản

**Package**: `bdspropb.AssetService`

**Purpose**: CRUD và operations trên Asset (tài sản thực tế đã mua).

### 3.1. CRUD Operations

#### 3.1.1. CreateAsset - Tạo Asset mới

```protobuf
rpc CreateAsset(SaveAssetRequest) returns (Response) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/asset"
        body: "*"
    };
}
```

**SaveAssetRequest**:
```protobuf
message SaveAssetRequest {
    uint64 id = 1;
    string name = 2;
    double area = 3;
    optional uint64 wardId = 4;
    optional uint64 districtId = 5;
    optional uint64 provinceId = 6;
    optional uint64 productId = 7;        // Link tới Product (nếu có)
    optional uint64 parentAssetId = 8;    // Parent Asset (nếu split)
    string splitMergeStatus = 9;
    double purchasePrice = 10;            // Giá mua (required)
    string purchaseDate = 11;
    uint32 legalStatus = 12;              // 10=Sổ đỏ, 20=Sổ hồng, 30=Chưa có
    bool archived = 13;
    uint32 rentStatus = 14;
    string address = 15;
    string description = 16;
    optional uint64 propertyTypeId = 18;
    repeated LegalItem legalItems = 19;   // Giấy tờ pháp lý
    uint32 ownerType = 20;
}
```

---

#### 3.1.2. UpdateAsset - Cập nhật Asset

```protobuf
rpc UpdateAsset(SaveAssetRequest) returns (Response) {
    option (google.api.http) = {
        put: "/v2/bdspro/v2/asset/{id}"
        body: "*"
    };
}
```

---

#### 3.1.3. DeleteAsset - Xóa Asset

```protobuf
rpc DeleteAsset(IdRequest) returns (Response) {
    option (google.api.http) = {
        delete: "/v2/bdspro/v2/asset/{id}"
    };
}
```

---

### 3.2. Query APIs

#### 3.2.1. SearchInternal - Search Assets (internal)

```protobuf
rpc SearchInternal(AssetSearchRequest) returns (AssetSearchResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/asset/internal"
    };
}
```

**AssetSearchRequest**:
```protobuf
message AssetSearchRequest {
    string keyword = 1;
    uint32 type = 2;
    bool archived = 3;
}
```

**AssetSearchResponse**:
```protobuf
message AssetSearchResponse {
    repeated Asset data = 1;
    uint64 total = 2;
}
```

---

#### 3.2.2. GetDetail - Chi tiết Asset

```protobuf
rpc GetDetail(sharepb.IdRequest) returns (Asset) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/asset/detail/{id}"
    };
}
```

---

#### 3.2.3. SearchShare - Assets được share

```protobuf
rpc SearchShare(AssetSearchRequest) returns (AssetSearchResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/asset/search-share"
    };
}
```

---

#### 3.2.4. SearchSplit - Assets đã bị split

```protobuf
rpc SearchSplit(AssetSearchRequest) returns (AssetSearchResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/asset/split"
    };
}
```

---

### 3.3. Split & Merge Operations

#### 3.3.1. Split - Tách Asset thành nhiều Assets

```protobuf
rpc Split(SplitAssetRequest) returns (Asset) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/asset/split"
        body: "*"
    };
}
```

**SplitAssetRequest**:
```protobuf
message SplitAssetRequest {
    uint64 id = 1;              // ID Asset cha
    repeated Asset datas = 2;   // Danh sách Assets con
}
```

**Flow**:
1. Tạo các Asset con với `parentAssetId = id`
2. Set `splitMergeStatus = "split"` cho Asset cha
3. Create history record

---

#### 3.3.2. Merge - Gộp nhiều Assets thành 1

```protobuf
rpc Merge(MergeAssetRequest) returns (Asset) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/asset/merge"
        body: "*"
    };
}
```

**MergeAssetRequest**:
```protobuf
message MergeAssetRequest {
    repeated uint64 ids = 1;  // IDs các Assets cần gộp
    string new_name = 2;      // Tên Asset mới
    uint64 id = 3;            // ID Asset result (optional)
}
```

---

### 3.4. Other Operations

#### 3.4.1. Archived - Archive/Unarchive Asset

```protobuf
rpc Archived(ArchivedRequest) returns (Response) {
    option (google.api.http) = {
        put: "/v2/bdspro/v2/asset/archived"
    };
}
```

---

#### 3.4.2. GetPublish - Assets được publish của user

```protobuf
rpc GetPublish(IdRequest) returns (AssetSearchResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/asset/publish/{id}"  // profileId
    };
}
```

---

#### 3.4.3. CreateAssetOrganization - Tạo Asset cho Organization

```protobuf
rpc CreateAssetOrganization (SaveAssetRequest) returns (Response) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/asset/organization"
        body: "*"
    };
}
```

---

#### 3.4.4. CurrentOrganizationAssets - Assets của Organization

```protobuf
rpc CurrentOrganizationAssets (AssetSearchRequest) returns (AssetSearchResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/asset/organization"
    };
}
```

---

## 4. PostService - Quản lý Bài đăng

**Package**: `bdspropb.PostService`

**Purpose**: CRUD và operations trên Post (bài đăng tin rao).

### 4.1. Internal Query APIs

#### 4.1.1. GetPostByID - Lấy Post theo ID

```protobuf
rpc GetPostByID (GetPostByIDRequest) returns (Post);
```

---

#### 4.1.2. GetPostByIDs - Lấy nhiều Posts theo IDs

```protobuf
rpc GetPostByIDs(GetPostByIDsRequest) returns (GetPostByIDsResponse);
```

**Use Cases**: Service khác cần thông tin Posts để hiển thị.

---

#### 4.1.3. OwnerPost - Check ownership của Post

```protobuf
rpc OwnerPost(OwnerPostRequest) returns (OwnerPostResponse);
```

**OwnerPostRequest**:
```protobuf
message OwnerPostRequest {
    uint64 postId = 1;
    uint64 profileId = 2;
}
```

**OwnerPostResponse**:
```protobuf
message OwnerPostResponse {
    bool success = 1;  // true nếu user là owner
}
```

---

### 4.2. CRUD Operations

#### 4.2.1. CreatePost - Tạo Post mới

```protobuf
rpc CreatePost(PostSaveRequest) returns (Post) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/post"
        body: "*"
    };
}
```

**PostSaveRequest**:
```protobuf
message PostSaveRequest {
    uint64 productId = 1;         // Link tới Product (required)
    string title = 2;
    string content = 3;
    double price = 4;             // Giá đăng
    uint32 transactionType = 5;   // Bán/Thuê
    uint32 visibility = 6;        // Public/Internal/Private
    uint32 packageVisible = 7;    // 1=Thường, 2=VIP
    string expiredAt = 8;         // Ngày hết hạn
    uint32 numDate = 9;           // Số ngày đăng
    optional bool hidden = 10;
    repeated PostMedia mediaList = 11;
    uint32 status = 12;
    uint64 id = 13;               // ID (for update)
    uint32 priceType = 14;        // 1=Chính thức, 2=Thỏa thuận
}
```

---

#### 4.2.2. UpdatePost - Cập nhật Post

```protobuf
rpc UpdatePost(PostSaveRequest) returns (Post) {
    option (google.api.http) = {
        put: "/v2/bdspro/v2/post/{id}"
        body: "*"
    };
}
```

---

#### 4.2.3. Delete - Xóa Post

```protobuf
rpc Delete(IdRequest) returns (Response) {
    option (google.api.http) = {
        delete: "/v2/bdspro/v2/post/{id}"
    };
}
```

---

### 4.3. Query APIs

#### 4.3.1. GetDetail - Chi tiết Post

```protobuf
rpc GetDetail(IdRequest) returns (PostDetailResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/post/detail/{id}"
    };
}
```

**PostDetailResponse**:
```protobuf
message PostDetailResponse {
    Post post = 1;
    ProductResponse productInfo = 2;  // Thông tin Product
    Asset assetInfo = 3;               // Thông tin Asset (nếu có)
}
```

---

#### 4.3.2. GetPersonalPost - Posts cá nhân

```protobuf
rpc GetPersonalPost(PostSearchRequest) returns (PostListResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/post/personal"
    };
}
```

**PostSearchRequest**:
```protobuf
message PostSearchRequest {
    uint32 page = 1;
    uint32 size = 2;
    string content = 3;
    uint32 transactionType = 4;
    string text = 5;
    string title = 6;
    optional string expiredFrom = 7;
    optional string expiredTo = 8;
    optional bool hidden = 9;
    repeated uint64 status = 10;
    repeated uint64 types = 11;          // 1=Hiển thị, 2=Hết hạn, 3=Ẩn
    repeated uint64 packageVisibles = 12;
}
```

---

#### 4.3.3. GetGlobalPost - Posts công khai

```protobuf
rpc GetGlobalPost(PostSearchRequest) returns (PostListResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/post/global"
    };
}
```

---

#### 4.3.4. GetPublish - Posts được publish của user

```protobuf
rpc GetPublish(PostPublishSearch) returns (PostListResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/post/publish/{profileId}"
    };
}
```

---

### 4.4. Other Operations

#### 4.4.1. UpdateHidden - Ẩn/Hiện Post

```protobuf
rpc UpdateHidden(UpdateHiddenRequest) returns (Response) {
    option (google.api.http) = {
        put: "/v2/bdspro/v2/post/hidden"
        body: "*"
    };
}
```

**UpdateHiddenRequest**:
```protobuf
message UpdateHiddenRequest {
    uint64 postId = 1;
    bool hidden = 2;
}
```

---

#### 4.4.2. NewExpiredPost - Gia hạn Post

```protobuf
rpc NewExpiredPost(PostExpiredRequest) returns (Post) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/post/new-expired"
        body: "*"
    };
}
```

**PostExpiredRequest**:
```protobuf
message PostExpiredRequest {
    uint64 postId = 1;
    uint64 numDay = 2;  // Số ngày gia hạn
}
```

---

#### 4.4.3. CreatePostOrganization - Tạo Post cho Organization

```protobuf
rpc CreatePostOrganization (PostSaveRequest) returns (Post) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/post/organization"
        body: "*"
    };
}
```

---

#### 4.4.4. CurrentOrganizationPosts - Posts của Organization

```protobuf
rpc CurrentOrganizationPosts (PostSearchRequest) returns (PostListResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/post/organization"
    };
}
```

---

## 5. TransactionService - Quản lý Giao dịch

**Package**: `bdspropb.TransactionService`

**Purpose**: CRUD giao dịch mua bán/cho thuê.

### 5.1. CRUD Operations

#### 5.1.1. CreateTransaction - Tạo giao dịch mới

```protobuf
rpc CreateTransaction(CreateTransactionRequest) returns (Transaction) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/transaction/create"
        body: "*"
    };
}
```

**CreateTransactionRequest**:
```protobuf
message CreateTransactionRequest {
    uint32 ownerId = 1;
    optional uint32 organizationId = 2;
    double amount = 3;
    string currency = 4;
    string transactionName = 5;
    string description = 6;
    uint32 transactionType = 7;      // Bán/Thuê
    uint32 categoryId = 8;
    uint32 paymentMethodId = 9;
    optional uint32 relatedDealId = 10;
    optional uint32 relatedTransactionId = 11;
    string transactionDate = 12;
}
```

---

#### 5.1.2. UpdateTransaction - Cập nhật giao dịch

```protobuf
rpc UpdateTransaction(UpdateTransactionRequest) returns (Transaction) {
    option (google.api.http) = {
        put: "/v2/bdspro/v2/transaction/update/{id}"
        body: "*"
    };
}
```

---

#### 5.1.3. DeleteTransaction - Xóa giao dịch

```protobuf
rpc DeleteTransaction(DeleteTransactionRequest) returns (DeleteTransactionResponse) {
    option (google.api.http) = {
        delete: "/v2/bdspro/v2/transaction/delete/{id}"
    };
}
```

---

### 5.2. Query APIs

#### 5.2.1. ListTransactions - Danh sách giao dịch

```protobuf
rpc ListTransactions(ListTransactionsRequest) returns (ListTransactionsResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/transaction/list"
    };
}
```

**ListTransactionsRequest**:
```protobuf
message ListTransactionsRequest {
    optional string transactionType = 1;
    optional string approvalStatus = 2;
    optional uint32 transactionStatus = 3;
    optional string startDate = 4;
    optional string endDate = 5;
    optional int32 page = 6;
    optional int32 size = 7;
}
```

---

#### 5.2.2. GetTransactionDetail - Chi tiết giao dịch

```protobuf
rpc GetTransactionDetail(GetTransactionDetailRequest) returns (Transaction) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/transaction/detail/{id}"
    };
}
```

---

### 5.3. Approval Operations

#### 5.3.1. ApproveTransaction - Duyệt giao dịch

```protobuf
rpc ApproveTransaction(ApproveTransactionRequest) returns (Transaction) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/transaction/approve/{id}"
    };
}
```

---

#### 5.3.2. RejectTransaction - Từ chối giao dịch

```protobuf
rpc RejectTransaction(RejectTransactionRequest) returns (Transaction) {
    option (google.api.http) = {
        post: "/v2/bdspro/v2/transaction/reject/{id}"
    };
}
```

---

### 5.4. Reference Data

#### 5.4.1. GetTransactionTypes - Lấy danh sách loại giao dịch

```protobuf
rpc GetTransactionTypes(GetTransactionTypesRequest) returns (GetTransactionTypesResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/transaction/types"
    };
}
```

---

## 6. BdsproPublicService - Public APIs

**Package**: `bdspropb.BdsproPublicService`

**Purpose**: APIs công khai, không cần authentication (hoặc authentication nhẹ).

### 6.1. GetProductMarket - Products trên thị trường

```protobuf
rpc GetProductMarket(ProductMarketRequest) returns (ProductMarketResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/product/market"
    };
}
```

**ProductMarketRequest** (filter công khai):
```protobuf
message ProductMarketRequest {
    uint32 page = 1;
    uint32 size = 2;
    string sort = 3;
    string text = 4;
    uint64 provinceId = 5;
    uint64 districtId = 6;
    uint64 wardId = 7;
    int32 transactionType = 14;
    int32 saleStatus = 15;
    int32 rentStatus = 16;
    // ... price filters
}
```

**ProductMarketResponse** (denormalized data):
```protobuf
message ProductMarketResponse {
    repeated ProductMarket data = 1;
    int64 total = 2;
}
```

**ProductMarket** struct chứa tất cả thông tin denormalized (xem §19 trong bdspro.md).

---

### 6.2. Reference Data APIs

#### 6.2.1. GetPropertyType - Danh sách loại BĐS

```protobuf
rpc GetPropertyType(SearchQueryRequest) returns (ListItemResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/list/property-type"
    };
}
```

---

#### 6.2.2. GetRegion - Danh sách địa giới

```protobuf
rpc GetRegion(SearchQueryRequest) returns (ListItemResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/list/region"
    };
}
```

**SearchQueryRequest**:
```protobuf
message SearchQueryRequest {
    uint32 page = 1;
    uint32 size = 2;
    string text = 3;
    optional uint64 parentId = 4;  // Parent region (để lấy quận/phường)
}
```

---

#### 6.2.3. GetDocType - Danh sách loại giấy tờ

```protobuf
rpc GetDocType(SearchQueryRequest) returns (ListItemResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/list/doc-type"
    };
}
```

---

#### 6.2.4. GetProject - Danh sách dự án

```protobuf
rpc GetProject(SearchQueryRequest) returns (ListItemResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/list/project"
    };
}
```

---

#### 6.2.5. GetAmenity - Danh sách tiện ích

```protobuf
rpc GetAmenity(SearchQueryRequest) returns (ListItemResponse) {
    option (google.api.http) = {
        get: "/v2/bdspro/v2/list/amenity"
    };
}
```

---

## 7. Admin Services - APIs quản trị

BDSPro cũng có các Admin Services riêng:

1. **AdminProductService** - Admin APIs cho Product
2. **AdminAssetService** - Admin APIs cho Asset
3. **AdminPostService** - Admin APIs cho Post
4. **AdminProjectService** - Admin APIs cho Project
5. **AdminTransactionService** - Admin APIs cho Transaction

Các Admin APIs thường có quyền cao hơn, ví dụ:
- Xem tất cả records (không filter theo owner)
- Approve/Reject operations
- Bulk operations
- Analytics & Reports

---

## 8. Common Types & DTOs

### 8.1. MediaItem

```protobuf
message MediaItem {
    uint64 id = 1;
    string mediaUrl = 2;
    string mediaType = 3;      // "image" | "video"
    bool isMain = 4;            // Ảnh đại diện
    int32 order = 5;            // Thứ tự hiển thị
    int32 targetType = 6;       // Loại target
    uint64 targetId = 7;        // ID target
}
```

---

### 8.2. IdRequest (common)

```protobuf
message IdRequest {
    uint64 id = 1;
    uint32 page = 2;
    uint32 size = 3;
    string sort = 4;
    string text = 5;
    uint32 ownerType = 6;
    repeated uint64 ids = 7;
}
```

---

### 8.3. Response (common)

```protobuf
message Response {
    uint64 id = 1;
    string message = 2;
}
```

---

### 8.4. ArchivedRequest (common)

```protobuf
message ArchivedRequest {
    uint64 id = 1;
    bool archived = 2;
}
```

---

## 9. Quy tắc Call Internal APIs

### 9.1. Authentication & Context

**Tất cả APIs đều cần JWT token** trong metadata:

```go
// Client setup
ctx := context.Background()
md := metadata.Pairs("authorization", "Bearer "+token)
ctx = metadata.NewOutgoingContext(ctx, md)

// Call API
resp, err := client.GetDetail(ctx, &bdspropb.IdRequest{Id: 123})
```

**Context keys tự động**:
- `profileId`: Từ JWT
- `organizationId`: Từ JWT
- `roleId`: Từ JWT

---

### 9.2. Error Handling

APIs trả về gRPC Status codes:

```go
resp, err := client.CreateProduct(ctx, req)
if err != nil {
    st, ok := status.FromError(err)
    if ok {
        code := st.Code()        // codes.InvalidArgument, codes.NotFound, ...
        message := st.Message()
    }
}
```

**Common error codes**:
- `codes.InvalidArgument` (4001): Invalid request
- `codes.NotFound` (4002): Entity not found
- `codes.PermissionDenied` (4003): No permission
- `codes.Internal` (5001): Internal server error

---

### 9.3. Pagination Pattern

**Request**:
```protobuf
page = 1       // Page number (1-indexed)
size = 20      // Items per page (max 100)
```

**Response**:
```protobuf
data = [...]           // Array of items
total/totalElements = 150  // Total count
```

**Calculate**:
```go
totalPages := (total + size - 1) / size
hasNext := page < totalPages
```

---

### 9.4. Filter Pattern

**Multiple conditions = AND**:
```go
filter := &bdspropb.ProductSearch{
    PropertyTypeId: 1,  // AND loại BĐS = 1
    ProvinceId: 79,     // AND tỉnh = 79
    SaleStatus: 20,     // AND đang bán
}
```

**Empty/Zero value = Ignore filter**:
```go
filter := &bdspropb.ProductSearch{
    PropertyTypeId: 0,  // Không filter theo property type
    ProvinceId: 79,     // Filter theo tỉnh
}
```

---

### 9.5. Soft Delete Pattern

**Tất cả Delete APIs = Soft Delete**:
- Set `deleted_at = NOW()`
- Query tự động filter `deleted_at IS NULL`

**Permanent delete**: Chỉ admin mới có quyền.

---

## 10. Use Cases & Examples

### 10.1. Service khác cần thông tin Products

**Scenario**: News feed service cần hiển thị thông tin Product trong bài viết.

```go
// 1. Get Product by ID
client := bdspropb.NewProductServiceClient(conn)
product, err := client.GetDetail(ctx, &bdspropb.IdRequest{Id: productID})

// 2. Get multiple Products by IDs
products, err := client.GetDetailByIds(ctx, &bdspropb.IdRequest{
    Ids: []uint64{1, 2, 3, 4, 5},
})
```

---

### 10.2. Organization service cần thống kê

**Scenario**: Organization service cần đếm số BĐS của member.

```go
client := bdspropb.NewBdsproInternalServiceClient(conn)
resp, err := client.GetCountByOwner(ctx, &bdspropb.GetCountByOwnerRequest{
    OwnerOf: sharepb.OwnerOf_MEMBER,
    OwnerId: memberID,
})

fmt.Printf("Products: %d, Assets: %d, Posts: %d\n", 
    resp.TotalProduct, resp.TotalAsset, resp.TotalPost)
```

---

### 10.3. Tạo Product với Post cùng lúc

**Scenario**: User tạo Product và đăng tin luôn.

```go
client := bdspropb.NewProductServiceClient(conn)
resp, err := client.ProductNew(ctx, &bdspropb.ProductSaveRequest{
    Name: "Nhà phố Quận 1",
    Area: 120.5,
    ProvinceId: &provinceID,
    DistrictId: &districtID,
    WardId: &wardID,
    PriceData: &bdspropb.ProductPrice{
        SalePrice: &salePrice,
        Currency: "VND",
    },
    Post: &bdspropb.PostSaveWithProduct{
        Title: "Bán nhà phố Q1 giá tốt",
        Content: "...",
        Visibility: 10, // Public
        NumDate: 30,
    },
})
```

---

### 10.4. Split Product thành nhiều phần

**Scenario**: Tách đất lớn thành nhiều lô nhỏ.

```go
client := bdspropb.NewProductServiceClient(conn)
resp, err := client.ChildDevideProduct(ctx, &bdspropb.DevideChildRequest{
    ParentId: parentProductID,
    Childs: []*bdspropb.SaveChildRequest{
        {
            Name: "Lô A1",
            Area: 60.0,
            // ... other fields
        },
        {
            Name: "Lô A2",
            Area: 40.0,
            // ... other fields
        },
    },
})
```

---

### 10.5. Query Products có filter phức tạp

**Scenario**: Tìm nhà phố ở Q1, đang bán, giá 5-10 tỷ.

```go
client := bdspropb.NewProductServiceClient(conn)
resp, err := client.ProductMe(ctx, &bdspropb.ProductSearch{
    Page: 1,
    Size: 20,
    PropertyTypeId: 1,      // Nhà phố
    DistrictId: 760,        // Quận 1
    SaleStatus: 20,         // Đang bán
    SalePrice: 5000000000,  // Min 5 tỷ (simplified)
    // ... more filters
})

for _, product := range resp.Data {
    fmt.Printf("Product: %s - %s\n", product.Code, product.Name)
}
```

---

## 11. Best Practices

### 11.1. Khi nào dùng Internal APIs?

✅ **Nên dùng**:
- Service khác cần data từ bdspro
- Cross-service communication
- Data enrichment (lấy thông tin bổ sung)
- Analytics & Reporting

❌ **Không nên dùng**:
- Frontend gọi trực tiếp (dùng gateway)
- Có sẵn API trong service hiện tại

---

### 11.2. Performance Tips

**1. Batch queries khi có thể**:
```go
// ✅ Good: 1 call
products := client.GetDetailByIds(ctx, &IdRequest{Ids: ids})

// ❌ Bad: N calls
for _, id := range ids {
    product := client.GetDetail(ctx, &IdRequest{Id: id})
}
```

**2. Chỉ lấy fields cần thiết**:
- Dùng `GetDetailByIds` thay vì `GetDetail` nếu không cần full info
- Filter sớm, pagination đúng

**3. Cache khi phù hợp**:
- Reference data (PropertyType, Region, DocType)
- Rarely changed data

---

### 11.3. Error Handling Best Practices

```go
resp, err := client.CreateProduct(ctx, req)
if err != nil {
    st, ok := status.FromError(err)
    if !ok {
        // Non-gRPC error
        return err
    }
    
    switch st.Code() {
    case codes.InvalidArgument:
        // Validation error - show to user
        return fmt.Errorf("invalid input: %s", st.Message())
    case codes.NotFound:
        // Entity not found
        return fmt.Errorf("product not found")
    case codes.PermissionDenied:
        // No permission
        return fmt.Errorf("no permission")
    default:
        // Internal error
        return fmt.Errorf("internal error: %s", st.Message())
    }
}
```

---

### 11.4. Context & Timeout

**Luôn set timeout**:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := client.GetDetail(ctx, req)
```

**Pass metadata đúng**:
```go
md := metadata.Pairs(
    "authorization", "Bearer "+token,
    "x-request-id", requestID,
)
ctx = metadata.NewOutgoingContext(ctx, md)
```

---

## 12. Migration & Versioning

### 12.1. API Versioning

**Current**: `/v2/bdspro/v2/...`

**Breaking changes**: Tạo v3, maintain v2 trong 6 tháng.

**Non-breaking changes**: Update in-place.

---

### 12.2. Protobuf Compatibility

**✅ Safe changes**:
- Add new fields
- Add new methods
- Add new enum values

**❌ Breaking changes**:
- Remove fields
- Change field types
- Remove methods
- Change method signatures

---

## Tổng kết

BDSPro Internal Services cung cấp **comprehensive APIs** cho việc quản lý bất động sản:

### ✅ Services Coverage:

1. **BdsproInternalService**: Simple queries cho services khác
2. **ProductService**: Full CRUD + Split/Merge + AI Suggest
3. **AssetService**: Asset management + Split/Merge
4. **PostService**: Post management + Publishing
5. **TransactionService**: Transaction tracking + Approval
6. **BdsproPublicService**: Public market + Reference data
7. **Admin Services**: Admin operations

### ✅ Key Features:

- **Rich filtering**: Search với nhiều conditions
- **Batch operations**: GetByIds, BulkCreate, ...
- **Split/Merge**: Product & Asset operations
- **Access control**: Owner check, Sharing, Permissions
- **Soft delete**: Safe delete operations
- **Pagination**: Consistent pagination pattern
- **Error handling**: gRPC status codes chuẩn

### ✅ Khi nào gọi bdspro APIs:

- **Organization service**: Get member's BĐS count
- **News feed service**: Enrich posts với Product info
- **Analytics service**: Collect metrics
- **Notification service**: Notify về BĐS events
- **Search service**: Index BĐS data

File này phục vụ cho:
- ✅ Developers implement cross-service calls
- ✅ API documentation reference
- ✅ Integration testing
- ✅ Client SDK generation

