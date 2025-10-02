package dto

import _dto "common/domain/dto"

type AmenitySearchDTO struct {
	_dto.Pagable
	Text string `form:"text"`
}
