package entity

import (
	"time"

	"gorm.io/gorm"
)

// Product maps to the products table (scoped columns for phase-1 APIs).
type Product struct {
	ID                uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt         *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy         *uint64        `gorm:"column:created_by"`
	UpdatedBy         *uint64        `gorm:"column:updated_by"`
	ParentID          *uint64        `gorm:"column:parent_id"`
	AssetID           *uint64        `gorm:"column:asset_id"`
	PropertyID        *uint64        `gorm:"column:property_id;index"`
	ImageID           *uint64        `gorm:"column:image_id"`
	PropertyTypeID    *uint64        `gorm:"column:property_type_id"`
	DocTypeID         *uint64        `gorm:"column:doc_type_id"`
	DocTypeNote       string         `gorm:"column:doc_type_note"`
	ProjectID         *uint64        `gorm:"column:project_id"`
	Visibility        uint32         `gorm:"column:visibility;default:30"`
	Name              string         `gorm:"column:name;size:255"`
	Code              string         `gorm:"column:code;size:15"`
	CategoryID        *uint64        `gorm:"column:category_id"`
	Area              float64        `gorm:"column:area"`
	PositionURL       string         `gorm:"column:position_url"`
	Description       string         `gorm:"column:description;type:text"`
	LastPriceID       *uint64        `gorm:"column:last_price_id"`
	SourceType        uint32         `gorm:"column:source_type"`
	SourceStatus      uint32         `gorm:"column:source_status"`
	SourceContactID   *uint64        `gorm:"column:source_contact_id"`
	SourceContactNote string         `gorm:"column:source_contact_note;type:text"`
	SourceTime        *time.Time     `gorm:"column:source_time"`
	Priority          uint32         `gorm:"column:priority;default:20"`
	Note              string         `gorm:"column:note;type:text"`
	TransactionType   int32          `gorm:"column:transaction_type;default:10"`
	SaleStatus        uint32         `gorm:"column:sale_status;default:10"`
	SaleVisibility    uint32         `gorm:"column:sale_visibility;default:30"`
	RentStatus        uint32         `gorm:"column:rent_status;default:10"`
	RentVisibility    uint32         `gorm:"column:rent_visibility;default:30"`
	OwnerID           uint64         `gorm:"column:owner_id;index"`
	OwnerOf           uint32         `gorm:"column:owner_of;default:10"`
	ProvinceID        *uint64        `gorm:"column:province_id"`
	WardID            *uint64        `gorm:"column:ward_id"`
	Address           string         `gorm:"column:address;size:255"`
	GoogleMapLink     string         `gorm:"column:google_map_link;size:255"`
	Archived          bool           `gorm:"column:archived;default:false"`

	Price *ProductPrice `gorm:"foreignKey:ProductID;references:ID"`
}

func (Product) TableName() string { return "products" }

// ProductPrice maps to product_price.
type ProductPrice struct {
	ID                 uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt          *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index"`
	ProductID          *uint64        `gorm:"column:product_id;index"`
	Currency           string         `gorm:"column:currency"`
	SalePrice          *float64       `gorm:"column:sale_price;type:decimal(15,0)"`
	SaleCommission     *float64       `gorm:"column:sale_commission;type:decimal(15,0)"`
	SaleCommissionType *uint32        `gorm:"column:sale_commission_type"`
	Deposite           *float64       `gorm:"column:deposite;type:decimal(15,0)"`
	RentPrice          *float64       `gorm:"column:rent_price;type:decimal(15,0)"`
	RentCommission     *float64       `gorm:"column:rent_commission;type:decimal(15,0)"`
	RentCommissionType *uint32        `gorm:"column:rent_commission_type"`
	RentPaymentCycle   *uint32        `gorm:"column:rent_payment_cycle"`
	ChannelPrice       bool           `gorm:"column:channel_price;default:false"`
	ChangeNote         string         `gorm:"column:change_note;type:text"`
	ChangedBy          uint64         `gorm:"column:changed_by"`
}

func (ProductPrice) TableName() string { return "product_price" }

// ProductUser links a profile to a product for /product/me listing.
type ProductUser struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt       *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
	ProfileID       uint64         `gorm:"column:profile_id;index;not null"`
	ProductID       uint64         `gorm:"column:product_id;index;not null"`
	IsOwner         bool           `gorm:"column:is_owner;default:true"`
	RoleID          uint64         `gorm:"column:role_id;default:0"`
	OwnerOf         uint32         `gorm:"column:owner_of;default:10"`
	OriginProfileID *uint64        `gorm:"column:origin_profile_id;index"`
}

func (ProductUser) TableName() string { return "product_user" }
