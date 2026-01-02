package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

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
	WardID           *uint64            `json:"wardId,omitempty"`
	Ward             *WardV2            `gorm:"foreignKey:WardID;references:ID" json:"ward,omitempty"`
	ProvinceID       *uint64            `json:"provinceId" binding:"required"`
	Province         *ProvinceV2        `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
	ProductID        *uint64            `json:"productId"`
	ParentAssetID    *uint64            `json:"parentAssetId"`
	SplitMergeStatus string             `json:"splitMergeStatus"`
	PurchasePrice    float64            `json:"purchasePrice,omitempty" binding:"required"`
	PurchaseDate     *time.Time         `json:"purchaseDate,omitempty"`
	LegalStatus      enums.EDocType     `gorm:"default:10" json:"legalStatus,omitempty"`
	Archived         bool               `gorm:"default:false" json:"archived"`
	RentStatus       enums.EAssetStatus `gorm:"default:10" json:"rentStatus,omitempty"`
	Address          string             `json:"address,omitempty"`
	Description      string             `json:"description,omitempty"`
	Area             float64            `json:"area" binding:"required"`
	PropertyTypeId   *uint64            `json:"propertyTypeId"`
	PropertyType     *PropertyTypeItem  `gorm:"foreignKey:PropertyTypeId;references:ID" json:"propertyType,omitempty"`
	Product          *Product           `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
	LegalItems       []AssetLegal       `gorm:"foreignKey:AssetID;references:ID" json:"legalItems,omitempty"`
	ImageID          *uint64            `json:"imageId"` // ID của legal item đầu tiên (ảnh đại diện)
	OwnerID          *uint64            `json:"ownerId"`
	OwnerOf          enums.EOwnerOf     `json:"ownerOf"`
	Owner            *ProfileInfo       `gorm:"-" json:"owner,omitempty"`

	AssetExploitations []AssetExploitation `gorm:"foreignKey:AssetID;references:ID" json:"assetExploitations,omitempty"`
	// OwnerUserID      *uint64            `json:"ownerUserId"`
	// OwnerOrgID       *uint64            `json:"ownerOrgId"`

	RequestOrganizationID *uint64 `gorm:"-" json:"-"`
}

type AssetList struct {
	Asset
	PropertyTypeName string          `gorm:"column:property_type_name" json:"propertyTypeName"`
	ProvinceName     string          `gorm:"column:province_name" json:"provinceName"`
	WardName         string          `gorm:"column:ward_name" json:"wardName"`
	Permissions      string          `gorm:"column:permissions" json:"permissions"`
	LegalList        json.RawMessage `gorm:"column:legal_list" json:"legalList"` // raw json from query
	ImageUrl         string          `gorm:"column:image_url" json:"imageUrl"`   // URL của legal item đầu tiên (ảnh đại diện)
	// LegalItems       []LegalItem     `gorm:"-" json:"legalItems"`
}

// func (a *AssetList) TableName() string {
// 	return "assets"
// }

type AssetLog struct {
	CreatedBy        *uint64            `json:"-"`
	UpdatedBy        *uint64            `json:"-"`
	ID               uint64             `json:"id"`
	CreatedAt        *time.Time         `json:"-"`
	UpdatedAt        *time.Time         `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt     `json:"-"`
	Name             string             `json:"name,omitempty"`
	WardId           uint64             `json:"wardId"`
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
