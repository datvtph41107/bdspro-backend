package models

import (
	"bdspro/enums"
	_models "common/models"
)

// Product đại diện cho bảng sản phẩm

type Product struct {
	_models.BaseEntity
	ID                   uint64                  `gorm:"primaryKey;autoIncrement" json:"id"`
	AssetId              uint64                  `json:"assetId"`
	Status               enums.ProductStatus     `json:"status"`
	Visibility           enums.Scope             `json:"visibility" gorm:"default:1"` //phạm vi
	Name                 string                  `gorm:"size:255" json:"name"`
	Code                 string                  `gorm:"size:15" json:"code"`
	CategoryID           uint64                  `gorm:"not null" json:"category_id"`
	Area                 float64                 `json:"area"`
	Price                int64                   `json:"price"`
	Description          string                  `json:"description"`
	OrganizationID       uint64                  `json:"organization_id"`
	OrgStatus            int16                   `json:"org_status"`
	ProjectID            uint64                  `json:"project_id"`
	SourceType           int16                   `json:"source_type"`
	CustomerID           uint64                  `json:"customer_id"`
	ImportPrice          int64                   `json:"import_price"`
	CommissionPercentage float64                 `json:"commission_percentage"`
	CommissionAmount     float64                 `json:"commission_amount"`
	Note                 string                  `json:"note"`
	TransactionPrice     float64                 `json:"transaction_price"`
	GiaNhap              *float64                `gorm:"type:DECIMAL(15,0)" json:"gia_nhap,omitempty"`                                                              // Giá chủ nhà muốn net (nội bộ)
	SaleStatus           enums.ProductStatus     `gorm:"type:smallint;not null;default:1" json:"sale_status"`                                                       // 1=Chưa bán, 2=Đang bán, 3=Đã bán
	SalePriceTotal       *float64                `gorm:"type:DECIMAL(15,0)" json:"sale_price_total,omitempty"`                                                      // Giá bán niêm yết (công khai)
	SalePricePerM2       *float64                `gorm:"type:DECIMAL(15,2)" json:"sale_price_per_m2,omitempty"`                                                     // Giá bán/m² (công khai)
	SaleCommission       *float64                `gorm:"type:DECIMAL(10,2)" json:"sale_commission,omitempty"`                                                       // Hoa hồng bán (nội bộ)
	SaleCommissionType   *string                 `gorm:"type:VARCHAR(10);check:sale_commission_type IN ('percent', 'fixed')" json:"sale_commission_type,omitempty"` // "percent" hoặc "fixed"
	SaleVisibility       enums.ProductVisibility `gorm:"type:smallint;not null;default:1" json:"sale_visibility"`                                                   // 1=Riêng tư, 2=Nội bộ, 3=Công khai
	RentStatus           uint                    `gorm:"type:smallint;not null;default:1" json:"rent_status"`                                                       // 1=Chưa thuê, 2=Đang đăng tin, 3=Đang cho thuê
	RentPrice            *float64                `gorm:"type:DECIMAL(15,0)" json:"rent_price,omitempty"`                                                            // Giá thuê/tháng (công khai)
	RentCommission       *float64                `gorm:"type:DECIMAL(10,2)" json:"rent_commission,omitempty"`                                                       // Hoa hồng thuê (nội bộ)
	RentCommissionType   *string                 `gorm:"type:VARCHAR(10);check:rent_commission_type IN ('percent', 'fixed')" json:"rent_commission_type,omitempty"` // "percent" hoặc "fixed"
	RentPaymentCycle     *string                 `gorm:"type:VARCHAR(20)" json:"rent_payment_cycle,omitempty"`                                                      // Chu kỳ thanh toán ("tháng", "quý"...)
	RentVisibility       uint                    `gorm:"type:smallint;not null;default:1" json:"rent_visibility"`                                                   // 1=Riêng tư, 2=Nội bộ, 3=Công khai
	OwnerUserID          *uint64                 `gorm:"type:int8" json:"owner_user_id,omitempty"`                                                                  // Ai sở hữu (nếu đã bán)
	OwnerOrgID           *uint64                 `gorm:"type:int8" json:"owner_org_id,omitempty"`                                                                   // Tổ chức sở hữu (nếu cty mua)
	AssetID              *int64                  `gorm:"type:int8" json:"asset_id,omitempty"`                                                                       // Link tới tài sản (nếu có)
	GhiChuNoiBo          *string                 `gorm:"type:TEXT" json:"ghi_chu_noi_bo,omitempty"`                                                                 // Ghi chú nội bộ (môi giới)
	Archived             bool                    `gorm:"type:bool;not null;default:0" json:"archived"`                                                              // 0=chưa, 1=archive

	// CreatedUser *Profile           `gorm:"foreignKey:CreatedBy" json:"created_user,omitempty"`
	Attributes []AttributeProduct `gorm:"foreignKey:ProductID" json:"attributes,omitempty"`
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
	// Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	// Status     int16  `json:"status"`
	// Images               []ProductImage         `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	// ShareRequests        []ShareRequest         `gorm:"foreignKey:ProductID" json:"share_requests,omitempty"`
	// AssignMembers        []*OrganizationRequest `gorm:"many2many:product_member" json:"assign_members,omitempty"`

}
