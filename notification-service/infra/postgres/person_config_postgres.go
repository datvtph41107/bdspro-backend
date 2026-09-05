package postgres

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"notification/internal/domain"
)

// @bind: notification/internal/usecase.PersonConfigStore
type PersonConfigPostgres struct {
	DB *gorm.DB
}

func NewPersonConfigPostgres(db *gorm.DB) *PersonConfigPostgres {
	return &PersonConfigPostgres{DB: db}
}

func (r *PersonConfigPostgres) BatchUpsert(ctx context.Context, configs []*domain.PersonConfigEntity) ([]*domain.PersonConfigEntity, error) {
	if len(configs) == 0 {
		return []*domain.PersonConfigEntity{}, nil
	}

	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	for _, cfg := range configs {
		if cfg.ID > 0 {
			result := tx.Model(&domain.PersonConfigEntity{}).
				Where("id = ? AND user_id = ?", cfg.ID, cfg.UserID).
				Updates(map[string]interface{}{
					"key":        cfg.Key,
					"checked":    cfg.Checked,
					"is_default": cfg.IsDefault,
					"channel":    cfg.Channel,
				})
			if result.Error != nil {
				tx.Rollback()
				return nil, result.Error
			}

			if result.RowsAffected > 0 {
				continue
			}

			// Nếu không cập nhật được theo ID, reset ID để tạo mới theo user_id + key
			cfg.ID = 0
		}

		err := tx.Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "channel"}, {Name: "user_id"}, {Name: "key"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"checked":    cfg.Checked,
					"is_default": cfg.IsDefault,
				}),
			},
			clause.Returning{},
		).Create(cfg).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return configs, nil
}

func (r *PersonConfigPostgres) ListByUserID(ctx context.Context, userID uint64) ([]*domain.PersonConfigEntity, error) {
	if userID == 0 {
		return []*domain.PersonConfigEntity{}, nil
	}

	var configs []*domain.PersonConfigEntity
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("id asc").
		Find(&configs).Error
	if err != nil {
		return nil, err
	}

	return configs, nil
}
