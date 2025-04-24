package a_developer

import (
	"bdspro/v1/models"
	"common/case/crud2"
)

type DeveloperService struct {
	crud2.BaseService[models.Developer, *DeveloperRepo]
}

func NewDeveloperService(repo *DeveloperRepo) *DeveloperService {
	return &DeveloperService{
		BaseService: crud2.BaseService[models.Developer, *DeveloperRepo]{Repo: repo},
	}
}
