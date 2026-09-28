package dto

// SyncListRequest is the common page/timestamp query used by product/me and services.
type SyncListRequest struct {
	ID        uint64 `form:"id"`
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Text      string `form:"text"`
	Timestamp int64  `form:"timestamp"`
	Archived  *bool  `form:"archived"`
}

// SyncListByIDRequest binds path id plus sync list query (product detail, etc.).
type SyncListByIDRequest struct {
	ID        uint64 `path:"id"`
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Text      string `form:"text"`
	Timestamp int64  `form:"timestamp"`
	Archived  *bool  `form:"archived"`
}

// ToSyncListRequest maps bound path+query into the service DTO (includes ID).
func (r SyncListByIDRequest) ToSyncListRequest() SyncListRequest {
	return SyncListRequest{
		ID:        r.ID,
		Page:      r.Page,
		Size:      r.Size,
		Text:      r.Text,
		Timestamp: r.Timestamp,
		Archived:  r.Archived,
	}
}
