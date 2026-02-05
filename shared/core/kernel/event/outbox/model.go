package outbox

import (
	"encoding/json"
	"time"
)

type JSONMap map[string]interface{}

func (j JSONMap) Marshal() ([]byte, error)  { return json.Marshal(j) }
func (j *JSONMap) Unmarshal(b []byte) error { return json.Unmarshal(b, j) }

type OutboxEvent struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	EventType   string    `gorm:"type:varchar(255);index;not null"`
	Version     int       `gorm:"not null;default:1"`
	Payload     JSONMap   `gorm:"type:jsonb;not null"`
	TraceID     string    `gorm:"type:varchar(64);index"`
	Attempts    int       `gorm:"not null;default:0"`
	Processed   bool      `gorm:"not null;default:false;index"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	ProcessedAt *time.Time
}

func (OutboxEvent) TableName() string { return "outbox_events" }
