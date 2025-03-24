package ward

import (
	"bdspro/models"
	"common/case/crud"

	"gorm.io/gorm"
)

type WardRepo struct {
	crud.BaseRepo[models.Ward]
}

func NewWardRepo(db *gorm.DB) *WardRepo {
	return &WardRepo{
		BaseRepo: crud.BaseRepo[models.Ward]{DB: db},
	}
}
