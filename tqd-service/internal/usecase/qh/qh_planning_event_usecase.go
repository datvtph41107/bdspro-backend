package qh_usecase

import (
	_err "common/domain/err"
	"context"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/interface/repo"
)

type IQHPlanningEventUsecase interface {
	GetList(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh_domain.QHPlanningEvent, int64, error)
	GetRelations(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh_domain.QHPlanningRelationItem, int64, error)
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

func (u *qhPlanningEventUsecase) GetList(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh_domain.QHPlanningEvent, int64, error) {
	entities, total, err := u.repo.GetList(ctx, req)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}

	return entities, total, nil
}

func (u *qhPlanningEventUsecase) GetRelations(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh_domain.QHPlanningRelationItem, int64, error) {
	entities, total, err := u.repo.GetRelations(ctx, req)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return entities, total, nil
}
