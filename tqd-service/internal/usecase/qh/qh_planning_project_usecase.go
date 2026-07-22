package qh_usecase

import (
	_err "common/domain/err"
	"context"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/interface/repo"
)

type IQHPlanningProjectUsecase interface {
	Create(ctx context.Context, entity *qh_domain.QHPlanningProject) (*qh_domain.QHPlanningProject, error)
	Update(ctx context.Context, id uint64, entity *qh_domain.QHPlanningProject) (*qh_domain.QHPlanningProject, error)
	Delete(ctx context.Context, id uint64) (bool, error)
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningProject, error)
	GetList(ctx context.Context, req *qh_dto.ListPlanningProjectsRequest) ([]qh_domain.QHPlanningProject, int64, error)
}

type qhPlanningProjectUsecase struct {
	repo repo.IQHPlanningProjectRepo
}

// NewQHPlanningProjectUsecase @bind: internal/usecase/qh.IQHPlanningProjectUsecase
func NewQHPlanningProjectUsecase(repo repo.IQHPlanningProjectRepo) IQHPlanningProjectUsecase {
	return &qhPlanningProjectUsecase{
		repo: repo,
	}
}

func (u *qhPlanningProjectUsecase) Create(ctx context.Context, entity *qh_domain.QHPlanningProject) (*qh_domain.QHPlanningProject, error) {
	if err := u.repo.Create(ctx, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return entity, nil
}

func (u *qhPlanningProjectUsecase) Update(ctx context.Context, id uint64, entity *qh_domain.QHPlanningProject) (*qh_domain.QHPlanningProject, error) {
	if err := u.repo.Update(ctx, id, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể cập nhật dữ liệu",
		}
	}
	return entity, nil
}

func (u *qhPlanningProjectUsecase) Delete(ctx context.Context, id uint64) (bool, error) {
	if err := u.repo.Delete(ctx, id); err != nil {
		return false, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể xóa dữ liệu",
		}
	}
	return true, nil
}

func (u *qhPlanningProjectUsecase) GetByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningProject, error) {
	entity, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}
	return entity, nil
}

func (u *qhPlanningProjectUsecase) GetList(ctx context.Context, req *qh_dto.ListPlanningProjectsRequest) ([]qh_domain.QHPlanningProject, int64, error) {
	entities, total, err := u.repo.GetList(ctx, req)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return entities, total, nil
}
