package a_doc_type

import (
	"bdspro/models"
	"common/case/crud"

	"gorm.io/gorm"
)

type DocTypeRepo struct {
	crud.BaseRepo[models.DocType]
}

func NewDocTypeRepo(db *gorm.DB) *DocTypeRepo {
	return &DocTypeRepo{
		BaseRepo: crud.BaseRepo[models.DocType]{DB: db},
	}
}
