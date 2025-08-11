package _category

import (
	_models "common/models"
)

type ICategoryModel interface {
	TableName() string
}

type CategoryModel[M ICategoryModel] struct {
	_models.BaseEntity
	M        M    `gorm:"embedded" json:"data"`
	IsCustom bool `json:"isCustom"`
}

func (CategoryModel[M]) TableName() string {
	var t M
	return t.TableName()
}
