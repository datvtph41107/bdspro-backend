package domain

import (
	"time"

	"gorm.io/gorm"
)

// Task struct represents the tasks table
func (Task) TableName() string {
	return "tasks"
}

type Task struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	ProjectID   uint   `gorm:"not null;index"`
	Status      string `gorm:"size:50"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
