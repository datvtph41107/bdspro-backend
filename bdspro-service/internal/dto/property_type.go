package dto

import _dto "common/dto"

type PropertyTypeSearchDTO struct {
	_dto.Pagable
	Text string `form:"text"`
}
