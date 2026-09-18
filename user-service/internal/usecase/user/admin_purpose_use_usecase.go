package usecases

import (
	_crud "common/domain/crud"
	"user/internal/domain/user"
	"user/internal/interface/repo"
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
