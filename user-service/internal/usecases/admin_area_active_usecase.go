package usecases

import (
	_crud "common/domain/crud"
	_dto "common/domain/dto"
	"context"
	"user/internal/interface/repo"
	"user/internal/models"
)

type IAdminMainAreaUsecase interface {
	_crud.IBaseUsecase[models.MainAreaEntity]
	Search(ctx context.Context, text string, pagable _dto.IPagable) ([]models.MainAreaEntity, uint32, error)
}

type AdminMainAreaUsecase struct {
	_crud.BaseUsecase[models.MainAreaEntity, repo.IMainAreaRepo]
}

func NewAdminMainAreaUsecase(areaRepo repo.IMainAreaRepo) IAdminMainAreaUsecase {
	return &AdminMainAreaUsecase{
		BaseUsecase: _crud.BaseUsecase[models.MainAreaEntity, repo.IMainAreaRepo]{
			Repo: areaRepo,
		},
	}
}

func (u *AdminMainAreaUsecase) Search(ctx context.Context, text string, pagable _dto.IPagable) ([]models.MainAreaEntity, uint32, error) {
	return u.Repo.Search(ctx, text, pagable)
}
