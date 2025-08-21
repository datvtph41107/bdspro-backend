package usecase

import (
	"auth/internal/domain"
	"auth/internal/interface/repo"
	"common/case/crud"
)

type RoleGroupUsecase struct {
	crud.BaseUsecase[domain.RoleGroup, repo.RoleGroupRepository]
}

func NewRoleGroupUsecase(roleGroupRepo repo.RoleGroupRepository) *RoleGroupUsecase {
	return &RoleGroupUsecase{
		crud.BaseUsecase[domain.RoleGroup, repo.RoleGroupRepository]{
			Repo: roleGroupRepo,
		},
	}
}
