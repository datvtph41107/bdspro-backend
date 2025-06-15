package domain

import (
	_models "common/models"
	"time"
)

type Comment struct {
	_models.BaseEntity
	NewsFeedID uint64  `json:"news_feed_id"`
	UserID     uint64  `json:"user_id"`
	Content    string  `json:"content"`
	ParentID   *uint64 `json:"parent_id"`
	NumLike    uint32  `json:"num_like"`
	NumDisLike uint32  `json:"num_dislike"`
	NumReply   uint32  `json:"num_reply"`
	MediaUrl   string  `json:"media_url"`
	MediaType  string  `json:"media_type"`
}

func (Comment) TableName() string {
	return "comment"
}

type CommentQuery struct {
	_models.BaseEntity
	NewsFeedID uint64     `json:"news_feed_id"`
	UserID     uint64     `json:"user_id"`
	Content    string     `json:"content"`
	ParentID   *uint64    `json:"parent_id"`
	NumLike    uint32     `json:"num_like"`
	NumDisLike uint32     `json:"num_dislike"`
	NumReply   uint32     `json:"num_reply"`
	MediaUrl   string     `json:"media_url"`
	MediaType  string     `json:"media_type"`
	LikeAt     *time.Time `gorm:"column:like_at" json:"like_at"`
	DisLike    *bool      `gorm:"column:dis_like" json:"dis_like"`
}

func (CommentQuery) TableName() string {
	return "comment"
}
