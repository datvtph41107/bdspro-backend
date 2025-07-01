package user_model

import (
	_models "common/models"
	"time"
)

type Rate struct {
	_models.BaseEntity
	Anonymous  bool           `gorm:"not null" json:"anonymous"`
	TargetID   uint64         `gorm:"not null" json:"targetId" binding:"required"`
	TargetType uint8          `gorm:"default:10;not null" json:"targetType" binding:"required"` // 10: user 20: group 30: organization
	Score      uint8          `gorm:"not null" json:"score" binding:"required,gte=1,lte=5"`
	Comment    string         `gorm:"type:text" json:"comment,omitempty" binding:"required"`
	ParentID   *uint64        `json:"parentId,omitempty"`
	Attachs    []AttachEntity `gorm:"foreignKey:RateID" json:"attachs,omitempty"`
	RootID     *uint64        `json:"rootId,omitempty"`
	Hidden     bool           `gorm:"default:false" json:"hidden"`
	// CreatedUser UserItem       `gorm:"foreignKey:CreatedBy" json:"createdUser,omitempty"`
}

type RateItem struct {
	ID          uint64         `gorm:"primaryKey;column:id" json:"id"`
	Anonymous   bool           `gorm:"not null" json:"anonymous"`
	Score       uint8          `gorm:"not null" json:"score" binding:"required,gte=1,lte=5"`
	Comment     string         `gorm:"type:text" json:"comment,omitempty" binding:"required"`
	ParentID    *uint64        `json:"parentId,omitempty"`
	CreatedBy   *uint64        `gorm:"column:created_by" json:"-"`
	CreatedUser *UserItem      `gorm:"foreignKey:CreatedBy" json:"createdUser,omitempty"`
	CreatedAt   *time.Time     `gorm:"column:created_at" json:"createdAt"`
	Attachs     []AttachEntity `gorm:"foreignKey:RateID" json:"attachs,omitempty"`
}

func (RateItem) TableName() string {
	return "rates"
}

type RateStats struct {
	TotalReviews int64   `gorm:"column:total_reviews" json:"totalReviews"`
	AverageScore float64 `gorm:"column:average_score" json:"averageScore"`
	Star1Count   int64   `gorm:"column:star_1_count" json:"star1Count"`
	Star2Count   int64   `gorm:"column:star_2_count" json:"star2Count"`
	Star3Count   int64   `gorm:"column:star_3_count" json:"star3Count"`
	Star4Count   int64   `gorm:"column:star_4_count" json:"star4Count"`
	Star5Count   int64   `gorm:"column:star_5_count" json:"star5Count"`
}
