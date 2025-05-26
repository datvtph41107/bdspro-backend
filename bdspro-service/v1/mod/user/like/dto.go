package u_like

import (
	"bdspro/v1/models"
	_utils "common/utils"
	"time"
)

type LikeSaveRequest struct {
	ProductID       uint64                  `json:"productId" binding:"required"`
	Content         string                  `json:"content" binding:"required"`
	TransactionType uint                    `json:"transactionType" binding:"required"`
	Title           string                  `json:"title"`
	ExpiredAt       *time.Time              `json:"expiredAt"`
	Metadatas       models.ProductMediaItem `json:"metadatas"`
	// Status    uint   `json:"status"`
}

type UpdateHiddenRequest struct {
	LikeID uint64 `json:"productId" binding:"required"`
	Hidden bool   `json:"hidden" binding:"required"`
}

type ActionRequest struct {
	LikeID uint64 `json:"LikeId" binding:"required"`
	Action uint   `json:"action" binding:"required"`
}

type LikeUpdateRequest struct {
	Content         string                  `json:"content"`
	TransactionType uint                    `json:"transactionType"`
	Title           string                  `json:"title"`
	Metadatas       models.ProductMediaItem `json:"metadatas"`
	// Status    uint   `json:"status"`
}

type LikeShareRequest struct {
	LikeID  uint64 `json:"LikeId"`
	Content string `json:"content"`
	Title   string `json:"title"`
}

type LikeExpiredRequest struct {
	LikeID uint64 `json:"LikeId" binding:"required"`
	NumDay uint64 `json:"numDay" binding:"required"`
}

// type LikeInfo struct {
// }
type LikeDetail struct {
	Like        *models.Like    `json:"Like"`
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
