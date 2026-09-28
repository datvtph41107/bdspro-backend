package dto

// PropertyListRequest is the query for GET /v2/bdspro/v2/property/me.
type PropertyListRequest struct {
	ID        uint64 `form:"id"`
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Timestamp int64  `form:"timestamp"`
}

// PropertyDetailRequest is GET /v2/bdspro/v2/property/{id}/detail.
type PropertyDetailRequest struct {
	ID        uint64 `path:"id"`
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Timestamp int64  `form:"timestamp"`
}

// ToListRequest maps bound path+query into the service DTO.
func (r PropertyDetailRequest) ToListRequest() PropertyListRequest {
	return PropertyListRequest{
		ID:        r.ID,
		Page:      r.Page,
		Size:      r.Size,
		Timestamp: r.Timestamp,
	}
}
