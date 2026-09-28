package entity

import (
	"time"

	"gorm.io/gorm"
)

type FAQ struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy *uint64        `gorm:"column:created_by"`
	UpdatedBy *uint64        `gorm:"column:updated_by"`
	Question  string         `gorm:"column:question;type:text;not null"`
	Answer    string         `gorm:"column:answer;type:text;not null"`
	GroupKey  string         `gorm:"column:group_key;size:100;index"`
}

func (FAQ) TableName() string { return "faqs" }
