package dto

import _dto "common/domain/dto"

type DocTypeSearchDTO struct {
	_dto.Pagable
	Text string `form:"text"`
}
