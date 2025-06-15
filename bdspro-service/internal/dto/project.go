package dto

import _dto "common/dto"

type ProjectSearchDTO struct {
	_dto.Pagable
	Text string `form:"text"`
}
