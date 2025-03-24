package province

import (
	"bdspro/models"
	"common/case/crud"

	"gorm.io/gorm"
)

type ProvinceRepo struct {
	crud.BaseRepo[models.Province]
}

func NewProvinceRepo(db *gorm.DB) *ProvinceRepo {
	return &ProvinceRepo{
		BaseRepo: crud.BaseRepo[models.Province]{DB: db},
	}
}
