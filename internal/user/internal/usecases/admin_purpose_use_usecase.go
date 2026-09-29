package usecases

import (
	_crud "common/domain/crud"
	"user/internal/interface/repo"
	"user/internal/models"
)

type IAdminPurposeUseUsecase interface {
	_crud.IBaseUsecase[models.PurposeUseEntity]
}

type AdminPurposeUseUsecase struct {
	_crud.BaseUsecase[models.PurposeUseEntity, repo.IPurposeUseRepo]
}

func NewAdminPurposeUseUsecase(purposeRepo repo.IPurposeUseRepo) IAdminPurposeUseUsecase {
	return &AdminPurposeUseUsecase{
		BaseUsecase: _crud.BaseUsecase[models.PurposeUseEntity, repo.IPurposeUseRepo]{
			Repo: purposeRepo,
		},
	}
}
