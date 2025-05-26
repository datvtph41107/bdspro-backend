package u_product_child

import (
	"bdspro/v1/data"
	u_product_dto "bdspro/v1/mod/user/product/dto"
)

// type MergeAllProductChild struct {
// 	// u_product_dto.ProductSaveRequest
// 	ParentID *uint64 `json:"parentId" binding:"required"`
// }

type DevideChildRequest struct {
	ParentID *uint64                            `json:"parentId" binding:"required"`
	Childs   []u_product_dto.ProductSaveRequest `json:"childs" binding:"required,notemptyarray,dive"`
}

type MergeProductChildRequest struct {
	// u_product_dto.ProductSaveRequest
	ParentID *uint64   `json:"parentId" binding:"required"`
	ChildIDs *[]uint64 `json:"childIds" binding:"required"`
}
type SaveProductChildRequest struct {
	u_product_dto.ProductSaveRequest
	ParentID *uint64 `json:"parentId" binding:"required"`
	// Area     float64 `json:"area" binding:"required"`
	// Name     string  `json:"name" binding:"required"`
	// Name            string                  `gorm:"size:255" json:"name"`
	// Code            string                  `gorm:"size:15" json:"code"`
	// Area            float64                 `json:"area"`
	// Description     string                  `json:"description"`
	// Note            string                  `json:"note"`
	// PropertyTypeId  *uint64                 `json:"propertyTypeId"`
	// DocTypeId       *uint64                 `json:"docTypeId"`
	// AmenityIds      []uint64                `json:"amenityIds"`
	// ProjectId       *uint64                 `json:"projectId"`
	// ProvinceID      *uint64                 `json:"provinceId"`
	// DistrictID      *uint64                 `json:"districtId"`
	// WardID          *uint64                 `json:"wardId"`
	// TransactionType int                     `json:"transactionType"`
	// SaleStatus      enums.ProductStatus     `json:"saleStatus"`
	// SaleVisibility  enums.ProductVisibility `json:"saleVisibility"`
	// RentStatus      uint                    `json:"rentStatus"`     // 1=Chưa thuê, 2=Đang đăng tin, 3=Đang cho thuê
	// RentVisibility  uint                    `json:"rentVisibility"` // 1=Riêng tư, 2=Nội bộ, 3=Công khai
	// PriceDTO        *models.ProductPrice    `json:"price"`
	// MediaItemDTOs   []data.MediaItem        `json:"mediaItems"`
	// HouseInfoDTO    *data.HouseInfoUpdate   `json:"houseInfo"`
	// PrivateDataDTO  *models.ProductPrivate  `json:"productPrivate"`
	// GoogleMapLink   string                  `json:"googleMapLink"`
}

type ProductChildSearch struct {
	data.ProductSearchRequest
	// ParentID *uint64 `json:"parentId"`
	// Area     float64 `json:"area" binding:"required"`
	// Name     string  `json:"name" binding:"required"`
	// Name            string                  `gorm:"size:255" json:"name"`
	// Code            string                  `gorm:"size:15" json:"code"`
	// Area            float64                 `json:"area"`
	// Description     string                  `json:"description"`
	// Note            string                  `json:"note"`
	// PropertyTypeId  *uint64                 `json:"propertyTypeId"`
	// DocTypeId       *uint64                 `json:"docTypeId"`
	// AmenityIds      []uint64                `json:"amenityIds"`
	// ProjectId       *uint64                 `json:"projectId"`
	// ProvinceID      *uint64                 `json:"provinceId"`
	// DistrictID      *uint64                 `json:"districtId"`
	// WardID          *uint64                 `json:"wardId"`
	// TransactionType int                     `json:"transactionType"`
	// SaleStatus      enums.ProductStatus     `json:"saleStatus"`
	// SaleVisibility  enums.ProductVisibility `json:"saleVisibility"`
	// RentStatus      uint                    `json:"rentStatus"`     // 1=Chưa thuê, 2=Đang đăng tin, 3=Đang cho thuê
	// RentVisibility  uint                    `json:"rentVisibility"` // 1=Riêng tư, 2=Nội bộ, 3=Công khai
	// PriceDTO        *models.ProductPrice    `json:"price"`
	// MediaItemDTOs   []data.MediaItem        `json:"mediaItems"`
	// HouseInfoDTO    *data.HouseInfoUpdate   `json:"houseInfo"`
	// PrivateDataDTO  *models.ProductPrivate  `json:"productPrivate"`
	// GoogleMapLink   string                  `json:"googleMapLink"`
}

type ProductHistorySearch struct {
	Page uint `form:"page"`
	Size uint `form:"size"`
}
