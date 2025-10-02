package dto

import _dto "common/domain/dto"

type ProjectSearchDTO struct {
	_dto.Pagable
	Text string `form:"text"`
}
