package dto

import _dto "common/dto"

type DealCostSearchDTO struct {
	_dto.Pagable
	DealId uint64
}
