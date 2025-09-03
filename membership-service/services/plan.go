package services

import (
	"common/case/crud"

	"membership/models"
	"membership/repo"
)

type PlanService struct {
	crud.BaseUsecase[models.Plan, crud.ICrudRepo[models.Plan]]
}

func NewPlanService(repo *repo.PlanRepo) *PlanService {
	return &PlanService{
		BaseUsecase: crud.BaseUsecase[models.Plan, crud.ICrudRepo[models.Plan]]{Repo: repo},
	}
}
