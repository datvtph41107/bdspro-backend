package domain

import (
	_models "common/domain/entity"
	"fmt"
	"time"
	"tqd/internal/enums"
)

type OpenHour struct {
	_models.BaseEntity
	Name      string  `json:"name" gorm:"column:name;type:varchar(255)"`
	POIID     *uint64 `json:"poiId,omitempty" gorm:"column:poi_id;index"`
	DayOfWeek int     `json:"dayOfWeek" gorm:"column:day_of_week;not null"` // 2=Monday, 3=Tuesday, ..., 8=Sunday
	OpenTime  string  `json:"openTime" gorm:"column:open_time;type:time"`   // TIME format HH:MM:SS
	CloseTime string  `json:"closeTime" gorm:"column:close_time;type:time"` // TIME format HH:MM:SS
	Note      string  `json:"note" gorm:"column:note;type:text"`
	IsOpen    bool    `json:"isOpen" gorm:"column:is_open;default:true"`
	// lặp lại
	IsRecurring bool                 `json:"isRecurring" gorm:"column:is_recurring;default:true"`
	StartDate   *time.Time           `json:"startDate" gorm:"column:start_date"`
	EndDate     *time.Time           `json:"endDate" gorm:"column:end_date"`
	Status      enums.OpenHourStatus `json:"status" gorm:"column:status;default:0"`
	Type        enums.OpenHourType   `json:"type" gorm:"column:type;default:0"`
}

func (OpenHour) TableName() string {
	return "open_hours"
}

func (oh *OpenHour) Validate() error {
	if oh.POIID == nil || *oh.POIID == 0 {
		return fmt.Errorf("poi id is required")
	}
	if oh.DayOfWeek < 2 || oh.DayOfWeek > 8 {
		return fmt.Errorf("day of week must be between 2 (Monday) and 8 (Sunday)")
	}
	if oh.OpenTime == "" {
		return fmt.Errorf("open time is required")
	}
	if oh.CloseTime == "" {
		return fmt.Errorf("close time is required")
	}
	return nil
}
