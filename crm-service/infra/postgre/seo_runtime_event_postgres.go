package postgre

import (
	"context"

	impl "common/db"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/repo"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SeoRuntimeEventPostgres struct {
	db *gorm.DB
}

func NewSeoRuntimeEventPostgres(db *gorm.DB) repo.SeoRuntimeEventRepo {
	return &SeoRuntimeEventPostgres{db: db}
}

func (r *SeoRuntimeEventPostgres) CreateIdempotent(ctx context.Context, event *seo_domain.SeoRuntimeEvent) (bool, error) {
	if event == nil {
		return false, nil
	}
	result := impl.GetDB(ctx, r.db).WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "event_key"}},
		DoNothing: true,
	}).Create(event)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
