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
