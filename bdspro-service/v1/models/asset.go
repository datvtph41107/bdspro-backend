package models

import (
	"bdspro/internal/enums"
	_models "common/models"
	"time"

	"gorm.io/gorm"
)

type AssetSplit struct {
	Name             string             `json:"name,omitempty" binding:"required"`
	WardID           *uint64            `json:"wardId" binding:"required"`
	DistrictID       *uint64            `json:"districtId,omitempty"`
	ProvinceID       *uint64            `json:"provinceId,omitempty"`
	ProductID        *uint64            `json:"productId"`
	ParentAssetID    *uint64            `json:"parentAssetId"`
	SplitMergeStatus string             `json:"splitMergeStatus"`
	PurchasePrice    float64            `json:"purchasePrice,omitempty" binding:"required"`
	PurchaseDate     *time.Time         `json:"purchaseDate,omitempty"`
	LegalStatus      uint               `gorm:"default:10" json:"legalStatus,omitempty"`
	Archived         bool               `gorm:"default:false" json:"archived"`
	RentStatus       enums.EAssetStatus `gorm:"default:10" json:"rentStatus,omitempty"`
	Address          string             `json:"address,omitempty"`
	Description      string             `json:"description,omitempty"`
	Area             float64            `json:"area" binding:"required"`
	PropertyTypeId   *uint64            `json:"propertyTypeId"`
}

type Asset struct {
	// ID             uint64        `gorm:"primaryKey;autoIncrement" json:"id"`
	// OwnershipDocId uint64        `json:"ownership_doc_id"` // giấy tờ sở hữu
	// OwnershipDoc   *OwnershipDoc `gorm:"foreignKey:OwnershipDocId" json:"ownershipDoc,omitempty"`
	// ProvinceID     uint64        `json:"province_id"`
	// DistrictID     uint64        `json:"district_id"`
	// District       *District     `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	// WardID         uint64        `json:"ward_id"`
	// Ward           *Ward         `gorm:"foreignKey:WardID" json:"ward,omitempty"`
	// Address        string        `gorm:"size:255" json:"address"`
	// GoogleMapLink  string        `gorm:"size:255" json:"googleMapLink"`
	// DocStatus        uint64            `gorm:"default:1" json:"docStatus,omitempty"` // 1: sổ đỏ, 2: sổ hồng, 3: chưa có sổ
	// Price            float64           `json:"price" binding:"required"`
	// BuyedDate        *time.Time        `json:"buyedDate,omitempty"`
	_models.BaseEntity
	Name             string             `json:"name,omitempty" binding:"required"`
	WardID           *uint64            `json:"wardId" binding:"required"`
	Ward             *Region            `gorm:"foreignKey:WardID;references:ID" json:"ward,omitempty"`
	DistrictID       *uint64            `json:"districtId,omitempty"`
	District         *Region            `gorm:"foreignKey:DistrictID;references:ID" json:"district,omitempty"`
	ProvinceID       *uint64            `json:"provinceId,omitempty"`
	Province         *Region            `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
	ProductID        *uint64            `json:"productId"`
	OwnerUserID      *uint64            `json:"ownerUserId"`
	OwnerOrgID       *uint64            `json:"ownerOrgId"`
	ParentAssetID    *uint64            `json:"parentAssetId"`
	SplitMergeStatus string             `json:"splitMergeStatus"`
	PurchasePrice    float64            `json:"purchasePrice,omitempty" binding:"required"`
	PurchaseDate     *time.Time         `json:"purchaseDate,omitempty"`
	LegalStatus      uint               `gorm:"default:10" json:"legalStatus,omitempty"`
	Archived         bool               `gorm:"default:false" json:"archived"`
	RentStatus       enums.EAssetStatus `gorm:"default:10" json:"rentStatus,omitempty"`
	Address          string             `json:"address,omitempty"`
	Description      string             `json:"description,omitempty"`
	Area             float64            `json:"area" binding:"required"`
	PropertyTypeId   *uint64            `json:"propertyTypeId"`
	PropertyType     *PropertyTypeItem  `gorm:"foreignKey:PropertyTypeId;references:ID" json:"propertyType,omitempty"`
	Product          *Product           `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
}

type AssetList struct {
	Asset
	PropertyTypeName string `gorm:"column:property_type_name" json:"propertyTypeName"`
	ProvinceName     string `gorm:"column:province_name" json:"provinceName"`
	DistrictName     string `gorm:"column:district_name" json:"districtName"`
	WardName         string `gorm:"column:ward_name" json:"wardName"`
	Permissions      string `gorm:"column:permissions" json:"permissions"`
}

type AssetLog struct {
	CreatedBy        *uint64            `json:"-"`
	UpdatedBy        *uint64            `json:"-"`
	ID               uint64             `json:"id"`
	CreatedAt        *time.Time         `json:"-"`
	UpdatedAt        *time.Time         `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt     `json:"-"`
	Name             string             `json:"name,omitempty"`
	WardId           uint64             `json:"wardId"`
	DistrictId       uint64             `json:"districtId,omitempty"`
	ProvinceId       uint64             `json:"provinceId,omitempty"`
	Price            float64            `json:"price"`
	BuyedDate        *time.Time         `json:"buyedDate,omitempty"`
	DocStatus        uint64             `json:"docStatus,omitempty"` // 1: sổ đỏ, 2: sổ hồng, 3: chưa có sổ
	ProductID        *uint64            `json:"productId"`
	OwnerUserID      *uint64            `json:"ownerUserId"`
	OwnerOrgID       *uint64            `json:"ownerOrgId"`
	ParentAssetID    *uint64            `json:"parentAssetId"`
	SplitMergeStatus string             `json:"splitMergeStatus"`
	PurchasePrice    float64            `json:"purchasePrice,omitempty"`
	PurchaseDate     *time.Time         `json:"purchaseDate,omitempty"`
	LegalStatus      string             `json:"legalStatus,omitempty"`
	Archived         bool               `json:"archived"`
	RentStatus       enums.EAssetStatus `json:"rentStatus,omitempty"`
	Address          string             `json:"address,omitempty"`
	Description      string             `json:"description,omitempty"`
	Area             float64            `json:"area" binding:"required"`
}

func (Asset) TableName() string {
	return "assets"
}
