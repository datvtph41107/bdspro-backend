package dto

import (
	_dto "common/dto"
)

type ContactSearchDTO struct {
	_dto.Pagable
	Text string `form:"text"`
}
