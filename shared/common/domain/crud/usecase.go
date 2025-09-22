package _crud

import (
	"context"

	_err "common/domain/err"
	_dto "common/dto"
)

type IBaseUsecase[T any] interface {
	Create(c context.Context, entity *T) (*T, *_err.ErrorDTO)
	Update(c context.Context, id uint64, entity *T) (*T, *_err.ErrorDTO)
	Delete(c context.Context, id uint64) (bool, *_err.ErrorDTO)
	GetByID(c context.Context, id uint64) (*T, *_err.ErrorDTO)
	GetAll(c context.Context) ([]T, *_err.ErrorDTO)
	GetDetail(c context.Context, id uint64) (*T, *_err.ErrorDTO)
	GetList(c context.Context, pagable _dto.IPagable) ([]T, int64, *_err.ErrorDTO)
	GetListByGroupKey(c context.Context, groupKey uint32, pagable _dto.IPagable) ([]T, int64, *_err.ErrorDTO)
}

type BaseUsecase[T any, R ICrudRepo[T]] struct {
	Repo R
}

func (s *BaseUsecase[T, R]) Create(c context.Context, entity *T) (*T, *_err.ErrorDTO) {
	if err := s.Repo.Create(c, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}

	return entity, nil
}

func (s *BaseUsecase[T, R]) Update(c context.Context, id uint64, entity *T) (*T, *_err.ErrorDTO) {
	if err := s.Repo.Update(c, id, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể cập nhật dữ liệu",
		}
	}

	return entity, nil
}

func (s *BaseUsecase[T, R]) Delete(c context.Context, id uint64) (bool, *_err.ErrorDTO) {
	if err := s.Repo.Delete(c, id); err != nil {
		return false, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể xóa dữ liệu",
		}
	}

	return true, nil
}

func (s *BaseUsecase[T, R]) GetByID(c context.Context, id uint64) (*T, *_err.ErrorDTO) {
	entity, err := s.Repo.GetByID(c, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}

	return entity, nil
}

func (s *BaseUsecase[T, R]) GetAll(c context.Context) ([]T, *_err.ErrorDTO) {
	entities, err := s.Repo.GetAll(c)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}

	return entities, nil
}

func (s *BaseUsecase[T, R]) GetDetail(c context.Context, id uint64) (*T, *_err.ErrorDTO) {
	entity, err := s.Repo.GetDetail(c, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}

	return entity, nil
}

func (s *BaseUsecase[T, R]) GetList(c context.Context, pagable _dto.IPagable) ([]T, int64, *_err.ErrorDTO) {
	entities, total, err := s.Repo.GetList(c, pagable)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}

	return entities, total, nil
}
