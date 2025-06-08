package dto

import _dto "common/dto"

type PostSearch struct {
	_dto.Page
	Text string `json:"text" form:"text"`
}

type NewsFeedGlobalSearch struct {
	_dto.Page
	Text string `json:"text" form:"text"`
}
