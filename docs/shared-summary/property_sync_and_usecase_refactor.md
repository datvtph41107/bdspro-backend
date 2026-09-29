# Tổng kết: Property Sync API, SyncProvider & UserCreate Refactor

> **Mục đích**: Tài liệu cho AI đọc hiểu và thực hiện nhanh khi có yêu cầu tương tự.

---

## 1. Property Timestamp API & GetDetail Sync

### 1.1. API GetPropertyTimestamps

**Proto** (`shared/protobuf/schema/bdspro/property.proto`):
```protobuf
rpc GetPropertyTimestamps(sharepb.IdRequest) returns (PropertyTimestampsResponse) {
  option (google.api.http) = { get: "/v2/bdspro/v2/property/{id}/timestamp" };
};

message PropertyTimestampsResponse {
  map<string, int64> timestamps = 1;  // keys: "detail", "history"
}
```

**Handler** (`bdspro-service/infra/handler/property_handler.go`):
- Gọi `SyncProvider.MGet(ctx, redisKeys)` với 2 key:
  - `property:{profileId}:{propertyId}` → detail
  - `property:histories:{profileId}:{propertyId}` → history
- Trả `timestamps["detail"]` và `timestamps["history"]` cho FE so sánh cache

### 1.2. GetDetail dùng SyncProvider

- **Key format**: `property:{profileId}:{propertyId}` (lấy profileId từ context)
- **Flow**: `HasUpdated(key, timestamp)` → nếu `!updated` return empty → fetch DB → `PutTimestamp(key, updatedAt)`
- Giống pattern của `product_handler.GetDetail`

---

## 2. SyncProvider – GetKey & Constants

### 2.1. Hàm GetKey

**File**: `shared/common/provider/sync_provider.go`

```go
// GetKey trả về key Redis: {prefix}:{profileId}:{resourceId}
func (s *SyncProvider) GetKey(ctx context.Context, prefix string, resourceId uint64) string {
	profileId := _utils.GetProfileIdWithContext(ctx)
	return fmt.Sprintf("%s:%d:%d", prefix, profileId, resourceId)
}
```

### 2.2. Constants (Prefix keys)

```go
const (
	SyncKeyProduct               = "product"
	SyncKeyProductPriceHist      = "product:price_hist"
	SyncKeyProductDistrs         = "product:distrs"
	SyncKeyProductPosts          = "product:posts"
	SyncKeyProductDeals          = "product:deals"
	SyncKeyProductNotes          = "product:notes"
	SyncKeyProductAppointments   = "product:appointments"
	SyncKeyProductContacts       = "product:contacts"
	SyncKeyPropertyHistories     = "property:histories"
	SyncKeyPropertyMe            = "ppt:me"      // list property của user
	SyncKeyPropertyDetail        = "ppt:id"      // detail
)
```

### 2.3. Cách dùng trong handler

```go
key := s.SyncProvider.GetKey(ctx, _utils.SyncKeyProduct, req.Id)
redisKeys := []string{
	s.SyncProvider.GetKey(ctx, _utils.SyncKeyProduct, req.Id),
	s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductPriceHist, req.Id),
	// ...
}
```

---

## 3. Property UserCreate Refactor

### 3.1. Thứ tự tạo bảng

1. **PropertyIdentify** (BasicInfo) – luôn tạo trước
2. **PropertyLocation** (với `PropertyIdentifyID`)
3. **PropertyLandInfo**
4. **PropertyBuildingInfo**
5. **PropertyEdvidence**
6. **PropertyExternalRef**
7. **PropertyMedia**
8. **PropertyLineage** – tạo cuối, chứa ID của các bảng trên
9. **PropertyUser** – link lineage với owner
10. Product, Asset, PropertyRelation, Amenity, AssetLegal

### 3.2. Repos mới

**Interface** (`bdspro-service/internal/repo/property_repo.go`):
- `PropertyIdentifyRepo` – Create
- `PropertyLocationRepo` – Create
- `PropertyExternalRefRepo` – Create
- `PropertyEdvidenceRepo` – Create

**Postgres**:
- `infra/postgres/property_identify_postgres.go`
- `infra/postgres/property_location_postgres.go`
- `infra/postgres/property_external_ref_postgres.go`
- `infra/postgres/property_edvidence_postgres.go`

### 3.3. Logic UserCreate

- Nếu không có `req.BasicInfo` → tạo mặc định `PropertyIdentify{PID: originId, Version: 1}`
- Mọi entity con đều set `PropertyIdentifyID` trước khi `Create`
- `PropertyLineage` gán: `PropertyIdentifyID`, `LocationID`, `LandInfoID`, `BuildingInfoID`, `EdvidenceID`, `ExternalRefID`
- Luôn tạo `PropertyUser` sau `PropertyLineage`
- Sync timestamp list: `SyncProvider.PutTimestamp(ctx, GetKey(ctx, SyncKeyPropertyMe, 0), ...)`

### 3.4. Wire

- Thêm 4 postgres constructor vào `wire/wire.go`
- Inject vào `PropertyUsecase`

---

## 4. Quick reference

| Yêu cầu | File/Component |
|---------|----------------|
| Thêm sync key mới | `shared/common/provider/sync_provider.go` – constants |
| API timestamp mới | Proto → Handler dùng `GetKey` + `MGet` |
| Tạo entity có FK | Tạo entity cha trước → lấy ID → gán vào entity con |
| Thêm repo mới | `internal/repo/*.go` interface + `infra/postgres/*_postgres.go` + wire |

---

## 5. Migration

Đảm bảo các bảng sau tồn tại:
- `property_identify`
- `property_location`
- `property_external_ref`
- `property_edvidence`
