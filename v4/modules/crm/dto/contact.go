package dto

// ContactSearchRequest mirrors ContactSearchRequest query params for list/user.
type ContactSearchRequest struct {
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Text      string `form:"text"`
	Timestamp int64  `form:"timestamp"`
}

// ContactListUserRequest is GET /v2/crm/contact/list/user/{ownerId}.
type ContactListUserRequest struct {
	// ownerId 0 means current user (legacy FE contract for list/me alias).
	OwnerID   uint64 `path:"ownerId" pathAllowZero:"true"`
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Text      string `form:"text"`
	Timestamp int64  `form:"timestamp"`
}
