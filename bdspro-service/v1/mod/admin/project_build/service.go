package a_project_build

import (
	"bdspro/v1/models"
	"common/case/crud2"
)

type ProjectBuildService struct {
	crud2.BaseService[models.ProjectBuild, *ProjectBuildRepo]
}

func NewProjectBuildService(repo *ProjectBuildRepo) *ProjectBuildService {
	return &ProjectBuildService{
		BaseService: crud2.BaseService[models.ProjectBuild, *ProjectBuildRepo]{Repo: repo},
	}
}
