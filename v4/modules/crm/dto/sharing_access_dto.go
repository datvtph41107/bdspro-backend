package dto

import (
	_dto "common/domain/dto"
)

type SharingAccessSearchDTO struct {
	_dto.Pagable
	ContactId uint64 `form:"contactId"`
	Text      string `form:"text"`
}