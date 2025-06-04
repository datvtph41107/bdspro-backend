package dto

import _dto "common/dto"

type DocTypeSearchDTO struct {
	_dto.Page
	Text string `form:"text"`
}
