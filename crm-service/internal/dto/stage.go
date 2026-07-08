package dto

import (
	_dto "common/domain/dto"
)

type StageSearchDTO struct {
	_dto.Pagable
	StageName string `form:"stageName"`
}