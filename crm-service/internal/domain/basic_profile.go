package domain

import "time"

type BasicProfile struct {
	ProfileID  uint64
	FullName   string
	Avatar     string
	Phone      string
	Email      string
	VerifiedAt *time.Time

	EntityID uint64 `gorm:"column:entity_id" json:"entityId"`
	Type     uint32
}

func (BasicProfile) TableName() string {
	return "basic_profile"
}