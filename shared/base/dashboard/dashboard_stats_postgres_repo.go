package dashboard

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type DashboardStatsPostgresRepository struct {
	db *gorm.DB
}

func NewDashboardStatsPostgresRepository(db *gorm.DB) *DashboardStatsPostgresRepository {
	return &DashboardStatsPostgresRepository{db: db}
}

func (r *DashboardStatsPostgresRepository) Create(ctx context.Context, tableName string, stats *DashboardStats) (*DashboardStats, error) {
	err := r.db.Table(tableName).Create(stats).Error
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *DashboardStatsPostgresRepository) GetByTime(ctx context.Context, tableName string, timeReq time.Time) (*DashboardStats, error) {
	var stats DashboardStats
	dateOnly := timeReq.Truncate(24 * time.Hour)
	err := r.db.Table(tableName).Where("DATE(calculate_time) = DATE(?)", dateOnly).First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *DashboardStatsPostgresRepository) GetByTwoTime(ctx context.Context, tableName string, firstTime time.Time, secondTime time.Time) (*DashboardStats, *DashboardStats, error) {
	var firstStats DashboardStats
	var secondStats DashboardStats
	// Truncate times to date (remove time component)
	firstDateOnly := firstTime.Truncate(24 * time.Hour)
	secondDateOnly := secondTime.Truncate(24 * time.Hour)

	err := r.db.Table(tableName).Where("DATE(calculate_time) = DATE(?)", firstDateOnly).First(&firstStats).Error
	if err != nil {
		return nil, nil, err
	}
	err = r.db.Table(tableName).Where("DATE(calculate_time) = DATE(?)", secondDateOnly).First(&secondStats).Error
	if err != nil {
		return nil, nil, err
	}
	return &firstStats, &secondStats, nil
}

func (r *DashboardStatsPostgresRepository) UpdateOrCreate(ctx context.Context, tableName string, stats *DashboardStats) (*DashboardStats, error) {
	var existingStats DashboardStats
	err := r.db.Table(tableName).Where("calculate_time = ?", stats.CalculateTime).First(&existingStats).Error

	if err != nil {
		// Không tìm thấy, tạo mới
		err = r.db.Table(tableName).Create(stats).Error
		if err != nil {
			return nil, err
		}
		return stats, nil
	}

	// Cập nhật count
	existingStats.Count = stats.Count
	err = r.db.Table(tableName).Save(&existingStats).Error
	if err != nil {
		return nil, err
	}
	return &existingStats, nil
}
