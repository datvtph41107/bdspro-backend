package dto

import _dto "common/dto"

type DocTypeSearchDTO struct {
	_dto.Pagable
	Text string `form:"text"`
}
