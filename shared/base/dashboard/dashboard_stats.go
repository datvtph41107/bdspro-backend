package dashboard

import "time"

type DashboardStats struct {
	ID            uint32    `gorm:"column:id;primaryKey;autoIncrement"`
	CalculateTime time.Time `gorm:"column:calculate_time;not null"`
	Count         int64     `gorm:"column:count;not null;default:0"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
