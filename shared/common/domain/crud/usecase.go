package _crud

import (
	"context"

	_dto "common/domain/dto"
	_err "common/domain/err"
)

type IBaseUsecase[T any] interface {
	Create(c context.Context, entity *T) (*T, error)
	Update(c context.Context, id uint64, entity *T) (*T, error)
	Delete(c context.Context, id uint64) (bool, error)
	GetByID(c context.Context, id uint64) (*T, error)
	GetAll(c context.Context) ([]T, error)
	GetDetail(c context.Context, id uint64) (*T, error)
	GetList(c context.Context, pagable _dto.IPagable) ([]T, int64, error)
}

type BaseUsecase[T any, R ICrudRepo[T]] struct {
	Repo R
}

func (s *BaseUsecase[T, R]) Create(c context.Context, entity *T) (*T, error) {
	if err := s.Repo.Create(c, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}

	return entity, nil
}

func (s *BaseUsecase[T, R]) Update(c context.Context, id uint64, entity *T) (*T, error) {
	if err := s.Repo.Update(c, id, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể cập nhật dữ liệu",
		}
	}

	return entity, nil
}

func (s *BaseUsecase[T, R]) Delete(c context.Context, id uint64) (bool, error) {
	if err := s.Repo.Delete(c, id); err != nil {
		return false, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể xóa dữ liệu",
		}
	}

	return true, nil
}

func (s *BaseUsecase[T, R]) GetByID(c context.Context, id uint64) (*T, error) {
	entity, err := s.Repo.GetByID(c, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}

	return entity, nil
}

func (s *BaseUsecase[T, R]) GetAll(c context.Context) ([]T, error) {
	entities, err := s.Repo.GetAll(c)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}

	return entities, nil
}

func (s *BaseUsecase[T, R]) GetDetail(c context.Context, id uint64) (*T, error) {
	entity, err := s.Repo.GetDetail(c, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}

	return entity, nil
}

func (s *BaseUsecase[T, R]) GetList(c context.Context, pagable _dto.IPagable) ([]T, int64, error) {
	entities, total, err := s.Repo.GetList(c, pagable)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}

	return entities, total, nil
}
