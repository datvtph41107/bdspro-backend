package dto

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_dto "common/domain/dto"
	sharepb "pb/types/shared"
	"time"
)

// type MediaItem struct {
// 	MediaURL  string `json:"mediaUrl"`
// 	MediaType string `json:"mediaType"`
// 	IsMain    bool   `json:"isMain"`
// 	Order     int    `json:"order"`
// }

type ProductCacheVersionResponse struct {
	HasChanged    bool
	ProductIds    []uint64
	LastUpdatedAt int64
}

type HouseInfoUpdate struct {
	// house info
	NumBedroom  *int32  `json:"numBedroom,omitempty"`  // Số phòng ngủ (optional)
	NumBathroom *int32  `json:"numBathroom,omitempty"` // Số WC (optional)
	NumFloor    *int32  `json:"numFloor,omitempty"`    // Số tầng (optional)
	NumFront    *int32  `json:"numFront,omitempty"`    // Số mặt tiền
	NumCarPark  *int32  `json:"numCarPark,omitempty"`  // Số chỗ để ôtô
	Furniture   *string `json:"furniture,omitempty"`   // Nội thất: "full", "partial", "none"
	Orientation *string `json:"orientation,omitempty"` // Hướng nhà (Đông, Tây, Nam,...)

	OrientationHouse enums.EHouseOrient      `json:"orientationHouse,omitempty"` // Hướng nhà (Đông, Tây, Nam,...)
	CertificateHouse enums.EHouseCertificate `json:"certificateHouse,omitempty"` // Chứng nhận
	RoadWidth        *float64                `json:"roadWidth,omitempty"`        // Đường vào (m2)
	FrontWidth       *float64                `json:"frontWidth,omitempty"`       // Mặt tiền (m)
	BackWidth        *float64                `json:"backWidth,omitempty"`        // Hậu (m)
	Width            *float64                `json:"width,omitempty"`            // Rộng (m)
	Height           *float64                `json:"height,omitempty"`           // Dài (m)
}

type ProductSaveRequest struct {
	ParentID          *uint64                 `gorm:"size:255" json:"parentId"`
	Name              string                  `gorm:"size:255" json:"name" binding:"required"`
	Code              string                  `gorm:"size:15" json:"code"`
	Status            uint32                  `json:"status" binding:"required"`
	AreaLand          float64                 `json:"areaLand" binding:"required"`
	AreaTotal         float64                 `json:"areaTotal"`
	Description       string                  `json:"description"`
	Note              string                  `json:"note"`
	PropertyTypeId    *uint64                 `json:"propertyTypeId"`
	DocTypeId         *uint64                 `json:"docTypeId"`
	AmenityIds        []uint64                `json:"amenityIds"`
	ProjectId         *uint64                 `json:"projectId"`
	ProvinceID        *uint64                 `json:"provinceId"`
	WardID            *uint64                 `json:"wardId"`
	TransactionType   int32                   `json:"transactionType"`
	PriceData         *domain.ProductPrice    `json:"priceData" binding:"required"`
	MediaItems        []MediaItem             `json:"mediaItems"`
	HouseInfo         *HouseInfoUpdate        `json:"houseInfo"`
	PrivateData       *domain.ProductPrivate  `json:"productPrivate"`
	GoogleMapLink     string                  `json:"googleMapLink"`
	SourceType        uint32                  `json:"sourceType"`
	SourceContactId   *uint64                 `json:"sourceContactId,omitempty"`
	SourceContactNote string                  `json:"sourceContactNote,omitempty"`
	PostCreate        *PostSaveRequest        `json:"post"`
	AssetCreate       *AssetCreateWithProduct `json:"asset"`
	SaleVisibility    enums.EVisibility       `json:"saleVisibility"`
	RentVisibility    enums.EVisibility       `json:"rentVisibility"`
	Visibility        enums.EVisibility       `json:"visibility"`
	OwnerType         enums.EOwnerOf          `json:"ownerType"`
	OwnerId           uint64                  `json:"ownerId"`
	ImageId           *uint64                 `json:"imageId"`

	CertificateHouseNote string `json:"certificateHouseNote,omitempty"`
	// SaleStatus      enums.EProductStatus    `json:"saleStatus"`
	// RentStatus      enums.EProductStatus    `json:"rentStatus"`
	// SaleStatus      enums.ProductStatus     `json:"saleStatus"`
	// SaleVisibility  enums.ProductVisibility `json:"saleVisibility"`
	// RentStatus      uint                    `json:"rentStatus"`     // 1=Chưa thuê, 2=Đang đăng tin, 3=Đang cho thuê
	// RentVisibility  uint                    `json:"rentVisibility"` // 1=Riêng tư, 2=Nội bộ, 3=Công khai

	RequestOrganizationID *uint64 `json:"-"`
}

type UpdateProductRequest struct {
	ProductID   *uint64           `json:"productId"`
	Name        string            `json:"name"`
	Status      uint32            `json:"status"`
	Visibility  enums.EVisibility `json:"visibility"`
	Description string            `json:"description"`
	Note        string            `json:"note"`

	//
	PropertyTypeID   *uint64            `json:"propertyTypeId"`
	AreaLand         float64            `json:"area"`
	AreaTotal        float64            `json:"areaTotal"`
	Direction        enums.EHouseOrient `json:"direction"`
	BalconyDirection enums.EHouseOrient `json:"balconyDirection"`
	Bedrooms         *uint32            `json:"bedrooms,omitempty"`  // Số phòng ngủ
	Bathrooms        *uint32            `json:"bathrooms,omitempty"` // Số phòng tắm
	RoadWidth        *float64           `json:"roadWidth,omitempty"` // Đường vào (m2)
	Floors           *uint32            `json:"floors,omitempty"`    // Số tầng

	AmenityIDs []uint64 `json:"amenityIds"`

	Price           *float64 `json:"price"`
	ImportPrice     *float64 `json:"importPrice"`
	CommissionValue *float64 `json:"commissionValue"`
	CommissionType  uint32   `json:"commissionType"`
	ChannelPrice    bool     `json:"channelPrice"`
	ChangeNote      string   `json:"changeNote"`

	MediaItems []MediaItem `json:"mediaItems"`

	// Nguồn
	ResourceType      enums.ESourceType   `json:"sourceType"`
	SourceContactId   uint64              `json:"sourceContactId,omitempty"`
	SourceContactNote string              `json:"sourceContactNote,omitempty"`
	SourceStatus      enums.ESourceStatus `json:"sourceStatus"`
	Priority          enums.EPriority     `json:"priority"`
	SourceTagIDs      []uint64            `json:"sourceTagIds"`

	// dự án
	ProjectID         *uint64 `json:"projectId"`
	MiningMode        enums.EMiningMode
	MiningScope       string
	MiningDescription string
	SourceType        uint32
}

type ContactSourceProfileDTO struct {
	ProfileId uint64
}

type UpdateProductSourceRequest struct {
	ProductID         uint64              `json:"productId"`
	SourceType        enums.ESourceType   `json:"sourceType"`
	SourceContactId   uint64              `json:"sourceContactId,omitempty"`
	SourceContactNote string              `json:"sourceContactNote,omitempty"`
	SourceStatus      enums.ESourceStatus `json:"sourceStatus"`
	Priority          enums.EPriority     `json:"priority"`
	SourceTagIDs      []uint64            `json:"sourceTagIds"`
	// AssignedUserID    *uint64           `json:"assignedUserId,omitempty"`
}

type ProductSaveDetect struct {
	Content string `json:"content" form:"content"`
}

type ArchivedRequest struct {
	ProductId *uint64 `json:"productId" binding:"required"`
	Archived  bool    `json:"archived" binding:"required"`
}

type StatusUpdate struct {
	ProductID  *uint64                  `json:"productId" binding:"required"`
	SaleStatus enums.EProductSaleStatus `json:"saleStatus" binding:"required"`
	RentStatus enums.EProductRentStatus `json:"rentStatus" binding:"required"`
	ContactID  *uint64                  `json:"contactId"`
	Note       string                   `json:"note"`
	Amount     float64                  `json:"amount"`
	Phone      string                   `json:"phone"`
}

type VisibilityUpdate struct {
	ProductID  *uint64           `json:"productId" binding:"required"`
	Visibility enums.EVisibility `json:"visibility" binding:"required"`
}

type DepositeRequest struct {
	ProductID     *uint64 `json:"productId" binding:"required"`
	CustomerName  string  `json:"customerName" binding:"required,min=2"`
	CustomerPhone string  `json:"customerPhone" binding:"required,min=9"`
	Amount        uint64  `json:"amount" binding:"required"`
	Note          string  `json:"note"`
}

type CreateResponse struct {
	Product *domain.Product `json:"product"`
	Post    *domain.Post    `json:"post"`
	Asset   *domain.Asset   `json:"asset"`
}

// type ProductMarket struct {
// 	ID               int64            `gorm:"column:id" json:"id"`
// 	CreatedAt        time.Time        `gorm:"column:created_at" json:"createdAt"`
// 	UpdatedAt        time.Time        `gorm:"column:updated_at" json:"updatedAt"`
// 	Name             string           `gorm:"column:name" json:"name"`
// 	Code             string           `gorm:"column:code" json:"code"`
// 	CategoryID       *int64           `gorm:"column:category_id" json:"categoryId"`
// 	Area             *float64         `gorm:"column:area" json:"area"`
// 	Description      *string          `gorm:"column:description" json:"description"`
// 	TransactionType  *int16           `gorm:"column:transaction_type" json:"transactionType"`
// 	SaleStatus       int16            `gorm:"column:sale_status" json:"saleStatus"`
// 	SaleVisibility   int16            `gorm:"column:sale_visibility" json:"saleVisibility"`
// 	RentStatus       int16            `gorm:"column:rent_status" json:"rentStatus"`
// 	RentVisibility   int16            `gorm:"column:rent_visibility" json:"rentVisibility"`
// 	OwnerID          *int64           `gorm:"column:owner_id" json:"ownerId"`
// 	OwnerType        enums.EOwnerOf `json:"ownerType"`
// 	PropertyTypeID   *int64           `gorm:"column:property_type_id" json:"propertyTypeId"`
// 	DocTypeID        *int64           `gorm:"column:doc_type_id" json:"docTypeId"`
// 	PositionURL      *string          `gorm:"column:position_url" json:"positionUrl"`
// 	Address          *string          `gorm:"column:address" json:"address"`
// 	GoogleMapLink    *string          `gorm:"column:google_map_link" json:"googleMapLink"`
// 	AssetID          *int64           `gorm:"column:asset_id" json:"assetId"`
// 	ApartmentID      *int64           `gorm:"column:apartment_id" json:"apartmentId"`
// 	LastPriceID      *int64           `gorm:"column:last_price_id" json:"lastPriceId"`
// 	TransactionPrice *float64         `gorm:"column:transaction_price" json:"transactionPrice"`
// 	DeletedAt        *time.Time       `gorm:"column:deleted_at" json:"deletedAt"`
// 	Archived         bool             `gorm:"column:archived" json:"archived"`

// 	// House Info
// 	NumBedroom  *int64  `gorm:"column:num_bedroom" json:"numBedroom"`
// 	NumBathroom *int64  `gorm:"column:num_bathroom" json:"numBathroom"`
// 	NumFloor    *int64  `gorm:"column:num_floor" json:"numFloor"`
// 	Furniture   *string `gorm:"column:furniture" json:"furniture"`
// 	Orientation *string `gorm:"column:orientation" json:"orientation"`
// 	NumFront    *int64  `gorm:"column:num_front" json:"numFront"`
// 	NumCarPark  *int64  `gorm:"column:num_car_park" json:"numCarPark"`
// 	NumToilet   *int64  `gorm:"column:num_toilet" json:"numToilet"`

// 	// Amenities and Media
// 	AmenityIDs *string        `gorm:"column:amenity_ids" json:"amenityIds"`
// 	MediaURLs  pq.StringArray `gorm:"type:varchar[];column:media_urls" json:"mediaUrls" sql:"type:varchar[]"`
// 	MediaTypes pq.StringArray `gorm:"type:varchar[];column:media_types" json:"mediaTypes" sql:"type:varchar[]"`
// 	MainImage  *string        `gorm:"column:main_image" json:"mainImage"`

// 	// Region Names
// 	ProvinceName *string `gorm:"column:province_name" json:"provinceName"`
// 	DistrictName *string `gorm:"column:district_name" json:"districtName"`
// 	WardName     *string `gorm:"column:ward_name" json:"wardName"`

// 	// DocType and PropertyType Name
// 	DocTypeName      *string `gorm:"column:doc_type_name" json:"docTypeName"`
// 	PropertyTypeName *string `gorm:"column:property_type_name" json:"propertyTypeName"`

// 	// Latest Product Price Info
// 	Currency           *string  `gorm:"column:currency" json:"currency"`
// 	PriceOwner         *int64   `gorm:"column:price_owner" json:"priceOwner"`
// 	SalePrice          *float64 `gorm:"column:sale_price" json:"salePrice"`
// 	SaleCommission     *int64   `gorm:"column:sale_commission" json:"saleCommission"`
// 	SaleCommissionType *int64   `gorm:"column:sale_commission_type" json:"saleCommissionType"`
// 	RentPrice          *float64 `gorm:"column:rent_price" json:"rentPrice"`
// 	RentCommission     *float64 `gorm:"column:rent_commission" json:"rentCommission"`
// 	RentCommissionType *int64   `gorm:"column:rent_commission_type" json:"rentCommissionType"`
// 	RentPaymentCycle   *int64   `gorm:"column:rent_payment_cycle" json:"rentPaymentCycle"`
// 	ImportPrice        *int64   `gorm:"column:import_price" json:"importPrice"`
// 	OperatingCost      *float64 `gorm:"column:operating_cost" json:"operatingCost"`
// 	InternalNote       *int64   `gorm:"column:internal_note" json:"internalNote"`
// 	TargetProfit       *int64   `gorm:"column:target_profit" json:"targetProfit"`
// 	PrivateDocs        *float64 `gorm:"column:private_docs" json:"privateDocs"`
// 	Deposite           *float64 `gorm:"column:deposite" json:"deposite"`
// 	// OwnerOrgID       *int64     `gorm:"column:owner_org_id" json:"ownerOrgId"`

// }

type ProductSearch struct {
	_dto.Pagable
	OwnerID          *uint64 `json:"ownerId"`
	OwnerType        *uint64 `json:"ownerType"`
	ParentID         *uint64 `json:"parentId"`
	PropertyTypeID   *uint64 `json:"propertyTypeId"`
	DocTypeID        *uint64 `json:"docTypeId"`
	ProjectID        *uint64 `json:"projectId"`
	ProvinceID       *uint64 `json:"provinceId"`
	WardID           *uint64 `json:"wardId"`
	TransactionType  *int32  `json:"transactionType"`
	SaleStatus       *int32  `json:"saleStatus"`
	RentStatus       *int32  `json:"rentStatus"`
	SaleVisibility   *int32  `json:"saleVisibility"`
	RentVisibility   *int32  `json:"rentVisibility"`
	SourceType       *int32  `json:"sourceType"`
	TransactionPrice *int32  `json:"transactionPrice"`
	SalePrice        *int32  `json:"salePrice"`
	RentPrice        *int32  `json:"rentPrice"`
	SaleCommission   *int32  `json:"saleCommission"`
	RentCommission   *int32  `json:"rentCommission"`
	RentPaymentCycle *int32  `json:"rentPaymentCycle"`
	Archived         *bool   `json:"archived"`
}

type GroupProductSearchRequest struct {
	_dto.Pagable
	GroupID          uint64  `json:"groupId" binding:"required"`
	Text             string  `json:"text"`
	ParentID         *uint64 `json:"parentId"`
	PropertyTypeID   *uint64 `json:"propertyTypeId"`
	DocTypeID        *uint64 `json:"docTypeId"`
	ProjectID        *uint64 `json:"projectId"`
	ProvinceID       *uint64 `json:"provinceId"`
	WardID           *uint64 `json:"wardId"`
	TransactionType  *int32  `json:"transactionType"`
	SaleStatus       *int32  `json:"saleStatus"`
	RentStatus       *int32  `json:"rentStatus"`
	SaleVisibility   *int32  `json:"saleVisibility"`
	RentVisibility   *int32  `json:"rentVisibility"`
	SourceType       *int32  `json:"sourceType"`
	TransactionPrice *int32  `json:"transactionPrice"`
	SalePrice        *int32  `json:"salePrice"`
	RentPrice        *int32  `json:"rentPrice"`
	SaleCommission   *int32  `json:"saleCommission"`
	RentCommission   *int32  `json:"rentCommission"`
	RentPaymentCycle *int32  `json:"rentPaymentCycle"`
	Archived         *bool   `json:"archived"`
}

type UserProductSearchRequest struct {
	_dto.Pagable
	UserID           uint64  `json:"userId" binding:"required"`
	Text             string  `json:"text"`
	ParentID         *uint64 `json:"parentId"`
	PropertyTypeID   *uint64 `json:"propertyTypeId"`
	DocTypeID        *uint64 `json:"docTypeId"`
	ProjectID        *uint64 `json:"projectId"`
	ProvinceID       *uint64 `json:"provinceId"`
	WardID           *uint64 `json:"wardId"`
	TransactionType  *int32  `json:"transactionType"`
	SaleStatus       *int32  `json:"saleStatus"`
	RentStatus       *int32  `json:"rentStatus"`
	SaleVisibility   *int32  `json:"saleVisibility"`
	RentVisibility   *int32  `json:"rentVisibility"`
	SourceType       *int32  `json:"sourceType"`
	TransactionPrice *int32  `json:"transactionPrice"`
	SalePrice        *int32  `json:"salePrice"`
	RentPrice        *int32  `json:"rentPrice"`
	SaleCommission   *int32  `json:"saleCommission"`
	RentCommission   *int32  `json:"rentCommission"`
	RentPaymentCycle *int32  `json:"rentPaymentCycle"`
	Archived         *bool   `json:"archived"`
}

type ProductListResponse struct {
	ID              uint64              `gorm:"column:id" json:"id"`
	PropertyId      uint64              `gorm:"column:property_id" json:"propertyId"`
	ImageID         *uint64             `gorm:"column:image_id" json:"imageId"`
	ImageURL        string              `gorm:"column:image_url" json:"imageUrl"`
	Name            string              `gorm:"column:name" json:"name"`
	Code            string              `gorm:"column:code" json:"code"`
	Area            float64             `gorm:"column:area" json:"area"`
	AvailableArea   float64             `gorm:"column:available_area" json:"availableArea"`
	PropertyTypeId  uint64              `gorm:"column:property_type_id" json:"propertyTypeId"`
	DocTypeId       uint64              `gorm:"column:doc_type_id" json:"docTypeId"`
	ProjectId       uint64              `gorm:"column:project_id" json:"projectId"`
	ProjectName     string              `gorm:"column:project_name" json:"projectName"`
	ProvinceId      uint64              `gorm:"column:province_id" json:"provinceId"`
	ProvinceName    string              `gorm:"column:province_name" json:"provinceName"`
	WardId          uint64              `gorm:"column:ward_id" json:"wardId"`
	WardName        string              `gorm:"column:ward_name" json:"wardName"`
	TransactionType int32               `gorm:"column:transaction_type" json:"transactionType"`
	SaleStatus      uint32              `gorm:"column:sale_status" json:"saleStatus"`
	RentStatus      uint32              `gorm:"column:rent_status" json:"rentStatus"`
	SaleVisibility  uint32              `gorm:"column:sale_visibility" json:"saleVisibility"`
	RentVisibility  uint32              `gorm:"column:rent_visibility" json:"rentVisibility"`
	SourceType      uint32              `gorm:"column:source_type" json:"sourceType"`
	Price           domain.ProductPrice `gorm:"-" json:"price"`

	PropertyTypeName string `gorm:"column:property_type_name" json:"propertyTypeName"`

	NumBedroom  *int32                  `gorm:"column:num_bedroom" json:"numBedroom"`
	NumBathroom *int32                  `gorm:"column:num_bathroom" json:"numBathroom"`
	Address     *sharepb.AddressV3Proto `gorm:"-" json:"address"`
	// Internal fields for address mapping
	AddressDetail string  `gorm:"column:address" json:"-"`
	DistrictID    *uint64 `gorm:"column:district_id" json:"-"`
	DistrictName  string  `gorm:"column:district_name" json:"-"`

	PriceID             uint64  `gorm:"column:price_id" json:"-"`
	PriceCurrency       string  `gorm:"column:price_currency" json:"-"`
	PriceSalePrice      float64 `gorm:"column:price_sale_price" json:"-"`
	PriceSaleCommission float64 `gorm:"column:price_sale_commission" json:"-"`
	PriceSaleCommType   uint32  `gorm:"column:price_sale_commission_type" json:"-"`
	PriceDeposite       float64 `gorm:"column:price_deposite" json:"-"`
	PriceRentPrice      float64 `gorm:"column:price_rent_price" json:"-"`
	PriceRentCommission float64 `gorm:"column:price_rent_commission" json:"-"`
	PriceRentCommType   uint32  `gorm:"column:price_rent_commission_type" json:"-"`
	PriceRentPayCycle   uint32  `gorm:"column:price_rent_payment_cycle" json:"-"`

	// Statistics fields
	PublicListingCount int64 `gorm:"-" json:"publicListingCount,omitempty"` // Số tin đăng công khai

	// BuildingInfo fields (from PropertyBuildingInfo)
	BuildingInfo *domain.PropertyBuildingInfo `gorm:"-" json:"buildingInfo,omitempty"`
	// Internal fields for building info mapping
	BuildingType      *uint32  `gorm:"column:building_type" json:"-"`
	ConstructionArea  *float64 `gorm:"column:construction_area" json:"-"`
	FloorArea         *float64 `gorm:"column:floor_area" json:"-"`
	Floors            *uint32  `gorm:"column:floors" json:"-"`
	BuildingBedrooms  *uint32  `gorm:"column:building_bedrooms" json:"-"`
	BuildingBathrooms *uint32  `gorm:"column:building_bathrooms" json:"-"`
	Direction         *uint32  `gorm:"column:direction" json:"-"`
	BalconyDirection  *uint32  `gorm:"column:balcony_direction" json:"-"`

	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deletedAt"`

	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updatedAt"`
	ShareCount       uint32    `gorm:"column:share_count" json:"shareCount,omitempty"`             // Số lượt chia sẻ
	PostCount        uint32    `gorm:"column:post_count" json:"postCount,omitempty"`               // Số tin đăng
	AssetCount       uint32    `gorm:"column:asset_count" json:"assetCount,omitempty"`             // Số tài sản
	DealCount        uint32    `gorm:"column:deal_count" json:"dealCount,omitempty"`               // Số thương vụ
	AppointmentCount uint32    `gorm:"column:appointment_count" json:"appointmentCount,omitempty"` // Số lịch hẹn
	ContactCount     uint32    `gorm:"column:contact_count" json:"contactCount,omitempty"`         // Số người quan tâm

	// Avatar field
	AvatarURL       string `gorm:"column:avatar_url" json:"avatarUrl"`
	AvatarThumbURL  string `gorm:"column:avatar_thumb_url" json:"avatarThumbUrl"`
	AvatarMediaType string `gorm:"column:avatar_media_type" json:"avatarMediaType"`
	AvatarSortOrder int32  `gorm:"column:avatar_sort_order" json:"avatarSortOrder"`
}

type ProductSummaryDTO struct {
	TotalProducts   int64   `json:"totalProducts"`
	SellingProducts int64   `json:"sellingProducts"`
	TotalValue      float64 `json:"totalValue"`
}
