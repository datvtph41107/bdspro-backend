package dto

type NewsFeedOfUserDTO struct {
	PostID uint64 `json:"postId" binding:"required"`
	UserID uint64 `json:"userId" binding:"required"`
}

type NewsFeedOfGroupDTO struct {
	PostID  uint64 `json:"postId" binding:"required"`
	GroupID uint64 `json:"groupId" binding:"required"`
}
