package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
	"time"

	"github.com/lib/pq"
)

// Product đại diện cho bảng sản phẩm

type Product struct {
	_models.BaseEntity
	ID              uint64               `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	ParentId        *uint64              `gorm:"type:int8" json:"parentId,omitempty"`
	DepositeID      *uint64              `gorm:"type:int8" json:"depositeId,omitempty"`
	Deposite        *Deposite            `gorm:"foreignKey:DepositeID;references:ID" json:"deposite,omitempty"`
	AssetId         *uint64              `gorm:"type:int8" json:"assetId,omitempty"`
	PropertyTypeId  *uint64              `json:"propertyTypeId,omitempty"`
	PropertyType    *PropertyTypeItem    `gorm:"foreignKey:PropertyTypeId;references:ID" json:"propertyType,omitempty"`
	DocTypeId       *uint64              `json:"docTypeId,omitempty"`
	DocType         *DocTypeItem         `gorm:"foreignKey:DocTypeId;references:ID" json:"docType,omitempty"`
	Amenities       []AmenityItem        `gorm:"many2many:product_amenity;foreignKey:ID;joinForeignKey:ProductID;References:ID;joinReferences:AmenityID" json:"amenities,omitempty"`
	Visibility      enums.EVisibility    `json:"visibility,omitempty" gorm:"default:10"` //phạm vi
	Name            string               `gorm:"size:255" json:"name,omitempty"`
	Code            string               `gorm:"size:15" json:"code,omitempty"`
	CategoryID      *uint64              `json:"categoryId,omitempty"`
	Area            float64              `json:"area"`
	PositionUrl     string               `json:"positionUrl,omitempty"`
	Description     string               `json:"description,omitempty"`
	OrganizationID  *uint64              `json:"organizationId,omitempty"`
	OrgStatus       int16                `json:"orgStatus,omitempty"`
	LastPriceID     *uint64              `json:"lastPriceId,omitempty"`
	Price           *ProductPrice        `gorm:"foreignKey:LastPriceID;references:ID" json:"price,omitempty"`
	ApartmentID     *uint64              `json:"apartmentId,omitempty"`
	Apartment       *Apartment           `gorm:"foreignKey:ApartmentID;references:ID" json:"apartment,omitempty"`
	ApartmentAttr   *ApartmentAttribute  `gorm:"-" json:"apartmentAttribute,omitempty"`
	Build           *ProjectBuild        `gorm:"-" json:"build,omitempty"`
	Project         *Project             `gorm:"-" json:"project,omitempty"`
	Developer       *Developer           `gorm:"-" json:"developer,omitempty"`
	SourceType      int16                `json:"sourceType,omitempty"`
	CustomerID      *uint64              `json:"customerId,omitempty"`
	PrivateData     *ProductPrivate      `json:"privateInfo,omitempty"`
	Note            string               `json:"note,omitempty"`
	TransactionType uint                 `gorm:"type:smallint;default:10" json:"transactionType"`
	SaleStatus      enums.EProductStatus `gorm:"type:smallint;not null;default:10" json:"saleStatus"`     // 1=Chưa bán, 2=Đang bán, 3=Đã bán
	SaleVisibility  enums.EVisibility    `gorm:"type:smallint;not null;default:10" json:"saleVisibility"` // 1=Riêng tư, 2=Nội bộ, 3=Công khai
	RentStatus      enums.EProductStatus `gorm:"type:smallint;not null;default:40" json:"rentStatus"`     // 1=Chưa thuê, 2=Đang đăng tin, 3=Đang cho thuê
	RentVisibility  enums.EVisibility    `gorm:"type:smallint;not null;default:10" json:"rentVisibility"` // 1=Riêng tư, 2=Nội bộ, 3=Công khai
	OwnerUserID     *uint64              `gorm:"type:int8" json:"ownerUserId,omitempty"`                  // Ai sở hữu (nếu đã bán)
	OwnerType       enums.EOwnerType     `gorm:"type:smallint;not null;default:10" json:"ownerType"`      // 1=Cá nhân, 2=Tổ chức
	ProvinceID      *uint64              `json:"provinceId,omitempty"`
	Province        *Region              `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
	DistrictID      *uint64              `json:"districtId,omitempty"`
	District        *Region              `gorm:"foreignKey:DistrictID;references:ID" json:"district,omitempty"`
	WardID          *uint64              `json:"wardId,omitempty"`
	Ward            *Region              `gorm:"foreignKey:WardID;references:ID" json:"ward,omitempty"`
	Address         string               `gorm:"size:255" json:"address,omitempty"`
	GoogleMapLink   string               `gorm:"size:255" json:"googleMapLink,omitempty"`
	Attributes      []AttributeProduct   `gorm:"foreignKey:ProductID" json:"attributes,omitempty"`
	MediaList       []ProductMediaItem   `gorm:"foreignKey:ProductID" json:"mediaList,omitempty"`
	HouseInfo       *HouseInfo           `gorm:"foreignKey:ProductID" json:"houseInfo,omitempty"`
	Archived        bool                 `gorm:"type:bool;not null;default:false" json:"archived"` // 0=chưa, 1=archive
	// OwnerOrgID      *uint64              `gorm:"type:int8" json:"ownerOrgId,omitempty"`                   // Tổ chức sở hữu (nếu cty mua)
	CreatedUser *ProfileInfo `gorm:"foreignKey:CreatedBy;references:ProfileId" json:"createdUser,omitempty"`
	UpdatedUser *ProfileInfo `gorm:"foreignKey:UpdatedBy;references:ProfileId" json:"updatedUser,omitempty"`
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

type ProductMarket struct {
	ID               int64      `gorm:"column:id" json:"id"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	Name             string     `gorm:"column:name" json:"name"`
	Code             string     `gorm:"column:code" json:"code"`
	CategoryID       *int64     `gorm:"column:category_id" json:"categoryId"`
	Area             *float64   `gorm:"column:area" json:"area"`
	Description      *string    `gorm:"column:description" json:"description"`
	TransactionType  *int16     `gorm:"column:transaction_type" json:"transactionType"`
	SaleStatus       int16      `gorm:"column:sale_status" json:"saleStatus"`
	SaleVisibility   int16      `gorm:"column:sale_visibility" json:"saleVisibility"`
	RentStatus       int16      `gorm:"column:rent_status" json:"rentStatus"`
	RentVisibility   int16      `gorm:"column:rent_visibility" json:"rentVisibility"`
	OwnerUserID      *int64     `gorm:"column:owner_user_id" json:"ownerUserId"`
	OwnerOrgID       *int64     `gorm:"column:owner_org_id" json:"ownerOrgId"`
	PropertyTypeID   *int64     `gorm:"column:property_type_id" json:"propertyTypeId"`
	DocTypeID        *int64     `gorm:"column:doc_type_id" json:"docTypeId"`
	PositionURL      *string    `gorm:"column:position_url" json:"positionUrl"`
	Address          *string    `gorm:"column:address" json:"address"`
	GoogleMapLink    *string    `gorm:"column:google_map_link" json:"googleMapLink"`
	AssetID          *int64     `gorm:"column:asset_id" json:"assetId"`
	ApartmentID      *int64     `gorm:"column:apartment_id" json:"apartmentId"`
	LastPriceID      *int64     `gorm:"column:last_price_id" json:"lastPriceId"`
	TransactionPrice *float64   `gorm:"column:transaction_price" json:"transactionPrice"`
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"deletedAt"`
	Archived         bool       `gorm:"column:archived" json:"archived"`

	// House Info
	NumBedroom  *int64  `gorm:"column:num_bedroom" json:"numBedroom"`
	NumBathroom *int64  `gorm:"column:num_bathroom" json:"numBathroom"`
	NumFloor    *int64  `gorm:"column:num_floor" json:"numFloor"`
	Furniture   *string `gorm:"column:furniture" json:"furniture"`
	Orientation *string `gorm:"column:orientation" json:"orientation"`
	NumFront    *int64  `gorm:"column:num_front" json:"numFront"`
	NumCarPark  *int64  `gorm:"column:num_car_park" json:"numCarPark"`
	NumToilet   *int64  `gorm:"column:num_toilet" json:"numToilet"`

	// Amenities and Media
	AmenityIDs *string        `gorm:"column:amenity_ids" json:"amenityIds"`
	MediaURLs  pq.StringArray `gorm:"type:varchar[];column:media_urls" json:"mediaUrls" sql:"type:varchar[]"`
	MediaTypes pq.StringArray `gorm:"type:varchar[];column:media_types" json:"mediaTypes" sql:"type:varchar[]"`
	MainImage  *string        `gorm:"column:main_image" json:"mainImage"`

	// Region Names
	ProvinceName *string `gorm:"column:province_name" json:"provinceName"`
	DistrictName *string `gorm:"column:district_name" json:"districtName"`
	WardName     *string `gorm:"column:ward_name" json:"wardName"`

	// DocType and PropertyType Name
	DocTypeName      *string `gorm:"column:doc_type_name" json:"docTypeName"`
	PropertyTypeName *string `gorm:"column:property_type_name" json:"propertyTypeName"`

	// Latest Product Price Info
	Currency           *string  `gorm:"column:currency" json:"currency"`
	PriceOwner         *int64   `gorm:"column:price_owner" json:"priceOwner"`
	SalePrice          *float64 `gorm:"column:sale_price" json:"salePrice"`
	SaleCommission     *int64   `gorm:"column:sale_commission" json:"saleCommission"`
	SaleCommissionType *int64   `gorm:"column:sale_commission_type" json:"saleCommissionType"`
	RentPrice          *float64 `gorm:"column:rent_price" json:"rentPrice"`
	RentCommission     *float64 `gorm:"column:rent_commission" json:"rentCommission"`
	RentCommissionType *int64   `gorm:"column:rent_commission_type" json:"rentCommissionType"`
	RentPaymentCycle   *int64   `gorm:"column:rent_payment_cycle" json:"rentPaymentCycle"`
	ImportPrice        *int64   `gorm:"column:import_price" json:"importPrice"`
	OperatingCost      *float64 `gorm:"column:operating_cost" json:"operatingCost"`
	InternalNote       *int64   `gorm:"column:internal_note" json:"internalNote"`
	TargetProfit       *int64   `gorm:"column:target_profit" json:"targetProfit"`
	PrivateDocs        *float64 `gorm:"column:private_docs" json:"privateDocs"`
	Deposite           *float64 `gorm:"column:deposite" json:"deposite"`
}

func (ProductMarket) TableName() string {
	return "product_market"
}
