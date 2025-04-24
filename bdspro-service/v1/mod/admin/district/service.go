package district

import (
	"bdspro/v1/models"
	"common/case/crud"
)

type DistrictService struct {
	crud.BaseService[models.District]
}

func NewDistrictService(repo *DistrictRepo) *DistrictService {
	return &DistrictService{
		BaseService: crud.BaseService[models.District]{Repo: repo.BaseRepo},
	}
}
