package crud

import (
	"context"
	"fmt"

	_dto "common/domain/dto"
	_errors "common/errors"
)

type IBaseUsecase[T any] interface {
	Create(c context.Context, entity *T) (*T, error)
	Update(c context.Context, id uint64, entity *T) (*T, error)
	Delete(c context.Context, id uint64) (bool, error)
	GetByID(c context.Context, id uint64) (*T, error)
	GetAll(c context.Context) ([]T, error)
	GetDetail(c context.Context, id uint64) (*T, error)
	GetList(c context.Context, pagable _dto.IPagable) ([]T, int64, error)
	GetListByGroupKey(c context.Context, groupKey uint32, pagable _dto.IPagable) ([]T, int64, error)
}

type BaseUsecase[T any, R ICrudRepo[T]] struct {
	Repo R
}

func (s *BaseUsecase[T, R]) Create(c context.Context, entity *T) (*T, error) {
	if err := s.Repo.Create(c, entity); err != nil {
		return nil, fmt.Errorf("create record: %w", err)
	}
	return entity, nil
}

func (s *BaseUsecase[T, R]) Update(c context.Context, id uint64, entity *T) (*T, error) {
	if err := s.Repo.Update(c, id, entity); err != nil {
		return nil, fmt.Errorf("update record %d: %w", id, err)
	}
	return entity, nil
}

func (s *BaseUsecase[T, R]) Delete(c context.Context, id uint64) (bool, error) {
	if err := s.Repo.Delete(c, id); err != nil {
		return false, fmt.Errorf("delete record %d: %w", id, err)
	}
	return true, nil
}

func (s *BaseUsecase[T, R]) GetByID(c context.Context, id uint64) (*T, error) {
	entity, err := s.Repo.GetByID(c, id)
	if err != nil {
		return nil, _errors.ReturnError(_errors.DataNotFound, _errors.WithCause(err))
	}
	return entity, nil
}

func (s *BaseUsecase[T, R]) GetAll(c context.Context) ([]T, error) {
	entities, err := s.Repo.GetAll(c)
	if err != nil {
		return nil, fmt.Errorf("list records: %w", err)
	}
	return entities, nil
}

func (s *BaseUsecase[T, R]) GetDetail(c context.Context, id uint64) (*T, error) {
	entity, err := s.Repo.GetDetail(c, id)
	if err != nil {
		return nil, _errors.ReturnError(_errors.DataNotFound, _errors.WithCause(err))
	}
	return entity, nil
}

func (s *BaseUsecase[T, R]) GetList(c context.Context, pagable _dto.IPagable) ([]T, int64, error) {
	entities, total, err := s.Repo.GetList(c, pagable)
	if err != nil {
		return nil, 0, fmt.Errorf("list paged records: %w", err)
	}
	return entities, total, nil
}
