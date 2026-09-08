package dto

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_dto "common/domain/dto"
	"time"
)

type Region struct {
	ID     uint64
	Name   string
	Code   uint32
	Active bool
}

// ==================== PropertySearchDTO (UPDATED) ====================
// PropertySearchDTO chứa các tham số tìm kiếm và lọc danh sách Property
// Đã bổ sung các trường filter mới từ UI (scope, identified, nationalVerified, địa chỉ, thực thể liên quan, thời gian)
type PropertySearchDTO struct {
	_dto.Pagable
	Text           string  `json:"text,omitempty"`
	PropertyTypeID *uint64 `json:"propertyTypeId,omitempty"`
	LocationID     *uint64 `json:"locationId,omitempty"` // deprecated, giữ để tương thích
	ProjectID      *uint64 `json:"projectId,omitempty"`
	RecordStatus   string  `json:"recordStatus,omitempty"` // "active", "inactive", "archived"
	SourceType     *uint32 `json:"sourceType,omitempty"`
	// NEW: Các filter mới từ UI
	Scope            *string `json:"scope,omitempty"`            // "private", "shared", "public"
	Identified       *bool   `json:"identified,omitempty"`       // true: đã định danh, false: chưa
	NationalVerified *bool   `json:"nationalVerified,omitempty"` // true: đã xác thực QG

	// NEW: Địa chỉ hành chính
	ProvinceID *uint64 `json:"provinceId,omitempty"`
	DistrictID *uint64 `json:"districtId,omitempty"`
	WardID     *uint64 `json:"wardId,omitempty"`

	// NEW: Lọc theo sự tồn tại của thực thể liên quan
	HasAsset   *bool `json:"hasAsset,omitempty"`
	HasProduct *bool `json:"hasProduct,omitempty"`
	HasListing *bool `json:"hasListing,omitempty"`

	// NEW: Khoảng thời gian
	CreatedFrom *time.Time `json:"createdFrom,omitempty"`
	CreatedTo   *time.Time `json:"createdTo,omitempty"`
	UpdatedFrom *time.Time `json:"updatedFrom,omitempty"`
	UpdatedTo   *time.Time `json:"updatedTo,omitempty"`
	Owned       bool       `json:"owned"`
	RegionID    *uint64    `json:"regionId"`

	RequestUserId *uint64 `json:"-"`
	Sync          bool
}

// ==================== PropertyWithCounts (UPDATED) ====================
// PropertyWithCounts chứa Property cùng với các thông tin tổng hợp từ JOIN và subquery
// Đã bổ sung các trường mới từ domain (Identifier, NationalID, NationalIDVerified, Scope, DistrictID, BuildingID...)
// và các trường tên (districtName, blockName, buildingName) phục vụ hiển thị
type PropertyWithCounts struct {
	ID             uint64    `gorm:"column:id"`
	Title          string    `gorm:"column:title"`
	AddressDetail  string    `gorm:"column:address_detail"`
	ProvinceID     *uint64   `gorm:"column:province_id"`
	DistrictID     *uint64   `gorm:"column:district_id"`
	WardID         *uint64   `gorm:"column:ward_id"`
	ProjectID      *uint64   `gorm:"column:project_id"`
	PropertyTypeID *uint64   `gorm:"column:property_type_id"`
	Scope          int16     `gorm:"column:scope"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	CreatedBy      uint64    `gorm:"column:created_by"`
	AreaTotal      float64   `gorm:"column:area_total"`
	SourceType     uint32    `gorm:"column:source_type"`
	StatusRecord   uint32    `gorm:"column:record_status"`
	Visibility     uint32    `gorm:"column:visibility"`

	// ===== NAME (alias) =====
	ProvinceName     string `gorm:"column:province_name"`
	DistrictName     string `gorm:"column:district_name"`
	WardName         string `gorm:"column:ward_name"`
	ProjectName      string `gorm:"column:project_name"`
	PropertyTypeName string `gorm:"column:property_type_name"`

	// ===== BUILDING INFO =====
	BuildingType     *int32   `gorm:"column:building_type"`
	ConstructionArea *float64 `gorm:"column:construction_area"`
	FloorArea        *float64 `gorm:"column:floor_area"`
	Floors           *int32   `gorm:"column:floors"`
	Bedrooms         *int32   `gorm:"column:bedrooms"`
	Bathrooms        *int32   `gorm:"column:bathrooms"`
	Direction        *int32   `gorm:"column:direction"`
	BalconyDirection *int32   `gorm:"column:balcony_direction"`

	// ===== AVATAR =====
	AvatarURL      string `gorm:"column:avatar_url"`
	AvatarThumbURL string `gorm:"column:avatar_thumb_url"`

	// ===== COUNTS =====
	AssetCount   int64 `gorm:"column:asset_count"`
	ProductCount int64 `gorm:"column:product_count"`
	ListingCount int64 `gorm:"column:listing_count"`
}

type CreatePropertyProductRequest struct {
	Lineage      *domain.PropertyLineage      `json:"lineage"`
	Info         *domain.PropertyInfo         `json:"info"`
	LandInfo     *domain.PropertyLandInfo     `json:"landInfo"`
	BuildingInfo *domain.PropertyBuildingInfo `json:"buildingInfo"`
	ExternalRef  *domain.PropertyExternalRef  `json:"externalRef"`
	Edvidence    *domain.PropertyEdvidence    `json:"edvidence"`
	MediaList    []domain.PropertyMedia       `json:"mediaList"`
	Location     *domain.PropertyLocation     `json:"location"`

	LegalInfo  *domain.AssetLegal `json:"legalInfo"`
	AmenityIds []uint64           `json:"amenityIds"`
	TagIDs     []uint64           `json:"tagIds"`
	Post       *domain.Post       `json:"post"`
	Product    *domain.Product    `json:"product"`
	Asset      *domain.Asset      `json:"asset"`
	Note       string             `json:"note"`
}

type ApplyPriceUpdateRequest struct {
	ProductID    uint64                   `json:"productId" validate:"required"`
	Price        float64                  `json:"price,omitempty"`
	PriceStatus  enums.ProductPriceStatus `json:"priceStatus,omitempty"`
	ChangeNote   string                   `json:"changeNote,omitempty"`
	ChannelPrice bool                     `json:"channelPrice"`
}

type ApplyPriceUpdateResponse struct {
	ProductID  uint64    `json:"productId"`
	OldPrice   *float64  `json:"oldPrice,omitempty"`
	NewPrice   *float64  `json:"newPrice,omitempty"`
	PricePerM2 *float64  `json:"pricePerM2,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// Amenities
type TagDTO struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type PropertyTagDTO struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Icon string `json:"icon"`
}

type TagSearchDTO struct {
	_dto.Pagable
	Keyword string
	Type    string
}

type UpdatePropertyTagsRequest struct {
	PropertyID uint64   `json:"propertyId"`
	TagIDs     []uint64 `json:"tagIds"`
	SourceType uint32   `json:"sourceType"`
}

type PropertyAmenityDTO struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
	Checked bool   `json:"checked"`
}

type UpdatePropertyAmenitiesRequest struct {
	PropertyID uint64   `json:"propertyId"`
	AmenityIDs []uint64 `json:"amenityIds"`
	SourceType uint32   `json:"sourceType"`
}

type ArchivePropertiesDTO struct {
	IDs      []uint64 `json:"ids"`
	Archived bool     `json:"archived"`
}

type HiddenPropertiesDTO struct {
	IDs    []uint64 `json:"ids"`
	Hidden bool     `json:"hidden"`
}

type ArchiveProductsDTO struct {
	IDs      []uint64 `json:"ids"`
	Archived bool     `json:"archived"`
}
