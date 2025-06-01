package u_post_dto

import (
	"bdspro/internal/enums"
	"bdspro/v1/models"
	"time"
)

type PostSaveRequest struct {
	ProductID       uint64                  `json:"productId" binding:"required"`
	Title           string                  `json:"title"`
	Content         string                  `json:"content" binding:"required"`
	Price           float64                 `json:"postPrice"`
	TransactionType uint                    `json:"transactionType" binding:"required"`
	Visibility      enums.EVisibility       `json:"visibility" binding:"required"`
	PackageVisible  uint                    `json:"packageVisible"`
	ExpiredAt       *time.Time              `json:"expiredAt"`
	NumDate         int                     `json:"numDate"`
	Hidden          bool                    `json:"hidden"`
	Metadatas       models.ProductMediaItem `json:"metadatas"`
	// Status    uint   `json:"status"`
}

type PostSaveWithProduct struct {
	ProductID       uint64                  `json:"productId"`
	Title           string                  `json:"title"`
	Content         string                  `json:"content"`
	Price           float64                 `json:"postPrice"`
	TransactionType uint                    `json:"transactionType"`
	Visibility      enums.EVisibility       `json:"visibility"`
	PackageVisible  uint                    `json:"packageVisible"`
	ExpiredAt       *time.Time              `json:"expiredAt"`
	NumDate         int                     `json:"numDate"`
	Hidden          bool                    `json:"hidden"`
	Metadatas       models.ProductMediaItem `json:"metadatas"`
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
	Title           string                  `json:"title"`
	Content         string                  `json:"content"`
	Price           float64                 `json:"postPrice"`
	TransactionType uint                    `json:"transactionType"`
	PackageVisible  uint                    `json:"packageVisible"`
	Visibility      enums.EVisibility       `json:"visibility"`
	NumDate         int                     `json:"numDate"`
	Metadatas       models.ProductMediaItem `json:"metadatas"`
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

// type PostInfo struct {
// }
type PostDetail struct {
	Post        *models.Post    `json:"post"`
	ProductInfo *models.Product `json:"product"`
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
