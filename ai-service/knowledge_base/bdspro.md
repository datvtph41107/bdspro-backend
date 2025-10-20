# BDSPro Domain Knowledge - Chi tiết Domain Entities

Tài liệu này phân tích chi tiết tất cả các domain entities trong **bdspro-service** - service quản lý bất động sản (Real Estate Management).

---

## 1. Product Domain - Sản phẩm BĐS

**Mục đích**: Quản lý thông tin sản phẩm bất động sản (listings) - đại diện cho BĐS được đăng bán/cho thuê.

### 1.1. Product Entity

**Table**: `products`

**Embed**: `_models.BaseEntity` (ID, CreatedAt, UpdatedAt, DeletedAt, CreatedBy, UpdatedBy)

### 1.2. Core Fields

#### Thông tin cơ bản:
```go
Name           string   // Tên sản phẩm
Code           string   // Mã sản phẩm (SP000001, SP000002...)
Area           float64  // Diện tích (m²)
Description    string   // Mô tả chi tiết
Note           string   // Ghi chú
```

#### Phân loại sản phẩm:
```go
PropertyTypeId *uint64        // ID loại BĐS (nhà, đất, căn hộ...)
PropertyType   *PropertyTypeItem  // Relation: loại BĐS
DocTypeId      *uint64        // ID loại giấy tờ pháp lý
DocType        *DocTypeItem   // Relation: sổ đỏ/hồng/chưa có sổ
CategoryID     *uint64        // Danh mục sản phẩm
```

#### Liên kết với entities khác:
```go
ParentId       *uint64        // Sản phẩm cha (nếu là sản phẩm con)
AssetId        *uint64        // Liên kết tới tài sản (Asset)
ApartmentID    *uint64        // Liên kết tới căn hộ
Apartment      *Apartment     // Relation
PostIds        []uint64       // Danh sách bài đăng (Post IDs)
```

### 1.3. Pricing Fields

```go
LastPriceID    *uint64        // ID của bảng giá hiện tại
Price          *ProductPrice  // Relation: thông tin giá
```

**Lưu ý**: Giá được lưu trong bảng riêng `product_price` để tracking lịch sử thay đổi giá.

### 1.4. Location Fields - Địa chỉ

```go
ProvinceID     *uint64   // ID tỉnh/thành
Province       *Region   // Relation
DistrictID     *uint64   // ID quận/huyện
District       *Region   // Relation
WardID         *uint64   // ID phường/xã
Ward           *Region   // Relation
Address        string    // Địa chỉ chi tiết
GoogleMapLink  string    // Link Google Maps
PositionUrl    string    // URL vị trí

// Computed fields (không lưu DB)
ProvinceName   string `gorm:"-"`
DistrictName   string `gorm:"-"`
WardName       string `gorm:"-"`
```

### 1.5. Transaction Status - Trạng thái giao dịch

#### Sale (Bán):
```go
SaleStatus        enums.EProductSaleStatus  // Trạng thái bán
SaleVisibility    enums.EVisibility         // Phạm vi hiển thị khi bán
SaleTransactionID *uint64                   // ID giao dịch bán

// Enum values:
// - EProductNotSold (10): Chưa bán
// - EProductSelling (20): Đang bán
// - EProductSold    (30): Đã bán
```

#### Rent (Cho thuê):
```go
RentStatus        enums.EProductRentStatus  // Trạng thái cho thuê
RentVisibility    enums.EVisibility         // Phạm vi hiển thị khi cho thuê
RentTransactionID *uint64                   // ID giao dịch thuê

// Enum values:
// - EProductNotRent (10): Chưa cho thuê
// - EProductRenting (20): Đang cho thuê
// - EProductRented  (30): Đã cho thuê
```

#### Transaction Type:
```go
TransactionType   uint  // 10=Bán, 20=Thuê (gorm default:10)
```

### 1.6. Ownership & Visibility - Quyền sở hữu

```go
OwnerID        uint64         // ID người sở hữu
OwnerOf        enums.EOwnerOf // Loại owner (Member/Group/Organization)
Visibility     enums.EVisibility // Phạm vi hiển thị (Public/Internal/Private)

// Visibility enums:
// - EVisiblePublic   (10): Công khai
// - EVisibleInternal (20): Nội bộ
// - EVisiblePrivate  (30): Riêng tư
```

### 1.7. Additional Information

#### Media:
```go
MediaList      []ProductMediaItem  // Danh sách ảnh/video
```

#### House Info (cho nhà ở):
```go
HouseInfo      *HouseInfo  // Thông tin nhà: số phòng, WC, tầng...
```

#### Amenities (tiện ích):
```go
Amenities      []AmenityItem  // Many-to-Many: gần trường, chợ, bệnh viện...
```

#### Private Data (dữ liệu nội bộ):
```go
PrivateData    *ProductPrivate  // Giá nhập, hoa hồng, ghi chú nội bộ
```

#### Source & Archive:
```go
SourceType     enums.ESourceType  // Nguồn: import, manual, API...
Archived       bool               // Đã lưu trữ? (default: false)
```

### 1.8. Relations Summary

**Product** có quan hệ với:
- `PropertyType` (1:1) - Loại BĐS
- `DocType` (1:1) - Loại giấy tờ
- `Region` (1:3) - Tỉnh, Quận, Phường
- `ProductPrice` (1:1) - Giá hiện tại
- `ProductPrivate` (1:1) - Thông tin nội bộ
- `Apartment` (1:1) - Căn hộ (nếu có)
- `HouseInfo` (1:1) - Thông tin nhà
- `ProductMediaItem` (1:N) - Ảnh/video
- `AmenityItem` (N:N) - Tiện ích xung quanh
- `Asset` (N:1) - Tài sản gốc

---

## 2. Asset Domain - Tài sản thực tế

**Mục đích**: Quản lý tài sản BĐS thực sự mà công ty/cá nhân **sở hữu** (đã mua).

**Sự khác biệt với Product**:
- **Product**: BĐS để bán/cho thuê (listing) - có thể là của người khác
- **Asset**: BĐS mà mình đã mua và sở hữu (portfolio management)

### 2.1. Asset Entity

**Table**: `assets`

**Embed**: `_models.BaseEntity`

### 2.2. Core Fields

```go
Name             string   // Tên tài sản
Area             float64  // Diện tích (m²)
Address          string   // Địa chỉ chi tiết
Description      string   // Mô tả
```

### 2.3. Location Fields

```go
ProvinceID       *uint64   // Tỉnh/thành
Province         *Region
DistrictID       *uint64   // Quận/huyện
District         *Region
WardID           *uint64   // Phường/xã (required)
Ward             *Region
```

### 2.4. Purchase Information - Thông tin mua

```go
PurchasePrice    float64       // Giá mua (required)
PurchaseDate     *time.Time    // Ngày mua
```

### 2.5. Legal Status - Trạng thái pháp lý

```go
LegalStatus      enums.EDocType  // Sổ đỏ/hồng/chưa có (default: 10)
LegalItems       []AssetLegal    // Relation: Danh sách giấy tờ pháp lý
```

**AssetLegal** bao gồm:
- `DocumentName`: Tên giấy tờ
- `DocumentURL`: Link file scan
- `DocumentType`: Loại giấy tờ
- `IssuedDate`: Ngày cấp
- `ExpiryDate`: Ngày hết hạn
- `Description`: Mô tả

### 2.6. Asset Status - Trạng thái tài sản

```go
RentStatus       enums.EAssetStatus  // Trạng thái cho thuê (default: 10)
Archived         bool                // Đã lưu trữ (default: false)

// EAssetStatus values:
// - EAssetStatusOwning  (10): Đang sở hữu
// - EAssetStatusRenting (20): Đang cho thuê
// - EAssetStatusSold    (30): Đã bán
// - EAssetStatusSelling (40): Đang bán
```

### 2.7. Split/Merge Tracking - Tách/Gộp đất

```go
ParentAssetID    *uint64   // ID tài sản cha (nếu tách từ tài sản khác)
SplitMergeStatus string    // Trạng thái tách/gộp
```

**Use case**: Khi tách 1 thửa đất lớn thành nhiều thửa nhỏ, hoặc gộp nhiều thửa nhỏ thành 1 thửa lớn.

### 2.8. Link to Product

```go
ProductID        *uint64   // Link tới Product (nếu Asset này đang được rao bán)
Product          *Product
```

### 2.9. Property Type

```go
PropertyTypeId   *uint64
PropertyType     *PropertyTypeItem  // Loại BĐS
```

### 2.10. Ownership

```go
OwnerID          *uint64         // ID người sở hữu
OwnerType        enums.EOwnerOf  // Member/Group/Organization
Owner            *ProfileInfo    // Thông tin owner (computed)
```

### 2.11. Asset Exploitation - Khai thác tài sản

```go
AssetExploitations []AssetExploitation  // Danh sách hợp đồng khai thác
```

**AssetExploitation** bao gồm:
- `ContractName`: Tên hợp đồng
- `TypeID`: Loại khai thác (cho thuê, quảng cáo...)
- `ContractValue`: Giá trị hợp đồng
- `Customer`: Tên khách hàng
- `CyclePay`: Chu kỳ thanh toán
- `IncomeAmount`: Số tiền thu được
- `StartDate`, `EndDate`: Thời gian
- `FileContract`: File hợp đồng
- `Status`: Trạng thái

### 2.12. Asset List View

**AssetList** struct: View để query danh sách assets với thông tin computed:

```go
type AssetList struct {
    Asset
    PropertyTypeName string          // Tên loại BĐS
    ProvinceName     string          // Tên tỉnh
    DistrictName     string          // Tên quận
    WardName         string          // Tên phường
    Permissions      string          // Quyền truy cập
    LegalList        json.RawMessage // Raw JSON các giấy tờ
}
```

---

## 3. Post Domain - Bài đăng

**Mục đích**: Quản lý các bài đăng (tin rao) trên hệ thống - đại diện cho việc **công khai** Product ra thị trường.

**Quan hệ**: 1 Product có thể có nhiều Posts (đăng nhiều lần, các gói VIP khác nhau)

### 3.1. Post Entity

**Table**: `posts`

**Embed**: `_models.BaseEntity`

### 3.2. Core Fields

```go
Code            string         // Mã bài đăng (PO000001...)
Title           string         // Tiêu đề
Content         string         // Nội dung mô tả
ProductId       uint64         // Link tới Product (required)
Product         *Product       // Relation
```

### 3.3. Transaction Type

```go
TransactionType enums.TransactionType  // Bán hay Cho thuê
```

### 3.4. Pricing

```go
Price           float64   // Giá đăng trên tin (có thể khác giá trong Product)
PriceType       uint32    // 1: Chính thức, 2: Thỏa thuận
```

### 3.5. Post Status - Trạng thái bài đăng

```go
Status            enums.EPostStatus  // Trạng thái cơ bản (default: 20)

// EPostStatus values:
// - EPostActive   (10): Đang hoạt động
// - EPostPending  (20): Chờ duyệt
// - EPostApproved (30): Đã duyệt
// - EPostRejected (40): Từ chối
// - EPostHidden   (50): Ẩn
```

```go
TransactionStatus enums.EPostVisibleStatus  // Trạng thái giao dịch (computed)
VisibleStatus     enums.EPostVisibleStatus  // Trạng thái hiển thị (computed)

// EPostVisibleStatus values:
// - EPostTransactionSelling (100): Đang bán
// - EPostTransactionSold    (110): Đã bán
// - EPostTransactionRenting (200): Đang cho thuê
// - EPostTransactionRented  (210): Đã cho thuê
```

### 3.6. Visibility & Package

```go
Visibility      enums.EVisibility  // Phạm vi hiển thị (Public/Internal/Private)
PackageVisible  uint               // Gói tin: 1=Thường, 2=VIP
Hidden          bool               // Ẩn bài đăng (default: false)
```

### 3.7. Duration & Expiry

```go
NumDate         int          // Số ngày đăng
ExpiredAt       *time.Time   // Ngày hết hạn
```

### 3.8. Publishing

```go
PublishedAt     *time.Time   // Thời gian xuất bản
PublishedBy     uint64       // Người xuất bản
```

### 3.9. Ownership

```go
OwnerID         uint64         // ID người sở hữu
OwnerOf         enums.EOwnerOf // Member/Group/Organization (default: 10)
```

### 3.10. Engagement Metrics

```go
Like            uint64   // Số lượt thích
Comment         uint64   // Số bình luận
NumView         uint64   // Số lượt xem
```

### 3.11. Media

```go
MediaList       []PostMediaEntity  // Danh sách ảnh/video của bài đăng
```

**Lưu ý**: Media của Post có thể khác với Media của Product (chọn lọc ảnh đẹp hơn để đăng).

### 3.12. PostItem View

**PostItem** struct: View để query danh sách posts với thông tin computed:

```go
type PostItem struct {
    Post
    PostPrice             float64
    Address               uint64
    ProvinceName          string
    DistrictName          string
    WardName              string
    MediaListRaw          json.RawMessage
    Status                enums.EPostStatus
    TransactionStatus     enums.TransactionStatus
    VisibleStatus         enums.EPostVisibleStatus
    TransactionTypeName   string  // Tên loại giao dịch
    TransactionStatusName string  // Tên trạng thái giao dịch
    VisibleStatusName     string  // Tên trạng thái hiển thị
}
```

---

## 4. Project Domain - Dự án BĐS

**Mục đích**: Quản lý các dự án bất động sản (chung cư, khu đô thị...) do chủ đầu tư phát triển.

### 4.1. Project Entity

**Table**: `projects`

**Embed**: `_models.BaseEntity`

### 4.2. Core Fields

```go
Name        string   // Tên dự án
Description string   // Mô tả dự án
```

### 4.3. Developer

```go
DeveloperID uint64         // ID chủ đầu tư
Developer   *DeveloperItem // Relation
```

### 4.4. Builds - Các tòa nhà

```go
Builds      []ProjectBuildItem  // Danh sách các tòa/block
```

**ProjectBuild** Entity:
- **Table**: `project_build`
- **Fields**:
  ```go
  Name         string   // Tên tòa (A, B, C...)
  NumFloor     uint     // Số tầng
  NumApartment uint     // Tổng số căn hộ
  Note         string   // Ghi chú
  ProjectID    uint64   // Link tới Project
  ```
- **Relations**:
  ```go
  Project      *ProjectItem        // Relation
  Apartments   []ApartmentItem     // Danh sách căn hộ
  Attributes   []ApartmentAttrItem // Các loại căn hộ (2PN, 3PN...)
  ```

### 4.5. ProjectItem

**ProjectItem** struct: Lightweight version dùng cho relations:

```go
type ProjectItem struct {
    ID   uint64
    Name string
}
```

---

## 5. Apartment Domain - Căn hộ

**Mục đích**: Quản lý từng căn hộ cụ thể trong một tòa nhà của dự án.

### 5.1. Apartment Entity

**Table**: `apartment`

**Embed**: `_models.BaseEntity`

### 5.2. Core Fields

```go
Name            string  // Mã căn hộ (A-10-01)
Note            string  // Ghi chú
Floor           *int    // Tầng
Ordinal         int     // Thứ tự (not null)
```

### 5.3. Attribute

```go
ApartmentAttrID uint64  // ID loại căn hộ (2PN, 3PN...)
```

**ApartmentAttribute** định nghĩa loại căn hộ:
- Diện tích
- Số phòng ngủ
- Số WC
- Layout
- Giá bán dự kiến

### 5.4. Build Reference

```go
BuildID         uint64  // Thuộc tòa nào
```

### 5.5. Status

```go
Status          uint    // Trạng thái: Đã bán/Còn trống/Đang đặt cọc...
Archived        uint    // Đã lưu trữ
```

### 5.6. ApartmentItem

**ApartmentItem** struct: Lightweight version dùng cho relations:

```go
type ApartmentItem struct {
    ID              uint64
    Name            string
    Floor           *int
    Ordinal         int
    BuildID         uint64
    Status          uint
    Archived        uint
    ApartmentAttrID uint64
}
```

---

## 6. ProductPrice Domain - Giá sản phẩm

**Mục đích**: Tracking lịch sử thay đổi giá của Product. Mỗi lần cập nhật giá tạo 1 record mới.

### 6.1. ProductPrice Entity

**Table**: `product_price`

**Embed**: `_models.BaseEntity`

### 6.2. Core Fields

```go
ProductID        *uint64   // Link tới Product
Currency         string    // Đơn vị tiền tệ (VND, USD...)
```

### 6.3. Sale Price - Giá bán

```go
SalePrice        *float64  // Giá bán niêm yết (DECIMAL(15,0))
SaleCommission   *float64  // Hoa hồng bán (DECIMAL(15,0))
Deposite         *float64  // Tiền đặt cọc (DECIMAL(15,0))
```

### 6.4. Rent Price - Giá thuê

```go
RentPrice        *float64  // Giá thuê/tháng (DECIMAL(15,0))
RentCommission   *float64  // Hoa hồng thuê (DECIMAL(15,0))
RentPaymentCycle *uint32   // Chu kỳ thanh toán (tháng, quý, năm)
```

**Lưu ý**: 
- Tất cả giá đều là pointer để có thể NULL
- DECIMAL(15,0) để tránh lỗi làm tròn với số lớn
- Mỗi Product có `LastPriceID` trỏ tới giá hiện tại

---

## 7. ProductPrivate Domain - Thông tin nội bộ

**Mục đích**: Lưu thông tin nội bộ/bí mật của Product (chỉ agent/admin xem được).

### 7.1. ProductPrivate Entity

**Table**: `product_private`

**Embed**: `_models.BaseEntity`

### 7.2. Core Fields

```go
ProductId     uint64    // Link 1:1 tới Product
```

### 7.3. Internal Pricing

```go
ImportPrice   *float64  // Giá nhập/giá chủ nhà muốn net
OperatingCost *float64  // Chi phí vận hành (DECIMAL(15,0))
TargetProfit  *float64  // Lợi nhuận mục tiêu
```

### 7.4. Internal Note

```go
InternalNote  string    // Ghi chú nội bộ môi giới
```

### 7.5. Private Documents

```go
PrivateDocs   string    // Danh sách URL tài liệu nội bộ (text)
```

**Helper methods**:
```go
func (p *ProductPrivate) SetPrivateDocs(docs []string) {
    p.PrivateDocs = strings.Join(docs, ",")
}

func (p *ProductPrivate) GetPrivateDocs() []string {
    if p.PrivateDocs == "" {
        return []string{}
    }
    return strings.Split(p.PrivateDocs, ",")
}
```

**Lưu ý**: PrivateDocs lưu dạng CSV trong DB, convert khi đọc/ghi.

---

## 8. HouseInfo Domain - Thông tin nhà

**Mục đích**: Lưu thông tin chi tiết về nhà ở (số phòng, tầng, hướng...).

### 8.1. HouseInfo Entity

**Table**: `house_info`

**Embed**: `_models.BaseEntityNotId` (không có ID riêng, dùng ProductID làm PK)

### 8.2. Core Fields

```go
ProductID   uint64  // Primary Key & Foreign Key (not null)
```

### 8.3. Room Information

```go
NumBedroom  *int32  // Số phòng ngủ
NumBathroom *int32  // Số WC/phòng tắm
NumToilet   *int32  // Số toilet riêng
NumFloor    *int32  // Số tầng
```

### 8.4. Parking & Front

```go
NumFront    *int32  // Số mặt tiền
NumCarPark  *int32  // Số chỗ đậu xe ô tô
```

### 8.5. Furniture & Orientation

```go
Furniture   *string  // Nội thất: "full", "partial", "none" (varchar(50))
Orientation *string  // Hướng nhà: Đông, Tây, Nam, Bắc... (varchar(50))
```

**Lưu ý**: Tất cả fields đều là pointer (optional) vì không phải nhà nào cũng có đủ thông tin.

---

## 9. Region Domain - Địa giới hành chính

**Mục đích**: Quản lý cây địa giới: Tỉnh > Quận/Huyện > Phường/Xã.

### 9.1. Region Entity

**Table**: `region`

**Embed**: `_models.BaseEntity`

### 9.2. Core Fields

```go
Name     string   // Tên (TP.HCM, Quận 1, Phường Bến Nghé...)
Code     uint     // Mã số hành chính
CodeName string   // Tên không dấu (tphcm, quan1, phuong-ben-nghe)
```

### 9.3. Hierarchy

```go
Level    uint     // Cấp: 1=Tỉnh, 2=Quận/Huyện, 3=Phường/Xã
ParentID *uint64  // ID của cấp cha
Unit     string   // Đơn vị: "Thành phố", "Quận", "Phường"...
```

**Cách query**:
- Lấy tất cả Tỉnh: `WHERE level = 1`
- Lấy Quận của Tỉnh X: `WHERE level = 2 AND parent_id = X`
- Lấy Phường của Quận Y: `WHERE level = 3 AND parent_id = Y`

---

## 10. PropertyType Domain - Loại BĐS

**Mục đích**: Danh mục các loại bất động sản (Nhà riêng, Căn hộ, Đất nền, Biệt thự...).

### 10.1. PropertyType Entity

**Table**: `property_type`

**Embed**: `_models.BaseEntity`

### 10.2. Core Fields

```go
ID     uint64  // Primary Key
Name   string  // Tên loại BĐS (varchar(100), not null)
Active bool    // Còn sử dụng? (default: true)
```

### 10.3. PropertyTypeItem

**PropertyTypeItem** struct: Lightweight version dùng cho relations:

```go
type PropertyTypeItem struct {
    ID   uint64
    Name string
}
```

**Các loại phổ biến**:
- Nhà riêng
- Căn hộ/Chung cư
- Đất nền
- Biệt thự
- Nhà mặt phố
- Shophouse
- Văn phòng
- Kho xưởng

---

## 11. DocType Domain - Loại giấy tờ pháp lý

**Mục đích**: Danh mục các loại giấy tờ sở hữu BĐS.

### 11.1. DocType Entity

**Table**: `doc_type`

**Embed**: `_models.BaseEntity`

### 11.2. Core Fields

```go
ID     uint64  // Primary Key
Name   string  // Tên loại giấy tờ (varchar(100), not null)
Active bool    // Còn sử dụng? (default: true)
```

### 11.3. DocTypeItem

**DocTypeItem** struct: Lightweight version dùng cho relations:

```go
type DocTypeItem struct {
    ID   uint64
    Name string
}
```

**Các loại phổ biến**:
- Sổ đỏ (Sổ đỏ chính chủ)
- Sổ hồng
- Giấy tờ hợp lệ khác
- Chưa có sổ
- Hợp đồng mua bán
- Giấy phép xây dựng

---

## 12. Amenity Domain - Tiện ích xung quanh

**Mục đích**: Danh mục các tiện ích gần Product (trường học, chợ, bệnh viện...).

### 12.1. Amenity Entity

**Table**: `amenity`

**Embed**: `_models.BaseEntity`

### 12.2. Core Fields

```go
Name   string  // Tên tiện ích (varchar(100), unique, not null)
Active bool    // Còn hiển thị? (default: true)
```

### 12.3. Many-to-Many Relation

```go
// Product <-> Amenity: Many-to-Many
type ProductAmenity struct {
    ProductID uint64  // Primary Key
    AmenityID uint64  // Primary Key
}
```

**Table**: `product_amenity`

### 12.4. AmenityItem

**AmenityItem** struct: Lightweight version:

```go
type AmenityItem struct {
    ID   uint64
    Name string
}
```

**Các amenity phổ biến**:
- Gần trường học
- Gần chợ/siêu thị
- Gần bệnh viện
- Gần trung tâm thương mại
- Gần công viên
- Gần bến xe/ga tàu
- An ninh tốt
- Khu vực yên tĩnh

---

## 13. Developer Domain - Chủ đầu tư

**Mục đích**: Quản lý thông tin các chủ đầu tư dự án BĐS.

### 13.1. Developer Entity

**Table**: `developer`

**Embed**: `_models.BaseEntity`

### 13.2. Core Fields

```go
Name          string      // Tên chủ đầu tư
Slug          string      // URL-friendly slug (size: 15)
Description   string      // Mô tả
```

### 13.3. Branding

```go
LogoUrl       string  // URL logo
BackgroundUrl string  // URL ảnh nền/cover
```

### 13.4. Contact Information

```go
Address       string      // Địa chỉ công ty (size: 15)
Phone         string      // Số điện thoại (size: 15)
Email         string      // Email (size: 255)
Website       string      // Website
```

### 13.5. Foundation

```go
FoundedAt     *time.Time  // Ngày thành lập
Status        uint        // Trạng thái hoạt động
```

### 13.6. DeveloperItem

**DeveloperItem** struct: Lightweight version:

```go
type DeveloperItem struct {
    ID     uint64
    Name   string
    Status uint
}
```

---

## 14. ProductMedia Domain - Media của sản phẩm

**Mục đích**: Quản lý ảnh/video của Product.

### 14.1. ProductMedia Entity

**Table**: `product_media`

**Embed**: `_models.BaseEntity`

### 14.2. Core Fields

```go
ProductID uint64  // Link tới Product (not null)
MediaURL  string  // URL ảnh/video (varchar(500), not null)
MediaType string  // "image" hoặc "video" (varchar(20), not null)
```

### 14.3. Main Image

```go
IsMain    bool  // Ảnh đại diện? (default: false)
```

**Lưu ý**: Chỉ có 1 media có `IsMain = true` cho mỗi Product.

### 14.4. Ordering

```go
Order     int   // Thứ tự hiển thị
```

### 14.5. ProductMediaItem

**ProductMediaItem** struct: Lightweight version:

```go
type ProductMediaItem struct {
    ID        uint64
    ProductID uint64 `gorm:"-"` // Không select
    MediaURL  string `json:"mediaUrl"`
    MediaType string `json:"mediaType"`
    IsMain    bool   `json:"isMain"`
    Order     int    `json:"sort_order"`
}
```

---

## 15. Sharing Domain - Chia sẻ quyền truy cập

**Mục đích**: Quản lý việc chia sẻ Product/Asset/Post cho user/group/org khác.

### 15.1. Sharing Entity

**Table**: `sharings`

**Embed**: `_models.BaseEntity`

### 15.2. Target - Đối tượng được chia sẻ

```go
TargetID    uint64            // ID đối tượng (Product/Asset/Post)
TargetType  enums.ETargetType // Loại đối tượng
```

### 15.3. Owner - Người sở hữu

```go
OwnerID     uint64         // ID người sở hữu
OwnerType   enums.EOwnerOf // Member/Group/Organization
```

### 15.4. Share To - Chia sẻ cho ai

```go
ShareID     uint64  // ID người được chia sẻ
ShareType   string  // Loại: "member", "group", "organization"
```

### 15.5. Permissions - Quyền hạn

```go
Permissions string  // Danh sách quyền (CSV string)
```

**Helper methods**:
```go
func (p *Sharing) SetFields(fields []string) {
    p.Permissions = strings.Join(fields, ",")
}

func (p *Sharing) GetFields() []string {
    if p.Permissions == "" {
        return []string{}
    }
    return strings.Split(p.Permissions, ",")
}
```

**Permissions có thể**:
- `view`: Xem
- `edit`: Sửa
- `delete`: Xóa
- `share`: Chia sẻ tiếp

**Example**: `"view,edit"` → Được xem và sửa

---

## 16. PostMedia Domain - Media của bài đăng

**Mục đích**: Quản lý ảnh/video của Post (có thể khác với ProductMedia).

### 16.1. PostMediaEntity

**Table**: `post_media`

**Embed**: `_models.BaseEntity`

### 16.2. Core Fields

```go
PostID    uint64  // Link tới Post
MediaURL  string  // URL ảnh/video
MediaType string  // "image" hoặc "video"
IsMain    bool    // Ảnh đại diện
Order     int     // Thứ tự
```

**Lưu ý**: Tương tự ProductMedia nhưng link tới Post thay vì Product.

---

## 17. AssetLegal Domain - Giấy tờ pháp lý tài sản

**Mục đích**: Lưu trữ thông tin các giấy tờ pháp lý của Asset.

### 17.1. AssetLegal Entity

**Table**: `asset_legals`

**Embed**: `_models.BaseEntity`

### 17.2. Core Fields

```go
AssetID             uint64      // Link tới Asset (required)
DocumentName        string      // Tên giấy tờ
DocumentURL         string      // URL file scan
DocumentType        string      // Loại giấy tờ (sổ đỏ, sổ hồng...)
```

### 17.3. Dates

```go
IssuedDate          *time.Time  // Ngày cấp
ExpiryDate          *time.Time  // Ngày hết hạn
```

### 17.4. Description

```go
Description         string  // Mô tả chi tiết
```

### 17.5. Split/Merge Tracking

```go
RelatedSplitMergeID *uint64  // ID liên quan đến tách/gộp đất
```

**Use case**: Khi tách đất, giấy tờ mới sẽ link về ID của split/merge transaction.

---

## 18. AssetExploitation Domain - Khai thác tài sản

**Mục đích**: Quản lý các hợp đồng khai thác tài sản (cho thuê, quảng cáo, hợp tác...).

### 18.1. AssetExploitation Entity

**Table**: `asset_exploitation`

**Embed**: `_models.BaseEntity`

### 18.2. Core Fields

```go
AssetID       uint64  // Link tới Asset (required)
ContractName  string  // Tên hợp đồng (required)
```

### 18.3. Type

```go
TypeID        *uint64        // ID loại khai thác (required)
AssetType     *AssetCostType // Relation: cho thuê, quảng cáo, hợp tác...
```

### 18.4. Contract Value

```go
ContractValue float64  // Giá trị hợp đồng (required)
IncomeAmount  float64  // Doanh thu thực tế
```

### 18.5. Customer

```go
Customer      string  // Tên khách hàng (required)
```

### 18.6. Payment Cycle

```go
CyclePay      uint  // Chu kỳ thanh toán (tháng, quý, năm)
```

### 18.7. Duration

```go
StartDate     *time.Time  // Ngày bắt đầu
EndDate       *time.Time  // Ngày kết thúc
```

### 18.8. Documents

```go
FileContract  string      // URL file hợp đồng
Description   string      // Mô tả chi tiết
```

### 18.9. Status & Ownership

```go
Status        uint           // Trạng thái (default: 10)
OwnerID       uint64         // Người sở hữu
OwnerType     enums.EOwnerOf // Member/Group/Organization (default: 10)
```

### 18.10. AssetExploitationQuery

**AssetExploitationQuery** struct: View với thông tin computed:

```go
type AssetExploitationQuery struct {
    _models.BaseEntity
    AssetID       uint64
    ContractName  string
    TypeID        *uint64
    TypeName      string      // Tên loại (computed)
    ContractValue float64
    Customer      string
    CyclePay      uint
    IncomeAmount  float64
    StartDate     *time.Time
    EndDate       *time.Time
    Description   string
    FileContract  string
    Status        uint
}
```

---

## 19. ProductMarket View - Market View cho Product

**Mục đích**: View để hiển thị danh sách Product trên thị trường với thông tin đầy đủ (denormalized).

### 19.1. ProductMarket View

**Table**: `product_market` (materialized view hoặc joined query)

### 19.2. Core Product Info

```go
ID               int64
CreatedAt        time.Time
UpdatedAt        time.Time
Name             string
Code             string
Area             *float64
Description      *string
Address          *string
GoogleMapLink    *string
CategoryID       *int64
```

### 19.3. Transaction Status

```go
TransactionType  *int32
SaleStatus       int16
SaleVisibility   int16
RentStatus       int16
RentVisibility   int16
TransactionPrice *float64
```

### 19.4. Links

```go
PropertyTypeID   *int64
DocTypeID        *int64
PositionURL      *string
AssetID          *int64
ApartmentID      *int64
LastPriceID      *int64
OwnerID          *int64
```

### 19.5. House Info (denormalized)

```go
NumBedroom  *int64
NumBathroom *int64
NumFloor    *int64
Furniture   *string
Orientation *string
NumFront    *int64
NumCarPark  *int64
NumToilet   *int64
```

### 19.6. Media (aggregated)

```go
AmenityIDs *string         // CSV: "1,2,5"
MediaURLs  pq.StringArray  // Array: ["url1", "url2"]
MediaTypes pq.StringArray  // Array: ["image", "video"]
MainImage  *string         // URL ảnh chính
```

### 19.7. Region Names (denormalized)

```go
ProvinceName *string  // Tên tỉnh (computed)
DistrictName *string  // Tên quận (computed)
WardName     *string  // Tên phường (computed)
```

### 19.8. Type Names (denormalized)

```go
DocTypeName      *string  // Tên loại giấy tờ (computed)
PropertyTypeName *string  // Tên loại BĐS (computed)
```

### 19.9. Latest Price Info (denormalized)

```go
Currency           *string
SalePrice          *float64
SaleCommission     *int64
SaleCommissionType *int64
RentPrice          *float64
RentCommission     *float64
RentCommissionType *int64
RentPaymentCycle   *int64
ImportPrice        *int64
OperatingCost      *float64
InternalNote       *int64
TargetProfit       *int64
PrivateDocs        *float64
Deposite           *float64
```

### 19.10. Flags

```go
DeletedAt        *time.Time
Archived         bool
```

**Lưu ý**: View này dùng để query performance tốt hơn, không cần join nhiều bảng.

---

## 20. Domain Relationships Summary

### 20.1. Core Relationships

```
Product (1) --- (1) ProductPrice      [LastPriceID]
Product (1) --- (1) ProductPrivate    [ProductId]
Product (1) --- (1) HouseInfo         [ProductID]
Product (1) --- (N) ProductMedia      [ProductID]
Product (N) --- (N) Amenity           [product_amenity table]
Product (1) --- (1) Apartment         [ApartmentID]
Product (N) --- (1) Asset             [AssetId]
Product (N) --- (1) PropertyType      [PropertyTypeId]
Product (N) --- (1) DocType           [DocTypeId]
Product (N) --- (1) Region            [ProvinceID, DistrictID, WardID]

Post (N) --- (1) Product              [ProductId]
Post (1) --- (N) PostMedia            [PostID]

Asset (1) --- (N) AssetLegal          [AssetID]
Asset (1) --- (N) AssetExploitation   [AssetID]
Asset (N) --- (1) PropertyType        [PropertyTypeId]
Asset (N) --- (1) Region              [ProvinceID, DistrictID, WardID]

Project (1) --- (N) ProjectBuild      [ProjectID]
Project (N) --- (1) Developer         [DeveloperID]

ProjectBuild (1) --- (N) Apartment    [BuildID]
ProjectBuild (N) --- (1) Project      [ProjectID]

Region (Tree) --- Self                [ParentID]

Sharing (N) --- (1) Product/Asset/Post [TargetID + TargetType]
```

### 20.2. Ownership Pattern

Tất cả các entities chính đều có:
```go
OwnerID    uint64         // ID người/nhóm/tổ chức sở hữu
OwnerType  enums.EOwnerOf // 10=Member, 20=Group, 30=Organization
```

### 20.3. Visibility Pattern

Product và Post có:
```go
Visibility  enums.EVisibility  // 10=Public, 20=Internal, 30=Private
```

### 20.4. Status Pattern

- **Product**: `SaleStatus`, `RentStatus` (Not/Selling/Sold, Not/Renting/Rented)
- **Asset**: `RentStatus` (Owning/Renting/Sold/Selling)
- **Post**: `Status` (Active/Pending/Approved/Rejected/Hidden)

### 20.5. Audit Pattern

Tất cả entities embed `BaseEntity` có:
```go
ID         uint64
CreatedAt  time.Time
UpdatedAt  time.Time
DeletedAt  *time.Time  // Soft delete
CreatedBy  *uint64     // Auto từ context
UpdatedBy  *uint64     // Auto từ context
```

---

## 21. Common Query Patterns

### 21.1. Get Products for Sale (công khai)

```go
db.Where("sale_status = ? AND sale_visibility = ? AND deleted_at IS NULL", 
    enums.EProductSelling, enums.EVisiblePublic).
    Preload("PropertyType").
    Preload("Province").
    Preload("District").
    Preload("Ward").
    Preload("MediaList").
    Find(&products)
```

### 21.2. Get Products by Owner

```go
db.Where("owner_id = ? AND owner_of = ? AND deleted_at IS NULL", 
    ownerID, enums.EOwnerOfMember).
    Find(&products)
```

### 21.3. Get Assets with Legal Docs

```go
db.Preload("LegalItems").
    Preload("PropertyType").
    Preload("Province").
    Where("owner_id = ? AND deleted_at IS NULL", ownerID).
    Find(&assets)
```

### 21.4. Get Posts by Product

```go
db.Where("product_id = ? AND status IN (?) AND deleted_at IS NULL", 
    productID, []int{enums.EPostActive, enums.EPostApproved}).
    Preload("MediaList").
    Order("published_at DESC").
    Find(&posts)
```

### 21.5. Get Region Hierarchy

```go
// Lấy tất cả tỉnh
db.Where("level = 1 AND deleted_at IS NULL").Find(&provinces)

// Lấy quận của tỉnh
db.Where("level = 2 AND parent_id = ? AND deleted_at IS NULL", provinceID).Find(&districts)

// Lấy phường của quận
db.Where("level = 3 AND parent_id = ? AND deleted_at IS NULL", districtID).Find(&wards)
```

### 21.6. Get Project with Builds and Apartments

```go
db.Preload("Developer").
    Preload("Builds.Apartments").
    Preload("Builds.Attributes").
    Where("id = ? AND deleted_at IS NULL", projectID).
    First(&project)
```

### 21.7. Get Product Market View (denormalized)

```go
db.Table("product_market").
    Where("deleted_at IS NULL AND archived = false").
    Where("sale_status = ?", enums.EProductSelling).
    Where("sale_visibility = ?", enums.EVisiblePublic).
    Order("created_at DESC").
    Limit(20).
    Find(&productMarkets)
```

---

## 22. Enum Values Reference

### 22.1. EVisibility - Phạm vi hiển thị

```go
EVisiblePublic   = 10  // Công khai (mọi người xem)
EVisibleInternal = 20  // Nội bộ (tổ chức xem)
EVisiblePrivate  = 30  // Riêng tư (chỉ owner xem)
```

### 22.2. EOwnerOf - Loại owner

```go
EOwnerOfMember      = 10  // Thành viên (cá nhân trong org)
EOwnerOfGroup       = 20  // Nhóm
EOwnerOfOrgnization = 30  // Tổ chức
EOwnerOfUser        = 40  // Người dùng (ngoài org)
```

### 22.3. EProductSaleStatus - Trạng thái bán

```go
EProductNotSold = 10  // Chưa bán
EProductSelling = 20  // Đang bán
EProductSold    = 30  // Đã bán
```

### 22.4. EProductRentStatus - Trạng thái cho thuê

```go
EProductNotRent = 10  // Chưa cho thuê
EProductRenting = 20  // Đang cho thuê
EProductRented  = 30  // Đã cho thuê
```

### 22.5. EAssetStatus - Trạng thái tài sản

```go
EAssetStatusOwning  = 10  // Đang sở hữu
EAssetStatusRenting = 20  // Đang cho thuê
EAssetStatusSold    = 30  // Đã bán
EAssetStatusSelling = 40  // Đang bán
```

### 22.6. EPostStatus - Trạng thái bài đăng

```go
EPostActive   = 10  // Đang hoạt động
EPostPending  = 20  // Chờ duyệt
EPostApproved = 30  // Đã duyệt
EPostRejected = 40  // Từ chối
EPostHidden   = 50  // Ẩn
```

### 22.7. EPostVisibleStatus - Trạng thái giao dịch post

```go
EPostTransactionSelling = 100  // Đang bán
EPostTransactionSold    = 110  // Đã bán
EPostTransactionRenting = 200  // Đang cho thuê
EPostTransactionRented  = 210  // Đã cho thuê
```

---

## 23. Business Logic Notes

### 23.1. Product vs Asset

**Khi nào tạo Product?**
- Khi có BĐS cần rao bán/cho thuê (có thể là của người khác)
- Product có thể không liên kết tới Asset (listing thuần túy)

**Khi nào tạo Asset?**
- Khi công ty/cá nhân **mua** BĐS về sở hữu
- Asset là tài sản thực tế, cần quản lý chi phí, khai thác, pháp lý

**Liên kết**:
- Asset có thể tạo Product để bán lại (`Product.AssetId`)
- Nhiều Product có thể cùng link tới 1 Asset (bán nhiều lần)

### 23.2. Product vs Post

**Product**: 
- Là thông tin BĐS (master data)
- Không đổi nhiều theo thời gian

**Post**:
- Là tin đăng (ad campaign)
- 1 Product có nhiều Posts (đăng lại nhiều lần, gói VIP khác nhau)
- Có giá đăng, thời gian, package
- Có engagement metrics (like, view, comment)

### 23.3. Price History

**Mỗi lần thay đổi giá**:
1. Tạo record mới trong `product_price`
2. Update `Product.LastPriceID` trỏ tới record mới
3. Record cũ vẫn giữ (history)

### 23.4. Split/Merge Asset

**Tách đất** (Split):
1. Asset cha: Set `SplitMergeStatus = "split"`
2. Tạo nhiều Asset con với `ParentAssetID = asset_cha.ID`
3. Asset con inherit thông tin từ cha
4. Legal docs mới link về split transaction

**Gộp đất** (Merge):
1. Nhiều Asset con: Set `SplitMergeStatus = "merged"`
2. Tạo 1 Asset cha mới
3. Legal docs mới link về merge transaction

### 23.5. Sharing Logic

**Khi share Product/Asset/Post**:
1. Tạo record trong `sharings`
2. Set `TargetID`, `TargetType`, `ShareID`, `ShareType`
3. Set `Permissions`: "view", "edit", "delete", "share"
4. Người được share có thể truy cập theo permissions

**Check quyền**:
```go
// Trong repo/usecase
query := db.Where("(owner_id = ? AND owner_of = ?) OR id IN (SELECT target_id FROM sharings WHERE share_id = ? AND target_type = ? AND deleted_at IS NULL)", 
    ownerID, ownerOf, userID, "product")
```

### 23.6. Visibility Logic

**Public (10)**:
- Hiển thị trên website/app công khai
- Mọi người xem được

**Internal (20)**:
- Chỉ thành viên trong Organization xem được
- Không hiển thị ra ngoài

**Private (30)**:
- Chỉ Owner và người được share xem được
- Không hiển thị cho ai khác

---

## 24. Data Consistency Rules

### 24.1. Product Consistency

**Rule 1**: Nếu `Product.AssetId` != null, Asset phải tồn tại và chưa bị xóa.

**Rule 2**: Nếu `Product.LastPriceID` != null, ProductPrice phải tồn tại.

**Rule 3**: HouseInfo chỉ có khi PropertyType là "Nhà riêng" hoặc "Biệt thự".

**Rule 4**: Apartment chỉ có khi PropertyType là "Căn hộ/Chung cư".

**Rule 5**: `SaleStatus` và `RentStatus` không thể cùng lúc = "Sold" và "Rented".

### 24.2. Post Consistency

**Rule 1**: `Post.ProductId` phải trỏ tới Product tồn tại.

**Rule 2**: Nếu Product đã bán (`SaleStatus = Sold`), không tạo Post mới với `TransactionType = Bán`.

**Rule 3**: `ExpiredAt` phải > `PublishedAt`.

### 24.3. Asset Consistency

**Rule 1**: `PurchasePrice` > 0.

**Rule 2**: Nếu `ParentAssetID` != null, Asset cha phải tồn tại và có `SplitMergeStatus = "split"`.

**Rule 3**: `RentStatus = Renting` thì phải có ít nhất 1 AssetExploitation active.

### 24.4. Region Consistency

**Rule 1**: Level 1 (Tỉnh) phải có `ParentID = null`.

**Rule 2**: Level 2 (Quận) phải có `ParentID` trỏ tới Tỉnh (Level 1).

**Rule 3**: Level 3 (Phường) phải có `ParentID` trỏ tới Quận (Level 2).

---

## 25. Common Validation Rules

### 25.1. Product Validation

```go
// Trong BeforeSave hook hoặc Validator

// 1. Required fields
if product.Name == "" {
    return errors.New("name is required")
}
if product.Area <= 0 {
    return errors.New("area must be > 0")
}

// 2. Status logic
if product.SaleStatus == enums.EProductSelling && product.SaleVisibility == 0 {
    return errors.New("sale_visibility required when selling")
}

// 3. Price consistency
if product.LastPriceID != nil {
    var price ProductPrice
    db.First(&price, product.LastPriceID)
    if price.ID == 0 {
        return errors.New("invalid price reference")
    }
}
```

### 25.2. Asset Validation

```go
// 1. Required fields
if asset.Name == "" {
    return errors.New("name is required")
}
if asset.PurchasePrice <= 0 {
    return errors.New("purchase_price must be > 0")
}
if asset.Area <= 0 {
    return errors.New("area must be > 0")
}
if asset.WardID == nil {
    return errors.New("ward_id is required")
}

// 2. Split/Merge logic
if asset.ParentAssetID != nil {
    var parent Asset
    db.First(&parent, asset.ParentAssetID)
    if parent.SplitMergeStatus != "split" {
        return errors.New("parent asset must have split status")
    }
}
```

### 25.3. Post Validation

```go
// 1. Required fields
if post.Title == "" {
    return errors.New("title is required")
}
if post.ProductId == 0 {
    return errors.New("product_id is required")
}
if post.NumDate <= 0 {
    return errors.New("num_date must be > 0")
}

// 2. Product exists and valid
var product Product
db.First(&product, post.ProductId)
if product.ID == 0 {
    return errors.New("product not found")
}

// 3. Transaction type match
if post.TransactionType == TransactionTypeSale && product.SaleStatus != EProductSelling {
    return errors.New("product not available for sale")
}
```

---

## 26. Transaction Domain - Giao dịch BĐS

**Mục đích**: Quản lý các giao dịch mua bán/cho thuê bất động sản trong hệ thống.

### 26.1. Transaction Entity

**Table**: `transactions`

**Embed**: `_models.BaseEntity`

### 26.2. Core Fields

```go
TransactionName      string   // Tên giao dịch (varchar(255), not null)
Description          string   // Mô tả chi tiết (text)
TransactionDate      time.Time // Ngày giao dịch (timestamp)
```

### 26.3. Ownership

```go
OwnerId              uint64         // ID người sở hữu (int8, not null)
OwnerType            enums.EOwnerOf // Member/Group/Organization (default: 10)
OrganizationId       *uint64        // ID tổ chức (int8)
```

### 26.4. Amount & Currency

```go
Amount               float64  // Số tiền (decimal(20,2), not null)
Currency             string   // Đơn vị tiền tệ (varchar(10), default: 'VND')
```

### 26.5. Transaction Type

```go
TransactionType      enums.TransactionType  // Loại giao dịch (default: 10)

// Enum values:
// - TransactionTypeSale (10): Bán
// - TransactionTypeRent (20): Thuê
```

### 26.6. Transaction Status

```go
TransactionStatus    enums.TransactionStatus  // Trạng thái (default: 10)

// Enum values:
// - TransactionStatusDraft     (10): Nháp
// - TransactionStatusDeposite  (20): Đặt cọc
// - TransactionStatusSigned    (30): Đã ký hợp đồng
// - TransactionStatusCompleted (40): Hoàn thành
// - TransactionStatusCanceled  (50): Hủy
```

### 26.7. Approval Status

```go
ApprovalStatus       enums.ApprovalStatus  // Trạng thái duyệt (default: 'pending')
ApprovedBy           *uint32               // Người duyệt (int4)
ApprovalDate         *time.Time            // Ngày duyệt (timestamp)
```

### 26.8. Category & Payment Method

```go
CategoryId           uint32  // Danh mục (int4)
PaymentMethodId      uint32  // Phương thức thanh toán (int4)
```

### 26.9. Related Entities

```go
ProductId            *uint64   // Link tới Product (int8)
Product              *Product  // Relation
ContactId            *uint64   // Link tới Contact/Customer (int8)
```

### 26.10. Related Deals & Transactions

```go
RelatedDealId        *uint32  // Deal liên quan (int4)
RelatedTransactionId *uint32  // Giao dịch liên quan (int4)
```

**Use case**: Link các giao dịch liên quan (ví dụ: thuê dài hạn có nhiều giao dịch thanh toán).

### 26.11. Deposite Information

```go
DepositeAmount       *float64  // Số tiền đặt cọc (decimal(20,2))
DepositeNote         *string   // Ghi chú đặt cọc (text)
```

### 26.12. Status & Flags

```go
Status               uint32  // Trạng thái tổng quát (default: 10)
IsDeleted            bool    // Đã xóa? (default: false)
```

### 26.13. Helper Methods

```go
func (t *Transaction) GetTransactionTypeName() string {
    switch t.TransactionType {
    case enums.TransactionTypeSale:
        return "Bán"
    case enums.TransactionTypeRent:
        return "Thuê"
    default:
        return "Không xác định"
    }
}

func (t *Transaction) GetTransactionStatusName() string {
    switch t.TransactionStatus {
    case enums.TransactionStatusDraft:
        return "Nháp"
    case enums.TransactionStatusDeposite:
        return "Đặt cọc"
    case enums.TransactionStatusSigned:
        return "Đã ký"
    case enums.TransactionStatusCompleted:
        return "Hoàn thành"
    case enums.TransactionStatusCanceled:
        return "Hủy"
    default:
        return "Không xác định"
    }
}
```

### 26.14. Transaction Flow

**Typical Sale Flow**:
1. **Draft** → Tạo giao dịch nháp
2. **Deposite** → Khách đặt cọc
3. **Signed** → Ký hợp đồng
4. **Completed** → Hoàn tất thanh toán, bàn giao
5. **Canceled** (optional) → Hủy giao dịch

**Typical Rent Flow**:
1. **Draft** → Tạo giao dịch thuê
2. **Deposite** → Đặt cọc + tiền thuê tháng đầu
3. **Signed** → Ký hợp đồng thuê
4. **Completed** → Bàn giao nhà
5. Tạo các giao dịch thanh toán tiền thuê hàng tháng (RelatedTransactionId)

---

## 27. ProductHistory Domain - Lịch sử sản phẩm

**Mục đích**: Tracking lịch sử thay đổi của Product (tách/gộp/chuyển đổi).

### 27.1. ProductHistory Entity

**Table**: `product_history`

**Embed**: `_models.BaseEntity`

### 27.2. Core Fields

```go
ID          uint64                // Primary Key
Action      enums.EProductHistory // Loại hành động (int)
```

**EProductHistory Actions**:
- Split (tách sản phẩm)
- Merge (gộp sản phẩm)
- Divide (chia nhỏ)
- Convert (chuyển đổi)

### 27.3. Parent Product

```go
ParentID    *uint64       // ID sản phẩm cha
Parent      *ProductInfo  // Relation: thông tin sản phẩm cha
```

### 27.4. Target Product

```go
ProductID   *uint64       // ID sản phẩm đích (kết quả)
Target      *ProductInfo  // Relation: sản phẩm kết quả
```

### 27.5. Child Products

```go
ChildIds    pq.Int64Array   // Mảng ID các sản phẩm con (bigint[])
Childs      []ProductInfo   // Relation: danh sách sản phẩm con (computed)
Areas       pq.Float64Array // Mảng diện tích tương ứng (double[])
```

**Use case**: Khi tách 1 sản phẩm lớn thành nhiều sản phẩm nhỏ.

### 27.6. Audit

```go
CreatedUser *ProfileInfo  // Người thực hiện (from CreatedBy)
```

### 27.7. ProductInfo Helper Struct

```go
type ProductInfo struct {
    ID          uint64   // Primary Key
    Name        string   // Tên (size: 255)
    Code        string   // Mã (size: 15)
    CategoryID  *uint64  // Danh mục
    Area        float64  // Diện tích
    LastPriceID *uint64  // Giá hiện tại
}
```

**Lưu ý**: ProductInfo là lightweight version của Product dùng cho history.

### 27.8. Use Cases

#### Split Product (Tách):
```go
// Tách 1 Product 100m² thành 2 Products: 60m² và 40m²
{
    Action:    Split,
    ParentID:  &productA.ID,
    ChildIds:  [productB.ID, productC.ID],
    Areas:     [60.0, 40.0],
}
```

#### Merge Products (Gộp):
```go
// Gộp 2 Products thành 1 Product mới
{
    Action:    Merge,
    ChildIds:  [productA.ID, productB.ID],
    ProductID: &productC.ID,  // Product mới
    Areas:     [60.0, 40.0],
}
```

---

## 28. Access Control Domains - Quản lý quyền truy cập

### 28.1. AssetUser - User có quyền trên Asset

**Mục đích**: Quản lý quyền của user/member đối với Asset.

**Table**: `asset_user`

**Embed**: `_models.BaseEntity`

#### Fields:
```go
ProfileID uint64  // ID user/member (not null, indexed)
AssetID   uint64  // ID asset (not null, indexed)
IsOwner   bool    // Là chủ sở hữu? (default: false)
RoleID    uint64  // ID vai trò (default: 0)
```

**Composite Index**: (ProfileID, AssetID) - Unique constraint

---

### 28.2. AssetOrganization - Organization có quyền trên Asset

**Mục đích**: Quản lý quyền của organization đối với Asset.

**Table**: `asset_organization`

**Embed**: `_models.BaseEntity`

#### Fields:
```go
OrganizationID uint64  // ID organization (not null, indexed)
AssetID        uint64  // ID asset (not null, indexed)
IsOwner        bool    // Tổ chức sở hữu? (default: false)
RoleID         uint64  // ID vai trò (default: 0)
```

---

### 28.3. ProductUser - User có quyền trên Product

**Table**: `product_user`

**Embed**: `_models.BaseEntity`

#### Fields:
```go
ProfileID uint64  // ID user (not null, indexed)
ProductID uint64  // ID product (not null, indexed)
IsOwner   bool    // Là chủ sở hữu? (default: false)
RoleID    uint64  // ID vai trò (default: 0)
```

---

### 28.4. ProductOrganization - Organization có quyền trên Product

**Table**: `product_organization`

**Embed**: `_models.BaseEntity`

#### Fields:
```go
OrganizationID uint64  // ID organization (not null, indexed)
ProductID      uint64  // ID product (not null, indexed)
IsOwner        bool    // Tổ chức sở hữu? (default: false)
RoleID         uint64  // ID vai trò (default: 0)
```

---

### 28.5. PostUser - User có quyền trên Post

**Table**: `post_user`

**Embed**: `_models.BaseEntity`

#### Fields:
```go
ProfileID uint64  // ID user (not null, indexed)
PostID    uint64  // ID post (not null, indexed)
IsOwner   bool    // Là chủ sở hữu? (default: false)
RoleID    uint64  // ID vai trò (default: 0)
```

---

### 28.6. PostOrganization - Organization có quyền trên Post

**Table**: `post_organization`

**Embed**: `_models.BaseEntity`

#### Fields:
```go
OrganizationID uint64  // ID organization (not null, indexed)
PostID         uint64  // ID post (not null, indexed)
IsOwner        bool    // Tổ chức sở hữu? (default: false)
RoleID         uint64  // ID vai trò (default: 0)
```

---

### 28.7. Access Control Pattern

**Query quyền của user**:
```go
// Check user có quyền trên Product không
var access ProductUser
db.Where("profile_id = ? AND product_id = ? AND deleted_at IS NULL", 
    profileID, productID).
    First(&access)

if access.ID == 0 {
    return errors.New("no access")
}
```

**Query tất cả Products user có quyền**:
```go
var productIDs []uint64
db.Model(&ProductUser{}).
    Where("profile_id = ? AND deleted_at IS NULL", profileID).
    Pluck("product_id", &productIDs)

var products []Product
db.Where("id IN (?) AND deleted_at IS NULL", productIDs).
    Find(&products)
```

**RoleID mapping**:
- `0`: No specific role (guest/viewer)
- `1`: Viewer - Xem
- `2`: Editor - Sửa
- `3`: Manager - Quản lý
- `4`: Owner - Chủ sở hữu (IsOwner = true)

---

## 29. SharingAccess Domain - Chia sẻ chi tiết

**Mục đích**: Quản lý việc chia sẻ quyền truy cập với kiểm soát fields và commission.

### 29.1. SharingAccess Entity

**Table**: `sharing_access`

**Embed**: `_models.BaseEntity`

### 29.2. Domain & Target

```go
DomainID    uint64              // ID đối tượng
Domain      enums.EDomainAccess // Loại domain (default: 10)

// EDomainAccess values:
// - Product (10)
// - Asset (20)
// - Post (30)
// - Transaction (40)
```

### 29.3. From (Người chia sẻ)

```go
FromType    enums.EOwnerOf  // Loại người chia sẻ (default: 10)
FromId      uint64          // ID người chia sẻ
```

### 29.4. To (Người nhận)

```go
ToType      enums.EOwnerOf  // Loại người nhận
ToId        uint64          // ID người nhận
```

### 29.5. Fields & Permissions

```go
Fields      string   // Các fields được xem (CSV string)
Permissions string   // Quyền hạn (CSV string)
```

**Helper methods**:
```go
func (p *SharingAccess) SetFields(fields []string) {
    p.Fields = strings.Join(fields, ",")
}

func (p *SharingAccess) GetFields() []string {
    if p.Fields == "" {
        return []string{}
    }
    return strings.Split(p.Fields, ",")
}
```

**Fields example**: `"name,address,area,price"` → Chỉ xem được các field này

**Permissions example**: `"view,edit"` → Có quyền xem và sửa

### 29.6. Commission

```go
Commission  float32  // Hoa hồng (% hoặc số tiền)
```

**Use case**: Khi chia sẻ Product cho agent khác, set commission họ được hưởng khi bán thành công.

### 29.7. ProductAccessQuery View

```go
type ProductAccessQuery struct {
    ID         uint64   // ID sharing
    DomainID   uint64   // Product ID
    Fields     string   // Các fields được xem
    Commission float32  // Hoa hồng
    ToId       uint64   // ID người được chia sẻ
    ToType     uint     // Loại: 1=cá nhân, 2=nhóm, 3=tổ chức
    ToName     string   // Tên người được chia sẻ (computed)
    Avatar     string   // Avatar (computed)
}
```

**Use case**: Query để hiển thị danh sách người được chia sẻ Product.

---

## 30. AssetShare Domain - Chia sẻ Asset đơn giản

**Mục đích**: Chia sẻ Asset với permissions đơn giản (không có fields control).

### 30.1. AssetShare Entity

**Table**: `asset_shares`

**No BaseEntity** - struct đơn giản

### 30.2. Core Fields

```go
AssetID     uint64  // ID asset (required)
TargetID    uint64  // ID người được chia sẻ (required)
TypeShare   uint32  // Loại chia sẻ (required)
Permissions string  // Quyền hạn (CSV string)
```

### 30.3. Helper Methods

```go
func (p *AssetShare) SetFields(fields []string) {
    p.Permissions = strings.Join(fields, ",")
}

func (p *AssetShare) GetFields() []string {
    if p.Permissions == "" {
        return []string{}
    }
    return strings.Split(p.Permissions, ",")
}
```

**Lưu ý**: Đơn giản hơn SharingAccess, không có commission và fields control.

---

## 31. RecordHistory Domain - Lịch sử thay đổi

**Mục đích**: Tracking lịch sử thay đổi của tất cả các loại records.

### 31.1. RecordHistory Entity

**Table**: `record_history`

**Embed**: `_models.BaseEntity`

### 31.2. Core Fields

```go
RecordID   *uint64              // ID record bị thay đổi
RecordType enums.ERecordType    // Loại record
ActionType enums.EHistoryAction // Loại hành động
Content    pq.StringArray       // Nội dung thay đổi (string[])
```

### 31.3. ERecordType Enum

```go
// Các loại records có thể track:
- Product
- Asset
- Post
- Transaction
- Project
- Apartment
```

### 31.4. EHistoryAction Enum

```go
// Các hành động:
- Create
- Update
- Delete
- Split
- Merge
- Share
- Transfer
```

### 31.5. Content Array

**Content** chứa mảng string mô tả thay đổi:

```go
content := []string{
    "Changed price from 5,000,000,000 to 5,500,000,000",
    "Updated status from Selling to Sold",
    "Added media: image1.jpg, image2.jpg",
}
```

**PostgreSQL Array**: Sử dụng `pq.StringArray` để lưu mảng trong PostgreSQL.

---

## 32. PaymentMethod Domain - Phương thức thanh toán

**Mục đích**: Danh mục các phương thức thanh toán.

### 32.1. PaymentMethod Entity

**Table**: `payment_methods` (inferred)

**No BaseEntity** - struct đơn giản

### 32.2. Core Fields

```go
Id         uint32               // Primary Key
Name       string               // Tên phương thức
MethodCode string               // Mã code (unique)
Status     enums.PaymentStatus  // Trạng thái
CreatedAt  time.Time            // Ngày tạo
UpdatedAt  time.Time            // Ngày cập nhật
```

### 32.3. Common Payment Methods

- **Cash**: Tiền mặt
- **BankTransfer**: Chuyển khoản
- **CreditCard**: Thẻ tín dụng
- **EWallet**: Ví điện tử (Momo, ZaloPay...)
- **Installment**: Trả góp
- **Mixed**: Kết hợp nhiều phương thức

---

## 33. AssetIncomeType Domain - Loại thu nhập từ tài sản

**Mục đích**: Danh mục các loại thu nhập khi khai thác Asset.

### 33.1. AssetIncomeType Entity

**Table**: `asset_income_types` (inferred)

**Embed**: `_models.BaseEntity`

### 33.2. Core Fields

```go
TypeName string  // Tên loại thu nhập
IsCustom bool    // Loại tự tạo? (user-defined)
```

### 33.3. Common Income Types

**System-defined (IsCustom = false)**:
- **Rent**: Thu nhập từ cho thuê
- **Advertising**: Thu nhập từ quảng cáo
- **Parking**: Thu nhập từ bãi đậu xe
- **Service**: Thu nhập dịch vụ

**User-defined (IsCustom = true)**:
- User có thể tạo các loại thu nhập riêng

---

## 34. Task Domain - Công việc

**Mục đích**: Quản lý công việc liên quan đến dự án.

### 34.1. Task Entity

**Table**: `tasks`

**No BaseEntity** - sử dụng GORM standard

### 34.2. Core Fields

```go
ID          uint           // Primary Key
Title       string         // Tiêu đề (size: 255, not null)
Description string         // Mô tả (text)
ProjectID   uint           // Link tới Project (not null, indexed)
Status      string         // Trạng thái (size: 50)
CreatedAt   time.Time      // Ngày tạo
UpdatedAt   time.Time      // Ngày cập nhật
DeletedAt   gorm.DeletedAt // Soft delete (indexed)
```

### 34.3. Task Status

Common statuses:
- `"pending"`: Chờ xử lý
- `"in_progress"`: Đang thực hiện
- `"completed"`: Hoàn thành
- `"cancelled"`: Hủy

**Lưu ý**: Status là string để flexible, có thể custom.

---

## 35. Additional Domains Summary

### 35.1. Deposite Domain (Deprecated)

**File**: `deposite.go` - **Đã bị comment out**

Thông tin đặt cọc giờ được lưu trong **Transaction** entity:
```go
Transaction.DepositeAmount
Transaction.DepositeNote
```

---

### 35.2. Customer Domain

**File**: `customer.go` - **Timeout khi đọc**

Có thể chứa thông tin khách hàng, nhưng cần verify lại.

---

### 35.3. Intent Domain

**File**: `intent.go` - **Timeout khi đọc**

Có thể là ý định mua/thuê của khách hàng.

---

### 35.4. Attribute Domains

**Files**: 
- `attribute.go`
- `attribute_value.go`
- `attribute_product.go`
- `apartment_attribute.go`

Các domain này quản lý thuộc tính động của Product/Apartment. Cần phân tích thêm.

---

### 35.5. Cost Domains

**Files**:
- `asset_cost.go`
- `cost_type.go`
- `cost_document.go`
- `income_document.go`

Quản lý chi phí và thu nhập từ Asset. Cần phân tích thêm.

---

## 36. Updated Relationships Summary

### 36.1. Transaction Relationships

```
Transaction (N) --- (1) Product      [ProductId]
Transaction (N) --- (1) Contact      [ContactId]
Transaction (N) --- (1) PaymentMethod [PaymentMethodId]
Transaction (N) --- (1) Transaction  [RelatedTransactionId] (self-referencing)
```

### 36.2. History Relationships

```
ProductHistory (N) --- (1) Product   [ParentID, ProductID]
RecordHistory  (N) --- (1) Any       [RecordID + RecordType] (polymorphic)
```

### 36.3. Access Control Relationships

```
AssetUser          (N) --- (1) Asset   [AssetID]
AssetUser          (N) --- (1) User    [ProfileID]

AssetOrganization  (N) --- (1) Asset  [AssetID]
AssetOrganization  (N) --- (1) Org    [OrganizationID]

ProductUser        (N) --- (1) Product [ProductID]
ProductOrganization(N) --- (1) Product [ProductID]

PostUser           (N) --- (1) Post   [PostID]
PostOrganization   (N) --- (1) Post   [PostID]
```

### 36.4. Sharing Relationships

```
SharingAccess (N) --- (1) Product/Asset/Post/Transaction [DomainID + Domain]
AssetShare    (N) --- (1) Asset                          [AssetID]
```

---

## 37. Updated Query Patterns

### 37.1. Get Transactions by Product

```go
db.Where("product_id = ? AND deleted_at IS NULL", productID).
    Preload("Product").
    Order("transaction_date DESC").
    Find(&transactions)
```

### 37.2. Get Product History

```go
db.Where("parent_id = ? OR product_id = ? AND deleted_at IS NULL", 
    productID, productID).
    Preload("Parent").
    Preload("Target").
    Preload("CreatedUser").
    Order("created_at DESC").
    Find(&histories)
```

### 37.3. Check User Access to Asset

```go
// Check AssetUser
var assetUser AssetUser
db.Where("profile_id = ? AND asset_id = ? AND deleted_at IS NULL", 
    profileID, assetID).
    First(&assetUser)

if assetUser.ID > 0 {
    // Has access
}
```

### 37.4. Get Shared Products

```go
db.Table("sharing_access").
    Where("to_id = ? AND to_type = ? AND domain = ? AND deleted_at IS NULL", 
        userID, enums.EOwnerOfMember, enums.EDomainProduct).
    Pluck("domain_id", &productIDs)

var products []Product
db.Where("id IN (?)", productIDs).Find(&products)
```

### 37.5. Get Record History

```go
db.Where("record_id = ? AND record_type = ? AND deleted_at IS NULL", 
    recordID, enums.ERecordTypeProduct).
    Order("created_at DESC").
    Find(&histories)
```

---

## 38. Business Logic Notes - Extended

### 38.1. Transaction Lifecycle

**Create Transaction**:
1. Tạo Transaction với status = Draft
2. Link tới Product và Contact
3. Set DepositeAmount nếu có đặt cọc

**Update to Deposite**:
1. Khách đặt cọc
2. Update TransactionStatus = Deposite
3. Tạo payment record

**Complete Transaction**:
1. Update TransactionStatus = Completed
2. Update Product.SaleStatus = Sold (nếu bán)
3. Update Product.SaleTransactionID = transaction.ID

### 38.2. Product Split Flow

**Split Product**:
1. Tạo nhiều Product mới (childs)
2. Tạo ProductHistory với:
   - Action = Split
   - ParentID = product cũ
   - ChildIds = array IDs products mới
   - Areas = array diện tích tương ứng
3. Archive Product cũ (`Archived = true`)

### 38.3. Access Control Logic

**Check quyền**:
1. Check Owner: `Product.OwnerID == userID && Product.OwnerOf == EOwnerOfMember`
2. Check ProductUser: `EXISTS(ProductUser WHERE ProfileID = userID AND ProductID = productID)`
3. Check Sharing: `EXISTS(SharingAccess WHERE ToId = userID AND DomainID = productID)`

**Priority**: Owner > ProductUser > SharingAccess

### 38.4. History Tracking

**Khi thay đổi quan trọng**:
1. Tạo RecordHistory
2. Set RecordID, RecordType
3. Set ActionType (Create/Update/Delete...)
4. Set Content array mô tả thay đổi

**Content format**:
```go
[]string{
    "Field: OldValue → NewValue",
    "price: 5,000,000,000 → 5,500,000,000",
}
```

---

## Tổng kết (Updated)

BDSPro service quản lý **toàn bộ lifecycle** của bất động sản với **38 domains chính**:

### Core Domains (Entities chính):
1. ✅ **Product** (§1): Sản phẩm BĐS - listings để bán/cho thuê
2. ✅ **Asset** (§2): Tài sản thực tế đã mua và sở hữu
3. ✅ **Post** (§3): Bài đăng tin rao công khai
4. ✅ **Transaction** (§26): Giao dịch mua bán/cho thuê
5. ✅ **Project** (§4): Dự án BĐS (chung cư, khu đô thị)
6. ✅ **ProjectBuild** (§4): Tòa nhà trong dự án
7. ✅ **Apartment** (§5): Căn hộ cụ thể

### Support Domains (Entities hỗ trợ):
8. ✅ **ProductPrice** (§6): Giá & lịch sử giá sản phẩm
9. ✅ **ProductPrivate** (§7): Thông tin nội bộ/bí mật
10. ✅ **ProductMedia** (§14): Ảnh/video sản phẩm
11. ✅ **HouseInfo** (§8): Thông tin nhà (phòng, tầng, hướng)
12. ✅ **PostMedia** (§16): Ảnh/video bài đăng
13. ✅ **AssetLegal** (§17): Giấy tờ pháp lý tài sản
14. ✅ **AssetExploitation** (§18): Hợp đồng khai thác tài sản

### Reference Data (Danh mục):
15. ✅ **Region** (§9): Địa giới hành chính (Tỉnh/Quận/Phường)
16. ✅ **PropertyType** (§10): Loại BĐS (Nhà/Đất/Căn hộ...)
17. ✅ **DocType** (§11): Loại giấy tờ (Sổ đỏ/Sổ hồng...)
18. ✅ **Amenity** (§12): Tiện ích xung quanh
19. ✅ **Developer** (§13): Chủ đầu tư
20. ✅ **PaymentMethod** (§32): Phương thức thanh toán
21. ✅ **AssetIncomeType** (§33): Loại thu nhập từ tài sản

### Access Control (Quản lý quyền):
22. ✅ **AssetUser** (§28.1): Quyền user trên Asset
23. ✅ **AssetOrganization** (§28.2): Quyền org trên Asset
24. ✅ **ProductUser** (§28.3): Quyền user trên Product
25. ✅ **ProductOrganization** (§28.4): Quyền org trên Product
26. ✅ **PostUser** (§28.5): Quyền user trên Post
27. ✅ **PostOrganization** (§28.6): Quyền org trên Post

### Sharing (Chia sẻ):
28. ✅ **Sharing** (§15): Chia sẻ đơn giản
29. ✅ **SharingAccess** (§29): Chia sẻ chi tiết với fields & commission
30. ✅ **AssetShare** (§30): Chia sẻ Asset đơn giản

### History & Tracking (Lịch sử):
31. ✅ **ProductHistory** (§27): Lịch sử tách/gộp/chuyển đổi Product
32. ✅ **RecordHistory** (§31): Lịch sử thay đổi tất cả records

### View & Query (Views tối ưu):
33. ✅ **ProductMarket** (§19): View denormalized cho market
34. ✅ **AssetList** (§2.12): View cho danh sách Assets
35. ✅ **PostItem** (§3.12): View cho danh sách Posts
36. ✅ **ProductAccessQuery** (§29.7): View quyền truy cập Product

### Others (Khác):
37. ✅ **Task** (§34): Công việc liên quan dự án
38. ⏳ **Attribute Domains**: Thuộc tính động (chưa phân tích)
39. ⏳ **Cost Domains**: Chi phí tài sản (chưa phân tích)
40. ⏳ **Customer/Intent**: Khách hàng & ý định (chưa phân tích)

---

### Thiết kế Pattern (Clean Architecture):

✅ **BaseEntity Pattern**:
- Embed `_models.BaseEntity` → có ID, timestamps, CreatedBy, UpdatedBy
- Soft delete tự động (DeletedAt)
- Audit trail tự động

✅ **Ownership Pattern**:
```go
OwnerID   uint64         // ID người/nhóm/tổ chức
OwnerOf   enums.EOwnerOf // 10=Member, 20=Group, 30=Organization
```

✅ **Visibility Pattern**:
```go
Visibility  enums.EVisibility  // 10=Public, 20=Internal, 30=Private
```

✅ **Status Pattern**:
- Mỗi domain có enums riêng cho status
- Helper methods để get tên tiếng Việt

✅ **Relationship Pattern**:
- Foreign Key + Relation struct
- Preload để eager loading
- Indexes trên FK columns

✅ **Access Control Pattern**:
- Owner check → User access → Sharing
- RoleID: 0=Guest, 1=Viewer, 2=Editor, 3=Manager, 4=Owner

✅ **History Tracking Pattern**:
- RecordHistory với polymorphic relation
- ProductHistory cho split/merge operations
- Content array để mô tả chi tiết

---

### Quy tắc khi code cho bdspro-service:

🔒 **Security**:
- Luôn check quyền: Owner → Access → Sharing
- Filter theo organizationId khi cần
- Validate visibility trước khi return data

📊 **Query**:
- LUÔN check `deleted_at IS NULL`
- Use Preload cho relations
- Apply indexes trên FK và filter columns

🔄 **Transaction**:
- Wrap critical operations trong transaction
- Update status của related entities
- Create history records

✅ **Validation**:
- Validate required fields
- Check status transitions
- Validate business rules (price > 0, dates...)

📝 **Consistency**:
- Update Product.LastPriceID khi tạo ProductPrice mới
- Update Product.SaleStatus khi Transaction completed
- Archive parent khi split Product/Asset

🏷️ **Enum Usage**:
- Dùng enums thay vì magic numbers
- Helper methods để get display names
- Validate enum values

---

### Key Business Flows:

**1. Product Sale Flow**:
```
Create Product → Create Price → Create Post → 
Transaction (Draft → Deposite → Signed → Completed) → 
Update Product Status (Sold)
```

**2. Asset Management Flow**:
```
Purchase Asset → Add Legal Docs → 
Create Exploitation Contracts → Track Income → 
(Optional) Create Product to sell
```

**3. Sharing Flow**:
```
Check Owner → Create SharingAccess/AssetShare → 
Set Permissions & Commission → Notify recipient
```

**4. Split Product Flow**:
```
Create Child Products → Create ProductHistory → 
Archive Parent Product → Transfer relations
```

---

File knowledge base này phục vụ cho:
- ✅ AI Assistant trả lời câu hỏi về bdspro domain
- ✅ Developers mới tìm hiểu cấu trúc dữ liệu
- ✅ Code review và validation rules
- ✅ Database schema design reference
- ✅ API design và business logic implementation

