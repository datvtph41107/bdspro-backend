package common

import (
	"bdspro/internal/enums"
	_utils "common/utils"
	"context"
)

type CrudUsecase[T any, D DTO] struct {
	Repo IBaseRepo[T, D]
}

func (uc *CrudUsecase[T, D]) Search(c context.Context, ownerType enums.EOwnerOf, dto D) ([]T, int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	return uc.Repo.Search(c, profileId, ownerType, dto)
}

func (uc *CrudUsecase[T, D]) GetByID(c context.Context, id uint64) (*T, error) {
	if err := uc.CheckPermission(c, id); err != nil {
		return nil, err
	}
	return uc.Repo.GetByID(c, id)
}

func (uc *CrudUsecase[T, D]) Detail(c context.Context, id uint64) (*T, error) {
	if err := uc.CheckPermission(c, id); err != nil {
		return nil, err
	}
	return uc.Repo.Detail(c, id)
}

func (uc *CrudUsecase[T, D]) Create(c context.Context, entity *T) error {
	return uc.Repo.Create(c, entity)
}

func (uc *CrudUsecase[T, D]) Update(c context.Context, id uint64, entity *T) error {
	if err := uc.CheckPermission(c, id); err != nil {
		return err
	}
	return uc.Repo.Update(c, id, entity)
}

func (uc *CrudUsecase[T, D]) Delete(c context.Context, id uint64) error {
	if err := uc.CheckPermission(c, id); err != nil {
		return err
	}
	return uc.Repo.Delete(c, id)
}

func (uc *CrudUsecase[T, D]) CheckPermission(c context.Context, id uint64) error {
	return nil
}
