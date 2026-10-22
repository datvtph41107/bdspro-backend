package entity

import (
	"time"

	"gorm.io/gorm"
)

// Asset maps to the assets table (owner-held real estate asset).
// ORM::TABLE
type Asset struct {
	ID               uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt        *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy        *uint64        `gorm:"column:created_by"`
	UpdatedBy        *uint64        `gorm:"column:updated_by"`
	Name             string         `gorm:"column:name"`
	WardID           *uint64        `gorm:"column:ward_id"`
	ProvinceID       *uint64        `gorm:"column:province_id"`
	ProductID        *uint64        `gorm:"column:product_id;index"`
	ParentAssetID    *uint64        `gorm:"column:parent_asset_id"`
	SplitMergeStatus string         `gorm:"column:split_merge_status"`
	PurchasePrice    float64        `gorm:"column:purchase_price"`
	PurchaseDate     *time.Time     `gorm:"column:purchase_date"`
	LegalStatus      uint32         `gorm:"column:legal_status;default:10"`
	Archived         bool           `gorm:"column:archived;default:false"`
	RentStatus       uint32         `gorm:"column:rent_status;default:10"`
	Address          string         `gorm:"column:address"`
	Description      string         `gorm:"column:description"`
	Area             float64        `gorm:"column:area"`
	PropertyTypeID   *uint64        `gorm:"column:property_type_id"`
	ImageID          *uint64        `gorm:"column:image_id"`
	OwnerID          *uint64        `gorm:"column:owner_id;index"`
	OwnerOf          uint32         `gorm:"column:owner_of;default:10"`

	LegalItems []AssetLegal `gorm:"foreignKey:AssetID;references:ID"`
}

func (Asset) TableName() string { return "assets" }

// AssetLegal is a legal document attached to an asset.
// ORM::TABLE
type AssetLegal struct {
	ID                  uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt           *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt           *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy           *uint64        `gorm:"column:created_by"`
	UpdatedBy           *uint64        `gorm:"column:updated_by"`
	AssetID             uint64         `gorm:"column:asset_id;index;not null"`
	DocumentName        string         `gorm:"column:document_name"`
	DocumentURL         string         `gorm:"column:document_url"`
	DocumentType        string         `gorm:"column:document_type"`
	IssuedDate          *time.Time     `gorm:"column:issued_date"`
	ExpiryDate          *time.Time     `gorm:"column:expiry_date"`
	Description         string         `gorm:"column:description"`
	RelatedSplitMergeID *uint64        `gorm:"column:related_split_merge_id"`
}

func (AssetLegal) TableName() string { return "asset_legals" }
