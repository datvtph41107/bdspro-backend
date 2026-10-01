//go:build ignore

package postgres

import (
	"context"

	impl "common/db"
	seo_domain "userv4/modules/crm/entity/seo"
	repo "userv4/modules/crm/repository"

	"gorm.io/gorm"
)

type SeoCountPostgres struct {
	db *gorm.DB
}

func NewSeoCountPostgres(db *gorm.DB) repo.SeoCountRepo {
	return &SeoCountPostgres{db: db}
}

func (r *SeoCountPostgres) FlushClicks(ctx context.Context, logs []*seo_domain.SeoCountLog) error {
	if len(logs) == 0 {
		return nil
	}
	return impl.GetDB(ctx, r.db).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, row := range logs {
			if row == nil || row.SeoID == 0 || row.Clicked <= 0 {
				continue
			}
			if err := tx.Exec(
				`INSERT INTO seo_count_log (seo_id, clicked, "timestamp") VALUES (?, ?, ?)`,
				row.SeoID, row.Clicked, row.Timestamp,
			).Error; err != nil {
				return err
			}
			if err := tx.Exec(`
				INSERT INTO seo_count (seo_id, count, updated_at)
				VALUES (?, ?, ?)
				ON CONFLICT (seo_id) DO UPDATE SET
					count = seo_count.count + EXCLUDED.count,
					updated_at = EXCLUDED.updated_at
			`, row.SeoID, row.Clicked, row.Timestamp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
