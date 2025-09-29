package repo

import (
	"common/case/crud"
	_db "common/db"
	"context"

	"membership/models"

	"gorm.io/gorm"
)

type PlanRepo struct {
	crud.CrudRepo[models.Plan]
	*_db.TransactionRepo
}

func NewPlanRepo(db *gorm.DB) *PlanRepo {
	return &PlanRepo{
		CrudRepo: crud.CrudRepo[models.Plan]{TransactionRepo: _db.NewTransactionRepo(db)},
	}
}

func (p *PlanRepo) ListForUser() ([]models.Plan, error) {
	var entities []models.Plan

	err := p.TransactionRepo.GetDB(context.Background()).Debug().Where("deleted_at IS NULL and visibility = true").Find(&entities).Error
	return entities, err
}
