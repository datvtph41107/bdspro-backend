package dto

import _dto "common/dto"

type CostTypeSearchDTO struct {
	_dto.Pagable
	TypeName string
	CostType uint32
}
