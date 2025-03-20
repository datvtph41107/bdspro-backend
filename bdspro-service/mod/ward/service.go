package ward

import (
	"bdspro/models"
	"common/case/crud"
)

type WardService struct {
	crud.BaseService[models.Ward]
}

func NewWardService(repo *WardRepo) *WardService {
	return &WardService{
		BaseService: crud.BaseService[models.Ward]{Repo: repo.BaseRepo},
	}
}
