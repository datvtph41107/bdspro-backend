package entity

import (
	_models "common/models"
	"organization/internal/enums"
)

type CostType struct {
	_models.BaseEntity
	TypeName    string
	IsCustom    bool
	CostType    enums.ECostType
	Description string
}
