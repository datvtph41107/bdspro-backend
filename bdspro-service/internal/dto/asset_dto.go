package dto

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_dto "common/dto"
	"time"
)

type AssetCreateWithProduct struct {
	// Name          string            `json:"name,omitempty" binding:"required"`
	PurchasePrice float64             `json:"purchasePrice,omitempty"`
	PurchaseDate  *time.Time          `json:"purchaseDate,omitempty"`
	LegalStatus   uint                `gorm:"default:10" json:"legalStatus,omitempty"`
	RentStatus    enums.EAssetStatus  `gorm:"default:10" json:"rentStatus,omitempty"`
	Description   string              `json:"description,omitempty"`
	LegalItems    []domain.AssetLegal `json:"legalItems"`
	// Archived      bool              `gorm:"default:false" json:"archived"`
	// PropertyTypeId *uint64           `json:"propertyTypeId"`
	// Address       string            `json:"address,omitempty"`
	// Area           float64           `json:"area" binding:"required"`
}

type AssetSearchRequest struct {
	_dto.Pagable
	ProvinceID      *uint64  `form:"provinceId"`
	DistrictID      *uint64  `form:"districtId"`
	WardID          *uint64  `form:"wardId"`
	Text            string   `form:"text"`
	PropertyTypeIds []uint64 `form:"propertyTypeIds"`
	RentStatus      []uint64 `form:"rentStatus"`
	Archived        *bool    `form:"archived"`
	OnlyParent      *bool    `form:"onlyParent"`
	ParentID        *uint64  `form:"parentId"`
	Keyword         string   `json:"keyword"`
	Type            uint32   `json:"type"`

	RequestUserId         *uint64 `form:"-"`
	RequestOrganizationId *uint64 `form:"-"`
}

type AssetSaveRequest struct {
	ID               uint64      `json:"id"`
	Name             string      `json:"name"`
	Area             float64     `json:"area"`
	WardID           *uint64     `json:"ward_id"`
	DistrictID       *uint64     `json:"district_id"`
	ProvinceID       *uint64     `json:"province_id"`
	ProductID        *uint64     `json:"product_id"`
	ParentAssetID    *uint64     `json:"parent_asset_id"`
	SplitMergeStatus string      `json:"split_merge_status"`
	PurchasePrice    float64     `json:"purchase_price"`
	PurchaseDate     string      `json:"purchase_date"`
	LegalStatus      uint32      `json:"legal_status"`
	Archived         bool        `json:"archived"`
	RentStatus       uint32      `json:"rent_status"`
	Address          string      `json:"address"`
	Description      string      `json:"description"`
	PropertyTypeID   *uint64     `json:"property_type_id"`
	LegalItems       []LegalItem `json:"legal_items"`
	OwnerType        uint32      `json:"owner_type"`
}

type MergeAssetDTO struct {
	IDs []uint64 `json:"ids" binding:"required,notemptyarray,dive"`
}

type SplitAssetDTO struct {
	ID    uint64         `json:"id" binding:"required"`
	Datas []domain.Asset `json:"datas" binding:"required,notemptyarray,dive"`
}

type AssetArchivedDTO struct {
	ID       uint64
	Archived bool
}

type AssetWithChildCountDTO struct {
	domain.AssetList
	NumChild int64 `json:"numChild"`
}

type LegalItem struct {
	ID             uint64 `json:"id"`
	AssetID        uint64 `json:"asset_id"`
	DocumentType   string `json:"document_type"`
	DocumentNumber string `json:"document_number"`
	IssueDate      string `json:"issue_date"`
	ExpiryDate     string `json:"expiry_date"`
	Status         uint32 `json:"status"`
}

type MergeAssetRequest struct {
	IDs     []uint64 `json:"ids"`
	NewName string   `json:"new_name"`
	ID      uint64   `json:"id"`
}

type AssetData struct {
	Name           string  `json:"name"`
	Area           float64 `json:"area"`
	Address        string  `json:"address"`
	Description    string  `json:"description"`
	PurchasePrice  float64 `json:"purchase_price"`
	PurchaseDate   string  `json:"purchase_date"`
	LegalStatus    uint32  `json:"legal_status"`
	PropertyTypeID *uint64 `json:"property_type_id"`
}
