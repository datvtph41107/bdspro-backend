package models

import (
	"time"

	"gorm.io/gorm"
)

// Project struct represents the projects table
// TableName sets the table name explicitly
func (Project) TableName() string {
	return "projects"
}

type Project struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
