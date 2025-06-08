package postgre

import (
	"context"
	"social/internal/domain"

	"gorm.io/gorm"
)

// @bind: social/internal/repo.ReportRepo
type PostgreReport struct {
	db *gorm.DB
}

func NewPostgreReport(db *gorm.DB) *PostgreReport {
	return &PostgreReport{db: db}
}

func (r *PostgreReport) Create(ctx context.Context, report *domain.Report) (*domain.Report, error) {
	if err := r.db.WithContext(ctx).Create(report).Error; err != nil {
		return nil, err
	}
	return report, nil
}

func (r *PostgreReport) ExistByNewsFeedIdAndUserId(ctx context.Context, newsFeedId uint64, userId uint64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Report{}).
		Where("news_feed_id = ? AND user_id = ? AND deleted_at IS NULL", newsFeedId, userId).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
