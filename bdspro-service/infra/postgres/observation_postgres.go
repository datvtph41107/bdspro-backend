package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type PostgreObservationRepo struct {
	DB *gorm.DB
}

func NewPostgreObservationRepo(db *gorm.DB) repo.ObservationRepository {
	return &PostgreObservationRepo{
		DB: db,
	}
}

func (r *PostgreObservationRepo) Insert(
	ctx context.Context,
	o *domain.Observation,
) error {
	return GetDB(ctx, r.DB).Create(o).Error
}

func (r *PostgreObservationRepo) GetProductStats(
	ctx context.Context,
	productID uint64,
) (*dto.ProductStats, error) {
	var result dto.ProductStats
	err := GetDB(ctx, r.DB).
		Table("observations").
		Select(`
			entity_id as product_id,
			count(distinct profile_id) as view_count,
			sum(duration) as total_duration
		`).
		Where("entity_type = 'product'").
		Where("entity_id = ?", productID).
		Where("action = 'view'").
		Group("entity_id").
		Scan(&result).Error

	return &result, err
}
