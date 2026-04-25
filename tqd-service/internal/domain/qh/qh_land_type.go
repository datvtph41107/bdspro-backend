package qh_domain

import _models "common/domain/entity"

type QHLandType struct {
	_models.BaseEntity
	Name string
	Code string `column:"varchar(10)"`
}

func (QHLandType) TableName() string {
	return "qh_land_type"
}
