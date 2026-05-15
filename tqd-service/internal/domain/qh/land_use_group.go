package qh_domain

import "time"

type LandUseGroup struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	Code           string    `gorm:"type:varchar(50);uniqueIndex" json:"code"`
	Name           string    `gorm:"type:varchar(100)" json:"name"`
	Description    string    `gorm:"type:text" json:"description"`
	CanBuild       bool      `gorm:"default:false" json:"canBuild"`
	Priority       int       `gorm:"default:50" json:"priority"`
	BuildCondition string    `gorm:"type:text" json:"buildCondition"`
	Color          string    `gorm:"type:varchar(10)" json:"color"`
	IsActive       bool      `gorm:"default:true" json:"isActive"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	UpdatedBy      string    `gorm:"type:varchar(100)" json:"updatedBy"`
}

func (LandUseGroup) TableName() string {
	return "land_use_groups"
}
