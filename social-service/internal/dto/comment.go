package dto

import _dto "common/domain/dto"

type CommentRequest struct {
	_dto.Pagable
	NewsFeedID uint64  `json:"news_feed_id"`
	ParentID   *uint64 `json:"parent_id"`
}
