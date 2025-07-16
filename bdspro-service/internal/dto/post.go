package dto

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_dto "common/dto"
	"time"
)

type PostLinkSearch struct {
	_dto.Pagable
	ProductID *uint64 `json:"productId"`
}

type PostSearchRequest struct {
	_dto.Pagable
	Content         string     `form:"content"`
	TransactionType uint32     `form:"transactionType"`
	Text            string     `form:"text"`
	Title           string     `form:"title"`
	ExpiredFrom     *time.Time `form:"expiredFrom"`
	ExpiredTo       *time.Time `form:"expiredTo"`
	Hidden          *bool      `form:"hidden"`
	Status          []uint64   `form:"status"`
	Types           []uint64   `form:"types"` // 1: hiển thị, 2: hết hạn; 3: tạm ẩn
	PackageVisibles []uint64   `form:"packageVisibles"`
	ProductID       *uint64    `form:"productId"`
}

type PostDetail struct {
	Post        *domain.PostItem `json:"post"`
	ProductInfo *domain.Product  `json:"product"`
	AssetInfo   *domain.Asset    `json:"asset"`
}

type PostSaveRequest struct {
	ProductID       uint64                `json:"productId" binding:"required"`
	Title           string                `json:"title"`
	Content         string                `json:"content" binding:"required"`
	Price           float64               `json:"postPrice"`
	TransactionType enums.TransactionType `json:"transactionType" binding:"required"`
	PriceType       uint32                `json:"priceType"`
	Visibility      enums.EVisibility     `json:"visibility" binding:"required"`
	PackageVisible  uint                  `json:"packageVisible"`
	ExpiredAt       *time.Time            `json:"expiredAt"`
	NumDate         int                   `json:"numDate"`
	Hidden          bool                  `json:"hidden"`
	Medatadatas     []PostMediaItem       `json:"mediaList"`
	// Status    uint   `json:"status"`
}

type PostSaveWithProduct struct {
	ProductID       uint64                    `json:"productId"`
	Title           string                    `json:"title"`
	Content         string                    `json:"content"`
	Price           float64                   `json:"postPrice"`
	TransactionType uint                      `json:"transactionType"`
	Visibility      enums.EVisibility         `json:"visibility"`
	PackageVisible  uint                      `json:"packageVisible"`
	ExpiredAt       *time.Time                `json:"expiredAt"`
	NumDate         int                       `json:"numDate"`
	Hidden          bool                      `json:"hidden"`
	Metadatas       []domain.ProductMediaItem `json:"metadatas"`
	// Status    uint   `json:"status"`
}

type UpdateHiddenRequest struct {
	PostID uint64 `json:"postId" binding:"required"`
	Hidden bool   `json:"hidden"`
}

type ActionRequest struct {
	PostID uint64 `json:"postId" binding:"required"`
	Action uint   `json:"action" binding:"required"`
}

type PostUpdateRequest struct {
	Title           string            `json:"title"`
	Content         string            `json:"content"`
	Price           float64           `json:"postPrice"`
	TransactionType uint32            `json:"transactionType"`
	PackageVisible  uint              `json:"packageVisible"`
	Visibility      enums.EVisibility `json:"visibility"`
	PriceType       uint32            `json:"priceType"`
	NumDate         int               `json:"numDate"`
	MediaList       []PostMediaItem   `json:"mediaList"`
}

type PostShareRequest struct {
	PostID  uint64 `json:"postId"`
	Content string `json:"content"`
	Title   string `json:"title"`
}

type PostExpiredRequest struct {
	PostID uint64 `json:"postId" binding:"required"`
	NumDay uint64 `json:"numDay" binding:"required"`
}

type SearchRequest struct {
	Size            int        `form:"size,default=20"`
	Page            int        `form:"page,default=0"`
	Content         string     `form:"content"`
	TransactionType uint       `form:"transactionType"`
	Text            string     `form:"text"`
	Title           string     `form:"title"`
	ExpiredFrom     *time.Time `form:"expiredFrom"`
	ExpiredTo       *time.Time `form:"expiredTo"`
	Hidden          *bool      `form:"hidden"`
	Status          []uint64   `form:"status"`
	Types           []uint64   `form:"types"` // 1: hiển thị, 2: hết hạn; 3: tạm ẩn
	PackageVisibles []uint64   `form:"packageVisibles"`
}

type PostPublishSearch struct {
	_dto.Pagable
}
