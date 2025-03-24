package a_property_type

import (
	"bdspro/models"
	"common/case/crud"

	"gorm.io/gorm"
)

type PropertyTypeRepo struct {
	crud.BaseRepo[models.PropertyType]
}

func NewPropertyTypeRepo(db *gorm.DB) *PropertyTypeRepo {
	return &PropertyTypeRepo{
		BaseRepo: crud.BaseRepo[models.PropertyType]{DB: db},
	}
}
