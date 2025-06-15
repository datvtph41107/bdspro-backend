package dto

import (
	_dto "common/dto"
)

type StageSearchDTO struct {
	_dto.Pagable
	StageName string `form:"stageName"`
}
