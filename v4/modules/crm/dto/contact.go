package dto

// ContactSearchRequest mirrors ContactSearchRequest query params for list/user.
type ContactSearchRequest struct {
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Text      string `form:"text"`
	Timestamp int64  `form:"timestamp"`
}
