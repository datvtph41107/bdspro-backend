package dto

// type MergeAllProductChild struct {
// 	// u_product_dto.ProductSaveRequest
// 	ParentID *uint64 `json:"parentId" binding:"required"`
// }

type DevideChildRequest struct {
	ParentID *uint64              `json:"parentId" binding:"required"`
	Childs   []ProductSaveRequest `json:"childs" binding:"required,notemptyarray,dive"`
}

type MergeProductChildRequest struct {
	// u_product_dto.ProductSaveRequest
	ParentID *uint64   `json:"parentId" binding:"required"`
	ChildIDs *[]uint64 `json:"childIds" binding:"required"`
}
type SaveProductChildRequest struct {
	ProductSaveRequest
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
	// PriceDTO        *domain.ProductPrice    `json:"price"`
	// MediaItemDTOs   []data.MediaItem        `json:"mediaItems"`
	// HouseInfoDTO    *data.HouseInfoUpdate   `json:"houseInfo"`
	// PrivateDataDTO  *domain.ProductPrivate  `json:"productPrivate"`
	// GoogleMapLink   string                  `json:"googleMapLink"`
}

type ProductChildSearch struct {
	ProductSearchRequest
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
	// PriceDTO        *domain.ProductPrice    `json:"price"`
	// MediaItemDTOs   []data.MediaItem        `json:"mediaItems"`
	// HouseInfoDTO    *data.HouseInfoUpdate   `json:"houseInfo"`
	// PrivateDataDTO  *domain.ProductPrivate  `json:"productPrivate"`
	// GoogleMapLink   string                  `json:"googleMapLink"`
}
