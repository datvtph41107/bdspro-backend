package dto

import _dto "common/dto"

type CommentRequest struct {
	_dto.Page
	NewsFeedID uint64  `json:"news_feed_id"`
	ParentID   *uint64 `json:"parent_id"`
}
