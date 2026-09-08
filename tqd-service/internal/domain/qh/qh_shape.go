package qh_domain

import _models "common/domain/entity"

type QHShape struct {
	_models.BaseEntity
	Name string
}

func (QHShape) TableName() string {
	return "qh_shape"
}
