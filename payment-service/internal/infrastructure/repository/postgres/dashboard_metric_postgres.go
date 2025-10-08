package postgres

import (
	"context"
	"errors"
	"payment/internal/domain/entity"
	"payment/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

type DashboardMetricPostgresRepository struct {
	db *gorm.DB
}

func NewDashboardMetricPostgresRepository(db *gorm.DB) repository.DashboardMetricRepository {
	return &DashboardMetricPostgresRepository{db: db}
}

func (r *DashboardMetricPostgresRepository) CreateOrUpdate(ctx context.Context, metric *entity.DashboardMetric) (*entity.DashboardMetric, error) {
	if err := r.db.WithContext(ctx).Where("time = ?", metric.Time).First(&entity.DashboardMetric{}).Error; err != nil {
		if err := r.db.WithContext(ctx).Create(metric).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			return nil, err
		}
		return metric, nil
	}
	if err := r.db.WithContext(ctx).Model(&entity.DashboardMetric{}).Where("time = ?", metric.Time).Updates(metric).Error; err != nil {
		return nil, err
	}
	return metric, nil
}

func (r *DashboardMetricPostgresRepository) GetDashboardMetrics(ctx context.Context, fromDate time.Time, toDate time.Time) ([]*entity.DashboardMetric, error) {
	var metrics []*entity.DashboardMetric
	if err := r.db.WithContext(ctx).Where("time BETWEEN ? AND ?", fromDate, toDate).Find(&metrics).Error; err != nil {
		return nil, err
	}
	return metrics, nil
}
