package rest

import _models "common/models"

type IRestModel interface {
	TableName() string
}

type RestModel[M IRestModel] struct {
	_models.BaseEntity
	M M `gorm:"embedded" json:"data"`
}

func (RestModel[M]) TableName() string {
	var t M
	return t.TableName()
}
