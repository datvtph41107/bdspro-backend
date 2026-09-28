package usecases

import (
	_crud "common/domain/crud"
	"userv4/internal/domain/user"
	"userv4/internal/interface/repo"
)

type IAdminPurposeUseUsecase interface {
	_crud.IBaseUsecase[user.PurposeUseEntity]
}

type AdminPurposeUseUsecase struct {
	_crud.BaseUsecase[user.PurposeUseEntity, repo.IPurposeUseRepo]
}

func NewAdminPurposeUseUsecase(purposeRepo repo.IPurposeUseRepo) IAdminPurposeUseUsecase {
	return &AdminPurposeUseUsecase{
		BaseUsecase: _crud.BaseUsecase[user.PurposeUseEntity, repo.IPurposeUseRepo]{
			Repo: purposeRepo,
		},
	}
}
