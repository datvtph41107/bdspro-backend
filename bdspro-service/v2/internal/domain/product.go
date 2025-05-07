package domain

import (
	"bdspro/enums"
	_models "common/models"
)

// Product đại diện cho bảng sản phẩm

type Product struct {
	_models.BaseEntity
	ID              uint64                  `gorm:"primaryKey;autoIncrement" json:"id"`
	ParentId        *uint64                 `gorm:"type:int8" json:"parentId"`
	DepositeID      *uint64                 `gorm:"type:int8" json:"depositeId"`
	AssetId         *uint64                 `gorm:"type:int8" json:"assetId"`
	PropertyTypeId  *uint64                 `json:"propertyTypeId"`
	PropertyType    *PropertyTypeItem       `gorm:"foreignKey:PropertyTypeId;references:ID" json:"propertyType,omitempty"`
	DocTypeId       *uint64                 `json:"docTypeId"`
	DocType         *DocTypeItem            `gorm:"foreignKey:DocTypeId;references:ID" json:"docType,omitempty"`
	Amenities       []AmenityItem           `gorm:"many2many:product_amenity;foreignKey:ID;joinForeignKey:ProductID;References:ID;joinReferences:AmenityID" json:"amenities,omitempty"`
	Visibility      enums.Scope             `json:"visibility" gorm:"default:1"` //phạm vi
	Name            string                  `gorm:"size:255" json:"name"`
	Code            string                  `gorm:"size:15" json:"code"`
	CategoryID      *uint64                 `json:"categoryId"`
	Area            float64                 `json:"area"`
	PositionUrl     string                  `json:"positionUrl"`
	Description     string                  `json:"description"`
	OrganizationID  *uint64                 `json:"organizationId"`
	OrgStatus       int16                   `json:"orgStatus"`
	LastPriceID     *uint64                 `json:"lastPriceId"`
	Price           *ProductPrice           `gorm:"foreignKey:ID;references:ProductID" json:"price"`
	ApartmentID     *uint64                 `json:"apartmentId"`
	Apartment       *Apartment              `gorm:"foreignKey:ApartmentID;references:ID" json:"apartment"`
	ApartmentAttr   *ApartmentAttribute     `gorm:"-" json:"apartmentAttribute"`
	Build           *ProjectBuild           `gorm:"-" json:"build"`
	Project         *Project                `gorm:"-" json:"project"`
	Developer       *Developer              `gorm:"-" json:"developer"`
	SourceType      int16                   `json:"sourceType"`
	CustomerID      *uint64                 `json:"customerId"`
	PrivateData     *ProductPrivate         `json:"privateInfo"`
	Note            string                  `json:"note"`
	TransactionType int                     `gorm:"type:smallint;default:1" json:"transactionType"`
	SaleStatus      enums.ProductStatus     `gorm:"type:smallint;not null;default:1" json:"saleStatus"`     // 1=Chưa bán, 2=Đang bán, 3=Đã bán
	SaleVisibility  enums.ProductVisibility `gorm:"type:smallint;not null;default:3" json:"saleVisibility"` // 1=Riêng tư, 2=Nội bộ, 3=Công khai
	RentStatus      enums.ProductStatus     `gorm:"type:smallint;not null;default:3" json:"rentStatus"`     // 1=Chưa thuê, 2=Đang đăng tin, 3=Đang cho thuê
	RentVisibility  enums.ProductVisibility `gorm:"type:smallint;not null;default:1" json:"rentVisibility"` // 1=Riêng tư, 2=Nội bộ, 3=Công khai
	OwnerUserID     *uint64                 `gorm:"type:int8" json:"ownerUserId,omitempty"`                 // Ai sở hữu (nếu đã bán)
	OwnerOrgID      *uint64                 `gorm:"type:int8" json:"ownerOrgId,omitempty"`                  // Tổ chức sở hữu (nếu cty mua)
	ProvinceID      *uint64                 `json:"provinceId"`
	Province        *Region                 `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
	DistrictID      *uint64                 `json:"districtId"`
	District        *Region                 `gorm:"foreignKey:DistrictID;references:ID" json:"district,omitempty"`
	WardID          *uint64                 `json:"wardId"`
	Ward            *Region                 `gorm:"foreignKey:WardID;references:ID" json:"ward,omitempty"`
	Address         string                  `gorm:"size:255" json:"address"`
	GoogleMapLink   string                  `gorm:"size:255" json:"googleMapLink"`
	Attributes      []AttributeProduct      `gorm:"foreignKey:ProductID" json:"attributes,omitempty"`
	MediaList       []ProductMediaItem      `gorm:"foreignKey:ProductID" json:"mediaList,omitempty"`
	HouseInfo       *HouseInfo              `gorm:"foreignKey:ProductID" json:"houseInfo,omitempty"`
	CreatedUser     *ProfileInfo            `gorm:"foreignKey:CreatedBy;references:ProfileId" json:"createdUser,omitempty"`
	UpdatedUser     *ProfileInfo            `gorm:"foreignKey:UpdatedBy;references:ProfileId" json:"updatedUser,omitempty"`
	Archived        bool                    `gorm:"type:bool;not null;default:false" json:"archived"` // 0=chưa, 1=archive
	// SalePriceTotal       *float64                `gorm:"type:DECIMAL(15,0)" json:"sale_price_total,omitempty"`                                                      // Giá bán niêm yết (công khai)
	// SalePricePerM2       *float64                `gorm:"type:DECIMAL(15,2)" json:"sale_price_per_m2,omitempty"`                                                     // Giá bán/m² (công khai)
	// SaleCommission       *float64                `gorm:"type:DECIMAL(10,2)" json:"sale_commission,omitempty"`                                                       // Hoa hồng bán (nội bộ)
	// SaleCommissionType   *string                 `gorm:"type:VARCHAR(10);check:sale_commission_type IN ('percent', 'fixed')" json:"sale_commission_type,omitempty"` // "percent" hoặc "fixed"
	// RentPrice  *float64 `gorm:"type:DECIMAL(15,0)" json:"rent_price,omitempty"`      // Giá thuê/tháng (công khai)
	// RentCommission       *float64                `gorm:"type:DECIMAL(10,2)" json:"rent_commission,omitempty"`                                                       // Hoa hồng thuê (nội bộ)
	// RentCommissionType   *string                 `gorm:"type:VARCHAR(10);check:rent_commission_type IN ('percent', 'fixed')" json:"rent_commission_type,omitempty"` // "percent" hoặc "fixed"
	// GhiChuNoiBo    *string            `gorm:"type:TEXT" json:"ghi_chu_noi_bo,omitempty"`              // Ghi chú nội bộ (môi giới)
	// ImportPrice   int64               `json:"importPrice"`
	// CommissionPercentage float64             `json:"commissionPercentage"`
	// CommissionAmount     float64             `json:"commissionAmount"`
	// Status         enums.ProductStatus `json:"status"`
	// Price          int64               `json:"price"`
	// ProjectID            *uint64                 `json:"projectId"`
	// Project              *ProjectItem            `gorm:"foreignKey:ProjectID;references:ID" json:"project,omitempty"`
	// TransactionPrice float64 `json:"transactionPrice"`
	// GiaNhap          *float64            `gorm:"type:DECIMAL(15,0)" json:"gia_nhap,omitempty"`        // Giá chủ nhà muốn net (nội bộ)
	// CreatedBy            uint64              `json:"created_by"`
	// CreatedAt            time.Time           `gorm:"autoCreateTime" json:"created_at"`
	// UpdatedAt            time.Time           `gorm:"autoUpdateTime" json:"updated_at"`
	// DeletedAt            gorm.DeletedAt      `gorm:"index" json:"-"`
	// Category       *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	// Organizations        []*Organization        `gorm:"many2many:product_organization" json:"organizations,omitempty"`
	// Customer             *Customer              `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	// Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	// Status     int16  `json:"status"`
	// Images               []ProductImage         `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	// ShareRequests        []ShareRequest         `gorm:"foreignKey:ProductID" json:"share_requests,omitempty"`
	// AssignMembers        []*OrganizationRequest `gorm:"many2many:product_member" json:"assign_members,omitempty"`
}

func (Product) TableName() string {
	return "products"
}

type ProductItem struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:255" json:"name"`
	Code string `gorm:"size:15" json:"code"`
}

func (ProductItem) TableName() string {
	return "products"
}
