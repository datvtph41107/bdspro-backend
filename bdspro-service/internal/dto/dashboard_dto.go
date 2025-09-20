package dto

import sharepb "pb/types/shared"

type CountByOwnerResponse struct {
	TotalProduct  uint32 `json:"totalProduct"`
	TotalAsset    uint32 `json:"totalAsset"`
	TotalNewsfeed uint32 `json:"totalNewsfeed"`
	TotalProject  uint32 `json:"totalProject"`
}

type CountByOwnerRequest struct {
	OwnerOf sharepb.OwnerOf `json:"ownerOf"`
	OwnerId uint64          `json:"ownerId"`
}
