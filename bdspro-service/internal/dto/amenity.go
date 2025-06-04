package dto

import _dto "common/dto"

type AmenitySearchDTO struct {
	_dto.Page
	Text string `form:"text"`
}
