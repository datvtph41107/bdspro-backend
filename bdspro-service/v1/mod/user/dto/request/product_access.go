package user_dto

import (
	"bdspro/internal/enums"
	_dto "common/dto"
)

type ProductAccessSearch struct {
	_dto.Page
	ProductID  uint64            `form:"productId"`
	TargetType enums.ETargetType `form:"targetType"`
}
