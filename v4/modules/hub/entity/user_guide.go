package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserGuide struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt   *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy   *uint64        `gorm:"column:created_by"`
	UpdatedBy   *uint64        `gorm:"column:updated_by"`
	Title       string         `gorm:"column:title;type:text;not null"`
	Description string         `gorm:"column:description;type:text;not null"`
	GroupKey    string         `gorm:"column:group_key;type:text;not null;index"`
	Key         string         `gorm:"column:key;index"`
	Mode        int32          `gorm:"column:mode;not null;default:10"`
	Steps       []UserGuideStep `gorm:"foreignKey:UserGuideID;references:ID"`
}

func (UserGuide) TableName() string { return "tb_user_guide" }

type UserGuideStep struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt   *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
	UserGuideID uint64         `gorm:"column:user_guide_id;index;not null"`
	StepOrder   int            `gorm:"column:step_order;not null"`
	Image       string         `gorm:"column:image;type:text"`
	Content     string         `gorm:"column:content;type:text;not null"`
}

func (UserGuideStep) TableName() string { return "user_guide_steps" }
