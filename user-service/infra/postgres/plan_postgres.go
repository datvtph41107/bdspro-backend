package postgres

import (
	_db "common/db"
	_provider "common/provider"
	"context"
	"user/internal/interface/repo"
	"user/internal/models"
)

type PlanRepo struct {
	_provider.CrudRepo[models.Plan]
}

func NewPlanRepo(db *_db.TransactionRepo) repo.IPlanRepo {
	repo := &PlanRepo{}
	repo.Init(repo, db)
	return repo
}

func (r *PlanRepo) GetVisiblePlans(ctx context.Context) ([]models.Plan, error) {
	var entities []models.Plan
	err := r.TransactionRepo.GetDB(ctx).Where("deleted_at IS NULL AND visibility = true").Find(&entities).Error
	return entities, err
}
