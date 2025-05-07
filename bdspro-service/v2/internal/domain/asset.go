package domain

import (
	"bdspro/enums"
	_models "common/models"
	"time"
)

type Asset struct {
	_models.BaseEntity
	Name             string            `json:"name" binding:"required"`
	WardId           uint64            `json:"wardId" binding:"required"`
	DistrictId       uint64            `json:"districtId"`
	ProvinceId       uint64            `json:"provinceId"`
	Price            float64           `json:"price" binding:"required"`
	BuyedDate        *time.Time        `json:"buyedDate"`
	DocStatus        uint64            `json:"docStatus"` // 1: sổ đỏ, 2: sổ hồng, 3: chưa có sổ
	ProductID        *uint64           `json:"productId"`
	OwnerUserID      *uint64           `json:"ownerUserId"`
	OwnerOrgID       *uint64           `json:"ownerOrgId"`
	ParentAssetID    *uint64           `json:"parentAssetId"`
	SplitMergeStatus string            `json:"splitMergeStatus"`
	PurchasePrice    float64           `json:"purchasePrice"`
	PurchaseDate     *time.Time        `json:"purchaseDate"`
	LegalStatus      string            `json:"legalStatus"`
	Archived         bool              `gorm:"default:false" json:"archived"`
	RentStatus       enums.AssetStatus `gorm:"default:1" json:"rentStatus"`
	Area             float64           `json:"area" binding:"required"`
	Address          string            `json:"address"`
	Description      string            `json:"description"`
}

func (Asset) TableName() string {
	return "assets"
}
