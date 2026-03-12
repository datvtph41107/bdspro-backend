package dto

import (
	"bdspro/internal/domain"
	"common/pkg/optional" // thay bằng path thật
)

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
	AreaRegions  *UpdateAreaRegionsDTO
	LineageMeta  *UpdateLineageMetaDTO
}

// ─── Sub-commands ─────────────────────────────────────────────────────────────

type UpdateInfoDTO struct {
	Title        *string
	Note         *string
	UnitCode     *string
	Identifier   *string
	Level        *string
	Scope        *uint32
	SourceType   *uint32
	LegalStatus  *uint32
	RecordStatus *uint32

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
	BuildStatus      *uint32
	BuildingType     *uint32
	Direction        *uint32
	BalconyDir       *uint32
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

type UpdateAreaRegionsDTO struct {
	AreaRegionIDs []uint64
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

type UpdatePropertyCommand struct {
	LineageID uint64

	Info         *InfoUpdate
	Location     *LocationUpdate
	LandInfo     *LandInfoUpdate
	BuildingInfo *BuildingUpdate
	Evidence     *EvidenceUpdate
	Media        *MediaUpdate
	Amenities    *AmenitiesUpdate
	ExternalRefs *ExternalRefsUpdate
	AreaRegions  *AreaRegionsUpdate
	LineageMeta  *LineageMetaUpdate
}

type InfoUpdate struct {
	Title           optional.String
	Note            optional.String
	UnitCode        optional.String
	Identifier      optional.String
	Level           optional.String
	Scope           optional.Uint32
	SourceType      optional.Uint32
	LegalStatus     optional.Uint32
	RecordStatus    optional.Uint32
	AvatarID        optional.Uint64
	PropertyTypeID  optional.Uint64
	ProjectID       optional.Uint64
	OriginProfileID optional.Uint64
}

type LocationUpdate struct {
	AddressDetail optional.String
	MapURL        optional.String
	Latitude      optional.Float64
	Longitude     optional.Float64
	RegionID      optional.Uint64
	ProvinceID    optional.Uint64
	DistrictID    optional.Uint64
	WardID        optional.Uint64
}

type LandInfoUpdate struct {
	Basic       *LandInfoBasicUpdate
	Area        *LandInfoAreaUpdate
	Dimension   *LandInfoDimensionUpdate
	Expiry      *LandInfoExpiryUpdate
	Note        optional.String
	LandNote    optional.String
	PurposeUsed optional.Int32
}

type LandInfoBasicUpdate struct {
	DocumentNo   optional.String
	IssuringAuth optional.String
	DocumentType optional.Uint32
	Plot         optional.Uint32
	Sheet        optional.Uint32
}

type LandInfoAreaUpdate struct {
	AreaTotal optional.Float64
	AreaLand  optional.Float64
	AreaPlant optional.Float64
}

type LandInfoDimensionUpdate struct {
	FrontWidth  optional.Float64
	Depth       optional.Float64
	StreetWidth optional.Float64
}

type LandInfoExpiryUpdate struct {
	ExpiredLand  optional.String // ISO datetime
	ExpiredPlant optional.String
}

type BuildingUpdate struct {
	Note             optional.String
	BuildStatus      optional.Uint32
	BuildingType     optional.Uint32
	Direction        optional.Uint32
	BalconyDir       optional.Uint32
	AreaActual       optional.Float64
	AreaFloor        optional.Float64
	AreaConstruction optional.Float64
	Floors           optional.Uint32
	RoomNumber       optional.Uint32
	Bedrooms         optional.Uint32
	Bathrooms        optional.Uint32
}

type EvidenceUpdate struct {
	Title       optional.String
	Description optional.String
	FileID      optional.Uint64
}

type MediaUpdate struct {
	Items optional.Slice[MediaItemUpdate]
}

type MediaItemUpdate struct {
	ID        optional.Uint64
	MediaType optional.String
	MediaURL  optional.String
	ThumbURL  optional.String
	SortOrder optional.Int32
	IsCover   optional.Bool
}

type AmenitiesUpdate struct {
	AmenityIDs optional.Slice[uint64]
}

type ExternalRefsUpdate struct {
	Items optional.Slice[ExternalRefItemUpdate]
}

type ExternalRefItemUpdate struct {
	ID            optional.Uint64
	ExternalRefID optional.Uint64
	SourceSystem  optional.Int32
	SourceCode    optional.String
	Confidence    optional.Uint32
	SyncStatus    optional.Int32
	Note          optional.String
}

type AreaRegionsUpdate struct {
	AreaRegionIDs optional.Slice[uint64]
}

type LineageMetaUpdate struct {
	NationalID         optional.String
	LineageStatus      optional.Int32
	VerifiedNationalAt optional.String
}
