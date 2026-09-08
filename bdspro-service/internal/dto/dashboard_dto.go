package dto

import (
	sharepb "pb/types/shared"
	"time"
)

type CountByOwnerResponse struct {
	TotalProduct  uint32 `json:"totalProduct"`
	TotalAsset    uint32 `json:"totalAsset"`
	TotalPost     uint32 `json:"totalPost"`
	TotalNewsfeed uint32 `json:"totalNewsfeed"`
	TotalProject  uint32 `json:"totalProject"`
}

type CountByOwnerRequest struct {
	OwnerOf sharepb.OwnerOf `json:"ownerOf"`
	OwnerId uint64          `json:"ownerId"`
}

type CountPostByTimeResponse struct {
	Date  time.Time `json:"date"`
	Count uint32    `json:"count"`
}

type CountOfUserResponse struct {
	TotalPosted            uint32 `json:"totalPosted"`
	TotalInterestedProduct uint32 `json:"totalInterestedProduct"`
	TotalSavedProduct      uint32 `json:"totalSavedProduct"`
}
