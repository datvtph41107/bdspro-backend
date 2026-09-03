package postgresjob

import "time"

type model struct {
	ID       string `gorm:"column:id;type:varchar(80);primaryKey"`
	ReportID uint64 `gorm:"column:report_id;not null;uniqueIndex"`
	UserID   uint64 `gorm:"column:user_id;not null;index"`

	Operation   string `gorm:"column:operation;type:varchar(128);not null"`
	OperationID string `gorm:"column:operation_id;type:varchar(128);not null;index"`
	CommandKey  string `gorm:"column:command_key;type:varchar(128);not null;index"`

	Status   string `gorm:"column:status;type:varchar(20);not null;index:idx_report_jobs_pending"`
	Attempts int    `gorm:"column:attempts;not null;default:0"`

	AvailableAt  time.Time  `gorm:"column:available_at;type:timestamptz;not null;index:idx_report_jobs_pending"`
	LockedAt     *time.Time `gorm:"column:locked_at;type:timestamptz"`
	LockedBy     string     `gorm:"column:locked_by;type:varchar(100);not null;default:''"`
	ClaimVersion int64      `gorm:"column:claim_version;not null;default:0"`
	LastError    string     `gorm:"column:last_error;type:text;not null;default:''"`

	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
}

func (model) TableName() string {
	return "report_jobs"
}
