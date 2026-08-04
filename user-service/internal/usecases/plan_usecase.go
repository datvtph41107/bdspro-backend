package usecases

import (
	_crud "common/domain/crud"
	_err "common/domain/err"
	"context"
	"user/internal/catalogshadow"
	"user/internal/interface/repo"
	"user/internal/models"
)

type PlanUsecase struct {
	_crud.BaseUsecase[models.Plan, repo.IPlanRepo]
	planRepo       repo.IPlanRepo
	shadowObserver catalogshadow.Observer
}

type IPlanUsecase interface {
	_crud.IBaseUsecase[models.Plan]
	GetVisiblePlans(ctx context.Context) ([]models.Plan, *_err.ErrorDTO)
}

func NewPlanUsecase(
	planRepo repo.IPlanRepo,
	shadowObserver catalogshadow.Observer,
) IPlanUsecase {
	return &PlanUsecase{
		BaseUsecase: _crud.BaseUsecase[models.Plan, repo.IPlanRepo]{
			Repo: planRepo,
		},
		planRepo:       planRepo,
		shadowObserver: shadowObserver,
	}
}

func (u *PlanUsecase) GetVisiblePlans(ctx context.Context) ([]models.Plan, *_err.ErrorDTO) {
	plans, err := u.planRepo.GetVisiblePlans(ctx)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể lấy danh sách gói",
		}
	}
	if u.shadowObserver != nil {
		u.shadowObserver.ObserveVisiblePlans(ctx, plans)
	}
	return plans, nil
}
