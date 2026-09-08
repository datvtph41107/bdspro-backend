package qh_domain

import _models "common/domain/entity"

type QHDirection struct {
	_models.BaseEntity
	Name string
}

func (QHDirection) TableName() string {
	return "qh_direction"
}
