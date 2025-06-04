package dto

import _dto "common/dto"

type ProjectSearchDTO struct {
	_dto.Page
	Text string `form:"text"`
}
