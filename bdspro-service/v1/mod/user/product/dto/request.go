package u_product_dto

import (
	"bdspro/enums"
	"bdspro/v1/models"
)

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
	ParentID        *uint64                `gorm:"size:255" json:"parentId"`
	Name            string                 `gorm:"size:255" json:"name" binding:"required"`
	Code            string                 `gorm:"size:15" json:"code"`
	Area            float64                `json:"area" binding:"required"`
	Description     string                 `json:"description"`
	Note            string                 `json:"note"`
	PropertyTypeId  *uint64                `json:"propertyTypeId"`
	DocTypeId       *uint64                `json:"docTypeId"`
	AmenityIds      []uint64               `json:"amenityIds"`
	ProjectId       *uint64                `json:"projectId"`
	ProvinceID      *uint64                `json:"provinceId"`
	DistrictID      *uint64                `json:"districtId"`
	WardID          *uint64                `json:"wardId"`
	TransactionType int                    `json:"transactionType"`
	PriceDTO        *models.ProductPrice   `json:"price" binding:"required"`
	MediaItemDTOs   []MediaItem            `json:"mediaItems"`
	HouseInfoDTO    *HouseInfoUpdate       `json:"houseInfo"`
	PrivateDataDTO  *models.ProductPrivate `json:"productPrivate"`
	GoogleMapLink   string                 `json:"googleMapLink"`
	// SaleStatus      enums.ProductStatus     `json:"saleStatus"`
	// SaleVisibility  enums.ProductVisibility `json:"saleVisibility"`
	// RentStatus      uint                    `json:"rentStatus"`     // 1=Chưa thuê, 2=Đang đăng tin, 3=Đang cho thuê
	// RentVisibility  uint                    `json:"rentVisibility"` // 1=Riêng tư, 2=Nội bộ, 3=Công khai
}
type ProductSaveDetect struct {
	Content string `json:"content"`
}

type ArchivedRequest struct {
	ProductId *uint64 `json:"productId" binding:"required"`
	Archived  bool    `json:"archived" binding:"required"`
}

type StatusUpdate struct {
	ProductID *uint64             `json:"productId" binding:"required"`
	Status    enums.ProductStatus `json:"status" binding:"required"`
}

type VisibilityUpdate struct {
	ProductID  *uint64                `json:"productId" binding:"required"`
	Visibility enums.VisibilityStatus `json:"visibility" binding:"required"`
}

type DepositeRequest struct {
	ProductID     uint64 `json:"productId" binding:"required"`
	CustomerName  string `json:"customerName" binding:"required,min=2"`
	CustomerPhone string `json:"customerPhone" binding:"required,min=9"`
	Amount        uint64 `json:"amount" binding:"required"`
	Note          string `json:"note"`
}
