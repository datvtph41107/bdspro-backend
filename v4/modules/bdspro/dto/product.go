package dto

// SyncListRequest is the common page/timestamp query used by product/me and property/me.
type SyncListRequest struct {
	ID        uint64 `form:"id" uri:"id"`
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Text      string `form:"text"`
	Timestamp int64  `form:"timestamp"`
	Archived  *bool  `form:"archived"`
}
