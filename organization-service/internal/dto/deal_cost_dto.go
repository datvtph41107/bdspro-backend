package dto

import _dto "common/domain/dto"

type DealCostSearchDTO struct {
	_dto.Pagable
	DealId uint64
}
