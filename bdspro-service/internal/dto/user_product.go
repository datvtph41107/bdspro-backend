package dto

import (
	"bdspro/internal/domain"
	user_dto "bdspro/internal/dto/user"
	"bdspro/internal/enums"
	"time"

	"github.com/lib/pq"
)

// type MediaItem struct {
// 	MediaURL  string `json:"mediaUrl"`
// 	MediaType string `json:"mediaType"`
// 	IsMain    bool   `json:"isMain"`
// 	Order     int    `json:"order"`
// }

type HouseInfoUpdate struct {
	// house info
	NumBedroom  *int    `json:"numBedroom,omitempty"`  // Số phòng ngủ (optional)
	NumBathroom *int    `json:"numBathroom,omitempty"` // Số WC (optional)
	NumFloor    *int    `json:"numFloor,omitempty"`    // Số tầng (optional)
	NumFront    *int    `json:"numFront,omitempty"`    // Số mặt tiền
	NumCarPark  *int    `json:"numCarPark,omitempty"`  // Số chỗ để ôtô
	Furniture   *string `json:"furniture,omitempty"`   // Nội thất: "full", "partial", "none"
	Orientation *string `json:"orientation,omitempty"` // Hướng nhà (Đông, Tây, Nam,...)
}

type ProductSaveRequest struct {
	ParentID        *uint64                          `gorm:"size:255" json:"parentId"`
	Name            string                           `gorm:"size:255" json:"name" binding:"required"`
	Code            string                           `gorm:"size:15" json:"code"`
	Area            float64                          `json:"area" binding:"required"`
	Description     string                           `json:"description"`
	Note            string                           `json:"note"`
	PropertyTypeId  *uint64                          `json:"propertyTypeId"`
	DocTypeId       *uint64                          `json:"docTypeId"`
	AmenityIds      []uint64                         `json:"amenityIds"`
	ProjectId       *uint64                          `json:"projectId"`
	ProvinceID      *uint64                          `json:"provinceId"`
	DistrictID      *uint64                          `json:"districtId"`
	WardID          *uint64                          `json:"wardId"`
	TransactionType int                              `json:"transactionType"`
	PriceDTO        *domain.ProductPrice             `json:"price" binding:"required"`
	MediaItemDTOs   []MediaItem                      `json:"mediaItems"`
	HouseInfoDTO    *HouseInfoUpdate                 `json:"houseInfo"`
	PrivateDataDTO  *domain.ProductPrivate           `json:"productPrivate"`
	GoogleMapLink   string                           `json:"googleMapLink"`
	SourceType      uint                             `json:"sourceType"`
	SaleStatus      *enums.EProductStatus            `json:"saleStatus"`
	RentStatus      *enums.EProductStatus            `json:"rentStatus"`
	PostCreate      *PostSaveWithProduct             `json:"post"`
	AssetCreate     *user_dto.AssetCreateWithProduct `json:"asset"`
	SaleVisibility  *enums.EVisibility               `json:"saleVisibility"`
	RentVisibility  *enums.EVisibility               `json:"rentVisibility"`
	// SaleStatus      enums.ProductStatus     `json:"saleStatus"`
	// SaleVisibility  enums.ProductVisibility `json:"saleVisibility"`
	// RentStatus      uint                    `json:"rentStatus"`     // 1=Chưa thuê, 2=Đang đăng tin, 3=Đang cho thuê
	// RentVisibility  uint                    `json:"rentVisibility"` // 1=Riêng tư, 2=Nội bộ, 3=Công khai
}

type ProductSaveDetect struct {
	Content string `json:"content" form:"content"`
}

type ArchivedRequest struct {
	ProductId *uint64 `json:"productId" binding:"required"`
	Archived  bool    `json:"archived" binding:"required"`
}

type StatusUpdate struct {
	ProductID *uint64              `json:"productId" binding:"required"`
	Status    enums.EProductStatus `json:"status" binding:"required"`
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
