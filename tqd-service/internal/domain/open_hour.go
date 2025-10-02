package domain

import (
	_models "common/domain/entity"
	"tqd/internal/enums"
)

// OpenHour represents the operating hours
type OpenHour struct {
	_models.BaseEntity
	Name      string               `json:"name" gorm:"size:255"`
	POIID     *uint64              `json:"poiId,omitempty" gorm:"index"`
	DayOfWeek int                  `json:"dayOfWeek" gorm:"not null"`  // 2=Monday, 3=Tuesday, ..., 7=Saturday, 8=Sunday
	OpenTime  string               `json:"openTime" gorm:"type:time"`  // TIME format
	CloseTime string               `json:"closeTime" gorm:"type:time"` // TIME format
	Note      string               `json:"note" gorm:"type:text"`
	IsOpen    bool                 `json:"isOpen" gorm:"default:true"`
	Status    enums.OpenHourStatus `json:"status" gorm:"default:0"`
	Type      enums.OpenHourType   `json:"type" gorm:"default:0"`
}

// TableName specifies the table name for OpenHour
func (OpenHour) TableName() string {
	return "open_hour"
}
