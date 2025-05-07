package u_post

import (
	"bdspro/v1/models"
	_utils "common/utils"
	"time"
)

type PostSaveRequest struct {
	ProductID       uint64                  `json:"productId" binding:"required"`
	Content         string                  `json:"content" binding:"required"`
	TransactionType uint                    `json:"transactionType" binding:"required"`
	Title           string                  `json:"title"`
	ExpiredAt       *time.Time              `json:"expiredAt"`
	Metadatas       models.ProductMediaItem `json:"metadatas"`
	// Status    uint   `json:"status"`
}

type UpdateHiddenRequest struct {
	PostID uint64 `json:"productId" binding:"required"`
	Hidden bool   `json:"hidden" binding:"required"`
}

type ActionRequest struct {
	PostID uint64 `json:"postId" binding:"required"`
	Action uint   `json:"action" binding:"required"`
}

type PostUpdateRequest struct {
	Content         string                  `json:"content"`
	TransactionType uint                    `json:"transactionType"`
	Title           string                  `json:"title"`
	Metadatas       models.ProductMediaItem `json:"metadatas"`
	// Status    uint   `json:"status"`
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
	_utils.Page
	Content         string     `form:"content"`
	TransactionType uint       `form:"transactionType"`
	Title           string     `form:"title"`
	ExpiredFrom     *time.Time `form:"expiredFrom"`
	ExpiredTo       *time.Time `form:"expiredTo"`
	Hidden          *bool      `form:"hidden"`
}
