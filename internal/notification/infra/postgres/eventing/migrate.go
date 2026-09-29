package eventing

import (
	"time"

	"gorm.io/gorm"
)

type paymentEventInboxRow struct {
	EventID       string    `gorm:"column:event_id;type:varchar(192);primaryKey"`
	EventType     string    `gorm:"column:event_type;type:varchar(128);not null"`
	SchemaVersion int       `gorm:"column:schema_version;not null"`
	OrderID       uint64    `gorm:"column:order_id;not null"`
	PayloadHash   string    `gorm:"column:payload_hash;type:char(64);not null"`
	ReceivedAt    time.Time `gorm:"column:received_at;type:timestamptz;not null;default:now()"`
}

func (paymentEventInboxRow) TableName() string { return "payment_event_inbox" }

type eventNotificationRow struct {
	ID            uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	SourceEventID string    `gorm:"column:source_event_id;type:varchar(192);not null;uniqueIndex"`
	OrderID       uint64    `gorm:"column:order_id;not null"`
	SubjectKind   string    `gorm:"column:subject_kind;type:varchar(32);not null;index:idx_event_notifications_subject,priority:1"`
	SubjectID     string    `gorm:"column:subject_id;type:varchar(128);not null;index:idx_event_notifications_subject,priority:2"`
	Kind          string    `gorm:"column:kind;type:varchar(96);not null"`
	Title         string    `gorm:"column:title;type:varchar(255);not null"`
	Body          string    `gorm:"column:body;type:text;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamptz;not null;default:now();index:idx_event_notifications_subject,priority:3,sort:desc"`
}

func (eventNotificationRow) TableName() string { return "event_notifications" }

func AutoMigrate(database *gorm.DB) error {
	return database.AutoMigrate(
		&paymentEventInboxRow{},
		&eventNotificationRow{},
		&deliveryRow{},
	)
}
