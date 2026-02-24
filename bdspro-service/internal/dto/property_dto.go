package dto

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_dto "common/domain/dto"
	"time"
)

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
}

// ==================== PropertyWithCounts (UPDATED) ====================
// PropertyWithCounts chứa Property cùng với các thông tin tổng hợp từ JOIN và subquery
// Đã bổ sung các trường mới từ domain (Identifier, NationalID, NationalIDVerified, Scope, DistrictID, BuildingID...)
// và các trường tên (districtName, blockName, buildingName) phục vụ hiển thị
type PropertyWithCounts struct {
	domain.Property // Kế thừa tất cả trường từ domain.Property (đã bao gồm các trường mới)

	// Các trường counts (từ subquery)
	ProductCount uint32 `gorm:"column:product_count" json:"productCount"`
	AssetCount   uint32 `gorm:"column:asset_count" json:"assetCount"`
	PostCount    uint32 `gorm:"column:post_count" json:"postCount"`

	// NEW: Thông tin địa chỉ (từ JOIN)
	ProvinceName string `gorm:"column:province_name" json:"provinceName,omitempty"`
	DistrictName string `gorm:"column:district_name" json:"districtName,omitempty"`
	WardName     string `gorm:"column:ward_name" json:"wardName,omitempty"`

	// NEW: Thông tin dự án (từ JOIN)
	ProjectName  string `gorm:"column:project_name" json:"projectName,omitempty"`
	BlockName    string `gorm:"column:block_name" json:"blockName,omitempty"`
	BuildingName string `gorm:"column:building_name" json:"buildingName,omitempty"`

	// Thông tin loại BĐS
	PropertyTypeName string `gorm:"column:property_type_name" json:"propertyTypeName,omitempty"`

	// Thông tin BuildingInfo (từ property_building_info)
	BuildingType     *uint32  `gorm:"column:building_type" json:"buildingType,omitempty"`
	ConstructionArea *float64 `gorm:"column:construction_area" json:"constructionArea,omitempty"`
	FloorArea        *float64 `gorm:"column:floor_area" json:"floorArea,omitempty"`
	Floors           *uint32  `gorm:"column:floors" json:"floors,omitempty"`
	Bedrooms         *uint32  `gorm:"column:bedrooms" json:"bedrooms,omitempty"`
	Bathrooms        *uint32  `gorm:"column:bathrooms" json:"bathrooms,omitempty"`
	Direction        *uint32  `gorm:"column:direction" json:"direction,omitempty"`
	BalconyDirection *uint32  `gorm:"column:balcony_direction" json:"balconyDirection,omitempty"`

	// Thông tin Avatar (từ property_media)
	AvatarURL       string `gorm:"column:avatar_url" json:"avatarUrl,omitempty"`
	AvatarThumbURL  string `gorm:"column:avatar_thumb_url" json:"avatarThumbUrl,omitempty"`
	AvatarMediaType string `gorm:"column:avatar_media_type" json:"-"`
	AvatarSortOrder int32  `gorm:"column:avatar_sort_order" json:"-"`
}

type CreatePropertyProductRequest struct {
	Property     *domain.Property             `json:"property"`
	LandInfo     *domain.PropertyLandInfo     `json:"landInfo"`
	BuildingInfo *domain.PropertyBuildingInfo `json:"buildingInfo"`
	LegalInfo    *domain.AssetLegal           `json:"legalInfo"`
	MediaList    []domain.PropertyMedia       `json:"mediaList"`
	AmenityIds   []uint64                     `json:"amenityIds"`
	TagIDs       []uint64                     `json:"tagIds"`
	Post         *domain.Post                 `json:"post"`
	Product      *domain.Product              `json:"product"`
	Asset        *domain.Asset                `json:"asset"`
	Note         string                       `json:"note"`
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
}

type ArchivePropertiesDTO struct {
	IDs      []uint64 `json:"ids"`
	Archived bool     `json:"archived"`
}

type ArchiveProductsDTO struct {
	IDs      []uint64 `json:"ids"`
	Archived bool     `json:"archived"`
}
