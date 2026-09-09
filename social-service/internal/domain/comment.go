package domain

import (
	_models "common/models"
	"errors"
	"time"
)

var ErrCommentNewsFeedUnavailable = errors.New("bài viết không tồn tại hoặc bị giới hạn bình luận")

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
	LikedAt    *time.Time `gorm:"column:liked_at" json:"liked_at"`
	DisLike    *bool      `gorm:"column:dis_like" json:"dis_like"`
}

func (CommentQuery) TableName() string {
	return "comment"
}
