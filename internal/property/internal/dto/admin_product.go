package dto

import (
	_dto "common/domain/dto"
	"time"
)

// AdminProductSearchResponse response cho admin search products
type AdminProductSearchResponse struct {
	_dto.Pagable
	Data       []AdminProductItem `json:"data"`
	Total      int64              `json:"total"`
	TotalPages uint32             `json:"totalPages"`
}

// AdminProductItem thông tin sản phẩm cho admin view
type AdminProductItem struct {
	ID              uint64     `json:"id"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	Address         string     `json:"address"`
	Province        string     `json:"province"`
	District        string     `json:"district"`
	Ward            string     `json:"ward"`
	Area            float64    `json:"area"`
	Price           float64    `json:"price"`
	PriceUnit       string     `json:"priceUnit"`
	Status          uint32     `json:"status"`
	StatusText      string     `json:"statusText"`
	OwnerName       string     `json:"ownerName"`
	OwnerEmail      string     `json:"ownerEmail"`
	OwnerPhone      string     `json:"ownerPhone"`
	CreatedAt       string     `json:"createdAt"`
	UpdatedAt       string     `json:"updatedAt"`
	ApprovedAt      *time.Time `json:"approvedAt"`
	RejectedAt      *time.Time `json:"rejectedAt"`
	RejectionReason string     `json:"rejectionReason"`
	Images          []string   `json:"images"`
	PropertyType    string     `json:"propertyType"`
	DocumentType    string     `json:"documentType"`
	Amenities       []string   `json:"amenities"`
	IsActive        bool       `json:"isActive"`
	IsArchived      bool       `json:"isArchived"`
	ViewCount       uint64     `json:"viewCount"`
	LikeCount       uint64     `json:"likeCount"`
	ShareCount      uint64     `json:"shareCount"`
}
