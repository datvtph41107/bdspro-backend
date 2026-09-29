package dto

import _dto "common/domain/dto"

type PropertyTypeSearchDTO struct {
	_dto.Pagable
	Text string `form:"text"`
}
