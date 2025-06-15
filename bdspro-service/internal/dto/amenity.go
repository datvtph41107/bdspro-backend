package dto

import _dto "common/dto"

type AmenitySearchDTO struct {
	_dto.Pagable
	Text string `form:"text"`
}
