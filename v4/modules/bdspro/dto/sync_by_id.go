package dto

// SyncByIDRequest is the common page/size/timestamp query for product-tab lists.
type SyncByIDRequest struct {
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Timestamp int64  `form:"timestamp"`
	Text      string `form:"text"`
	Keyword   string `form:"keyword"`
	Action    uint32 `form:"action"`
	ActorID   uint64 `form:"actorId"`
}

// SyncByIDPathRequest binds a path id (or productId) plus SyncByID query fields.
type SyncByIDPathRequest struct {
	ID        uint64 `path:"id|productId"`
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Timestamp int64  `form:"timestamp"`
	Text      string `form:"text"`
	Keyword   string `form:"keyword"`
	Action    uint32 `form:"action"`
	ActorID   uint64 `form:"actorId"`
}

// ToSyncByIDRequest strips path-bound id for service calls that take id separately.
func (r SyncByIDPathRequest) ToSyncByIDRequest() SyncByIDRequest {
	return SyncByIDRequest{
		Page:      r.Page,
		Size:      r.Size,
		Timestamp: r.Timestamp,
		Text:      r.Text,
		Keyword:   r.Keyword,
		Action:    r.Action,
		ActorID:   r.ActorID,
	}
}
