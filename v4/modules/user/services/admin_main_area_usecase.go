package services

import (
	_crud "common/domain/crud"
	_dto "common/domain/dto"
	"context"
	"userv4/internal/interface/repo"
	"userv4/modules/user/entity"
)

type IAdminMainAreaUsecase interface {
	_crud.IBaseUsecase[entity.MainAreaEntity]
	Search(ctx context.Context, text string, pagable _dto.IPagable) ([]entity.MainAreaEntity, uint32, error)
}

type AdminMainAreaUsecase struct {
	_crud.BaseUsecase[entity.MainAreaEntity, repo.IMainAreaRepo]
}

func NewAdminMainAreaUsecase(areaRepo repo.IMainAreaRepo) IAdminMainAreaUsecase {
	return &AdminMainAreaUsecase{
		BaseUsecase: _crud.BaseUsecase[entity.MainAreaEntity, repo.IMainAreaRepo]{
			Repo: areaRepo,
		},
	}
}

func (u *AdminMainAreaUsecase) Search(ctx context.Context, text string, pagable _dto.IPagable) ([]entity.MainAreaEntity, uint32, error) {
	return u.Repo.Search(ctx, text, pagable)
}
