package repo

import (
	"common/case/crud"
	"membership/models"

	"gorm.io/gorm"
)

type PlanRepo struct {
	crud.BaseRepo[models.Plan]
}

func NewPlanRepo(db *gorm.DB) *PlanRepo {
	return &PlanRepo{
		BaseRepo: crud.BaseRepo[models.Plan]{DB: db},
	}
}

func (p *PlanRepo) ListForUser() ([]models.Plan, error) {
	var entities []models.Plan

	err := p.DB.Debug().Where("deleted_at IS NULL and visibility = true").Find(&entities).Error
	return entities, err
}
