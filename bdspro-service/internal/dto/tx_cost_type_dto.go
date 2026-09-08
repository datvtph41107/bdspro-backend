package dto

import _dto "common/domain/dto"

type TxCostTypeSearchDTO struct {
	_dto.Pagable
	TypeName string
	CostType uint32
}
