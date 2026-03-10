package dto

import "bdspro/internal/domain"

// ─────────────────────────────────────────────────────────────────────────────
// UpdatePropertyDTO — command object duy nhất Handler → Usecase.
//
// Quy ước con trỏ:
//   - Sub-domain nil  = domain này KHÔNG được update
//   - Sub-domain !nil = domain này cần upsert
//   - Field *T nil    = field không thay đổi
//   - Field *T != nil = field cần update (kể cả khi value = zero/null)
//
// ─────────────────────────────────────────────────────────────────────────────
type UpdatePropertyDTO struct {
	LineageID uint64

	// FieldMask: danh sách path được phép update.
	// Rỗng = update tất cả field có trong payload.
	// Có giá trị = chỉ update đúng path được liệt kê.
	FieldMask []string

	Info         *UpdateInfoDTO
	Location     *UpdateLocationDTO
	LandInfo     *UpdateLandInfoDTO
	BuildingInfo *UpdateBuildingDTO
	Evidence     *UpdateEvidenceDTO
	Media        *UpdateMediaDTO
	Amenities    *UpdateAmenityDTO
	ExternalRefs *UpdateExternalRefDTO
	LineageMeta  *UpdateLineageMetaDTO
}

// ─── Sub-commands ─────────────────────────────────────────────────────────────

type UpdateInfoDTO struct {
	Title        *string
	Note         *string
	UnitCode     *string
	Identifier   *string
	Level        *string
	Scope        *int32
	SourceType   *int32
	LegalStatus  *int32
	RecordStatus *int32

	// *uint64 nil = bỏ qua; *uint64 value=0 = clear FK về NULL
	AvatarID        *uint64
	PropertyTypeID  *uint64
	ProjectID       *uint64
	OriginProfileID *uint64
}

type UpdateLocationDTO struct {
	AddressDetail *string
	MapURL        *string
	Latitude      *float64
	Longitude     *float64
	RegionID      *uint64
	ProvinceID    *uint64
	DistrictID    *uint64
	WardID        *uint64
}

// UpdateLandInfoDTO gom nhóm theo sub-group để align với proto và field_mask.
type UpdateLandInfoDTO struct {
	Basic     *LandInfoBasicDTO
	Area      *LandInfoAreaDTO
	Dimension *LandInfoDimensionDTO
	Expiry    *LandInfoExpiryDTO

	Note        *string
	LandNote    *string
	PurposeUsed *int32
}

type LandInfoBasicDTO struct {
	DocumentNo   *string
	IssuringAuth *string
	DocumentType *uint32
	Plot         *uint32
	Sheet        *uint32
}

type LandInfoAreaDTO struct {
	AreaTotal *float64
	AreaLand  *float64
	AreaPlant *float64
}

type LandInfoDimensionDTO struct {
	FrontWidth  *float64
	Depth       *float64
	StreetWidth *float64
}

type LandInfoExpiryDTO struct {
	// *string nil = bỏ qua; "" = clear về NULL; ISO8601 = set giá trị
	ExpiredLand  *string
	ExpiredPlant *string
}

type UpdateBuildingDTO struct {
	Note             *string
	BuildStatus      *int32
	BuildingType     *int32
	Direction        *int32
	BalconyDir       *int32
	AreaActual       *float64
	AreaFloor        *float64
	AreaConstruction *float64
	Floors           *uint32
	RoomNumber       *uint32
	Bedrooms         *uint32
	Bathrooms        *uint32
}

type UpdateEvidenceDTO struct {
	Title       *string
	Description *string
	FileID      *uint64 // nil = bỏ qua; 0 = clear
}

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

type UpdateAmenityDTO struct {
	AmenityIDs []uint64 // replace toàn bộ
}

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

type UpdateLineageMetaDTO struct {
	NationalID         *string
	LineageStatus      *int32
	VerifiedNationalAt *string // nil = bỏ qua; "" = clear; ISO8601 = set
}

// ─── Result ───────────────────────────────────────────────────────────────────

type UpdatePropertyResult struct {
	LineageID uint64
	UpdatedAt string

	// nil nếu domain không được update trong request này
	Info         *domain.PropertyInfo
	Location     *domain.PropertyLocation
	LandInfo     *domain.PropertyLandInfo
	BuildingInfo *domain.PropertyBuildingInfo
	Evidence     *domain.PropertyEdvidence
	Media        []domain.PropertyMedia
	Amenities    []domain.AmenityItem
	ExternalRefs []*domain.PropertyExternalRef
}
