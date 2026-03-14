package dto

import (
	"common/pkg/fieldmask"
)

// ─────────────────────────────────────────────────────────────────────────────
// UpdatePropertyDTO — command object duy nhất Handler → Usecase.
//
// Quy ước con trỏ toàn hệ:
//
//	Sub-domain:
//	  nil  = client không gửi domain này = KHÔNG update domain đó
//	  !nil = client gửi domain = cần upsert
//
//	Field *string / *uint32 / *int32 / *float64 / *bool:
//	  nil  = absent = KHÔNG update field
//	  &v   = present = UPDATE field (kể cả "" / 0 / false đều hợp lệ)
//
//	Field *uint64 (FK nullable):
//	  nil  = absent = KHÔNG update
//	  &0   = CLEAR FK → set NULL trong DB
//	  &>0  = set id
//
//	Field *string (DateString):
//	  nil  = absent = KHÔNG update
//	  &""  = CLEAR date → set NULL trong DB
//	  &RFC3339 = set date
//
//	Mask:
//	  Khởi tạo 1 lần tại Handler từ req.FieldMask.Paths.
//	  Truyền xuống qua DTO, tiêu thụ tại patch.Builder trong Usecase.
//	  Zero value (FieldMask{}) = AllowAll — không bao giờ nil.
//
// ─────────────────────────────────────────────────────────────────────────────
type UpdatePropertyDTO struct {
	LineageID uint64

	// Mask: value object bất biến, khởi tạo 1 lần tại handler.
	// fieldmask.FieldMask{} = AllowAll (backward compat khi client không gửi mask)
	Mask fieldmask.FieldMask

	Info         *UpdateInfoDTO
	Location     *UpdateLocationDTO
	LandInfo     *UpdateLandInfoDTO
	BuildingInfo *UpdateBuildingDTO
	Evidence     *UpdateEvidenceDTO
	Media        *UpdateMediaDTO       // list upsert — không dùng Mask
	Amenities    *UpdateAmenityDTO     // REPLACE list — không dùng Mask
	ExternalRefs *UpdateExternalRefDTO // list upsert — không dùng Mask
	AreaRegions  *UpdateAreaRegionsDTO // REPLACE list — không dùng Mask
	LineageMeta  *UpdateLineageMetaDTO // admin only — gateway đã strip

	Personalization *UpdatePersonalizationDTO
}

// ─── UpdateInfoDTO ────────────────────────────────────────────────────────────
type UpdateInfoDTO struct {
	Title      *string
	Note       *string
	UnitCode   *string
	Identifier *string
	Level      *string

	Scope        *uint32
	SourceType   *uint32
	LegalStatus  *uint32
	RecordStatus *uint32

	// FK: nil=absent | &0=CLEAR→NULL | &>0=set
	AvatarID        *uint64
	PropertyTypeID  *uint64
	ProjectID       *uint64
	OriginProfileID *uint64
}

// ─── UpdateLocationDTO ───────────────────────────────────────────────────────
type UpdateLocationDTO struct {
	AddressDetail *string
	MapURL        *string

	Latitude  *float64
	Longitude *float64

	// FK: nil=absent | &0=CLEAR→NULL | &>0=set
	RegionID   *uint64
	ProvinceID *uint64
	DistrictID *uint64
	WardID     *uint64
}

// ─── UpdateLandInfoDTO ───────────────────────────────────────────────────────
// Sub-group nil = sub-group không được gửi = không update cả group
type UpdateLandInfoDTO struct {
	DocumentNo   *string
	IssuringAuth *string
	DocumentType *uint32
	Plot         *uint32
	Sheet        *uint32

	AreaTotal *float64
	AreaLand  *float64
	AreaPlant *float64

	FrontWidth  *float64
	Depth       *float64
	StreetWidth *float64

	ExpiredLand  *string
	ExpiredPlant *string

	Note        *string
	LandNote    *string
	PurposeUsed *int32
}

// ─── UpdateBuildingDTO ───────────────────────────────────────────────────────
type UpdateBuildingDTO struct {
	Note *string

	BuildStatus  *uint32
	BuildingType *uint32
	Direction    *uint32
	BalconyDir   *uint32

	AreaActual       *float64
	AreaFloor        *float64
	AreaConstruction *float64

	Floors     *uint32
	RoomNumber *uint32
	Bedrooms   *uint32
	Bathrooms  *uint32
}

// ─── UpdateEvidenceDTO ───────────────────────────────────────────────────────
type UpdateEvidenceDTO struct {
	Title       *string
	Description *string
	// FK: nil=absent | &0=CLEAR→NULL | &>0=set
	FileID *uint64
}

// ─── UpdateMediaDTO ───────────────────────────────────────────────────────────
// List upsert — Mask không áp dụng, tự quản lý qua ID.
// ID=0 → INSERT | ID>0 → UPDATE record đó.
type UpdateMediaDTO struct {
	Items []UpdateMediaItemDTO
}

type UpdateMediaItemDTO struct {
	ID        uint64 // 0 = insert mới
	MediaType string
	MediaURL  string
	ThumbURL  string
	SortOrder int32
	IsCover   bool
}

// ─── UpdateAmenityDTO ─────────────────────────────────────────────────────────
// REPLACE toàn bộ list — Mask không áp dụng.
// AmenityIDs = [] → xóa hết amenities.
type UpdateAmenityDTO struct {
	AmenityIDs []uint64
}

// ─── UpdateAreaRegionsDTO ─────────────────────────────────────────────────────
// REPLACE toàn bộ list — Mask không áp dụng.
type UpdateAreaRegionsDTO struct {
	AreaRegionIDs []uint64
}

// ─── UpdateExternalRefDTO ─────────────────────────────────────────────────────
// List upsert — Mask không áp dụng, tự quản lý qua ID.
// ID=0 → INSERT | ID>0 → UPDATE record đó.
//
// Scalar fields (không *T) vì external_ref là upsert hoàn toàn:
// khi UPDATE, toàn bộ field được ghi lại — không cần partial update.
type UpdateExternalRefDTO struct {
	Items []UpdateExternalRefItemDTO
}

type UpdateExternalRefItemDTO struct {
	ID            uint64 // 0 = insert mới
	ExternalRefID uint64
	SourceSystem  int32
	SourceCode    string
	Confidence    uint32
	SyncStatus    int32
	Note          string
}

// ─── UpdateLineageMetaDTO ─────────────────────────────────────────────────────
// Admin/system only — gateway đã strip paths này cho non-admin.
type UpdateLineageMetaDTO struct {
	NationalID    *string
	LineageStatus *int32
	// nil=absent | ""=CLEAR→NULL | RFC3339=set
	VerifiedNationalAt *string
}

type UpdatePersonalizationDTO struct {
	ArchivedAt *string
	HiddenAt   *string
}

// ─── Result ───────────────────────────────────────────────────────────────────

type UpdatePropertyResult struct {
	LineageID uint64
	UpdatedAt string
}
