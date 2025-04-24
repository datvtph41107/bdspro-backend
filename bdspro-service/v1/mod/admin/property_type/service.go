package a_property_type

import (
	"bdspro/v1/models"
	"common/case/crud"
)

type PropertyTypeService struct {
	crud.BaseService[models.PropertyType]
}

func NewPropertyTypeService(repo *PropertyTypeRepo) *PropertyTypeService {
	return &PropertyTypeService{
		BaseService: crud.BaseService[models.PropertyType]{Repo: repo.BaseRepo},
	}
}
