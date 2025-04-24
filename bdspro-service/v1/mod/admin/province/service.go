package province

import (
	"bdspro/v1/models"
	"common/case/crud"
)

type ProvinceService struct {
	crud.BaseService[models.Province]
}

func NewProvinceService(repo *ProvinceRepo) *ProvinceService {
	return &ProvinceService{
		BaseService: crud.BaseService[models.Province]{Repo: repo.BaseRepo},
	}
}
