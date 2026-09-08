package qh_domain

import (
	_models "common/models"
)

type QHJurisdiction struct {
	_models.BaseEntity
	Name string `gorm:"type:varchar(255);not null" json:"name"`
}

func (QHJurisdiction) TableName() string {
	return "qh_jurisdictions"
}
