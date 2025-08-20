package admin_usecases

import (
	"bdspro/internal/domain"
	admin_repo "bdspro/internal/repo/admin"
	"common/case/crud"
)

type AdminProjectUsecase struct {
	crud.BaseService[domain.Project, admin_repo.AdminProjectRepo]
}

func NewAdminProjectUsecase(repo admin_repo.AdminProjectRepo) *AdminProjectUsecase {
	return &AdminProjectUsecase{
		BaseService: crud.BaseService[domain.Project, admin_repo.AdminProjectRepo]{
			Repo: repo,
		},
	}
}
