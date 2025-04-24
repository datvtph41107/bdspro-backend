package a_project

import (
	"bdspro/v1/models"
	"common/case/crud2"
)

type ProjectService struct {
	crud2.BaseService[models.Project, *ProjectRepo]
}

func NewProjectService(repo *ProjectRepo) *ProjectService {
	return &ProjectService{
		BaseService: crud2.BaseService[models.Project, *ProjectRepo]{
			Repo: repo,
		},
	}
}
