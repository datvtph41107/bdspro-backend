package dto

// BookmarkUsersListRequest is GET /v2/user/bookmark-users.
type BookmarkUsersListRequest struct {
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
}
