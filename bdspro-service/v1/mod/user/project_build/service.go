package u_project_build

import (
	"bdspro/v1/models"
	"common/case/crud2"
)

type UserProjectBuildService struct {
	crud2.BaseService[models.ProjectBuild, *UserProjectBuildRepo]
}

func NewUserProjectBuildService(repo *UserProjectBuildRepo) *UserProjectBuildService {
	return &UserProjectBuildService{
		BaseService: crud2.BaseService[models.ProjectBuild, *UserProjectBuildRepo]{Repo: repo},
	}
}
