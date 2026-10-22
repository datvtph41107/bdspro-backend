package services

import (
	_crud "common/domain/crud"
	"mdn4/internal/interface/repo"
	"mdn4/modules/user/domain/user"
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
