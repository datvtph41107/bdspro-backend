package u_asset_dto

import (
	"bdspro/internal/enums"
	"time"
)

type AssetCreateWithProduct struct {
	// Name          string            `json:"name,omitempty" binding:"required"`
	PurchasePrice float64           `json:"purchasePrice,omitempty"`
	PurchaseDate  *time.Time        `json:"purchaseDate,omitempty"`
	LegalStatus   uint              `gorm:"default:10" json:"legalStatus,omitempty"`
	RentStatus    enums.AssetStatus `gorm:"default:10" json:"rentStatus,omitempty"`
	Description   string            `json:"description,omitempty"`
	// Archived      bool              `gorm:"default:false" json:"archived"`
	// PropertyTypeId *uint64           `json:"propertyTypeId"`
	// Address       string            `json:"address,omitempty"`
	// Area           float64           `json:"area" binding:"required"`
}
