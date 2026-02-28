package domain

import "time"

type Region struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;size:255" json:"name"`
	Code      uint32    `gorm:"column:code;size:50" json:"code"`
	Active    bool      `gorm:"column:active;default:true;index" json:"active"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (Region) TableName() string {
	return "regions"
}
