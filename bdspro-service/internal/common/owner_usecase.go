package common

import (
	"bdspro/internal/enums"
	_utils "common/utils"
	"context"
)

type IOwnerUsecase[T any, D DTO] interface {
	Search(c context.Context, ownerType enums.EOwnerOf, dto D) ([]T, int64, error)
	GetByID(c context.Context, id uint64) (*T, error)
	Detail(c context.Context, id uint64) (*T, error)
	Create(c context.Context, entity *T) error
	Update(c context.Context, id uint64, entity *T) error
	Delete(c context.Context, id uint64) error
	CheckPermission(c context.Context, id uint64) error
}

type OwnerUsecase[T any, D DTO] struct {
	OwnerRepo IOwnerRepo[T, D]
	UC        IOwnerUsecase[T, D]
}

func (uc *OwnerUsecase[T, D]) Search(c context.Context, ownerType enums.EOwnerOf, dto D) ([]T, int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	return uc.OwnerRepo.Search(c, profileId, ownerType, dto)
}

func (uc *OwnerUsecase[T, D]) GetByID(c context.Context, id uint64) (*T, error) {
	err := uc.UC.CheckPermission(c, id)
	if err != nil {
		return nil, err
	}
	return uc.OwnerRepo.GetByID(c, id)
}

func (uc *OwnerUsecase[T, D]) Detail(c context.Context, id uint64) (*T, error) {
	err := uc.UC.CheckPermission(c, id)
	if err != nil {
		return nil, err
	}
	return uc.OwnerRepo.Detail(c, id)
}

func (uc *OwnerUsecase[T, D]) Create(c context.Context, entity *T) error {
	return uc.OwnerRepo.Create(c, entity)
}

func (uc *OwnerUsecase[T, D]) Update(c context.Context, id uint64, entity *T) error {
	err := uc.UC.CheckPermission(c, id)
	if err != nil {
		return err
	}
	return uc.OwnerRepo.Update(c, id, entity)
}

func (uc *OwnerUsecase[T, D]) Delete(c context.Context, id uint64) error {
	err := uc.UC.CheckPermission(c, id)
	if err != nil {
		return err
	}
	return uc.OwnerRepo.Delete(c, id)
}

func (uc *OwnerUsecase[T, D]) CheckPermission(c context.Context, id uint64) error {
	return nil
}
