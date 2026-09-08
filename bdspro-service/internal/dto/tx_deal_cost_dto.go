package dto

import _dto "common/domain/dto"

type TxDealCostSearchDTO struct {
	_dto.Pagable
	DealId uint64
}
