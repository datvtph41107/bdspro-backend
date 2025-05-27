package dto

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_dto "common/dto"
	"time"
)

type AssetCreateWithProduct struct {
	// Name          string            `json:"name,omitempty" binding:"required"`
	PurchasePrice float64            `json:"purchasePrice,omitempty"`
	PurchaseDate  *time.Time         `json:"purchaseDate,omitempty"`
	LegalStatus   uint               `gorm:"default:10" json:"legalStatus,omitempty"`
	RentStatus    enums.AssetStatus  `gorm:"default:10" json:"rentStatus,omitempty"`
	Description   string             `json:"description,omitempty"`
	LegalItems    []domain.LegalItem `json:"legalItems"`
	// Archived      bool              `gorm:"default:false" json:"archived"`
	// PropertyTypeId *uint64           `json:"propertyTypeId"`
	// Address       string            `json:"address,omitempty"`
	// Area           float64           `json:"area" binding:"required"`
}

type AssetSearchDTO struct {
	_dto.Page
	ProvinceID      *uint64  `form:"provinceId"`
	DistrictID      *uint64  `form:"districtId"`
	WardID          *uint64  `form:"wardId"`
	Text            string   `form:"text"`
	PropertyTypeIds []uint64 `form:"propertyTypeIds"`
	RentStatus      []uint64 `form:"rentStatus"`
	Archived        *bool    `form:"archived"`
	OnlyParent      *bool    `form:"onlyParent"`
	ParentID        *uint64  `form:"parentId"`
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

type AssetSaveRequest struct {
	domain.Asset
	LegalItems []domain.LegalItem `json:"legalItems"`
}
