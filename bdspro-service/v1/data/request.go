package data

import (
	"bdspro/v1/enums"
	"bdspro/v1/models"
	"time"
)

type NewProductRequest struct {
	Name        string               `gorm:"size:255" json:"name"`
	Code        string               `gorm:"size:15" json:"code"`
	Area        float64              `json:"area"`
	Price       *models.ProductPrice `json:"price"`
	Description string               `json:"description"`
	Note        string               `json:"note"`
}

type TextSearchRequest struct {
	Text string `form:"text"`
}
type ProductSearchRequest struct {
	Name             string     `form:"name"`
	Code             string     `form:"code"`
	Area             float64    `form:"area"`
	PriceFrom        uint64     `form:"priceFrom"`
	PriceTo          uint64     `form:"priceTo"`
	NumBedrooms      []uint     `form:"numBedrooms"`
	NumBathrooms     []uint     `form:"numBathrooms"`
	NumFloors        []uint     `form:"numFloors"`
	DocTypeIds       []uint64   `form:"docTypeIds"`
	PropertyTypeIds  []uint64   `form:"propertyTypeIds"`
	AmenityIds       []uint64   `form:"amenityIds"`
	SourceTypeIds    []uint     `form:"sourceTypeIds"`
	TransactionTypes []uint     `form:"transactionTypes"`
	SaleStatus       []uint     `form:"saleStatus"`
	RentStatus       []uint     `form:"rentStatus"`
	ProvinceIds      []uint64   `form:"provinceIds"`
	DistrictIds      []uint64   `form:"districtIds"`
	WardIds          []uint64   `form:"wardIds"`
	Description      string     `form:"description"`
	Note             string     `form:"note"`
	Timestamp        *time.Time `form:"timestamp"`
}

type MediaItem struct {
	MediaURL  string `json:"mediaUrl"`
	MediaType string `json:"mediaType"`
	IsMain    bool   `json:"isMain"`
	Order     int    `json:"order"`
}

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
	Name            string                  `gorm:"size:255" json:"name"`
	Code            string                  `gorm:"size:15" json:"code"`
	Area            float64                 `json:"area"`
	Description     string                  `json:"description"`
	Note            string                  `json:"note"`
	PropertyTypeId  *uint64                 `json:"propertyTypeId"`
	DocTypeId       *uint64                 `json:"docTypeId"`
	AmenityIds      []uint64                `json:"amenityIds"`
	ProjectId       *uint64                 `json:"projectId"`
	ProvinceID      *uint64                 `json:"provinceId"`
	DistrictID      *uint64                 `json:"districtId"`
	WardID          *uint64                 `json:"wardId"`
	TransactionType int                     `json:"transactionType"`
	SaleStatus      enums.ProductStatus     `json:"saleStatus"`
	SaleVisibility  enums.ProductVisibility `json:"saleVisibility"`
	RentStatus      uint                    `json:"rentStatus"`     // 1=Chưa thuê, 2=Đang đăng tin, 3=Đang cho thuê
	RentVisibility  uint                    `json:"rentVisibility"` // 1=Riêng tư, 2=Nội bộ, 3=Công khai
	PriceDTO        *models.ProductPrice    `json:"price"`
	MediaItemDTOs   []MediaItem             `json:"mediaItems"`
	HouseInfoDTO    *HouseInfoUpdate        `json:"houseInfo"`
	PrivateDataDTO  *models.ProductPrivate  `json:"productPrivate"`
	GoogleMapLink   string                  `json:"googleMapLink"`
}
type ProductSaveDetect struct {
	Content string `json:"content"`
}

//	type ElementRequest struct {
//		Floor   uint ``
//		Ordinal uint
//	}
//
//	type AttributeRequest struct {
//		Floor   uint ``
//		Ordinal uint
//	}
type ApartmentRequest struct {
	BuildID   uint64                    `json:"buildId"`
	Attribute models.ApartmentAttribute `json:"attribute"`
	Elements  []models.Apartment        `json:"elements"`
}

type ApartmentStatusRequest struct {
	BuildID uint64 `json:"buildId"`
	Data    struct {
		Status   int `json:"status"`
		Archived int `json:"archived"`
	} `json:"data"`
	Apartments []models.Apartment `json:"apartments"`
}

type RegionRequest struct {
	Text     string `form:"text"`
	Level    *int   `form:"level"`
	ParentID *uint  `form:"parentId"`
	Page     int    `form:"page"`
	PageSize int    `form:"size"`
}

type SearchName struct {
	Name     string `form:"name"`
	Page     int    `form:"page"`
	PageSize int    `form:"size"`
}
