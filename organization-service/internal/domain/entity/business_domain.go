package entity

import "time"

type BusinessDomain struct {
	ID          uint32 `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null;index"`
	Description *string
	Code        string    `gorm:"not null;uniqueIndex"`
	IsActive    bool      `gorm:"default:true;index"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	CreatedBy   uint32    `gorm:"index"`
	UpdatedBy   uint32
	DeletedAt   *time.Time `gorm:"index"`
	IsDeleted   bool       `gorm:"default:false"`
}

func (BusinessDomain) TableName() string {
	return "business_domains"
}
