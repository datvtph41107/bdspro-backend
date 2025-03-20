package district

import (
	"bdspro/models"
	"common/case/crud"

	"gorm.io/gorm"
)

type DistrictRepo struct {
	crud.BaseRepo[models.District]
}

func NewDistrictRepo(db *gorm.DB) *DistrictRepo {
	return &DistrictRepo{
		BaseRepo: crud.BaseRepo[models.District]{DB: db},
	}
}
