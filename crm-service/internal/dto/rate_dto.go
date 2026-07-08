package dto

import (
	_dto "common/domain/dto"
	"time"
)

type RateSearchDTO struct {
	_dto.Pagable
	OwnerId   uint64    `json:"ownerId"`
	ParentId  *uint64   `json:"parentId"`
	CreatedBy uint64    `json:"createdBy"`
	UpdatedBy uint64    `json:"updatedBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Registry  string    `json:"registry,omitempty"`
}

type RateCreateRequest struct {
	Anonymous bool     `json:"anonymous" binding:"required"`
	OwnerID   uint64   `json:"ownerId" binding:"required"`
	OwnerOf   uint32   `json:"ownerOf" binding:"required"`
	Score     uint8    `json:"score" binding:"required,gte=1,lte=5"`
	Comment   string   `json:"comment" binding:"required"`
	ParentID  *uint64  `json:"parentId,omitempty"`
	Attachs   []string `json:"attachs,omitempty"`
}

type RateUpdateRequest struct {
	Score   uint8    `json:"score" binding:"required,gte=1,lte=5"`
	Comment string   `json:"comment" binding:"required"`
	Attachs []string `json:"attachs,omitempty"`
}

type RateResponse struct {
	ID          uint64         `json:"id"`
	Anonymous   bool           `json:"anonymous"`
	OwnerID     uint64         `json:"ownerId"`
	OwnerOf     uint32         `json:"ownerOf"`
	Score       uint8          `json:"score"`
	Comment     string         `json:"comment"`
	ParentID    *uint64        `json:"parentId,omitempty"`
	CreatedBy   *uint64        `json:"createdBy,omitempty"`
	CreatedUser *UserItem      `json:"createdUser,omitempty"`
	CreatedAt   *time.Time     `json:"createdAt"`
	Attachs     []AttachEntity `json:"attachs,omitempty"`
}

type UserItem struct {
	ProfileID uint64 `json:"profileId"`
	FullName  string `json:"fullName"`
	Avatar    string `json:"avatar"`
}

type AttachEntity struct {
	ID       uint64 `json:"id"`
	FileURL  string `json:"fileUrl"`
	FileName string `json:"fileName"`
	FileType string `json:"fileType"`
}

type RateStatsResponse struct {
	TotalReviews int64   `json:"totalReviews"`
	AverageScore float64 `json:"averageScore"`
	Star1Count   int64   `json:"star1Count"`
	Star2Count   int64   `json:"star2Count"`
	Star3Count   int64   `json:"star3Count"`
	Star4Count   int64   `json:"star4Count"`
	Star5Count   int64   `json:"star5Count"`
}

type GetRateRequest struct {
	ID       uint64 `json:"id" binding:"required"`
	Registry string `json:"registry,omitempty"`
}

type UpdateHiddenRequest struct {
	ID       uint64 `json:"id" binding:"required"`
	Hidden   bool   `json:"hidden" binding:"required"`
	Registry string `json:"registry,omitempty"`
}