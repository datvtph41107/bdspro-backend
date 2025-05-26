package _owner

import (
	_models "common/models"
)

type IOwnerModel interface {
	TableName() string
}

type OwnerModel[M IOwnerModel] struct {
	_models.BaseEntity
	M       M      `gorm:"embedded" json:"data"`
	OwnerID uint64 `json:"ownerId"`
}

func (OwnerModel[M]) TableName() string {
	var t M
	return t.TableName()
}
