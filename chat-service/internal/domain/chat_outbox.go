package domain

import "time"

type ChatOutBox struct {
	ID            uint64 `gorm:"primaryKey;index:idx_outbox_status;priority:2"`
	AggregateID   uint64 `gorm:"not null"`
	AggregateType string `gorm:"varchar(50); not null"`
	EventType     string
	Payload       string
	Status        uint8 `gorm:"index:idx_outbox_status;priority:1"`
	CreatedAt     time.Time
	SentAt        time.Time
}
