package u_asset

import (
	"bdspro/v1/models"
)

type HistoryDTO struct {
	Page int `form:"page"`
	Size int `form:"size"`
}
type SearchDTO struct {
	Page              int      `form:"page"`
	Size              int      `form:"size"`
	PriceFrom         *int     `form:"priceFrom"`
	PriceTo           *int     `form:"priceTo"`
	WardId            *uint64  `form:"wardIds"`
	DistrictId        *uint64  `form:"districtIds"`
	ProvinceId        *uint64  `form:"provinceIds"`
	DocStatus         *uint64  `form:"docStatus"`
	ParentAssetID     *uint64  `form:"parentAssetId"`
	SplitMergeStatus  *string  `form:"splitMergeStatus"`
	PurchasePriceFrom *float64 `form:"purchasePriceFrom"`
	PurchasePriceTo   *float64 `form:"purchasePriceTo"`
	LegalStatus       *uint    `form:"legalStatus"`
	Archived          *bool    `form:"archived"`
	PropertyTypeIds   []uint64 `form:"propertyTypeIds" parser:"uint64s"`
	RentStatus        []uint64 `form:"rentStatus"`
	Text              string   `form:"text"`
	// RentStatus        *enums.AssetStatus `form:"rentStatus"`
	// ProvinceId        *uint64            `form:"provinceId"`
	// DistrictId        *uint64            `form:"districtId"`
	// WardId            *uint64            `form:"wardId"`

}

type ArchivedDTO struct {
	ID       uint64
	Archived bool
}

type MergeDTO struct {
	IDs []uint64 `json:"ids" binding:"required,notemptyarray,dive"`
	// Name          string            `json:"name" `
	// DocStatus     uint64            `json:"docStatus" `
	// PurchasePrice float64           `json:"purchasePrice"`
	// PurchaseDate  *time.Time        `json:"purchaseDate"`
	// LegalStatus   string            `json:"legalStatus"`
	// RentStatus    enums.AssetStatus `json:"rentStatus"`
	// ParentAssetID *uint64           `json:"parentAssetId"`
}

// type SplitItem struct {
// }

type SplitDTO struct {
	ID    uint64         `json:"id" binding:"required"`
	Datas []models.Asset `json:"datas" binding:"required,notemptyarray,dive"`
}

// --- swagger only
type SplitDTOSwagger struct {
	ID    uint64              `json:"id" binding:"required"`
	Datas []models.AssetSplit `json:"datas" binding:"required,notemptyarray,dive"`
}
