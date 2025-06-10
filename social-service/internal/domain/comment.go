package domain

import (
	_models "common/models"
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
