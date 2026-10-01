package qh_usecase

import (
	_err "common/domain/err"
	"context"
	"errors"

	qh "userv4/modules/tqd/entity/qh"
	qh_dto "userv4/modules/tqd/dto/qh"
	repo "userv4/modules/tqd/repository"

	"gorm.io/gorm"
)

type IQHPlanningEventUsecase interface {
	Create(ctx context.Context, entity *qh.QHPlanningEvent) (*qh.QHPlanningEvent, error)
	Update(ctx context.Context, id uint64, entity *qh.QHPlanningEvent) (*qh.QHPlanningEvent, error)
	Delete(ctx context.Context, id uint64) (bool, error)
	GetByID(ctx context.Context, id uint64) (*qh.QHPlanningEvent, error)
	GetList(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh.QHPlanningEvent, int64, error)
	GetRelations(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh.QHPlanningRelationItem, int64, error)
	CreateRelation(ctx context.Context, entity *qh.QHPlanningRelation) (*qh.QHPlanningRelationItem, error)
	UpdateRelation(ctx context.Context, id uint64, entity *qh.QHPlanningRelation) (*qh.QHPlanningRelationItem, error)
	DeleteRelation(ctx context.Context, id uint64) (bool, error)
	GetRelationByID(ctx context.Context, id uint64) (*qh.QHPlanningRelationItem, error)
}

type qhPlanningEventUsecase struct {
	repo repo.IQHPlanningEventRepo
}

// NewQHPlanningEventUsecase @bind: internal/usecase/qh.IQHPlanningEventUsecase
func NewQHPlanningEventUsecase(repo repo.IQHPlanningEventRepo) IQHPlanningEventUsecase {
	return &qhPlanningEventUsecase{
		repo: repo,
	}
}

func (u *qhPlanningEventUsecase) Create(ctx context.Context, entity *qh.QHPlanningEvent) (*qh.QHPlanningEvent, error) {
	if err := u.repo.Create(ctx, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	created, err := u.repo.GetByID(ctx, entity.ID)
	if err != nil {
		return entity, nil
	}
	return created, nil
}

func (u *qhPlanningEventUsecase) Update(ctx context.Context, id uint64, entity *qh.QHPlanningEvent) (*qh.QHPlanningEvent, error) {
	if _, err := u.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &_err.ErrorDTO{
				Code:    404,
				Message: "Không tìm thấy sự kiện",
			}
		}
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	if err := u.repo.Update(ctx, id, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể cập nhật dữ liệu",
		}
	}
	updated, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return entity, nil
	}
	return updated, nil
}

func (u *qhPlanningEventUsecase) Delete(ctx context.Context, id uint64) (bool, error) {
	if _, err := u.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &_err.ErrorDTO{
				Code:    404,
				Message: "Không tìm thấy sự kiện",
			}
		}
		return false, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	if err := u.repo.Delete(ctx, id); err != nil {
		return false, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể xóa dữ liệu",
		}
	}
	return true, nil
}

func (u *qhPlanningEventUsecase) GetByID(ctx context.Context, id uint64) (*qh.QHPlanningEvent, error) {
	entity, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &_err.ErrorDTO{
				Code:    404,
				Message: "Không tìm thấy sự kiện",
			}
		}
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return entity, nil
}

func (u *qhPlanningEventUsecase) GetList(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh.QHPlanningEvent, int64, error) {
	entities, total, err := u.repo.GetList(ctx, req)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}

	return entities, total, nil
}

func (u *qhPlanningEventUsecase) GetRelations(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh.QHPlanningRelationItem, int64, error) {
	entities, total, err := u.repo.GetRelations(ctx, req)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return entities, total, nil
}

func (u *qhPlanningEventUsecase) CreateRelation(ctx context.Context, entity *qh.QHPlanningRelation) (*qh.QHPlanningRelationItem, error) {
	if err := u.repo.CreateRelation(ctx, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	item, err := u.repo.GetRelationByID(ctx, entity.ID)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return item, nil
}

func (u *qhPlanningEventUsecase) UpdateRelation(ctx context.Context, id uint64, entity *qh.QHPlanningRelation) (*qh.QHPlanningRelationItem, error) {
	if _, err := u.repo.GetRelationByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &_err.ErrorDTO{
				Code:    404,
				Message: "Không tìm thấy quan hệ đồ án",
			}
		}
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	if err := u.repo.UpdateRelation(ctx, id, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể cập nhật dữ liệu",
		}
	}
	item, err := u.repo.GetRelationByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return item, nil
}

func (u *qhPlanningEventUsecase) DeleteRelation(ctx context.Context, id uint64) (bool, error) {
	if _, err := u.repo.GetRelationByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &_err.ErrorDTO{
				Code:    404,
				Message: "Không tìm thấy quan hệ đồ án",
			}
		}
		return false, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	if err := u.repo.DeleteRelation(ctx, id); err != nil {
		return false, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể xóa dữ liệu",
		}
	}
	return true, nil
}

func (u *qhPlanningEventUsecase) GetRelationByID(ctx context.Context, id uint64) (*qh.QHPlanningRelationItem, error) {
	item, err := u.repo.GetRelationByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &_err.ErrorDTO{
				Code:    404,
				Message: "Không tìm thấy quan hệ đồ án",
			}
		}
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return item, nil
}
