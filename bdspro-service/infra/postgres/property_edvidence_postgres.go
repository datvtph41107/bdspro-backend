package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type PostgrePropertyEdvidenceRepo struct {
	DB *gorm.DB
}

func NewPostgrePropertyEdvidenceRepo(db *gorm.DB) repo.PropertyEdvidenceRepo {
	return &PostgrePropertyEdvidenceRepo{DB: db}
}

func (r *PostgrePropertyEdvidenceRepo) Create(ctx context.Context, entity *domain.PropertyEdvidence) error {
	return GetDB(ctx, r.DB).Create(entity).Error
}
