package postgres

import (
	"context"
	"errors"

	"notification/internal/domain"
	usecase "notification/internal/usecase"

	"gorm.io/gorm"
)

// @bind: notification/internal/usecase.HistoryAuthStore
type HistoryAuthPostgres struct {
	db *gorm.DB
}

func NewHistoryAuthPostgres(db *gorm.DB) *HistoryAuthPostgres {
	return &HistoryAuthPostgres{db: db}
}

func (r *HistoryAuthPostgres) CreateHistoryAuth(ctx context.Context, entity *domain.HistoryAuthEntity) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *HistoryAuthPostgres) SearchHistoryAuth(ctx context.Context, params usecase.HistoryAuthSearch) ([]domain.HistoryAuthEntity, int64, error) {
	if params.UserID == 0 {
		return nil, 0, errors.New("userId is required")
	}

	query := r.db.WithContext(ctx).Model(&domain.HistoryAuthEntity{}).Where("user_id = ?", params.UserID)

	if params.OrganizationID != nil {
		query = query.Where("organization_id = ?", params.OrganizationID)
	}

	if params.ActionType != "" {
		query = query.Where("action_type = ?", params.ActionType)
	}

	if params.Channel != "" {
		query = query.Where("channel = ?", params.Channel)
	}

	if params.DeviceID != "" {
		query = query.Where("device_id = ?", params.DeviceID)
	}

	if params.Success != nil {
		query = query.Where("success = ?", *params.Success)
	}

	if params.FromDate != nil {
		query = query.Where("created_at >= ?", params.FromDate)
	}

	if params.ToDate != nil {
		query = query.Where("created_at <= ?", params.ToDate)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := params.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := params.Offset
	if offset < 0 {
		offset = 0
	}

	var histories []domain.HistoryAuthEntity
	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}
