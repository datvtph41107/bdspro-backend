package entity

// InteractiveEvent records a user interaction (screen view, click, etc.).
type InteractiveEvent struct {
	StartTime int64   `gorm:"not null;index"`
	EndTime   *int64  `gorm:"index"`
	Event     string  `gorm:"type:varchar(30);not null;index"`
	RefID     uint64  `gorm:"index"`
	Screen    *string `gorm:"type:varchar(40)"`
	Duration  *int64  `gorm:""`
	ProfileID *uint64 `gorm:"index"`
}

func (InteractiveEvent) TableName() string { return "interactive_events" }
