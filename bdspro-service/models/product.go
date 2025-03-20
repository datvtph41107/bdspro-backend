package models

import (
	"bdspro/enums"
	_models "common/models"
)

// Product đại diện cho bảng sản phẩm

type Product struct {
	_models.BaseEntity
	ID                   uint64              `gorm:"primaryKey;autoIncrement" json:"id"`
	AssetId              uint64              `json:"assetId"`
	Status               enums.ProductStatus `json:"status"`
	Visibility           enums.Scope         `json:"visibility" gorm:"default:1"` //phạm vi
	Name                 string              `gorm:"size:255" json:"name"`
	Code                 string              `gorm:"size:15" json:"code"`
	CategoryID           uint64              `gorm:"not null" json:"category_id"`
	Area                 float64             `json:"area"`
	Price                int64               `json:"price"`
	Description          string              `json:"description"`
	OrganizationID       uint64              `json:"organization_id"`
	OrgStatus            int16               `json:"org_status"`
	ProjectID            uint64              `json:"project_id"`
	SourceType           int16               `json:"source_type"`
	CustomerID           uint64              `json:"customer_id"`
	ImportPrice          int64               `json:"import_price"`
	CommissionPercentage float64             `json:"commission_percentage"`
	CommissionAmount     float64             `json:"commission_amount"`
	Note                 string              `json:"note"`
	TransactionPrice     float64             `json:"transaction_price"`
	CreatedUser          *Profile            `gorm:"foreignKey:CreatedBy" json:"created_user,omitempty"`
	Attributes           []AttributeProduct  `gorm:"foreignKey:ProductID" json:"attributes,omitempty"`
	// CreatedBy            uint64              `json:"created_by"`
	// CreatedAt            time.Time           `gorm:"autoCreateTime" json:"created_at"`
	// UpdatedAt            time.Time           `gorm:"autoUpdateTime" json:"updated_at"`
	// DeletedAt            gorm.DeletedAt      `gorm:"index" json:"-"`
	// ProvinceID           uint64              `json:"province_id"`
	// Province             *Province           `gorm:"foreignKey:ProvinceID" json:"province,omitempty"`
	// DistrictID           uint64              `json:"district_id"`
	// District             *District           `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	// WardID               uint64              `json:"ward_id"`
	// Ward                 *Ward               `gorm:"foreignKey:WardID" json:"ward,omitempty"`
	// Address              string              `gorm:"size:255" json:"address"`
	// GoogleMapLink        string              `gorm:"size:255" json:"google_map_link"`
	// Category       *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	// Organizations        []*Organization        `gorm:"many2many:product_organization" json:"organizations,omitempty"`
	// Customer             *Customer              `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	// Project    *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	// Status     int16  `json:"status"`
	// Images               []ProductImage         `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	// ShareRequests        []ShareRequest         `gorm:"foreignKey:ProductID" json:"share_requests,omitempty"`
	// AssignMembers        []*OrganizationRequest `gorm:"many2many:product_member" json:"assign_members,omitempty"`

}
