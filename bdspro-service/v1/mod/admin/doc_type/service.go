package a_doc_type

import (
	"bdspro/v1/models"
	"common/case/crud"
)

type DocTypeService struct {
	crud.BaseService[models.DocType]
}

func NewDocTypeService(repo *DocTypeRepo) *DocTypeService {
	return &DocTypeService{
		BaseService: crud.BaseService[models.DocType]{Repo: repo.BaseRepo},
	}
}
