package dto

import (
	"time"
	"user/internal/enums"
)

type ProfileDTO struct {
	ProfileID         uint64      `gorm:"primaryKey;autoIncrement"`
	FullName          string      `gorm:"size:255"`
	Avatar            string      `gorm:"size:500"`
	BackgroundImage   string      `gorm:"size:500"`
	Phone             string      `gorm:"size:15"`
	Phone2            string      `gorm:"size:15"`
	Email             string      `gorm:"size:255"`
	TaxCode           string      `gorm:"size:50"`
	Address           string      `gorm:"size:500"`
	Position          string      `gorm:"size:255"`
	DepartmentID      *uint64     `gorm:"column:department_id"`
	Birth             *time.Time  `gorm:"column:birth"`
	StartDate         *time.Time  `gorm:"column:start_date"`
	FrontIdentify     string      `gorm:"size:500"`
	BackIdentify      string      `gorm:"size:500"`
	Slogan            string      `gorm:"size:500"`
	Website           string      `gorm:"size:255"`
	Facebook          string      `gorm:"size:255"`
	Instagram         string      `gorm:"size:255"`
	Twitter           string      `gorm:"size:255"`
	Linkedin          string      `gorm:"size:255"`
	Youtube           string      `gorm:"size:255"`
	Introduction      string      `gorm:"type:text"`
	Visibility        uint8       `gorm:"default:1"`
	ProfileVisibility uint8       `gorm:"default:1"`
	StatusOnline      uint8       `gorm:"default:1"`
	TabDefault        uint8       `gorm:"default:1"`
	RoleRealEstate    uint8       `gorm:"default:10"`
	PlanID            *uint64     `gorm:"column:plan_id"`
	PlanAt            *time.Time  `gorm:"column:plan_at"`
	ReferralCode      string      `gorm:"size:255"`
	RoleType          enums.ERole `gorm:"default:10"`
	CreatedAt         time.Time   `gorm:"autoCreateTime"`
	UpdatedAt         time.Time   `gorm:"autoUpdateTime"`
	DeletedAt         *time.Time  `gorm:"column:deleted_at"`
}
