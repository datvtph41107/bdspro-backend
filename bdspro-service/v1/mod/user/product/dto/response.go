package u_product_dto

import (
	"bdspro/v1/models"
	"time"

	"github.com/lib/pq"
)

type CreateResponse struct {
	Product *models.Product `json:"product"`
	Post    *models.Post    `json:"post"`
	Asset   *models.Asset   `json:"asset"`
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
