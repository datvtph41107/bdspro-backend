package dto

import (
	_dto "common/dto"
)

type ContactSearchDTO struct {
	_dto.Page
	Text string `form:"text"`
}
