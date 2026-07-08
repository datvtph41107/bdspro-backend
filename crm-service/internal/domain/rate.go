package domain

import (
	_models "common/models"
	"time"
)

type Rate struct {
	_models.BaseEntity
	Anonymous  bool                   `gorm:"not null" json:"anonymous"`
	OwnerID    uint64                 `gorm:"not null" json:"ownerId" binding:"required"`
	OwnerOf    uint32                 `gorm:"default:10;not null" json:"ownerOf" binding:"required"` // 10: user 20: group 30: organization
	Score      uint8                  `gorm:"not null" json:"score" binding:"required,gte=1,lte=5"`
	Comment    string                 `gorm:"type:text" json:"comment,omitempty" binding:"required"`
	ParentID   *uint64                `json:"parentId,omitempty"`
	Attachs    []FeedbackAttachEntity `gorm:"foreignKey:RateID" json:"attachs,omitempty"`
	RootID     *uint64                `json:"rootId,omitempty"`
	Hidden     bool                   `gorm:"default:false" json:"hidden"`
	ApprovedAt *time.Time             `json:"approvedAt"`
	ApprovedBy *uint64                `json:"approvedBy"`
}

func (Rate) TableName() string {
	return "feedback_rates"
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

type FeedbackAttachEntity struct {
	_models.BaseEntity
	RateID   uint64 `gorm:"not null" json:"rateId"`
	FileURL  string `gorm:"type:varchar(500)" json:"fileUrl"`
	FileName string `gorm:"type:varchar(255)" json:"fileName"`
	FileType string `gorm:"type:varchar(50)" json:"fileType"`
}

func (FeedbackAttachEntity) TableName() string {
	return "feedback_rate_attachments"
}