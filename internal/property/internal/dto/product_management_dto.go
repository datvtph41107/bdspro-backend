package dto

// CreateProductForUserDTO DTO cho tạo product cho user
type CreateProductForUserDTO struct {
	Product   ProductSaveRequestDTO `json:"product"`
	ProfileID uint64                `json:"profileId"`
	IsOwner   bool                  `json:"isOwner"`
	RoleID    uint64                `json:"roleId"`
}

// CreateProductForOrgDTO DTO cho tạo product cho organization
type CreateProductForOrgDTO struct {
	Product        ProductSaveRequestDTO `json:"product"`
	OrganizationID uint64                `json:"organizationId"`
	IsOwner        bool                  `json:"isOwner"`
	RoleID         uint64                `json:"roleId"`
}

// ProductSaveRequestDTO DTO cho thông tin product cần tạo
type ProductSaveRequestDTO struct {
	ParentID        *uint64            `json:"parentId,omitempty"`
	Name            string             `json:"name"`
	Code            string             `json:"code"`
	Area            float64            `json:"area"`
	Description     string             `json:"description"`
	Note            string             `json:"note"`
	PropertyTypeID  *uint64            `json:"propertyTypeId,omitempty"`
	DocTypeID       *uint64            `json:"docTypeId,omitempty"`
	AmenityIDs      []uint64           `json:"amenityIds"`
	ProjectID       *uint64            `json:"projectId,omitempty"`
	ProvinceID      *uint64            `json:"provinceId,omitempty"`
	WardID          *uint64            `json:"wardId,omitempty"`
	TransactionType int32              `json:"transactionType"`
	PriceData       *ProductPriceDTO   `json:"priceData,omitempty"`
	MediaItems      []MediaItemDTO     `json:"mediaItems"`
	HouseInfo       *HouseInfoDTO      `json:"houseInfo,omitempty"`
	ProductPrivate  *ProductPrivateDTO `json:"productPrivate,omitempty"`
	GoogleMapLink   string             `json:"googleMapLink"`
	SaleStatus      uint32             `json:"saleStatus"`
	RentStatus      uint32             `json:"rentStatus"`
	SaleVisibility  uint32             `json:"saleVisibility"`
	RentVisibility  uint32             `json:"rentVisibility"`
	OwnerType       uint32             `json:"ownerType"`

	ID uint64 `json:"id"`

	SourceType        uint32  `json:"sourceType"`
	SourceContactId   *uint64 `json:"sourceContactId,omitempty"`
	SourceContactNote string  `json:"sourceContactNote,omitempty"`
}

// ProductPriceDTO DTO cho thông tin giá
type ProductPriceDTO struct {
	Currency           string  `json:"currency"`
	SalePrice          float64 `json:"salePrice"`
	SaleCommission     int64   `json:"saleCommission"`
	SaleCommissionType int64   `json:"saleCommissionType"`
	RentPrice          float64 `json:"rentPrice"`
	RentCommission     float64 `json:"rentCommission"`
	RentCommissionType int64   `json:"rentCommissionType"`
	RentPaymentCycle   int64   `json:"rentPaymentCycle"`
	ImportPrice        int64   `json:"importPrice"`
}

// MediaItemDTO DTO cho media
type MediaItemDTO struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

// HouseInfoDTO DTO cho thông tin nhà
type HouseInfoDTO struct {
	NumBedroom  *int64  `json:"numBedroom,omitempty"`
	NumBathroom *int64  `json:"numBathroom,omitempty"`
	NumFloor    *int64  `json:"numFloor,omitempty"`
	Furniture   *string `json:"furniture,omitempty"`
	Orientation *string `json:"orientation,omitempty"`
	NumFront    *int64  `json:"numFront,omitempty"`
	NumCarPark  *int64  `json:"numCarPark,omitempty"`
	NumToilet   *int64  `json:"numToilet,omitempty"`
}

// ProductPrivateDTO DTO cho thông tin riêng tư
type ProductPrivateDTO struct {
	Note string `json:"note"`
}

// CreateProductResponseDTO DTO cho response tạo product
type CreateProductResponseDTO struct {
	ProductID uint64 `json:"productId"`
	Message   string `json:"message"`
}
