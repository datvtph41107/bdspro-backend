package dto

import _dto "common/domain/dto"

type PostMediaItem struct {
	ID        uint64 `json:"id"`
	MediaURL  string `json:"mediaUrl"`
	MediaType string `json:"mediaType"`
	IsMain    bool   `json:"isMain"`
	Order     int    `json:"order"`
}

type PostMediaGetDTO struct {
	_dto.Pagable
	PostID uint64 `form:"postId"`
}
