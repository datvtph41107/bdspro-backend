package dto

import (
	_dto "common/dto"
)

type StageSearchDTO struct {
	_dto.Page
	StageName string `form:"stageName"`
}
