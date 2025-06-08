package domain

import (
	_models "common/models"
	"social/internal/enums"
)

type NewsFeed struct {
	_models.BaseEntity
	Title         string           `json:"title"`
	Content       string           `json:"content"`
	Image         string           `json:"image"`
	Link          string           `json:"link"`
	Visibility    enums.Visibility `gorm:"default:10" json:"visibility"`
	PostID        *uint64          `json:"post_id"`
	LikeNumber    int              `json:"like_number"`
	CommentNumber int              `json:"comment_number"`
	ShareNumber   int              `json:"share_number"`
	ViewNumber    int              `json:"view_number"`
}
