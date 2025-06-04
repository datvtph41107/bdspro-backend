package dto

import _dto "common/dto"

type PropertyTypeSearchDTO struct {
	_dto.Page
	Text string `form:"text"`
}
