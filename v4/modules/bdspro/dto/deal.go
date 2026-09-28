package dto

// DealSearchRequest mirrors bdspropb.DealSearchRequest query params.
type DealSearchRequest struct {
	Page           uint32  `form:"page"`
	Size           uint32  `form:"size"`
	Text           string  `form:"text"`
	OrganizationID *uint64 `form:"organizationId"`
	GroupID        *uint64 `form:"groupId"`
	OwnerID        *uint64 `form:"ownerId"`
	OwnerOf        uint32  `form:"ownerOf"`
	Timestamp      int64   `form:"timestamp"`
}
