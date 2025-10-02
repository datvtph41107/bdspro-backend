package dto

import (
	"bdspro/internal/enums"
	_dto "common/domain/dto"
	"time"
)

// CreatePostForUserDTO DTO cho tạo post cho user
type CreatePostForUserDTO struct {
	Post      PostSaveRequestDTO `json:"post"`
	ProfileID uint64             `json:"profileId"`
	IsOwner   bool               `json:"isOwner"`
	RoleID    uint64             `json:"roleId"`
}

// CreatePostForOrgDTO DTO cho tạo post cho organization
type CreatePostForOrgDTO struct {
	Post           PostSaveRequestDTO `json:"post"`
	OrganizationID uint64             `json:"organizationId"`
	IsOwner        bool               `json:"isOwner"`
	RoleID         uint64             `json:"roleId"`
}

// PostSaveRequestDTO DTO cho thông tin post cần tạo
type PostSaveRequestDTO struct {
	ProductID       uint64                `json:"productId" binding:"required"`
	ExpiredAt       *time.Time            `json:"expiredAt,omitempty"`
	Visibility      enums.EVisibility     `json:"visibility,omitempty"`
	Hidden          bool                  `json:"hidden"`
	Content         string                `json:"content,omitempty"`
	TransactionType enums.TransactionType `json:"transactionType,omitempty"`
	Title           string                `json:"title" binding:"required"`
	NumDate         int                   `json:"numDate,omitempty"`
	Type            uint                  `json:"type,omitempty"` // 1: normal, 2: feed, 3: share
	Price           float64               `json:"price,omitempty"`
	PriceType       uint32                `json:"priceType,omitempty"`      // 1: chính thức, 2: thỏa thuận
	PackageVisible  uint                  `json:"packageVisible,omitempty"` // 1: thường, 2: VIP
	OwnerType       enums.EOwnerOf        `json:"ownerType"`
	MediaItems      []PostMediaDTO        `json:"mediaItems,omitempty"`
	ID              uint64                `json:"id"`
}

// PostMediaDTO DTO cho media của post
type PostMediaDTO struct {
	URL      string `json:"url"`
	Type     string `json:"type"`
	Position int    `json:"position,omitempty"`
}

// CreatePostResponseDTO DTO cho response tạo post
type CreatePostResponseDTO struct {
	PostID  uint64 `json:"postId"`
	Message string `json:"message"`
}

// PostListRequestDTO DTO cho request danh sách tin đăng
type PostListRequestDTO struct {
	_dto.Pagable
	OrganizationID  *uint64                `json:"organizationId,omitempty"`
	ProfileID       *uint64                `json:"profileId,omitempty"`
	OwnerType       *enums.EOwnerOf        `json:"ownerType,omitempty"`
	ProductID       *uint64                `json:"productId,omitempty"`
	TransactionType *enums.TransactionType `json:"transactionType,omitempty"`
	Visibility      *enums.EVisibility     `json:"visibility,omitempty"`
	Hidden          *bool                  `json:"hidden,omitempty"`
	Type            *uint                  `json:"type,omitempty"`
	PackageVisible  *uint                  `json:"packageVisible,omitempty"`
	PublishedAtFrom *time.Time             `json:"publishedAtFrom,omitempty"`
	PublishedAtTo   *time.Time             `json:"publishedAtTo,omitempty"`
	ExpiredAtFrom   *time.Time             `json:"expiredAtFrom,omitempty"`
	ExpiredAtTo     *time.Time             `json:"expiredAtTo,omitempty"`
	Search          string                 `json:"search,omitempty"`
}

// PostListResponseDTO DTO cho response danh sách tin đăng
type PostListResponseDTO struct {
	_dto.Pagable
	Data []PostDetailDTO `json:"data"`
}

// PostDetailDTO DTO cho chi tiết tin đăng
type PostDetailDTO struct {
	ID              uint64                `json:"id"`
	ProductID       uint64                `json:"productId"`
	ExpiredAt       *time.Time            `json:"expiredAt,omitempty"`
	Visibility      enums.EVisibility     `json:"visibility,omitempty"`
	Hidden          bool                  `json:"hidden"`
	Content         string                `json:"content,omitempty"`
	TransactionType enums.TransactionType `json:"transactionType,omitempty"`
	Title           string                `json:"title"`
	NumDate         int                   `json:"numDate,omitempty"`
	Type            uint                  `json:"type,omitempty"`
	Price           float64               `json:"price,omitempty"`
	PriceType       uint32                `json:"priceType,omitempty"`
	PackageVisible  uint                  `json:"packageVisible,omitempty"`
	Like            uint64                `json:"numLike,omitempty"`
	Comment         uint64                `json:"numComment,omitempty"`
	NumView         uint64                `json:"numView,omitempty"`
	OwnerID         uint64                `json:"ownerId,omitempty"`
	OwnerType       enums.EOwnerOf        `json:"ownerType"`
	PublishedAt     *time.Time            `json:"publishedAt,omitempty"`
	PublishedBy     uint64                `json:"publishedBy,omitempty"`
	MediaList       []PostMediaDTO        `json:"mediaList,omitempty"`
	Product         *ProductBasicDTO      `json:"product,omitempty"`
	CreatedAt       *time.Time            `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time            `json:"updatedAt,omitempty"`
	Permissions     []string              `json:"permissions,omitempty"`
}

// ProductBasicDTO DTO cho thông tin cơ bản của sản phẩm
type ProductBasicDTO struct {
	ID           uint64  `json:"id"`
	Name         string  `json:"name,omitempty"`
	Code         string  `json:"code,omitempty"`
	Area         float64 `json:"area,omitempty"`
	Address      string  `json:"address,omitempty"`
	ProvinceName string  `json:"provinceName,omitempty"`
	DistrictName string  `json:"districtName,omitempty"`
	WardName     string  `json:"wardName,omitempty"`
}
