package services

import (
	"common/case/crud"
	"membership/models"
	"membership/repo"
)

type PlanService struct {
	crud.BaseService[models.Plan]
}

func NewPlanService(repo *repo.PlanRepo) *PlanService {
	return &PlanService{
		BaseService: crud.BaseService[models.Plan]{Repo: repo.BaseRepo},
	}
}
